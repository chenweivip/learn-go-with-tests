# 整数

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/integers)**

你应该对编程语言的整数操作早有经验，所以我们只需写一个add函数尝试一下。先写测试，创建integers/v1目录，其中创建一个叫[`adder_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/integers/v1/adder_test.go)的测试文件。

## 先写测试

```go
package integers

import "testing"

func TestAdder(t *testing.T) {
    sum := Add(2, 2)
    expected := 4

    if sum != expected {
        t.Errorf("expected '%d' but got '%d'", expected, sum)
    }
}
```

注意这里我们用的占位符是`%d`(之前用过`%q`)，这是专门用于对整数进行格式化的占位符。

注意，本次代码我们只演示整数操作，只需添加简单函数即可，不需要main主入口，也就不需要用main package。我们只需定义一个叫intergers的package，这个包用于组织整数操作/测试相关的代码。

**注意** Go语言中，一个目录里头只能有一个`package`，注意合理组织你的代码文件。[这里](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)有一个关于项目代码组织的参考和建议。

## 写代码逻辑

同样在integers/v1目录下，添加[`adder.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/integers/v1/adder.go)，下面是整数相加Add函数的定义，很简单：

```go
func Add(x, y int) int {
    return x + y
}
```

注意，在函数参数中，如果几个参数的类型相同，可以简写，例如`(x int, y int)`，可以简写成`(x, y int)。

另外，这次我们没有使用**具名返回值**，具名返回值常用在返回值的含义不太清楚的场合中，在本例中，Add函数返回值的含义是显而易见的。关于具名返回值的使用建议，可以参考[这个](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters)

在integers目录总，运行测试`go test`，确保测试通过。

## 重构

代码本身比较简单，没有多少需要改进的地方。

前面我们讲过，函数的具名返回参数会出现在文档中，同时也会出现在开发人员的文本编辑器中。

这点很有用，因为它可以提升代码的可读性。如果只需看代码中的类型签名和文档，就可以很容易读懂代码，显然对开发人员体验来说，是非常友好的。

可以在函数上添加注释型文档，这些文档会出现在Go Doc中。如果你看Go语言标准库的文档，也是这么来的。
```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
    return x + y
}
```

### 样例Examples

现在我们可以更进一步，制作golang语言代码样例[examples](https://blog.golang.org/examples)。你可以看到在golang语言标准库中有很多样例examples。

样例文档做法很多，比如在代码之外单独写readme文档，但是这种做法文档和代码很容易失去同步，因为开发人员在修改代码时，很容易忘却需要更新文档。

Go语言的样例文档可以像测试一样执行，所有文档和样例一般不会失去同步。

作为包测试族的一部分，样例会被编译(也可以被执行)。

和典型的测试一样，样例是在package下_test.go文件中的函数。在[`adder_test.go`](https://github.com/quii/learn-go-with-tests/blob/master/integers/v2/adder_test.go)中添加如下ExampleAdd函数：

```go
func ExampleAdd() {
    sum := Add(1, 5)
    fmt.Println(sum)
    // Output: 6
}
```

（注意在`adder_test.go`中需要添加`import "fmt"`导入fmt包，否则编译会通不过，强烈建议设置你的编译器，支持自动导入包，具体设置方法每种编译器都不一样，请自行研究。)

如果你修改了代码逻辑，但是样例没有同步修改，那么测试会失败。

现在可以运行package里头的测试族，`go test -v`，注意添加`-v`参数显示测试执行细节，我们可以看到样例函数也被执行了：

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

请注意，如果我们把注释 "//Output: 6"去掉，那么样例函数就不会被执行(虽然编译和测试可以通过)，你可以尝试一下。

添加样例example之后，它会出现在`godoc`中，提升你的代码用户体验。

可以试一下，运行`godoc -http=:6060`，然后浏览器访问`http://localhost:6060/pkg/`进行查看。

在godoc中，你可以看到你本地的`$GOPATH`中的所有包。比方说对于波波的机器，我在`http://localhost:6060/pkg/github.com/spring2go/learn-go-with-tests/integers/v2/`中，可以看到本章样例文档。

如果你想公开发布你的代码文档，你可以把文档发布在[godoc.org](https://godoc.org)。例如，本章的最终API发布在[https://godoc.org/github.com/quii/learn-go-with-tests/integers/v2](https://godoc.org/github.com/quii/learn-go-with-tests/integers/v2).



## 总结

本章我们学习到：

* 再次实践TDD流程，
* 整数相加，
* 写好文档，让用户更容易读懂你的代码，
* 如何书写样例Examples文档，它可以作为测试的一部分被执行。
