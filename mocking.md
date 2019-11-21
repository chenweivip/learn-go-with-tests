# Mocking

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/mocking)**

假设你得到一个新需求: 从3开始倒数到1，每个数打印一行(每次隔1秒)，数到0的时候打印"Go!"，然后退出。

```
3
2
1
Go!
```

为解决这个问题，我们会写一个函数`Countdown`，在`main`程序里头调用，如下:

```go
package main

func main() {
    Countdown()
}
```

虽然这是一个很简单的程序，但我们仍然会采用**增量式测试驱动**方法。

增量的意思是～每次都做一小步，但每小步都能生产出有用的软件。

如果每次都花大量时间开发代码，中间没有测试反馈，那么开发人员很容易掉入一个陷阱～看起来代码写得很快很多，但是让这些代码真正工作需要花费大量的hacking和调试时间，而且后续代码的可维护性差。**将需求分解为足够小的步骤，并且每一步都能生产出有用的软件，这种技能是非常重要的。**

下面是我们计划的分解和迭代步骤:

- 打印 3
- 打印 3, 2, 1 和 Go!
- 每行间隔1秒

## 先写测试

Our software needs to print to stdout and we saw how we could use DI to facilitate testing this in the DI section.
我们的软件要求输出到stdout。在之前的DI章节，我们学习过如何使用DI简化测试。

[countdown_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v1/countdown_test.go)

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    expected := "3"

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }
}
```

如果你对`buffer`还不熟悉，那么请先读[前面一章](dependency-injection.md)。

我们知道我们的`Countdown`函数要将结果写到某处，在Go语言中，`io.Writer`就是能够捕获这种功能的接口。

- 在`main`程序种，我们将结果输出到`os.Stdout`，这样用户就可以在终端上看到countdown的结果。
- 在测试中，我们将结果输出到`bytes.Buffer`，这样我们的测试就可以捕获并测试输出的结果。

## 写程序逻辑

```go
func Countdown(out *bytes.Buffer) {
    fmt.Fprint(out, "3")
}
```

我们使用了`fmt.Fprint`，它接受一个`io.Writer`接口(`*bytes.Buffer`遵循这个接口)，并将一个`string`写入到这个接口。现在测试可以通过。

## 重构

虽然`*bytes.Buffer`可以工作，我们最好使用更通用的接口。

```go
func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}
```

再次运行测试，应该还是可以通过。

下面是完整的主程序，我们在`main`中也调用了`Countdown`，这样我们的主程序也可以工作 ～ 我们小步行进，但是每一步都有可以工作的软件。

[countdown.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v1/countdown.go)

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}

func main() {
    Countdown(os.Stdout)
}
```

运行主程序`go run main.go`，确保主程序也可以工作。

这种测试驱动方法虽然看起来繁琐，但是我们建议对其它项目也都采用该方法。**每次实现一小个功能，让这个功能端到端能够工作，并且用测试覆盖这个功能**。

下面我们来实现打印2,1和"Go!"。

## 先写测试

有了测试代码的保护，我们可以继续迭代。我们不需要频繁停下来运行主程序校验功能，因为测试会确保我们的逻辑的正确的。

[countdown_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v2/countdown_test.go)

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

反引号(`)是另外一种创建字符串的语法，它支持字符串中包含新行

## 写程序逻辑

```go
func Countdown(out io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, "Go!")
}
```

我们用`for`循环向后计数(`i--`)，并用`fmt.Fprintln`将计数打印到`out`，每次输出都换行。最后，我们用`fmt.Fprint`输出"Go!"。

## 重构

我们可以把一些常量抽取出来:

[main.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v2/main.go)

```go
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, finalWord)
}
```

现在运行测试，可以得到期望的结果，但是我们的计数间隔1秒还没有实现。

在Go语言中，`time.Sleep`可以实现时间间隔，修改程序如下：

```go
func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

现在运行测试，也可以通过。

## Mocking

测试仍然可以通过，我们的软件也以预期方式工作，但是我们有一些问题:

- 我们的一个测试需要花费4秒钟运行！
    - 软件开发的前瞻性思维都强调快速反馈环的重要性。
    - **测试慢严重影响开发生产率**
    - 假设需求变得更复杂，需要更多测试。但是每次运行`Countdown`都要花费4秒钟，你能容忍吗？
- 我们还要测试程序的其它重要功能。

我们的程序依赖于`Sleep`ing，我们要将这种依赖抽取出来，这样，我们就可以在测试中控制这种依赖。

如果我们可以**mock**掉`time.Sleep`，那么我们就可以使用**依赖注入** ～ 用假的**spy**替代真实的`time.Sleep`，然后在spy中我们可以测试断言。

## 先写测试

我们将依赖定义为一个接口。这样，我们在`main`中可以使用真实的Sleeper，而在测试中用假的spy sleeper。虽然用了接口，但是我们的`Countdown`函数其实并不关心，并且我们还为调用方增加了灵活性。

```go
type Sleeper interface {
    Sleep()
}
```

我在`Countdown`函数中做了一些调整，`Countdown`函数本身并不负责sleep时间的长短 ～ 而是由函数的使用方决定。

