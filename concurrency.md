# 并发

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/concurrency)**

有这样一个场景：有个同事写了一个函数`CheckWebsites`，该函数检查一组URL的状态。

[CheckWebsites.go](https://github.com/spring2go/learn-go-with-tests/blob/master/concurrency/v1/CheckWebsites.go)

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        results[url] = wc(url)
    }

    return results
}
```

它返回一个map，键是URL，值是对应的检查状态的布尔值 - `true`表示响应正常，`false`表示响应不正常。

你还需要传入一个`WebsiteChecker`，它接受一个单一的URL，并返回一个状态检查布尔值。`CheckWebsites`函数内部使用`WebsiteChecker`检查所有的URLs的健康状态。

利用依赖 [dependency injection][DI]，你可以测试函数，但是不需要真正发起HTTP调用，并且测试稳定且快。

这个是我们已经写的测试：

[CheckWebsites_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/concurrency/v1/CheckWebsites_test.go)

```go
package concurrency

import (
    "reflect"
    "testing"
)

func mockWebsiteChecker(url string) bool {
    if url == "waat://furhurterwe.geds" {
        return false
    }
    return true
}

func TestCheckWebsites(t *testing.T) {
    websites := []string{
        "http://google.com",
        "http://blog.gypsydave5.com",
        "waat://furhurterwe.geds",
    }

    expect := map[string]bool{
        "http://google.com":          true,
        "http://blog.gypsydave5.com": true,
        "waat://furhurterwe.geds":    false,
    }

    got := CheckWebsites(mockWebsiteChecker, websites)

    if !reflect.DeepEqual(expect, got) {
        t.Fatalf("Wanted %v, got %v", want, got)
    }
}
```

这个函数已经上到生产上，用于检查上百个网站。但是你的同事开始陆续收到抱怨，说这个功能太慢了，所以老板找到你，要求优化这个函数的性能。

## 先写测试

我们用比对测试(benchmark)来测试`CheckWebsites`的性能，这样后面我们可以看到优化的效果。

[CheckWebsites_benchmark_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/concurrency/v1/CheckWebsites_benchmark_test.go)

```go
package concurrency

import (
    "testing"
    "time"
)

func slowStubWebsiteChecker(_ string) bool {
    time.Sleep(20 * time.Millisecond)
    return true
}

func BenchmarkCheckWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "a url"
    }

    for i := 0; i < b.N; i++ {
        CheckWebsites(slowStubWebsiteChecker, urls)
    }
}
```

这个比对测试使用100个urls对`CheckWebsites`进行测试，并使用了一个假的`WebsiteChecker`实现。`slowStubWebsiteChecker`是`WebsiteChecker`的一个实现，我们故意让它运行很慢。它用`time.Sleep`等20微秒，然后返回true。

现在运行比对测试 `go test -bench=.` (或者在Windows Powershell环境运行 `go test -bench=".")，你将看到如下输出：

```sh
pkg: github.com/spring2go/learn-go-with-tests/concurrency/v1
BenchmarkCheckWebsites-8               1        2282716423 ns/op
PASS
ok      github.com/spring2go/learn-go-with-tests/concurrency/v1 2.300s
```

比对测试结构，`CheckWebsites`平均运行2282716423纳秒 ～ 大约2.3秒。

下面我们来优化，让它跑更快一点。

### 写程序逻辑

现在我们来讲解并发(concurrency)，我们这里讲的并发是指'同事做多个事情'。其实我们在生活中也经常并发做做事的。

例如，早上起来我泡了一杯茶。我先用水壶烧水，当水还没烧开，我从冰箱取出牛奶，从柜子里取出茶叶，找到我最喜欢的杯子，将茶袋放入杯中，然后，当水壶烧开，我就将开水倒入茶杯。

在水壶还没有烧开的时候，我没有**干等**它烧开，然后再做其它事情，我是并行做的。

如果你能理解如何更快泡茶的原理，那么你就会知道如何让`CheckWebsites`跑得更快。当`CheckWebsites`发送一个请求到某个URL，在响应还没有回来之前，与其让计算机干等，不如让它同时发送下一个请求。

