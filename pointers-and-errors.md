# 指针(Pointer)和错误(Error)

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/pointers)**

上节我们学习了结构体struct，结构体可以把和某个概念相关的一组属性包装起来。

现在，你应该学会了使用结构体管理状态，暴露方法，让用户以受控的方式去访问或者改变结构体的状态。

金融科技和区块链技术是当前热点，所以我们将以银行系统为例来展开本章内容。

我们来创建一个`Wallet`结构体，它可以存比特币`Bitcoin`。

## 先写测试

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v0/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(10)

    got := wallet.Balance()
    expected := 10

    if got != expected {
        t.Errorf("got %d expected %d", got, expected)
    }
}
```

在[之前的例子中](./structs-methods-and-interfaces.md)中，我们通过字段名直接访问字段，但是我们对钱包的安全要求比较高，所以不希望将它的内部状态暴露出来，而是用方法控制对钱包的访问。

## 写程序代码

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v0/wallet.go)

```go
type Wallet struct {
	balance int
}

func (w Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w Wallet) Balance() int {
	return w.balance
}
```

如果上面的语法你还不熟悉，那么请先回上一章结构体(struct)好好阅读一遍。

在Go语言中，如果一个符号(如变量、类型或者函数名等)以小写字母打头，那么它表示私有的(private)，在它所在包之外是不可见的。在上面代码中，我们希望`Wallet`的方法(`Deposit`和`Balance`)对外可见，但是状态变量`balance`对外不可见。

注意之前学过，我们可以使用接收者(receiver)变量`w`来访问`Wallet`结构体中的内部`balance`字段。

现在我们运行测试，结果我们会得到如下错误:

`wallet_test.go:15: got 0 expected 10`

### 为什么出错????

很奇怪，我们的代码看起来应该可以正常工作。按理说，我们在`balance`变量上增加了一些存款，这个状态应该被`balance`保存起来了，然后通过`Balance`方法访问也应该返回被保存的值。

但是实际上在Go语言中，当你调用一个函数或者方法，参数是"值拷贝"的。当你调用`func (w Wallet) Deposit(amount int)`，这里的`w`是一个值拷贝～测试代码中创建的那个wallet实例的拷贝。

我们不想深入讲太多计算机科学的知识。简单讲，当你创建一个值 ～ 例如一个wallet，这个值存在内存的某处。你可以使用`&myVal`取地址运算来获取变量在内存中的地址。

在代码中添加一些打印输出语句:

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v0/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(10)

    got := wallet.Balance()

    fmt.Printf("address of balance in test is %v \n", &wallet.balance)

    expected := 10

    if got != expected {
        t.Errorf("got %d expected %d", got, expected)
    }
}
```

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v0/wallet.go)

```go
func (w Wallet) Deposit(amount int) {
    fmt.Printf("address of balance in Deposit is %v \n", &w.balance)
    w.balance += amount
}
```

注意，`\n`是转义字符，在输出内存地址后换行。

通过`&`取地址运算符，我们可以获取变量的地址，或者说指针。

再次运行测试:

```text
address of balance in Deposit is 0xc000016100 
address of balance in test is 0xc0000160f8 
```

You can see that the addresses of the two balances are different. So when we change the value of the balance inside the code, we are working on a copy of what came from the test. Therefore the balance in the test is unchanged.

你可以看到两个`balance`的地址是不一样的。所以，当我们在代码中改变`balance`的值，我们其实改的是来自测试中的一份拷贝，而测试中本身的`balance`没有变化。

我们可以用指针(pointer)来修复这个问题。所谓指针[pointer](https://gobyexample.com/pointers)，其实是变量的地址，通过指针传递(或者说地址传递)，我们就可以通过指针直接去修改某个变量的值。

修改代码:

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v1/wallet.go)

```go
func (w *Wallet) Deposit(amount int) {
    w.balance += amount
}

func (w *Wallet) Balance() int {
    return w.balance
}
```

修改后代码的区别在接收者(receiver)的类型，我们把`Wallet`改为`*Wallet`，读作"指向一个wallet的指针"。

再次运行测试，现在应该可以通过。

如果你有C语言背景，你会疑惑为什么测试可以通过？因为在函数里头我们都没有对指针进行反引用(deference)，指针反引用应该写成这样才对:

```go
func (w *Wallet) Balance() int {
    return (*w).balance
}
```

而我们之前是直接在指针上访问`balance`字段的。实际上，上面的代码使用`(*w)`也是完全合法的。但是，Go语言的作者认为这种语法太繁琐了，所以简化了，我们可以直接写`w.balance`，不需要先做反引用(dereference)。

指向结构体的指针甚至有一个自己的名称: 结构体指针(struct pointers)，并且它们是能够[自动反引用的](https://golang.org/ref/spec#Method_values).


## 重构

我们之前讲要创建的是比特币钱包，但我们的代码中还没有体现比特币，我们只是用了`int`来存储余额。

为比特币这个概念单独创建一个`struct`看起来有点重，`int`是可以正常工作，但是可读性不佳。

Go语言允许你基于现有类型创建新类型，语法如: `type MyName OriginalType`，这个也叫类型别名(type alias)。

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v1/wallet.go)

