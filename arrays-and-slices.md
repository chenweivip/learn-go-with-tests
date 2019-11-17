# 数组Array和切片Slice

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/arrays)**

数组可以存储多个元素，元素以特定顺序存储。Go语言中数组元素的类型必须相同，

我们经常需要对数组进行迭代操作。我们可以利用[之前的学过的 `for` 循环的知识](iteration.md)，来写一个`Sum`函数。`Sum`函数接受一个整型数组作为参数，然后计算并返回总和。

让我们继续练习TDD技术。

## 先写测试

[`sum_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v1/sum_test.go)

```go
package main

import "testing"

func TestSum(t *testing.T) {

    numbers := [5]int{1, 2, 3, 4, 5}

    got := Sum(numbers)
    expected := 15

    if got != expected {
        t.Errorf("got %d expected %d given, %v", got, expected, numbers)
    }
}
```

Go语言中的数组是**大小固定的**，声明数组变量的时候就要给定大小。有两种方式可以初始化一个数组:

* \[N\]type{value1, value2, ..., valueN} 例如 `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} 例如 `numbers := [...]int{1, 2, 3, 4, 5}`

在测试代码的错误消息中，把原始输入数组也打印出来，可以方便排查问题。这里我们用的占位符是`%v`，这个表示"缺省"格式，适用于任何类型，当然也适用于数组。

关于格式化字符串，可以参考[这里](https://golang.org/pkg/fmt/)


## 编写程序逻辑让测试通过

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v1/sum.go)

```go
func Sum(numbers [5]int) int {
    sum := 0
    for i := 0; i < 5; i++ {
        sum += numbers[i]
    }
    return sum
}
```

通过`array[index]`语法，可以获取数组指定索引位置的值。在上面的Sum函数中，我们用`for`循环对`numbers`数组迭代5次，将每个元素的值加到`sum`总和上。

运行测试`go test`，确保测试通过。

## 重构

我们可以引入 [`range`](https://gobyexample.com/range) 来重构一下代码:

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v2/sum.go):

```go
func Sum(numbers [5]int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

`range`也可以对数组进行迭代。每次调用 `range` 可以返回两个值，一个是索引index，另外一个是该索引位置的数组元素值。我们在index的位置用了 `_` [空位标识符](https://golang.org/doc/effective_go.html#blank)，表示忽略这个索引值，因为我们暂时不用，我们只用数组元素值。

### 数组和其类型

Go语言中的数组有一个特性，数组的大小是编码在其类型中的。如果一个函数期望的参数是 `[5]int`，而你传入的却是 `[4]int`，那么编译会通不过。就像一个函数期望 `int` 参数，而你却输入 `string` 参数， 两个是完全不同类型。

你可能觉得Go语言中的数组有固定长度限制，但是有些场景你不能限定长度，那该怎么办？Go语言中有**切片(slice)**，它是支持不定长度的集合类型。注意，Slice虽然是支持不定长度的集合类型，但是它也是有初始固定容量capacity的，后面我将进一步展开。

下一个需求，我们来实现对不定长度的集合元素求和。

## 先写测试

我们现在来实际使用 [slice type](https://golang.org/doc/effective_go.html#slices)，它是支持不定长的集合类型。语法和数组类似，声明时不需要设定长度，例如：

`mySlice := []int{1,2,3}` 而不是 `myArray := [3]int{1,2,3}`

[`sum_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v3/sum_test.go)

```go
func TestSum(t *testing.T) {

    t.Run("collection of 5 numbers", func(t *testing.T) {
        numbers := []int{1, 2, 3, 4, 5}

        got := Sum(numbers)
        expected := 15

        if got != expected {
            t.Errorf("got %d expected %d given, %v", got, expected, numbers)
        }
    })

    t.Run("collection of any size", func(t *testing.T) {
        numbers := []int{1, 2, 3}

        got := Sum(numbers)
        expected := 6

        if got != expected {
            t.Errorf("got %d expected %d given, %v", got, expected, numbers)
        }
    })

}
```

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v3/sum.go)

