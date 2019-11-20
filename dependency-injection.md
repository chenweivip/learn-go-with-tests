# 依赖注入

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/di/v1)**

学习本章之前，我们假定你已经阅读过结构体(structs)章节，因为要理解依赖注入必须先理解接口。

在编程社区里，大家对依赖注入有很多误解。希望本章可以澄清:

* 依赖注入并不一定需要框架
* 它也不会让你的设计变复杂
* 它有助于测试
* 它能帮助你写出通用的函数

我们想要写一个问候某人的函数，正如我们之前在hello-world章节做的那样，但是这次我们要测试的是**实际打印输出的部分**。

简单回忆一下，这个函数应该长成这样:

```go
func Greet(name string) {
    fmt.Printf("Hello, %s", name)
}
```

该如何测试呢？`fmt.Printf`会打印到控制台，但是我们很难捕获控制台的输出，然后用测试框架对其进行测试。

我们希望能够注入(**inject**)打印依赖，这里的注入其实就是传入。

我们的函数并不需要关心打印发生在哪里，或者是如何打印的，所以我们应当接受一个接口，而不是一个具体类型。

如果我们使用接口的话，我们就可以改变具体实现。在测试的时候，用一种可控的能够测试的实现。在真实环境中，再换成另外一种可以输出到控制台的实现。

如果你看下`fmt.Printf`的源码，你可以学习到这种让我们可以hook in具体实现的方式:

```go
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}
```

有意思的是，`Printf`底层只是调用了`Fprintf`，并传入了一个`os.Stdout`。

`os.Stdout`到底是什么？`Fprintf`期望传入的第一个参数到底是什么？

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

实际上是一个 `io.Writer`

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

随着你写的Go语言代码越来越多，你经常会碰到这个接口，因为它是一个通用接口，表示:"将数据写到某处"。

现在你知道，在底层你最终会使用`Writer`将我们的问候语发送到某处。现在我们可以利用这种抽象，让我们的代码变得易于测试和重用。

## 先写测试

[`di_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/di/v1/di_test.go)

```go
func TestGreet(t *testing.T) {
    buffer := bytes.Buffer{}
    Greet(&buffer,"Chris")

    got := buffer.String()
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

来自`bytes`包的`buffer`类型实现`Writer`接口。

在该测试中，我们让这个`buffer`成为我们的`Writer`，在调用`Greet`之后，我们就可以检查其中的内容。

## 写程序逻辑

类型`bytes.Buffer`是实现`Writer`接口的，这样我们就可以把问候语输出到一个`buffer`中。`fmt.Fprintf`和`fmt.Printf`类似，只是`fmt.Fprintf`接受的参数是一个`Writer`，而`fmt.Printf`缺省输出到标准控制台(stdout)。

```go
func Greet(writer *bytes.Buffer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

现在测试可以通过。

## 重构

上面的实现要求我们传一个`bytes.Buffer`类型的指针，这种做法技术上正确，但并不通用。

为了演示这个问题，可以尝试在主程序中调用一次`Greet`函数，传入一个`os.Stdout`:

```go
func main() {
    Greet(os.Stdout, "Elodie")
}
```

运行`go run di.go`会看到如下错误提示:

`./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet`

正如之前提到过的，`fmt.Fprintf`允许我们传入`io.Writer`接口，我们也知道`os.Stdout`和`bytes.Buffer`都实现这个接口。

如果我们修改代码，使用更通用的接口，那么我们的测试和主程序就都可以通过了。

[`di.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/di/v1/di.go)

```go
package main

import (
    "fmt"
    "os"
    "io"
)

func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
    Greet(os.Stdout, "Elodie")
}
```

## 关于io.Writer的更多内容

我们用`io.Writer`还可以将数据写到其它地方去吗？我们的`Greet`函数有多通用？

### Web服务器

运行下面的代码:

[`di.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/di/v2/di.go)

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
    Greet(w, "world")
}

func main() {
    http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
```

然后浏览器访问[http://localhost:5000](http://localhost:5000)，可以看到`Greet`函数输出在网页中。

后面章节我们会进一步讲HTTP服务器，所以目前不必纠结实现细节。

当你创建一个HTTP handler，入参有一个`http.ResponseWriter`和一个`http.Request`，其中`http.ResponseWriter`用于输出内容，`http.Request`用于获取用户请求。在我们的实现中，我们将问候语写入`http.ResponseWriter`实例。

你可能已经猜到了，`http.ResponseWriter`也是实现`io.Writer`接口的，所以我们可以在handler中重用`Greet`函数。

## 总结

刚开始我们的代码不太好测试，因为我们把数据写到不受我们控制的地方。

为了让代码易于测试，我们重构了代码，使用依赖注入的方式，让我们可以控制数据写到何处。**依赖注入(dependency injection)**让我们:

* **易于测试代码**，如果某个函数不易测试，通常是因为函数中硬编码了某种依赖或全局状态。如果我们的服务层使用了一个全局的数据库连接池，那么测试就不太容易，测起来也很慢。DI让我们能够注入数据库依赖(通过接口)，然后你可以mock掉某种在测试中无法控制的依赖。
* **关注分离(Separation of concerns)**，将**数据的去处**和**如何产生数据**两者进行解耦。如果你感觉一个方法/函数承担了太多职责(例如既产生数据也写到数据库，或者既处理HTTP请求也做业务逻辑处理)，那么DI可能是你需要利用的工具。
* **让代码易于重用**，首先是让代码可以在测试中重用，进一步如果想尝试某种新的实现，你就可以采用DI将新实现作为依赖进行注入。

### 关于Mock测试

后续章节我们会涉及Mocking。你可以用mock来替代代码功能，运行测试时，注入一个mock(假的)版本，这个mock是你可以控制和测试的。

### 关于Go语言标准库

对于`io.Writer`有一些熟悉之后，你就学会了在测试中可以使用`bytes.Buffer`作为`Writer`，然后在一个命令行程序或者一个Web服务器中，可以使用其它的`Writer`实现。

随着你对Go语言标准库的熟悉，你会越来越多的见到这类通用的接口。然后你可以在你的代码中，通过接口尽量重用标准库的功能，同时让你的软件更易于重用。

本章案例参考了书籍[Go程序设计语言](http://product.dangdang.com/25072202.html)，如果你想学习更多，推荐购买这本Go语言的经典书。