我们先创建一个让测试用的mock：

```go
type SpySleeper struct {
    Calls int
}

func (s *SpySleeper) Sleep() {
    s.Calls++
}
```

**Spy**是**mock**的一种，可以记录依赖是如何被使用的。Spy可以记录传入的参数，被调用了多少次，等等。在我们的案例中，我们跟踪`Sleep()`被调用了多少次，这样我们在测试中就可以校验。

更新测试注入我们的Spy依赖，并且断言sleep被调用了4次。

[countdown_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v3/countdown_test.go)

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}
    spySleeper := &SpySleeper{}

    Countdown(buffer, spySleeper)

    got := buffer.String()
    expected := `3
2
1
Go!`

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }

    if spySleeper.Calls != 4 {
        t.Errorf("not enough calls to sleeper, expected 4 got %d", spySleeper.Calls)
    }
}
```

## 写程序逻辑

修改`Countdown`函数，让其接受`Sleeper`接口，并且在其中调用`sleeper.Sleep()`：

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

这时，`main`程序编译会通不过，所以在`main`程序中，我们需要再创建一个真正的sleeper：

```go
type DefaultSleeper struct {}

func (d *DefaultSleeper) Sleep() {
    time.Sleep(1 * time.Second)
}
```

然后修改主调用程序：

[main.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v3/main.go)

```go
func main() {
    sleeper := &DefaultSleeper{}
    Countdown(os.Stdout, sleeper)
}
```

现在测试可以通过。

### 还有问题

有一个重要的逻辑我们还没有测试。

`Countdown` should sleep before each print, e.g:
`Countdown`在每次打印前应该先睡眠，例如:

- `Sleep`
- `Print N`
- `Sleep`
- `Print N-1`
- `Sleep`
- `Print Go!`
- etc

上面的测试仅仅断言`Countdown`里头有4次睡眠动作，但是那些睡眠动作的次序没有校验。

在写测试的过程中，如果你对测试不是100%确信，那么只要能举出反例就可以break测试！对`Countdown`做如下改变:

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
    }

    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

虽然实现是错误的，但是测试仍然可以通过。

我们要更新测试，仍然用spy可以校验程序的逻辑次序。

我们有两个不同的依赖，并且我们准备把它们的操作记录到一个list中。所以，我们为每个依赖创建一个spy。

```go
type CountdownOperationsSpy struct {
    Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}

