# 循环

**[本章代码](https://github.com/quii/learn-go-with-tests/tree/master/for)**

在Go语言中实现循环，你可以用关键字`for`。Go语言中没有`while`, `do`, `until`这些关键字，你只能用`for`，所谓少即是多，这反而是好事！

我们要实现一个函数，将某个字符重复5次输出，我们还是从先写测试开始。

## 先写测试

创建目录for/v1，其中创建测试文件`repeat_test.go`:

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}
```

## 写代码逻辑让测试通过

`for`的语法我们大都很熟悉，遵循类似C语言的语法。

在for/v1目录中创建代码文件[`repeat.go`](https://github.com/quii/learn-go-with-tests/blob/master/for/v1/repeat.go)。

```go
func Repeat(character string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

和其它语言如C/Java/JavaScript不同，Go语言中的`for`语句不需要用括号把三部分(初始值/判断/递增)括起来。并且，循环体的花括号{ }是必须的。下面这句你可能第一次看到：

```go
    var repeated string
```

之前我们用过 `:=` 可以同时声明和初始化变量。但是 `:=` 只是[同时做声明+初始化的简写](https://gobyexample.com/variables)，因为值的类型可以自动推导出来，所以可以简写。本例中，我们需要声明一个变量，无法自动推导类型，所以我们只能用非简写版本， var和string都不能少。`var` 也可以用来声明函数，这个我们后面会再讲。

现在运行测试，确保可以通过。

关于for循环的更多变体，可以参考[这里](https://gobyexample.com/for)。

## 重构

我们来对代码进行一点优化，引入一个常量，同时加法操作可以用 `+=` 运算符进行简化，创建v2版本[`repeat.go`](https://github.com/quii/learn-go-with-tests/blob/master/for/v2/repeat.go)：

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

`+=` 是加并且赋值操作符连写，将右边的操作数加到左边的操作数上，然后将值再赋给左边的操作数。这种方式也适用于其它类型，比如整数。

### 性能比对测试Benchmarking

Go语言中[内置支持性能比对测试](https://golang.org/pkg/testing/#hdr-Benchmarks)，并且方法和写测试非常类似。

在[`repeat_test.go`](https://github.com/quii/learn-go-with-tests/blob/master/for/vx/repeat_test.go)中添加性能比对测试函数`BenchmarkRepeat`：

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

可以看到这个函数的写法和测试很像。

这个函数带一个参数`b *testing.B`，它是一个钩子参数，可以让我们对接性能测试框架，比如函数中我们用到的 `b.N`，可以获取性能测试的循环次数。

当性能比对测试执行的时候，循环会运行 `b.N` 次，并且测试框架会测量执行时间。

你无需关心具体执行的次数，测试框架会帮你挑选最优值，保证在恰当的时间内运行足够测试。

现在运行命令 `go test -bench=.` 执行性能比对测试。(如果你在Windows Powershell终端下，请用`go test -bench="."`)： 

```text
goos: darwin
goarch: amd64
pkg: github.com/spring2go/learn-go-with-tests/for/vx
7439353           158 ns/op
PASS
```

`158 ns/op`的意思是，在我的本地机器上，Repeat函数平均需要花费136纳秒运行。这个性能是OK的！测试框架显示它运行了7439353次，再计算出平均运行时间。

**注意**，缺省情况下，性能比对测试以顺序方式运行。

## 扩展练习

* 修改Repeat函数，让调用方能够指定重复次数，然后修改代码和测试，并测试通过。
* 为你的函数添加`ExampleRepeat` 样例作为文档。
* 过一下[strings](https://golang.org/pkg/strings)标准库。找一两个你感兴趣的函数，然后为这些函数写测试(包括性能比对测试)。在标准库上多花时间学习是很有益的。

## 总结

* 进一步实践TDD，
* 学习 `for` 循环，
* 学习如何写性能比对测试benchmarks。