```go
func Sum(numbers []int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

运行测试 `go test`，确保通过。

## 重构

`Sum` 函数本身没有问题，但是我们也不能忽视测试代码的质量，大家看下下面的测试代码是否有优化空间？

```go
func TestSum(t *testing.T) {

    t.Run("collection of 5 numbers", func(t *testing.T) {
        numbers := []int{1, 2, 3, 4, 5}

        got := Sum(numbers)
        expected := 15

        if got != expected {
            t.Errorf("got %d expected %d given, %v", got, expected, numbers)
        }
    })

    t.Run("collection of any size", func(t *testing.T) {
        numbers := []int{1, 2, 3}

        got := Sum(numbers)
        expected := 6

        if got != expected {
            t.Errorf("got %d expected %d given, %v", got, expected, numbers)
        }
    })

}
```

有必要检视一下测试的实际价值。我们的目标并非测试越多越好，测试的数量只要能**覆盖代码的关键逻辑**即可。测试数量过多实际也会带来维护开销问题，记住**每个测试都有成本**。

在上面的测试中，我们看到其中有两个类似测试，这其实是有冗余的。如果某个大小的slice可以通过测试，那么基本上任意大小的slice都可以通过测试(在足够信任区间内)。

Go语言的测试工具内置支持[测试覆盖度检查](https://blog.golang.org/cover)，可以帮我们找出哪些代码没有被测试覆盖到。我想强调下，我们的目标并非100%测试覆盖率，这个工具只是帮我们了解代码的测试覆盖度。如果你严格遵循TDD方法，那么很可能你的测试覆盖率已经非常接近100%。

尝试运行:

`go test -cover`

可以看到输出:

```bash
PASS
coverage: 100.0% of statements
```

现在删除上面那个测试方法("collection of 5 numbers")，再运行测试覆盖检查，你会发现测试覆盖率还是100%，所以留一个测试就够了。

下面来一个新需求，添加一个称为 `SumAll` 的新函数，这个函数可以接受可变数量个slices作为参数，并且可以返回一个新的slice，其中每个元素的值是对应传入的每个slice的元素值总和。

例如:

`SumAll([]int{1,2}, []int{0,9})` 将返回 `[]int{3, 9}`

或者

`SumAll([]int{1,1,1})` 将返回 `[]int{3}`

## 先写测试

[`sum_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v3/sum_test.go):

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	expected := []int{3, 9}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v expected %v", got, expected)
	}
}
```

注意，这里对两个slice判断是否相等，不能再用 `got != expected`，可以用Go语言反射包中的 `!reflect.DeepEqual(got, expected` 进行比对，这个可以比对slice中的每一个元素(包括slice长度)。当然，你也可以用 `for` 循环迭代两个slice再比较每个元素，这会麻烦一点。

另外别忘了，你要在程序上方导入 `reflect` 反射包。

## 写程序逻辑

现在来定义 `SummAll` 函数:

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v4/sum.go)

```go
func SumAll(numbersToSum ...[]int) []int {
    lengthOfNumbers := len(numbersToSum)
    sums := make([]int, lengthOfNumbers)

    for i, numbers := range numbersToSum {
        sums[i] = Sum(numbers)
    }

    return sums
}
```

注意，Go函数中[可变数量参数](https://gobyexample.com/variadic-functions)的写法，`numbersToSum ...[]int`表示可变数量个slice(包括0个或者N个)，slice元素的类型是int。

另外，你也学习到可以使用 `make` 创建slice，语句 `sums := make([]int, lengthOfNumbers)` 表示创建一个整数型slice，并且赋值给 `sums`，slice的初始长度为lenthOfNumbers。

`len` 函数可以获取slice的长度，也可以获取array的长度。

slice的索引方式和array类似，`mySlice[N]` 可以获取第N个元素的值。

现在运行测试，确保正确通过。


## 重构

前面提到，slices是有初始固定容量(capacity)的。如果有一个容量为2的slice，但你却试图做赋值操作 `mySlice[10] = 1`，那么你会得到一个运行时(runtime)错误。有些场景下，我们刚开始并不明确数据量的大小，无法正确预估slice初始容量的大小，这个时候该怎么办？

这个时候，你可以用 `append` 函数，这个函数可以将元素值附加到一个slice尾部，如果slice的容量不够，它会自动创建一个容量更大的slice(并自动复制之前的元素)。

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v5/sum.go): 

```go
func SumAll(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        sums = append(sums, Sum(numbers))
    }

    return sums
}
```

采用这种实现方式，我们就不太需要关心容量问题。我们可以从一个空 `sums` slice开始，不断迭代 `numbersToSum` 并将 `Sum` 计算结果添加到 `sums` slice中。

下面又来一个新需求，我们要将 `SumAll` 改成 `SumAllTails` 。原来 `SumAll` 计算slice中所有元素的和，但 `SumAllTails` 只计算尾部元素的和，尾部元素是除第一个头部元素以外的其它所有元素。

## 先写测试

[`sum_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v6/sum_test.go)