const write = "write"
const sleep = "sleep"
```

`CountdownOperationsSpy`同时实现`io.Writer`和`Sleeper`，它将每次调用记录在一个slice中。在本次测试中，我们只关心操作的次序，所以我们将操作记录在一个操作名slice中就可以了。

现在可以在我们的测试族中添加子测试，这个测试校验睡眠和写入的次序。

```go
t.Run("sleep before every print", func(t *testing.T) {
    spySleepPrinter := &CountdownOperationsSpy{}
    Countdown(spySleepPrinter, spySleepPrinter)

    want := []string{
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
    }

    if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
        t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
    }
})
```

注意，我们之前改了`Countdown`的逻辑，现在需要调整回来，这样测试才能通过。

我们现在有两个`Sleeper`的spy实现，所以我们需要重构一下测试，让其中一个测输出的内容，另外一个测睡眠和输出操作的次序。最后，我们可以把第一个spy删掉，因为不需要了。

[countdown_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/mocking/v4/countdown_test.go)

```go
func TestCountdown(t *testing.T) {

    t.Run("prints 3 to Go!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Countdown(buffer, &CountdownOperationsSpy{})

        got := buffer.String()
        want := `3
2
1
Go!`

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

    t.Run("sleep before every print", func(t *testing.T) {
        spySleepPrinter := &CountdownOperationsSpy{}
        Countdown(spySleepPrinter, spySleepPrinter)

        want := []string{
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
        }

        if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
            t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
        }
    })
}
```

目前我们`Countdown`函数已经满足功能需求，并且逻辑正确。

## 将Sleeper变成可配置

最好将`Sleeper`变成可配置，这样我们在主程序中就可以调整睡眠时间。

### 先写测试

Let's first create a new type for `ConfigurableSleeper` that accepts what we need for configuration and testing.
我们先创建一个新类型`ConfigurableSleeper`:

```go
type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}
```

`duration`用于配置睡眠时间，`sleep`则可以传入一个sleep函数。`sleep`的签名和`time.Sleep`是一样的，这样我们在真实实现中就可以用`time.Sleep`，而在测试中用spy:

```go
type SpyTime struct {
    durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
    s.durationSlept = duration
}
```

有了这个spy，我们就可以为configurable sleeper创建一个新的测试。

```go
func TestConfigurableSleeper(t *testing.T) {
    sleepTime := 5 * time.Second

    spyTime := &SpyTime{}
    sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
    sleeper.Sleep()

    if spyTime.durationSlept != sleepTime {
        t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
    }
}
```

这个测试没有什么特别的，测试方式和之前的mock测试没有太大不同。

### 实现程序逻辑

主程序中，我们只需要为`ConfigurableSleeper`添加一个`Sleep`函数:

```go
func (c *ConfigurableSleeper) Sleep() {
    c.sleep(c.duration)
}
```

经过上面的调整，测试可以通过，那么我们为什么要花费力气把Sleeper变成可配置呢？下面会解释。

### 清理和重构

下一步，我们在main主程序中要实际使用`ConfigurableSleeper`:

```go
func main() {
    sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
```

现在运行测试和主程序，你可以看到结果和之前是一致的。

因为我们现在用了`ConfigurableSleeper`，现在可以删除掉`DefaultSleeper`实现了。现在我们有了一个更[通用](https://stackoverflow.com/questions/19291776/whats-the-difference-between-abstraction-and-generalization)的Sleeper，支持可配置睡眠的countdown功能。

## But isn't mocking evil?

你可能听说过mocking is evil。正如软件开发中的任何事物都可以是evil的，例如[DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)。

如果开发人员不能认真倾听测试的反馈，或者不重视重构，那么通常的结果是他们反而会反感测试。

当你要测某个功能的时候，如果你的mocking代码变得越来越复杂，或者需要mock掉很多功能，那么你应该检视你的代码设计，这通常是一个信号:

- 你将要测试的功能承担了太多的职责(因为需要mock掉太多的依赖)
  - 将功能进一步分解成模块，让它们职责单一
- 它的依赖太细粒度了
  - 思考是否可以将某些依赖整合为一个更有意义的模块
- 你的测试太过专注实现细节
  - 测试应该关注期望的行为，而非具体实现

通常，太多的mocking表明代码抽象太差。

**不少人认为的TDD的不足，其实是它的优势**，通常，代码很难测，其实是代码设计差的一个表现，换种说法，设计良好的代码更易于测试。

### 但是mock和测试并没有让我的开发变轻松！

你是否碰到过这样的场景？

- 你想做一些重构
- 但是重构需要改很多测试代码
- 你对TDD产生怀疑，然后在博客上写了一篇文章"Mocking considered harmful"

这种情况的出现，实际表明你测了太多的实现细节。测试应该关注期望的行为，除非实现细节对你的系统的运行很重要。

有时，到底测到什么程度不好把握，下面是一些建议:

- **重构的定义是:改变代码但是系统的行为不变**。如果你决定做一些重构，那么理论上，重构完了你可以直接提交代码，不需要改测试。所以写测试的时候要问自己：
  - 我测试的是系统行为，还是实现细节？
  - 如果我对这块代码做重构，那么我需要对测试做大调整吗？
- 虽然Go语言允许你测试私有函数，我建议尽量避免，因为私有函数是关于具体实现的。
- 我认为如果一个测试使用了超过3个mock，那么这是一个红色信号 ～ 需要花点时间重新思考你的设计。
- 谨慎使用spy。spy让你可以进入算法实现内部，这点有用，但也意味着测试代码和实现之间的一种紧耦合。**在使用spy的时候，确保你确实需要关注这些细节**。

软件开发中的规则总有例外，[Uncle Bob's的文章"When to mock"](https://8thlight.com/blog/uncle-bob/2014/05/10/WhenToMock.html)有一些不错的建议。

## 总结

### 进一步关于TDD

- 当你面对比较大的需求时，先将问题分解，分解为可实现的子问题，然后实现这些子问题。每个实现都是可以端到端工作的软件，并且每个实现都要用测试覆盖，测试反馈要快，小步快跑要远远好于"big bang"方法。
- 一旦你有了可以工作的小软件，你就容易在它基础上进行增量迭代开发，直到开发出你想要的最终软件。

> "When to use iterative development? You should use iterative development only on projects that you want to succeed."
> 
> "什么时候要用迭代式开发？如果你想要让项目成功的话，你就需要用迭代式开发"。

Martin Fowler.

### 关于Mocking

- **如果没有mocking，那么代码的很多重要部分就无法被测试覆盖**。在我们的案例中，我们就无法测试`Countdown`在每次输出之间有间隔睡眠时间，当然实际还有很多其它的例子。例如，要测试的系统对一个第三方服务有依赖调用(可能会失败)，或者要测试系统的某种特殊状态等，如果没有mocking就很难测试这些场景。
- 如果没有mock，那么仅仅只是测试一个简单的业务规则，你可能也需要搭建数据库和其它第三方依赖。然后你的测试就会很慢，导致**慢反馈环**。
- 因为需要搭建数据库或者Web服务才能测试，所以测试就容易不稳定，因为这些依赖的服务可能不稳定。

一旦开发人员学会了mocking这种技术，他们也倾向过度使用mock测试，去测试实现细节(how)，而不是期望行为(what)。因此，在测试前始终要先考虑清楚**测试的价值**，和对未来重构的影响。

本章我们只演示了**Spy**，它只是mock的一种。其实还有不同种类的mocks，[Uncle Bob有一篇易读的文章，解释不同的mock类型](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html)。在后续章节中，我们写的代码会依赖于其它代码提供数据，那时，我们会讲解**Stub**。