```go
type Bitcoin int

type Wallet struct {
    balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
    w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
    return w.balance
}
```

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v1/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(Bitcoin(10))

    got := wallet.Balance()

    expected := Bitcoin(10)

    if got != expected {
        t.Errorf("got %d expected %d", got, expected)
    }
}
```

采用这种类型别名语法创建`Bitcoin`，我们可以用`Bitcoin(999)`。

通过这种方式，我们相当于创建了一个新类型，并且我们可以在新类型上创建方法。如果你想在现有类型上添加一些领域特有的(domain specific)功能，那么这种方式就非常适合。

我们还可以在`Bitcoin`上实现 [Stringer](https://golang.org/pkg/fmt/#Stringer) 接口:

```go
type Stringer interface {
        String() string
}
```

该接口定义在`fmt`包中，实现该接口以后，你可以直接用`%s`格式化字符串，将类型实例的字符串表示打印出来。

```go
func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}
```

你可以看到，在类型别名上创建方法的语法，和直接在struct上创建方法是类似的。

下面我们也更新测试代码，在错误输出格式化字符串中，利用`String()`支持。

```go
    if got != expected {
        t.Errorf("got %s expected %s", got, expected)
    }
```

要看效果的话，把测试代码中的期望值改错就可以看到:

`wallet_test.go:18: got 10 BTC expected 20 BTC`

这样我们的测试输出更清晰了.

我们的下一个需求是实现`Withdraw`方法。

## 先写测试

`Withdraw`的逻辑和`Deposit`正好相反:

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v2/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}

        wallet.Deposit(Bitcoin(10))

        got := wallet.Balance()

        expected := Bitcoin(10)

        if got != expected {
            t.Errorf("got %s expected %s", got, expected)
        }
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}

        wallet.Withdraw(Bitcoin(10))

        got := wallet.Balance()

        expected := Bitcoin(10)

        if got != expected {
            t.Errorf("got %s expected %s", got, expected)
        }
    })
}
```

## 写代码逻辑

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v2/wallet.go)

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
    w.balance -= amount
}
```

现在测试可以通过。

## 重构

测试代码里头有一些重复，我们重构一下:

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v2/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    assertBalance := func(t *testing.T, wallet Wallet, expected Bitcoin) {
        t.Helper()
        got := wallet.Balance()

        if got != expected {
            t.Errorf("got %s expected %s", got, expected)
        }
    }

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

}
```

现在运行测试，确保程序仍然通过。

如果我们尝试超额提取会怎样？显然，我们是不能透支的。我们如何在`Withdraw`方法中把透支错误表达出来呢？

在Go语言中，如果你想表达错误，惯例是在函数中返回一个`err`，然后由调用者判断并采取后续动作。注意，不像Java语言，Go语言没有主动抛出异常的做法，只有返回`err`。

我们还是从测试开始。

## 先写测试

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)

    if err == nil {
        t.Error("expected an error but didn't get one")
    }
})
```

如果用户试图超额提取，我们希望`Withdraw`返回一个错误，并且balance保持不变。在超额提取情况下，测试判断err是否为nil，如果是nil，那么测试就失败并报错。

`nil`类似其它语言(如Java)中的`null`。`Withdraw`的返回值是`error`类型，`error`其实是一个接口。如果一个函数的入参或者返回值是接口，那么入参和返回值可以是`nil`的，所以`Withdraw`也可以返回`nil`。

和其它语言中的`null`一样，如果你试图访问`nil`上的值，系统就会抛**runtime panic**。所以你应该做`nil`检查，避免出现这种情况。

## 实现代码逻辑

[`wallet.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v3/wallet.go)

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("oh no")
    }

    w.balance -= amount
    return nil
}
```

记得在代码中导入`errors`包。

`errors.New`可以创建一个新的`error`实例，你可以给出一个错误消息.

## 重构

我们可以在测试代码中抽取出一个错误检查公共函数，让测试代码更清晰易读。

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v3/wallet_test.go)

```go
assertError := func(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Error("expected an error but didn't get one")
    }
}
```

测试用例更新:

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    wallet := Wallet{Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, Bitcoin(20))
    assertError(t, err)
})
```

前面我们返回的错误消息是"oh no"，这个我们后面还要优化的，因为这个消息对用户的提示作用不大。

假设错误最终返回到用户，让我们来更新一下测试，我们应该断言错误消息，而不只是检查错误存在。

## 先写测试

更新我们的`assertError`测试助手函数，传入一个`string`参数，用于比对:

```go
assertError := func(t *testing.T, got error, expected string) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but expected one")
    }

    if got.Error() != expected {
        t.Errorf("got %q, expected %q", got, expected)
    }
}
```