通常在Go语言中，当我们调用一个函数`doSomething()`，我们要等待调用返回(即便这个函数没有返回值，我们也需要等它返回)。我们说这个操作是**阻塞blocking**的 - 它执行完前我们必须等待。在Go语言中，一个不阻塞的操作将运行在一个单独的称为**goroutine**的**进程**中。

为了让Go启动一个新的goroutine，我们用`go doSometing()`这种语法，在函数调用前加一个关键字`go`。

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func() {
            results[url] = wc(url)
        }()
    }

    return results
}
```

启动goroutine的唯一方式是在函数调用之前加`go`，为了减少不必要的函数定义，我们经常使用**匿名函数**来启动一个goroutine。匿名函数从字面上看和普通函数差不多，只是没有函数名。上面的`for`循环体中就是一个匿名函数调用。

匿名函数有一些特性，上面的代码中可以展示两个。首先，它们在声明的同时就可以被执行 - 这就是为什么匿名函数的最后有一个`()`，这个相当于调用。其次，匿名函数可以访问它被声明时的词法作用域 - 在你声明匿名函数时所有可见的变量，在匿名函数的函数体中仍然可见。

上面的代码中，`url`、`wc`和`results`对匿名函数都可见。每次循环都会起一个goroutine，这些goroutine会将结果都存储到`results` map中。

但是当你运行 `go test`，你会得到类似如下错误:

```sh
--- FAIL: TestCheckWebsites (0.00s)
    CheckWebsites_test.go:31: Expected map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false], got map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/spring2go/learn-go-with-tests/concurrency/v2 0.015s

```

### 欢迎进入并发世界

你可能会得到不同的结果。你可能会得到一个panic消息(后面我们会讲)，不用担心，多运行几次你就会得到类似上面的结果。即使你多次运行还是得不到上面结果，那么先**假定**得到了上面的结果。

欢迎来到并发世界：如果你的程序的处理并发的方式不对，那么你可能会得到无法预测的结果。不用担心 - 这就是我们为什么要写测试的原因，它可以帮助我们以可预测的方式来处理并发。

### 回到我们的程序

`CheckWebsites`可能会返回一个空map，怎么回事？

真实原因在于，`for`循环所启动的那些goroutines，可能还来不及将结果添加到`results` map中，主函数`WebsiteChecker`就已经返回了，所以它对外返回了空map。

为了修复这个问题，我们可以让主函数`WebsiteChecker`等待一下，等其它goroutines做完它们的工作，然后再返回。等2秒钟应该够了，对不对？

```go
package concurrency

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func() {
            results[url] = wc(url)
        }()
    }

    time.Sleep(2 * time.Second)

    return results
}
```

现在运行测试，你可能会看到：

```sh
--- FAIL: TestCheckWebsites (0.00s)
    CheckWebsites_test.go:31: Expected map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false], got map[waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/spring2go/learn-go-with-tests/concurrency/v2 0.015s
```

还是有问题 - 为什么只有一个结果？你可能想通过延长等待时间的方式来修复 - 你可以这样尝试，但是没用。问题在于：每次`for`循环迭代，`url`变量会被重用 - 每次迭代`url`会获得一个新值。但是每一个goroutine对`url`变量都有引用(它们没有独立的`url`变量拷贝)。所以所有的goroutines得到的都是最后一次迭代的`url`值。这就是为什么结果显示只有最后一个url。

按如下方式修复:

[CheckWebsites.go](https://github.com/spring2go/learn-go-with-tests/blob/master/concurrency/v2/CheckWebsites.go)

```go
package concurrency

import (
    "time"
)

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func(u string) {
            results[u] = wc(u)
        }(url)
    }

    time.Sleep(2 * time.Second)

    return results
}
```

我们给匿名函数传一个参数 - `u` -，调用匿名函数的时候传入`url`作为参数，确保goroutine获得的`url`值是固定的 - 也就是那次迭代对应的某个url。`u`是`url`的一份拷贝，所以它不变。

现在如果你运气好的话，测试可以通过:

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        2.012s
```

