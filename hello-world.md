# Hello, World

**[本章所有代码](https://github.com/spring2go/learn-go-with-tests/tree/master/hello-world)**

大家常用Hello, world作为学习新编程语言的第一个例子，本课也不例外。

根据Go语言的惯例，请在如下位置创建目录：`$GOPATH/src/github.com/{your-user-id}/learn-go-with-tests/hello-world`，比方说在波波的机器上，创建指令为：`mkdir -p $GOPATH/src/github.com/spring2go/learn-go-with-tests/hello-world`。后续章节我们会沿用这个惯例。

在hello-world目录中创建v1版本[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v1/hello.go)文件：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

然后运行这个程序，打开一个终端窗口，进入hello-world目录，运行命令`go run hello.go`，可以看到Hello, world输出。


## 这个程序是如何工作的？

一个Go语言程序，需要一个主入口程序，这个主入口程序必须声明在`main`包中，并且必须在其中定义一个`main`函数。在Go语言中，包(package)是一种代码的组织单位。

`func`关键字用于定义一个函数，函数有名字和函数体。

`import "fmt"`的作用是导入名称为`fmt`的包，程序中我们用到了`fmt`包中的`Println`打印函数，用于输出Hello, world。

## 如何写测试？

下面我们要来测试这个程序，为了让这个程序容易测试，我们需要调整一下代码，把输出Hello, world的功能独立出来，修改后的v2版[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v2/hello.go)程序如下：

```go
package main

import "fmt"

func Hello() string {
    return "Hello, world"
}

func main() {
    fmt.Println(Hello())
}
```

我们单独创建一个名称为Hello的函数，函数声明以关键字`func`标示，并且需要标示说明这个函数的返回类型，这里是`string`，也就是说这个函数必须返回一个`string`类型的值。

下面我们可以来创建测试文件了，新建文件[`hello_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v2/hello_test.go)，这个程序用来测试`Hello`函数，代码如下：

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    expected := "Hello, world"

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }
}
```

直接在终端运行`go test`运行测试，这个测试应该可以通过，你可以尝试修改`expected`的值，故意让测试不通过。

**注意**：和其它语言如Java/C不同，Go语言中的每行语句后面是不需要加标点的。

可以看到，Go语言是内置支持测试的，不需要蛮烦安装其它测试工具。

### 测试规范

写测试和写函数类似，只需遵循一些规范：

-   文件名必须以`_test.go`结尾，例如`xxx_test.go`。
-   测试函数名必须以`Test`打头。
-   测试函数接受且仅接受一个参数`t *testing.T`。

关于`t *testing.T`这个参数，目前你只需要知道这是对接测试框架的钩子(hook)参数，有了这个参数，你可以调用测试框架的方法，比如在测试失败时调用`t.Fail()`，让测试框架处理失败。

下面有一些Go语言的语法说明：

#### `if`

If条件语句和其它语言类似，没必要多说。

#### 声明变量+赋值

`varName := value`是Go语言中的简写的变量声明+赋值。例如上面的`got := Hello()`，写全的话可以写成`var got string = Hello()`。Go语言的赋值语句具有自动类型推导能力，根据后面的赋值，Go语言可以自动推导出got是`string`类型，所以变量类型`string`可以省略，前面语句也可以写成var get = Hello()。把`var`再省略就可以写成got := Hello()。因为简写输入最少，我们在变量声明+赋值场合基本上都用简写。

#### `t.Errorf`

`t.Errorf`表示调用`t`上的_方法_`Errorf`，也就是测试失败时输出一个消息。`Errorf`中的`f`表示对参数进行格式化输出，可以把参数插入到格式化字符串的占位符(比如`%q`)部分，`%q`占位符参数可以把参数以双引号括起来。关于占位符的更多内容，请参考官方文档[fmt go doc](https://golang.org/pkg/fmt/#hdr-Printing)。


关于方法和函数的区别，后续我们会进一步展开。

### 关于Go doc

细致的文档是Go语言的一大特色。你可以在本地开启Go语言的文档服务，运行命令`godoc -http :8000`，然后浏览器访问[localhost:8000/pkg](http://localhost:8000/pkg)，就可以浏览本地安装的Go语言所支持的标准库。

大部分Go语言的标准库都有不错的文档，而且还带样例。通过浏览器访问[http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/)，你就可以详细查看Go语言所支持的测试功能。

如果你在本地无法运行`godoc`命令，那么可能你安装的是较新版本的Go语言(1.13之后)，新版本[不再自动包含`godoc`](https://golang.org/doc/go1.13#godoc)。你可以手工安装，运行命令`go get golang.org/x/tools/cmd/godoc`即可，注意保持网络可访问！

### 定制人名

下面我们要扩展一下程序的功能，不再简单输出Hello, world，而是能够根据给定的人名输出，例如，输入"Bobo"，就输出"Hello, Bobo“。

按照测试驱动开发方法学的要求，我们先写测试代码[hello_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v3/hello_test.go):

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Bobo")
    expected := "Hello, Bobo"

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }
}
```

然后运行测试 `go test`：

```text
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

显然，测试无法运行，程序通不过编译，因为我们还没有修改Hello函数支持这个测试，修改v3版[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v3/hello.go)，支持定制人名：

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

### 常量

Go语言中的常量定义方式如下：

```go
const englishHelloPrefix = "Hello, "
```

下面我们重构一下代码v4版[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v4/hello.go)

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    return englishHelloPrefix + name
}
```

所谓**重构(refactoring)**，就是在不改变程序功能逻辑的情况下，调整优化程序代码。在程序开发中，把经常用到的相同字符串常量化，是提升代码可读性和可维护性的一种最佳实践。经过上面的重构，再次运行测试，确保重构后程序逻辑正确。

## 新需求

下面我们要再次完善程序，当Hello函数输入为空字符串的时候，我们希望输出"Hello World"，而不是"Hello, "。

我们先写测试：

```go
func TestHello(t *testing.T) {

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Bobo")
        expected := "Hello, Bobo"

        if got != expected {
            t.Errorf("got %q expected %q", got, expected)
        }
    })

    t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"

        if got != expected {
            t.Errorf("got %q expected %q", got, expected)
        }
    })

}
```

细心学员会发现，上面的测试代码有冗余，我们可以通过引入子测试(subtest)来重构优化代码。所谓子测试，其实就是公共可重用的测试逻辑。

按如下方式重构测试代码：

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, expected string) {
        t.Helper()
        if got != expected {
            t.Errorf("got %q expected %q", got, expected)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        expected := "Hello, Chris"
        assertCorrectMessage(t, got, expected)
    })

    t.Run("empty string defaults to 'World'", func(t *testing.T) {
        got := Hello("")
        expected := "Hello, World"
        assertCorrectMessage(t, got, expected)
    })

}
```
上面的代码中，我们把断言逻辑抽取到一个子测试函数`assertCorrectMessage`中，这样可以提升重用度、代码可读和可维护性。Go语言支持在某个函数中再书写子函数(也叫闭包函数)，然后在函数中可以调用子函数。我们把参数`t *testing.T`传给`assertCorrectMessage`，这样就可以在子函数中访问测试框架的方法，比如错误输出。