然后更新调用者:

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)
    assertError(t, err, "cannot withdraw, insufficient funds")
})
```

我们在测试中引入了`t.Fatal`，如果这句被调用，测试将被终止。因为如果`got`是`nil`的话(也就是没有error)，那么就没必要做后续的断言。如果不用`t.Fatal`，测试会继续，然后会引发一个panic，因为后面语句会对nil指针操作。

## 写代码逻辑

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("cannot withdraw, insufficient funds")
    }

    w.balance -= amount
    return nil
}
```

现在测试可以通过。

## 重构

We have duplication of the error message in both the test code and the `Withdraw` code.

在测试代码和`Withdraw`函数中都有相同错误消息，这个是重复的。如果某个开发人员改了程序中的错误消息，那么测试中的错误消息也需要同步修改，这个很烦人。其实测试并不关心具体的错误消息，它只关心在超额取款的情况下，withdray方法需要返回某种有意义的错误。

在Go语言中，错误errors也是值，所以我们可以通过引入变量进行重构，让这个错误成为single source of truth。

```go
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return ErrInsufficientFunds
    }

    w.balance -= amount
    return nil
}
```

`var`关键字可以定义包内可见的全局变量。

重构之后，`Withdraw`函数清晰不少。

下面我们重构测试代码，不再硬编码错误消息，而是引用全局错误变量:

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t *testing.T, wallet Wallet, expected Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }
}

func assertError(t *testing.T, got error, expected error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but expected one")
    }

    if got != expected {
        t.Errorf("got %q, expected %q", got, expected)
    }
}
```

现在测试看起来也更清晰了。我把测试助手函数移到了主测试函数之后，这样开发人员在阅读代码的时候，可以从主测试开始看起，而不是先看测试助手函数。

Another useful property of tests is that they help us understand the _real_ usage of our code so we can make sympathetic code. We can see here that a developer can simply call our code and do an equals check to `ErrInsufficientFunds` and act accordingly.

测试还有一个好处～可以起到样例代码的作用。开发人员通过看测试代码，就可以理解该如何调用代码功能。例如，开发人员通过看针对`Withdraw`函数的测试，就会知道超额的情况下，`Withdraw`会返回错误`ErrInsufficientFunds`，然后他/她在代码中会预先检查这个错误，并做相应处理。

### 遗漏的错误检查

虽然Go编译器对我们有很大帮助，但只限语法词法检查，程序逻辑它帮不了忙，比方说，错误处理可能被疏忽。

我们的测试代码其实漏掉了一个场景，为了把它找出来，可以先通过终端安装`errcheck`这个工具，它是Go语言的程序代码错误检查工具(linter)之一。

`go get -u github.com/kisielk/errcheck`

然后，在代码目录中，运行`errcheck`，你应该看到类似如下输出：

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

这句话是说，这行代码的返回值(即便是nil)我们没有做检查处理。在我的计算机上，这行代码对应正常的取款场景，即便是正常流程，我们也要检查下`Withdraw`应该没有返回error(或者说返回应该是`nil`)。

下面是修复后的最终测试代码:

[`wallet_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/pointers/v4/wallet_test.go)

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
        assertNoError(t, err)
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t *testing.T, wallet Wallet, expected Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != expected {
        t.Errorf("got %s expected %s", got, expected)
    }
}

func assertNoError(t *testing.T, got error) {
    t.Helper()
    if got != nil {
        t.Fatal("got an error but didn't expected one")
    }
}

func assertError(t *testing.T, got error, expected error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but expceted one")
    }

    if got != expected {
        t.Errorf("got %s, expected %s", got, expected)
    }
}
```

## 总结

### 指针

* Go语言的函数/方法的入参采用的是值传递(也就是拷贝入参的值)，如果你想修改原数据的状态，你需要使用指针(pointer)传递，通过指针传递，你才能在函数/方法中修改这个指针指向的值。
* Go语言采用值传递在很多场合下是合适的，但有的时候你不希望用值传递，而是使用引用(reference)传递。引用传递的场合: 包含很大数据的结构体的场合，或者是你只需要一个实例的场合(例如数据库连接池)。

### nil

* 指针可以是nil
* 如果一个函数返回的是指针，那么你必须做nil检查，否则程序可能会抛出运行时异常，这种情况下编译器是帮不了你的。
* 用于表达一个可以为空的值

### 错误Errors

* 错误Errors用于表示调用函数/方法时的一种失败情况
* 在本章的测试中，我们得出结论：在测试中直接检查error中的字符串消息的做法不好。后面我们用更有意义的error常量进行了重构，同时提升了代码和测试的质量，也让用户更容易使用我们提供的API。
* 错误处理涉及很多方面，本章只是一个简介。后续章节我们还会涉及更多错误处理的策略。
* [不要只是检查错误，应该优雅地处理错误](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### 基于现有类型创建新类型

* 也称类型别名，可以扩展现有类型，添加领域特定(domain specific)的功能
* 可以实现接口

在编写Go语言程序的过程中，大部分时间你会和指针/错误打交道，所以你必须熟练使用这两者。所幸的是，如果你不小心搞错了，编译器通常会帮我们解决很多问题，你只需要花点时间阅读编译器的错误提示。
