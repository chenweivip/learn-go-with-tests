# 结构体struct、方法method和接口interface

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/structs)**

假设有一个需求，给定宽height和高width，计算矩形的周长。我们可以写一个函数`Perimeter(width float64, height float64)`，其中`float64`表示浮点数，例如`123.45`。

现在你对TDD方法应该很熟悉了。

## 先写测试

[shapes_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v1/shapes_test.go)

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    expected := 40.0

    if got != expected {
        t.Errorf("got %.2f expected %.2f", got, expected)
    }
}
```

注意新的字符串格式化占位符，`f`是浮点数占位符，`.2`表示打印两位小数。

## 写课程序代码

[shapes.go](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v1/shapes.go)

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}
```

很简单，对吧。现在我们再创建一个函数`Area(width, height float64)`，它可以返回矩形的面积。

你可以尝试先自己来实现，记得遵循TDD方法。

测试代码类似如下:

[shapes_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v2/shapes_test.go)

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    expected := 40.0

    if got != expected {
        t.Errorf("got %.2f expected %.2f", got, expected)
    }
}

func TestArea(t *testing.T) {
    got := Area(12.0, 6.0)
    expected := 72.0

    if got != expected {
        t.Errorf("got %.2f expected %.2f", got, expected)
    }
}
```

实现代码如下:

[shapes.go](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v2/shapes.go)

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
    return width * height
}
```

## 重构

上面的代码可以实现功能，但是代码本身和矩形没有直接关联。一个粗心的程序员可能会误传入一个三角形的宽度和高度，却没有意识到结果是错误的。

我们可以给函数更明确的名称，例如`RectangleArea`。一种更合理的做法是定义我们自己的`Rectangle`类型，通过这个类型封装矩形这个概念。

我们可以用**struct**来创建一个简单类型。[struct](https://golang.org/ref/spec#Struct_types)，简单理解，就是包含一组字段的一个结构体，可以用来存数据。

声明一个`Rectangle`结构体:

```go
type Rectangle struct {
    Width float64
    Height float64
}
```

我们使用`Rectangle`来重构测试:

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v3/shapes_test.go)

```go
func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{10.0, 10.0}
    got := Perimeter(rectangle)
    expected := 40.0

    if got != expected {
        t.Errorf("got %.2f expected %.2f", got, expected)
    }
}

func TestArea(t *testing.T) {
    rectangle := Rectangle{12.0, 6.0}
    got := Area(rectangle)
    expected := 72.0

    if got != expected {
        t.Errorf("got %.2f expected %.2f", got, expected)
    }
}
```