子函数中的`t.Helper()`方法告知测试框架在错误输出时，输出调用`assertCorrectMessage`语句的行号，而不是`assertCorrectMessage`子函数内的行号，这样可以方便开发人员跟踪问题。如果你还不理解`t.Helper()`的作用，可以故意修改测试让它失败，然后分别注释或者不注释`t.Helper()`这句，看看效果体会一下。

显然，现在就运行测试会失败，因为我们还没有调整Hello函数的逻辑，调整代码v5版本[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v5/hello.go)，添加一个`if`条件判断，如下

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

现在运行测试，确保测试可以通过。


### 测试驱动开发的纪律

测试驱动开发(Test Driven Development，简称TDD)，是一种现代敏捷软件开发方法学。一个典型的TDD流程包含如下步骤：

1. 写测试
2. 写程序逻辑
3. 运行测试，调整程序逻辑，直到测试通过
4. 重构(包括程序和测试）

TDD的核心逻辑是**缩短反馈环**，要点是每次都写少量程序逻辑，通过测试快速获得反馈。这种方法虽然看起来前期要多花一些时间写测试，但是可以提升代码质量和可维护性，中长期反而可以提升开发效率。尤其在你后续需要重构的时候，已有的测试代码可以保障你快速重构。如果你一开始忽略测试，虽然短期看可以更快更多写代码，但是随着代码越堆越多，长期代码不可维护，难以重构。

**注意**，实际开发中，第1～2步没有严格顺序要求，可以先写测试，再写程序逻辑，也可以倒过来，先写程序逻辑，再写测试，大部分程序员倾向后者。顺序并不重要，重要的是通过测试快速获取反馈。

## 再来一个新需求

下面再来一个新需求，我们的Hello world程序，除了要支持英文(并且缺省支持的是英文)，现在还要求支持中文。

添加一个测试：

```go
    t.Run("in Chinese", func(t *testing.T) {
        expected := Hello("波波", "Chinese")
        want := "你好, 波波"
        assertCorrectMessage(t, got, expected)
    })
```
修改v6版[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v6/hello.go)，支持中文(缺省英文)：

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == "Chinese" {
        return "你好, " + name
    }

    return englishHelloPrefix + name
}
```

确保测试可以通过。

再次通过重构优化，提取常量字符串：


```go
const chinese = "Chinese"
const englishHelloPrefix = "Hello, "
const chineseHelloPrefix = "你好, "

