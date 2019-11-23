# 选择Select

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/select)**

假设你接到一个需求，写一个函数`WebsiteRacer`，该函数接收两个URL，对这两个URL发起调用，看哪一个先返回。如果两者都不能在10秒以内返回，就返回一个`error`。

为了实现这个函数，我们要用到：

- `net/http` 发起HTTP调用。
- `net/http/httptest` 创建测试用HTTP服务器.
- goroutines.
- `select` 同步进程.

## 先写测试

先从最简单开始：

```go
func TestRacer(t *testing.T) {
    slowURL := "http://www.facebook.com"
    fastURL := "http://www.quii.co.uk"

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

这个测试不能满足我们的要求，但是我们不追求一开始就完美，而是循序渐进达成目标。

## 写程序逻辑

```go
func Racer(a, b string) (winner string) {
    startA := time.Now()
    http.Get(a)
    aDuration := time.Since(startA)

    startB := time.Now()
    http.Get(b)
    bDuration := time.Since(startB)

    if aDuration < bDuration {
        return a
    }

    return b
}
```

对于每个URL:

1. 在调用`URL`之前，我们先使用`time.Now()`记录开始时间。
2. 然后我们使用[`http.Get`](https://golang.org/pkg/net/http/#Client.Get)获取这个`URL`的内容。该函数返回一个[`http.Response`](https://golang.org/pkg/net/http/#Response)和一个`error`，但是目前我们并不关心这些值。
3. `time.Since`接受一个开始时间，然后返回当前时间和开始时间之间的差值，类型为`time.Duration`。

之后我们就比较两个间隔时间哪个更快。

### 问题

运行这个测试可能会通过，也可能不会通过。问题在于，为了测试我们的逻辑，我们必须测真实的站点。

基于HTTP进行测试的场景很多，所以Go语言标准库提供了工具，可以帮助我们进行测试。

在mocking和依赖注入章节，我们提到最好不要依赖外部服务来测试我们的代码，因为外部服务：

- 慢
- 不稳定
- 无法测边界条件

在标准库中，有一个[`net/http/httptest`](https://golang.org/pkg/net/http/httptest/)包，它方便我们创建mock HTTP server。

我修改测试使用mocks，这样我们就可以测受控和可靠的服务器。

```go
func TestRacer(t *testing.T) {

    slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))

    fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    expect := fastURL
    got := Racer(slowURL, fastURL)

    if got != expect {
        t.Errorf("got %q, expect %q", got, expect)
    }

    slowServer.Close()
    fastServer.Close()
}
```

语法稍微有点复杂，你需要花点时间消化下。

`httptest.NewServer`接受一个`http.HandlerFunc`类型的参数，这个参数是我们传入的一个**匿名函数**。

`http.HandlerFunc`是一个函数类型，完整声明是这样的：`type HandlerFunc func(ResponseWriter, *Request)`。

这些语句表达的是需要传入一个函数，这个函数接受一个`ResponseWriter`和一个`Request`，对于一个HTTP服务器来说，这些是必须的。

没有什么特别神奇的地方，**如果我们写一个真实的Go语言HTTP服务器，也是这么写的**。唯一的区别是，我们把`HandlerFunc`包裹在一个`httptest.NewServer`中，方便我们做测试。测试时，这个服务器会开启并监听在某个端口上，测试结束后，你可以关闭这个服务器和端口。

在两个mock服务器中，在接收到request时，我们让其中一个`time.Sleep`的时间长一点，另外一个则短一点。最后，两个服务器都通过`w.WriteHeader(http.StatusOK)`的方式返回`OK`响应。

现在再次运行测试，确保测试可以通过，而且测试会更快。建议你故意修改sleep时间，让测试不通过，这样可以让你进一步理解代码如何工作。

## 重构

我们的程序和测试中都有一些冗余，让我们来重构一下：

[racer.go](https://github.com/spring2go/learn-go-with-tests/blob/master/select/v1/racer.go)

```go
func Racer(a, b string) (winner string) {
    aDuration := measureResponseTime(a)
    bDuration := measureResponseTime(b)

    if aDuration < bDuration {
        return a
    }

    return b
}

func measureResponseTime(url string) time.Duration {
    start := time.Now()
    http.Get(url)
    return time.Since(start)
}
```

重构之后，我们的`Racer`代码更易于阅读。

[racer_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/select/v1/racer_test.go)

```go
func TestRacer(t *testing.T) {

    slowServer := makeDelayedServer(20 * time.Millisecond)
    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    expect := fastURL
    got := Racer(slowURL, fastURL)

    if got != expect {
        t.Errorf("got %q, expect %q", got, expect)
    }
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(delay)
        w.WriteHeader(http.StatusOK)
    }))
}
```

我们抽出一个创建mock服务器的函数`makeDelayedServer`，这样测试代码更简洁。

### `defer`

在一个函数调用的前面加一个`defer`关键字，那么这个函数会**在包含它的函数的最后**才执行。

有时我们需要清理资源，例如关闭文件，或者在我们的案例中，我们需要关闭服务器，让它释放端口资源。

我们期望资源清理的动作在函数的最后才执行，但是对资源清理的调用一般写在创建资源之后(而不是函数的最后)，这样代码容易阅读。

重构后的代码有很大改善，基于目前我们掌握的Go语言语法，目前的改善是合理的，但我们还可以让它变得更简单。

### 同步进程

- 既然Go语言支持并发，我们为什么要依次顺序测网站的性能呢？我们应该可以并发测。
- 对于**确切的响应时间**我们并不关心，我们只是想知道哪个先返回。

为了能够并发测，我们需要需要引入一种新的Go语言结构`select`，它可以方便我们对进程进行同步。

[racer.go](https://github.com/spring2go/learn-go-with-tests/blob/master/select/v2/racer.go)

```go
func Racer(a, b string) (winner string) {
    select {
    case <-ping(a):
        return a
    case <-ping(b):
        return b
    }
}