修改程序代码[`shapes.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v3/shapes.go)

```go
func Perimeter(rectangle Rectangle) float64 {
    return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
    return rectangle.Width * rectangle.Height
}
```

可以通过 `myStruct.field` 语法来访问结构体的字段。

通过给函数传一个`Rectangle`类型的参数，这个函数的作用会更加明确。但实际上还有更合理的做法，可以直接使用结构体struct来实现，我们后面会展开。

下一个需求是给圆形写一个`Area`函数。

## 先写测试

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := Area(rectangle)
        expected := 72.0

        if got != expected {
            t.Errorf("got %g expected %g", got, expected)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := Area(circle)
        expected := 314.1592653589793

        if got != expected {
            t.Errorf("got %g expected %g", got, expected)
        }
    })

}
```

可以看到，格式化占位符`f`可以用`g`替代，用`f`的话难以知道确切的小数位，而`g`可以在错误消息中显示完整的小数位([参考fmt选项](https://golang.org/pkg/fmt/))。

## 写程序代码

我们先定义`Circle`类型结构体:

```go
type Circle struct {
    Radius float64
}
```

然后我们实现计算圆形面积的函数，你可以尝试添加`Area(rectangle Rectangle)`函数:

```go
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

但是编译不通过，Go语言不允许你在同一块中重复声明`Area`函数:

`./shapes.go:20:32: Area redeclared in this block`

有两个办法解决这个问题:

* 我们可以将同名的函数声明在不同的包package中，但是这样做有点把事情搞复杂了。
* 我们也可以利用struct类型来定义[方法method](https://golang.org/ref/spec#Method_declarations)。

### 什么是方法?

虽然到目前为止我们只写过函数(function)，但其实我们已经用过一些方法(method)。之前我们调用`t.Errorf`，其实我们是在调用实例`t`(类型为`testing.T`)上的方法`Errorf`。

所谓方法，是带接收者(receiver)的一个函数。方法声明将方法名和方法体绑定起来，并且将这个方法关联到接收者的基础类型上。

方法和函数非常像，但调用方式不同，方法是通过对应实例调用的。你可以在任意地方调用函数，例如`Area(rectangle)`，但你只能在某个"事物"上调用方法。

来看具体例子，我们先修改测试，改为调用方法，后面我们再修改程序代码。

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v4/shapes_test.go)

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := rectangle.Area()
        expected := 72.0

        if got != expected {
            t.Errorf("got %g expected %g", got, expected)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := circle.Area()
        expected := 314.1592653589793

        if got != expected {
            t.Errorf("got %g expected %g", got, expected)
        }
    })

}
```

## 写程序代码

[`shapes.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v4/shapes.go)

```go
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Width + rectangle.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
```

方法声明的语法和函数很像，主要区别在方法接收者的语法: `func (receiverName RecieverType) MethodName(args)`。

当你调用某种类型的方法，可以通过`receiverName`这个变量获取对当前实例的引用。在很多其它语言(比如Java)中，是通过`this`这个接收者来获取当前实例的引用的。

在Go语言中，接收者变量的命名惯例是使用类型的第一个字母，并且小写。

```go
r Rectangle
```

注意，在Circle的`Area`函数中，我们引用了`math`包中的`PI`常量，记得要导入`math`包。

现在运行测试，确保测试通过。


## 重构

目前测试代码里头有重复，我们的两个测试方法的流程都类似: 创建一个形状实例，然后调用`Area()`方法计算面积，最后比对面积。

我们可以抽取公共测试逻辑`checkArea`，它接收一个形状(Shape)，这个形状可以是`Rectangle`，也可以是`Circle`，只要满足支持计算`Area()`即可。

在Go语言中，这个Shape可以用**接口interface**来实现。

在类似Go这样的静态类型语言中，[接口Interfaces](https://golang.org/ref/spec#Interface_types)是一种非常强大的概念，它允许我们创建一种类似具有范型能力的函数～这类函数可以接收不同的类型作为参数，它让我们可以创建高度解耦的代码，同时继续保持类型安全。

为了引入接口，我们先重构测试:

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v5/shapes_test.go)

```go
func TestArea(t *testing.T) {

    checkArea := func(t *testing.T, shape Shape, expected float64) {
        t.Helper()
        got := shape.Area()
        if got != expected {
            t.Errorf("got %g expected %g", got, expected)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })

}
```

`checkArea`是我们抽取出来的一个公共测试函数，它要求传入一个`Shape`，如果我们传入的不是一个Shape，那么编译器就会报错。

这个Shape到底长啥样？在Go语言中，我们只需要定义一个接口声明:

[`shapes.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v5/shapes.go)

```go
type Shape interface {
    Area() float64
}
```

就像我们之前创建`Rectangle`和`Circle`一样，我们再创建了一个新类型`Shape`，只不过这次我们用的是`interface`，而不是之前的`struct`。

现在运行测试，可以通过。

### 这个接口有点奇怪？

Go语言中的接口和其它语言中的接口很不一样。在其它语言中，你的类型必须显式地实现接口，就像`My type Foo implements interface Bar`这样。

但在我们的案例中:

* `Rectangle`有一个称为`Area`的方法，它返回一个`float64`类型的返回值，所以它满足`Shape`接口规范
* `Circle`也有一个称为`Area`的方法，它也返回一个`float64`类型的返回值，所以它也满足`Shape`接口规范
* `string`没有称为`Area`的方法，所以它不满足`Shape`接口规范
* 等等

在Go语言中，**接口解析是隐式的**。只要你传入的类型满足接口类型规范(具有接口要求的方法)，编译就会通过，它不要求显示声明。

### 解耦

注意，我们的测试公共函数`checkArea`并不关心传入的是一个`Rectangle` or `Circle` or `Triangle`。通过声明一个接口，这个函数就和具体的类型解耦了，它只需关心具体的操作逻辑。

接口规范声明支持哪些方法，具体类型只要具备同名方法就满足接口，这种方式在软件设计中非常重要，后续章节我们会讲解更多细节。

## 进一步重构

既然我们对struct已经有所理解，我们可以引入"表驱动测试"。

[表驱动测试(Table Driven Tests)](https://github.com/golang/go/wiki/TableDrivenTests)是一种测试方法，在对一组测试用例进行相同测试的时候，表驱动测试比较有用。

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v6/shapes_test.go)

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        expected  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.expected {
            t.Errorf("got %g expected %g", got, tt.expected)
        }
    }

}
```

在上面的测试中，我们声明了一个结构体切片(a slice of structs)，这个结构体是一个**匿名结构体**，具有两个字段，`shape`和`expected`，然后我们创建一个`Rectangle`和一个`Circle`，作为测试用例填充到切片中，最后将切片赋值给`areaTests`变量。

然后我们对`areaTests`切片进行迭代，使用结构体上的字段运行测试。

采用这种做法，开发人员只需要添加一个新的Shape结构体类型，实现`Area`方法，然后创建对应实例并添加到测试切片列表中，就可以进行测试。

表驱动测试是一种有用的测试方法，但是开发起开需要一些额外的投入。如果你需要对某个接口的不同实现进行测试，那么表驱动测试是一种比较合适的方法。

我们再添加一个形状~三角形(triangle)，来演示表驱动测试。

## 先写测试

为我们的新形状添加一个测试很简单，只需在测试列表中添加 `{Triangle{12, 6}, 36.0},`。

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v7/shapes_test.go)

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        expected  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
        {Triangle{12, 6}, 36.0},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.expected {
            t.Errorf("got %g expected %g", got, tt.expected)
        }
    }

}
```

## 添加程序逻辑

新建三角形Triangle结构体:

[`shapes.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v7/shapes.go)

```go
type Triangle struct {
    Base   float64
    Height float64
}

func (t Triangle) Area() float64 {
    return (t.Base * t.Height) * 0.5
}
```

注意，Triangle结构体必须具备`Area()`方法，并且返回值类型是`float64`，这样才能满足`Shape`接口规范要求，否则编译通不过。

运行测试，校验通过。

## 重构

到目前为止，我们的代码实现是可以的，但是测试方面还可以再提升。

看一下下面的代码:

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

这几行代码的可读性不佳，含义并不明显，或者说不太容易理解。

之前我们通过`MyStruct{val1, val2}`方式创建实例，但实际上你还可以命名字段。

再看下面的重构后的代码:

```go
        {shape: Rectangle{Width: 12, Height: 6}, expected: 72.0},
        {shape: Circle{Radius: 10}, expected: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, expected: 36.0},
```

In [Test-Driven Development by Example](https://g.co/kgs/yCzDLF) Kent Beck refactors some tests to a point and asserts:

在[Test-Driven Development by Example](https://g.co/kgs/yCzDLF)这本书中，Kent Beck在重构完一些测试代码后指出:

> The test speaks to us more clearly, as if it were an assertion of truth, **not a sequence of operations**
> 
> 测试应当浅显易懂，看上去就是断言一些易懂的事实，而不是系列难理解的操作


显然，重构后的代码更浅显易懂。

## 让测试输出更有意义

之前对`Triangle`的测试，如果`Area`函数逻辑不正确，那么错误输出可能类似如下:
`shapes_test.go:31: got 0.00 expected 36.00`.

我们知道这个错误和`Triangle`有关，那是因为我们正好在测它。但是如果我们的测试用例很多(比如超过20个)，然后其中一个有bug，如果错误输出不够明确的话，开发人员如何知道具体是哪个用例失败了呢？这个也是开发者体验问题，他们可能需要反复翻看代码才能具体定位哪个用例出错了。

我们可以把错误消息格式化字符串改为`%#v got %.2f expected %.2f`。`%#v`格式化字符串会把相关结构体及其字段都打印出来，这样开发人员就比较容易查看和定位问题。

为了进一步提升测试代码的可读性，我们可以把`expected`字段命名为更具描述性的字段如`hasArea`。

关于表驱动测试的最后一个技巧是使用 `t.Run`，并给测试用例命名。

通过将每个用例包裹在`t.Run`方法中，那么测试失败时会输出更清晰的错误消息，因为它会打印出用例名称，例如:


```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 expected 72.10
```

并且你还可以指定运行表中的某个用例 `go test -run TestArea/Rectangle`。

以下是重构后的最终测试代码:

[`shapes_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/structs/v8/shapes_test.go)

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        // using tt.name from the case to use it as the `t.Run` test name
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %g expected %g", tt.shape, got, tt.hasArea)
            }
        })

    }

}
```

## 总结

我们接触了更多的TDD实践，通过对基本的几何图形计算的改进，我们逐步了学习新的语言功能:

* 通过结构体struct来创建你自己的数据类型，它可以把一组相关数据包装起来，让你的代码意图更清晰
* 通过接口interfact，可以让函数接受遵循同一接口规范的不同类型作为输入\(参考[参数多态化parametric polymorphism](https://en.wikipedia.org/wiki/Parametric_polymorphism)\)
* 在数据类型上，可以添加方法来为类型添加功能，这些方法可以遵循某个接口规范
* 表驱动测试让你的测试断言更清晰，也让你的测试族更易于扩展和维护

**本章比较重要**，因为我们开始定义自己的类型了。在像Go这样的静态语言中，能够定制自己的类型是非常重要的，它让我们能够创建更大更复杂的软件系统，并且代码易于理解，模块化和测试。

接口是一种强大的解耦机制，让我们可以隔离和隐藏复杂性。在我们之前的公共测试函数中，测试代码并不需要确切知道传入的具体是哪种Shape类型，只需要能够调用实例的`Area`方法就可以了。

随着你对Go语言越来越熟悉，你会逐渐体会Go语言接口和标准库的强大能力。你会看到，标准库中大量定义和使用接口，通过让你的类型也实现标准库的接口，你就可以很快重用标准库的大量功能。