func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == chinese {
        return chineseHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

再运行时测试，确保重构正确。


### 支持法语

沿用之前测试驱动开发步骤

-   先写测试
-   修改程序逻辑
-   运行测试，调整程序逻辑，直到测试通过

修改v6版[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v6/hello.go)，Hello函数逻辑如下：

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == chinese {
        return chineseHelloPrefix + name
    }

    if language == french {
        return frenchHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

## `switch`语句

`if`语句过多，是程序复杂上升和可维护性下降的一个信号，我们可以通过`swich`语句进行重构，`switch`语句可以提升代码的可维护性和扩展性(假设我们后面要支持更多语言)。


```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    prefix := englishHelloPrefix

    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

重构完成，再次执行测试，确保通过。现在，学员应该对TDD方法有直观感受了，通过TDD，既可以保证我们代码的质量，同时可以提升我们的开发效率，当我们要实现新的需求，比如让我们的Hello, world支持一种新语言，我们更快速开发和交付功能。

### 最后一个重构

你可能觉得我们的Hello函数变得有点大了，好的，我们可以通过取出子函数进行重构，修改v8版本[`hello.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/hello-world/v8/hello.go):

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return
}
```

注意这里引入了一些新的Go语言语法：

1. 在函数`greetingPrefix`的签名中，我们使用了**具名返回值**`(prefix string)`，它会在函数中创建一个叫`prefix`的变量，另外：
	1. 这个变量缺省为“零”值，具体要看类型，`int`整型的话零值就是0，字符串的话零值就是空字符串“”。
	2. 函数返回时可以简写成`return`，相当于`return prefix`。
	3. 具名变量会显示在Go Doc中，可更清晰说明代码意图。
2. `switch`语句中，如果所有`case`语句都不匹配，就会走`default`分支。
3. `greetingPrefix`函数以小写字母打头，根据Go语言中的惯例，公共函数以大写字母打头，私有函数以小写字母打头。`greetingPrefix`是内部私有的，所以小写打头。


## 总结

没想到一个小小的`Hello, world`程序和TDD结合，可以衍生出这么多内容！现在你应该要理解和掌握：

### 一些Go语言语法

-   如何写测试
-   声明函数，包括函数参数和返回值
-   `if`, `const` 和 `switch` 的用法
-   声明变量和常量

### TDD流程和重要性

- 记住TDD的核心逻辑是**缩短反馈环**，要点是每次都写少量程序逻辑，通过测试快速获得反馈。短期看，TDD有一点开销，但是长期TDD可以显著提升软件质量和交付效率。TDD是开发人员必备技能。