但是如果你运气不好，你还是会碰到出错(这种错误更多出现在比对测试benchmark中，因为比对测试运行次数更多)：

```sh
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x132b923, 0x15)
        /usr/local/Cellar/go/1.13.4/libexec/src/runtime/panic.go:774 +0x72 fp=0xc000047718 sp=0xc0000476e8 pc=0x10305c2
runtime.mapassign_faststr(0x12c4920, 0xc000098d80, 0x132c8f5, 0x17, 0x0)
        /usr/local/Cellar/go/1.13.4/libexec/src/runtime/map_faststr.go:211 +0x417 fp=0xc000047780 sp=0xc000047718 pc=0x1014a47
github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites.func1(0x133dd80, 0xc000098d80, 0x132c8f5, 0x17)
        /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:15 +0x71 fp=0xc0000477c0 sp=0xc000047780 pc=0x1283801
runtime.goexit()
        /usr/local/Cellar/go/1.13.4/libexec/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc0000477c8 sp=0xc0000477c0 pc=0x105f471
created by github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites
        /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:14 +0x95

        ... 更多奇怪的错误 ...
```

错误很长令人害怕，没关系，深呼吸然后阅读stacktace：`fatal error: concurrent map writes`。有的时候，当我们运行测试，两个goroutines会试图同时写入`results` map。Go语言中的Maps不是并发安全的，所以当多个goroutine试图并发写入，会引发`fatal error`。

这里头有一个**竞争条件(race condition)**bug，当软件的输出依赖于一些我们无法控制的时序事件，就可能引发这种bug。因为我们无法控制每个goroutine在何时写入`results` map`，所以潜在可能有竞争条件bug。

Go语言内置支持[_race detector_][godoc_race_detector]，可以帮助我们查找竞争条件。为了启用这个功能，可以在测试时加`race`标记：`go test -race`。

你应该获得类似如下输出：

```sh
==================
WARNING: DATA RACE
Write at 0x00c00009cd80 by goroutine 10:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.13.4/libexec/src/runtime/map_faststr.go:202 +0x0
  github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites.func1()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:15 +0x82

Previous write at 0x00c00009cd80 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.13.4/libexec/src/runtime/map_faststr.go:202 +0x0
  github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites.func1()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:15 +0x82

Goroutine 10 (running) created at:
  github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:14 +0xb0
  github.com/spring2go/learn-go-with-tests/concurrency/v2.TestCheckWebsites()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites_test.go:28 +0x17f
  testing.tRunner()
      /usr/local/Cellar/go/1.13.4/libexec/src/testing/testing.go:909 +0x199

Goroutine 9 (finished) created at:
  github.com/spring2go/learn-go-with-tests/concurrency/v2.CheckWebsites()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:14 +0xb0
  github.com/spring2go/learn-go-with-tests/concurrency/v2.TestCheckWebsites()
      /Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites_test.go:28 +0x17f
  testing.tRunner()
      /usr/local/Cellar/go/1.13.4/libexec/src/testing/testing.go:909 +0x199
==================
```

输出很多很难读懂 - 但是第一行`WARNING: DATA RACE`是很明显的。继续看错误，我们看到两个goroutines试图同时写入map:

`Write at 0x00c00009cd80 by goroutine 10:`

两个试图同时写入一个内存地址：

`Previous write at 0x00c00009cd80 by goroutine 9`

上面我们还可以看到写入发生在哪个文件的哪行:

`/Users/william/go/src/github.com/spring2go/learn-go-with-tests/concurrency/v2/CheckWebsites.go:15`

所有信息都输出在终端上 - 你需要一点耐心，我们继续。

### 通道Channels

我们可以使用**通道channels**来协调goroutines，从而解决数据竞争问题。Channel是Go语言的一种数据结构，它可以同时用于接收和发送数据。Channel是可以协调进程之间协作的同步管道。

在本例中，我们想要协调父进程和每个goroutines(实际调用`WebsiteChecker`完成任务的进程)之间的通讯。

[CheckWebsites.go](https://github.com/spring2go/learn-go-with-tests/blob/master/concurrency/v3/CheckWebsites.go)

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
    string
    bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)
    resultChannel := make(chan result)

    for _, url := range urls {
        go func(u string) {
            resultChannel <- result{u, wc(u)}
        }(url)
    }

    for i := 0; i < len(urls); i++ {
        result := <-resultChannel
        results[result.string] = result.bool
    }

    return results
}
```