```go
func TestSumAllTails(t *testing.T) {
    got := SumAllTails([]int{1,2}, []int{0,9})
    expected := []int{2, 9}

    if !reflect.DeepEqual(got, expected) {
        t.Errorf("got %v expected %v", got, expected)
    }
}
```

## 写程序逻辑

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v6/sum.go)

```go
func SumAllTails(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        tail := numbers[1:]
        sums = append(sums, Sum(tail))
    }

    return sums
}
```

切片(Slice)可以被切分！语法就是 `slice[low:high]`，如果你省略 `low`， 那么它就截取[0 ~ high)(从0到high个，包含第0个，但不包含第high个)个元素；如果你省略 `high`，那么它就截取[high~length)个元素；low和high都省略的话，就是截取所有元素。**注意**，原切片和切分后获取的切片是共享存储的，改变一方的元素值，会影响另外一方。

在我们的代码中，我们用了 `numbers[1:]`，表示截取从第1个元素开始，到结尾的所有元素。注意Slice和Aarray一样，也是从0开始索引的。如果你对Slice切分还不熟悉，建议自己再写一些测试加深理解。

## 再次重构

如果我们给 `SumAllTails` 函数传一个空slice会怎样？空slice有尾部吗？Go语言会如何处理 `myEmptySlice[1:]`？

## 先写测试

```go
func TestSumAllTails(t *testing.T) {

    t.Run("make the sums of some slices", func(t *testing.T) {
        got := SumAllTails([]int{1,2}, []int{0,9})
        expected := []int{2, 9}

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v expected %v", got, expected)
        }
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        expected := []int{0, 9}

        if !reflect.DeepEqual(expected, want) {
            t.Errorf("got %v expected %v", got, expected)
        }
    })

}
```

## 尝试运行测试

```text
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```

注意，程序是通过编译的，这是一个运行时错误。编译错误是我们的朋友，它帮助我们检查程序的语法词法问题，运行时错误则比较危险，它直接影响到用户。

## 修复SumAllTails函数

[`sum.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v7/sum.go)

```go
func SumAllTails(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        if len(numbers) == 0 { // fix
            sums = append(sums, 0)
        } else {
            tail := numbers[1:]
            sums = append(sums, Sum(tail))
        }
    }

    return sums
}
```

## 重构

现在测试代码有一些冗余，我们可以重构一下，取出一个公共测试函数 `checkSums`:

[`sum_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/arrays/v7/sum_test.go)

```go
func TestSumAllTails(t *testing.T) {

    checkSums := func(t *testing.T, got, expected []int) {
        t.Helper()
        if !reflect.DeepEqual(got, expected) {
            t.Errorf("got %v want %v", got, want)
        }
    }

    t.Run("make the sums of tails of", func(t *testing.T) {
        got := SumAllTails([]int{1, 2}, []int{0, 9})
        expected := []int{2, 9}
        checkSums(t, got, expected)
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        expected := []int{0, 9}
        checkSums(t, got, expected)
    })

}
```

## 总结

本章我们学习了:

* 数组Arrays
* 切片Slices
  * 创建Slice的几种方式
  * Slice具有初始固定容量，可以用 `append` 添加元素，如果容量不够，它会自动创建新Slice。
  * 如何对Slice进行切分!
* `len` 函数可以获取一个数组或切片的长度
* 测试覆盖率(coverage)工具
* `reflect.DeepEqual` 的使用场合

我们之前演示的是整数型slices或arrays，其实slices/arrays的元素还可以是其它类型，甚至元素本身也可以是slices或arrays。如果需要，你可以这样定义 `[][]string`。

学习下 [关于切片的Go语言博客](https://blog.golang.org/go-slices-usage-and-internals)，进一步掌握slices。学习过程中，建议写些测试加深理解。

除了写测试外，另外一种方便的学习Go语言的方式是利用Go playground。通过playground，你可以尝试Go语言的大部分功能，你还可以分享代码。例如，[我写了一个关于slice的go playground，你可以尝试学习](https://play.golang.org/p/ICCWcRGIO68)

[另外一个例子](https://play.golang.org/p/bTrRmYfNYCp)，本例展示如果通过一个array获取一个slice，改变这个slice同时会改变原array的值。但我们可以复制array元素的方式创建一个slice，这样改变这个slice，不会影响原array的值。

[还有一个例子](https://play.golang.org/p/Poth8JS28sc)，本例展示在对一个巨大的slice进行切分之后，复制slice有什么好处。