func ping(url string) chan struct{} {
    ch := make(chan struct{})
    go func() {
        http.Get(url)
        close(ch)
    }()
    return ch
}
```

#### `ping`

我们创建了一个函数`ping`，它创建并返回一个`chan struct{}`。

在我们的案例中，我们并不**关心**发送到channel里头的是什么类型，**我们只是需要获得完成信号** ～ 关闭chennel就可以作为完成信号。

为何使用`struct{}`而不是另外一个类型如`bool`？因为从内存分配视角看，`chan struct{}`是最小的数据类型(实际不分配内存)，比`bool`还要小。因为我们只是需要一个完成信号，所以没必要分配内存。

在`ping`函数中，我启动一个goroutine，它会在完成`http.Get(url)`调用之后向channel发送一个信号。

##### 必须用 `make` 创建channels

注意我们必须用`make`创建channel，而不是如`var ch chan struct{}`。当你使用`var`，变量会被初始化为对应类型的"零"值。对于`string`，零值就是`""`，对于`int`，零值就是0，诸如此类。

对于channel，零值是`nil`，当你尝试向`nil`值的channel发送(`<-`)数据，它会永远阻塞，因为你无法向`nil` channel发送数据。

[你可以通过Go Playground校验这个行为](https://play.golang.org/p/IIbeAox5jKA)
#### 选择`select`

回忆一下之前的并发章节，你应该记得我们可以用`myVar := <-ch`的方式，等待从channel中获取值。这是一个**阻塞**调用，因为需要等待其它goroutine向channel中先发送值。

`select`可以让你在**多个**channels上同时等待。第一个有值的channel会赢，对应的这个`case`的代码将被执行。

在`Racer`函数中，我们通过两次调用`ping`创建了两个channel，每个channel对应一个`URL`。只要其中一个`ping`先写入channel，那么对应的`case`下的代码就会被执行，也就是返回快的`URL`（胜者)。

做了这些改变之后，我们的代码变得更清晰，实现也更简单。

### 超时

Our final requirement was to return an error if `Racer` takes longer than 10 seconds.
最后一个需求，如果`Racker`运行超过10秒还没决出胜者，那么就返回一个error。

## 先测测试

```go
t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
    serverA := makeDelayedServer(11 * time.Second)
    serverB := makeDelayedServer(12 * time.Second)

    defer serverA.Close()
    defer serverB.Close()

    _, err := Racer(serverA.URL, serverB.URL)

    if err == nil {
        t.Error("expected an error but didn't get one")
    }
})
```

我们让测试服务器运行慢一点，要超过10秒才返回，这样我们可以测试超时的场景。`Racer`现在返回两个值，胜者的URL(在测试中我们忽略)和一个`error`。

## 写程序逻辑

```go
func Racer(a, b string) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(10 * time.Second):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

在使用`select`时，`time.After`是一个非常有用的函数。如果不用`time.After`，我们的代码有可能永远阻塞，因为两个channels可能永远无法接收到值。用了`time.After`，可以保证我们的代码始终会返回，如果在规定时间内两个channels都没有拿到值，那么`time.After`会超时返回。

这样我们程序逻辑就很清晰：如果`a`或者`b`任意一个在规定时间内返回，那么先返回的是胜者；否则到达10秒，`time.After`会发出超时信号，程序返回一个`error`。

### 慢测试

对于这样一个逻辑很少的程序，测试要花10秒运行显然是很慢的。

我们可以让超时变得可配置。这样在测试中，我们把超时设短一点，在真实场景中，我们再把超时改回10秒。

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

下面是重构后的程序代码：

[racer.go](https://github.com/spring2go/learn-go-with-tests/blob/master/select/v3/racer.go)

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
    return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

重构的代码中，`Racer`使用缺省的(需求规定的)超时时间，`ConfigurableRacer`则可以配置超时时间，`Racer`内部间接使用`ConfigurableRacer`。程序的用户和第一个测试用例，会使用`Racer`，而超时测试用例会使用`ConfigurableRacer`。

下面是重构后的测试：

[racer_test).go](https://github.com/spring2go/learn-go-with-tests/blob/master/select/v2/racer_test.go)

```go
func TestRacer(t *testing.T) {

    t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
        slowServer := makeDelayedServer(20 * time.Millisecond)
        fastServer := makeDelayedServer(0 * time.Millisecond)

        defer slowServer.Close()
        defer fastServer.Close()

        slowURL := slowServer.URL
        fastURL := fastServer.URL

        expect := fastURL
        got, err := Racer(slowURL, fastURL)

        if err != nil {
            t.Fatalf("did not expect an error but got one %v", err)
        }

        if got != expect {
            t.Errorf("got %q, expect %q", got, expect)
        }
    })

    t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
        server := makeDelayedServer(25 * time.Millisecond)

        defer server.Close()

        _, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("expected an error but didn't get one")
        }
    })
}
```

在第一个测试用例中，我加了一个检查，确保对`Racer`的正常调用不会获得`error`。

## 总结

### `select`

- 让你可以在多个channels上进行等待。
- 有时候，你需要在`case`中加一个`time.After`，防止系统陷入永久阻塞状态。

### `httptest`

- 方便我们创建测试用HTTP服务，让测试可靠和可控。
- 和真实的`net/http`接口一致，学习成本很低。