除了`result` map，我们还引入了一个`resultChannel`，同样使用`make`进行创建。`chan result`是channel类型 - 一个可以容纳`result`的channel。`result`是我们引入的一个新类型，它有两个字段，一个string用来存要被`WebsiteChecker`所检查的url，一个bool用来存检查的结果。因为我们不需要具体的字段名，所以这两个字段都是匿名的，有时想不出好的命名，就可以采用匿名字段。

现在我们可以对urls进行迭代，在匿名函数中，我们使用通道的**发送表达式**，把每次检查的结果`result`(内含url和检查布尔结果)发送到通道`resultChannel`中，`<-`是发送到通道的操作符，发送时左边是通道，右边是要发送的值。

```go
// Send statement
resultChannel <- result{u, wc(u)}
```

下一个`for`循环从0到len(urls)再迭代一次。在循环体中，我们使用通道的**接收表达式**，将从通道中接收到的值赋给一个变量`result`。接收也是用`<-`操作符，但是这次两个操作数互换位置，通道在右边，要被赋值的变量在左边。

```go
// Receive expression
result := <-resultChannel
```

然后，我们用接收到的`result`更新`results` map。

通过将结果发送到一个通道channel中，我们就可以控制每次写入(和后面读取)的时序。虽然每个goroutine的操作都是并发进行的，但是channel在主函数和goroutines之间建立了一个同步管道，可以同步传递数据。

我们已经优化了代码，让它能运行更快，同时，我们也确保不能并行运行的部分能够正确顺序执行。我们用channels实现了多个进程之间的同步通讯。

现在我们可以运行比对测试benchmark:

```sh
pkg: github.com/spring2go/learn-go-with-tests/concurrency/v3
BenchmarkCheckWebsites-8              50          23166023 ns/op
PASS
ok      github.com/spring2go/learn-go-with-tests/concurrency/v3 1.199s
```
每次操作只用23166023纳秒，也就是0.023秒, 比老版本快了100倍. 优化成功！

## 总结

本章我们并没有讲很多TDD的内容。但其实在对`CheckWebsites`的不断重构中，我们一直都离不开测试(包括功能测试和性能比对benchmark测试)。由于这些测试的存在，我们可以大胆重构`CheckWebsites`，在确保功能正确的同时，还切实提升了其性能。

在性能优化过程中我们学习到：

- **goroutines**, Go语言中的基本并发单位，让我们可以同时检查多个站点。
- **匿名函数**, 让我们可以启动多个并发进程去检查网站。
- **通道channels**, 控制和协调不同进程间的通讯，避免**竞争条件race condition**
- **race detector**，帮助我们检查潜在有竞争条件的代码

### 关于性能优化

关于敏捷软件开发有这样一个说法(常常被误以为是Kent Beck讲的)：

> [Make it work, make it right, make it fast][wrf]
> 
> [先让它能工作，然后让它正确工作，再让它运行更快][wrf]

'work'的意思是说让测试通过，'right'指重构代码，'fast'指优化代码性能。我们只有在'make it work'和'make it right'之后再考虑‘make it fast'。本章案例中，我们拿到的代码就已经可以工作(只是太慢)，所以我们无需重构。但实际开发中，在另外两步还没有实现之前，我们绝不应该试图'make it fast'。

> [Premature optimization is the root of all evil][popt]
> 
> [过早的优化是万恶之源](popt]
> 
> -- Donald Knuth

[DI]: dependency-injection.md
[wrf]: http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast
[godoc_race_detector]: https://blog.golang.org/race-detector
[popt]: http://wiki.c2.com/?PrematureOptimization
