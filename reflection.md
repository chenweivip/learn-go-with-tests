# 反射

[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/reflection)

[来自Twitter的Go语言挑战题](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> 写一个函数`walk(x interface{}, fn func(string))`，该函数接受一个struct `x`和一个函数`fn`作为输入，`fn`将对struct `x`中的所有字符串字段进行调用。难度：递归级别。

为了解决这个问题，我们需要用到**反射reflection**。

> 编程语言中的反射指的是程序有能力检查自身的结构，主要通过类型检查。它是元编程的一种形式。反射属于高级编程主题，很对开发人员对反射感到困惑。

来自 [Go 语言博客: 反射](https://blog.golang.org/laws-of-reflection)

## 到底什么是`interface`?

Go语言支持类型安全(type-safety)，我们已经深有体会，例如声明函数只接受某种明确类型，如`string`，`int`，和我们自己定义的类型如`BankAccount`。

类型安全有很多好处，其中之一是一眼就可以看出支持的类型(相当于一种文档)，另外一个是编译器检查，如果传错类型，编译器会告诉我们。

但是你也可能碰到一种场景 ～ 要写一个函数，但是在编译期还不知道具体的参数类型。

Go语言允许你用`interface{}`这种类型来表达**任意**类型。

所以`walk(x interface{}, fn func(string))`可以接受任意类型的值作为`x`。

### 既然`interface`可以代表任何类型，为何还要定义具体类型的函数呢？

- 如果一个函数直接接受`interface`作为输入，那么它就会失去类型安全(type safety)检查。假如你想将`Foo.bar`(string类型)传给一个函数，结果却传了`Foo.baz`(int类型)，那会怎样？编译器无法给你错误提示。你也无法知道函数到底接受哪种类型的输入。如果函数明确声明接受某种类型的参数(例如`UserService`类型)，那你就不会轻易弄错。
- 作为写函数的作者，你就需要对输入的参数进行检查，检查输入的是什么类型，然后再对其做适当处理。这一般需要通过**反射**才能做到。这种做法比较复杂，代码比较难读，性能也会下降(因为需要做运行时类型检查）。

简言之，只有在真正需要时才考虑使用反射。

如果你想要多态函数(polymorphic functions)，先考虑你是否能够采用面向接口的设计方式(注意不是直接传`interface{}`)，这样你的函数就可以接受多个类型（只要这些类型实现你定义的接口)。

本次案例中，我们的函数要处理很多不同类型。和之前一样，我们将采用迭代方法，每支持一个功能，我们从先写测试开始，然后在此基础上不断重构，直到实现最终目标。

## 先写测试

我们先从只有一个字段的struct开始：

```go
func TestWalk(t *testing.T) {

    expected := "Bobo"
    var got []string

    x := struct {
        Name string
    }{expected}

    walk(x, func(input string) {
        got = append(got, input)
    })

    if len(got) != 1 {
        t.Errorf("wrong number of function calls, got %d expect %d", len(got), 1)
    }
}
```

`walk`函数接受一个匿名struct x，和一个匿名函数。struct x中只有一个字段，存储我们的期望字符串。匿名函数接受一个字符串输入，并将字符串添加到`got` slice中。刚开始我们先简单点，只检查`got`的长度是否满足期望，后面我们会细化进一步检查具体内容。

## 写程序逻辑

为了让上面的测试通过，我们可以给fn调用传任意字符串：

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

测试现在可以通过。下一步我们要具体断言`fn`被调用时，接受的是真正的字符串参数。

## 先写测试

修改测试，校验`fn`中接收到的字符串(在`got[0]`中)，和期望的字符串一致。

```go
if got[0] != expected {
    t.Errorf("got %q, want %q", got[0], expected)
}
```

## 完成程序逻辑

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v1/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)
    field := val.Field(0)
    fn(field.String())
}
```

上面的代码**既不安全也不完整**，但是还记得我们的TDD和增量驱动方法吗？我们的目标是小步行进，用最少代码先让它能工作，后面我们会不断重构优化。

我们需要使用反射来检查`x`，看它内部的属性。

[反射包](https://godoc.org/reflect)里头有一个函数`ValueOf`，它可以返回一个变量的`Value`。然后我们就可以检查这个值，包括它的字段(见代码下一行)。

然后，我们对传入的值做了一些乐观假设：

- 我们只检查第一个(也是唯一的一个)字段，如果一个字段都没有的话，那么就会导致panic。
- 然后我们对字段调用了`String()`方法，它将底层的值以string形式返回，如果底层的值无法以string形式返回，那么就会出错。

## 重构

对于简单的情况，我们的测试可以通过，但是我们知道目前的代码还有很多不足。

我会写更多测试，传入不同的值，让`fn`对不同值进行调用，然后检查`got` slice的值满足期望。

我们将测试重构为表驱动测试，这样方便我们继续测试新的场景。

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v1/reflection_test.go)

```go
func TestWalk(t *testing.T) {

    cases := []struct{
        Name string
        Input interface{}
        ExpectedCalls []string
    } {
        {
            "Struct with one string field",
            struct {
                Name string
            }{ "Bobo"},
            []string{"Bobo"},
        },
    }

    for _, test := range cases {
        t.Run(test.Name, func(t *testing.T) {
            var got []string
            walk(test.Input, func(input string) {
                got = append(got, input)
            })

            if !reflect.DeepEqual(got, test.ExpectedCalls) {
                t.Errorf("got %v, expect %v", got, test.ExpectedCalls)
            }
        })
    }
}
```

显然我们可以很容易添加新的场景，比如超过1个string字段的场景。

## 先写测试

在我们的测试用例中添加如下场景：

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v2/reflection_test.go)

```go
{
    "Struct with two string fields",
    struct {
        Name string
        City string
    }{"Bobo", "Shanghai"},
    []string{"Bobo", "Shanghai"},
}
```

## 调整程序逻辑

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v2/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i:=0; i<val.NumField(); i++ {
        field := val.Field(i)
        fn(field.String())
    }
}
```

`val`有一个方法`NumField`，可以返回这个值所有的字段数量。然后我们可以对所有字段进行迭代，并对每个字段调用`fn`方法。

## 完善

我们的程序有一个不足：`walk`假定每个字段都是`string`类型的，所以我们需要对代码进行完善。我们先写测试：

## 先写测试

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v3/reflection_test.go)

添加如下测试用例：

```go
{
    "Struct with non string field",
    struct {
        Name string
        Age  int
    }{"Bobo", 33},
    []string{"Bobo"},
},
```

## 调整程序逻辑

我们需要检查字段的类型是`string`：

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v3/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }
    }
}
```

可以通过字段的[`Kind`](https://godoc.org/reflect#Kind)检查其类型。

## 重构

现在的代码比之前要完善很多。

下一个场景，如果传入的`struct`不是扁平结构，而是嵌套结构的，怎么办？

## 先写测试

之前我们使用过匿名struct语法，很方便，这里我们可以继续使用：

```go
{
    "Nested fields",
    struct {
        Name string
        Profile struct {
            Age  int
            City string
        }
    }{"Bobo", struct {
        Age  int
        City string
    }{33, "Shanghai"}},
    []string{"Bobo", "Shanghai"},
},
```

值得注意的是，使用匿名struct语法之后，代码会比较难读。[社区有提议优化这种语法](https://github.com/golang/go/issues/12854)。

我们来重构一下，专门为这个场景创建一个新类型，然后在测试中引用这个类型。这对测试来说会引入一些复杂性 ～ 有些用于测试的代码在测试之外，但是读者可以通过初始化方式推断出`struct`的结构来。

在测试代码下方添加如下类型声明：

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v4/reflection_test.go)

```go
type Person struct {
    Name    string
    Profile Profile
}

type Profile struct {
    Age  int
    City string
}
```

调整测试代码，使用`Person`类型进行初始化，现在测试代码看起来会更清楚：

```go
{
    "Nested fields",
    Person{
        "Bobo",
        Profile{33, "Shanghai"},
    },
    []string{"Bobo", "Shanghai"},
},
```

## 调整程序逻辑

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v4/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }

        if field.Kind() == reflect.Struct {
            walk(field.Interface(), fn)
        }
    }
}
```

解决办法很简单，我们仍然检查字段的`Kind`，如果是`struct`类型，我们就对内嵌的`struct`再次调用`walk`(递归调用)。

## 重构

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

如果对某个值的比较超过两次，建议可以重构为`switch`方式，让代码易读也易于扩展。

下一步，如果传入的是一个指针会如何？

## 先写测试

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v5/reflection_test.go)

添加一个测试用例：

```go
{
    "Pointers to things",
    &Person{
        "Bobo",
        Profile{33, "Shanghai"},
    },
    []string{"Bobo", "Shanghai"},
},
```

## 测试通不过

```
=== RUN   TestWalk/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## 实现程序逻辑

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

不能在一个指针`Value`上使用`NumField`，对于指针类型，我们先要使用`Elem`取出底层的值。

## 重构

我们把从一个`interface{}`取出`reflect.Value`的动作重构为一个函数`getValue`。

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v5/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}

func getValue(x interface{}) reflect.Value {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    return val
}
```

这样做会让代码变多，但是这种重构在抽象级别上看是正确的。

- 先通过`getValue`从`x`中取出`reflect.Value`，这样我们不必关心是指针还是非指针。
- 再对字段进行迭代，根据其类型做相应处理。

下一步，我们要来考虑切片slice的场景。

## 先写测试

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v6/reflection_test.go)

```go
{
    "Slices",
    []Profile {
        {33, "Shanghai"},
        {34, "Beijing"},
    },
    []string{"Shanghai", "Beijing"},
},
```

## 实现程序逻辑

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    if val.Kind() == reflect.Slice {
        for i:=0; i< val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
        return
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

## 重构

上面的代码可以工作，但是质量不高。不必担心，我们的代码受测试保护，我们可以根据需要大胆重构。

如果你稍微抽象地思考一下，我们想让`walk`调用的是：

- struct上的每个字段
- slice中的每个值(未知类型)

我们的代码目前可以工作，但是抽象得不太好。我们先检查是否是slice（如果是的话，迭代执行完`walk`之后就返回)，然后我们再检查struct场景。

我们可以重构代码，先检查类型，再做具体工作`walk`：

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    switch val.Kind() {
    case reflect.Struct:
        for i:=0; i<val.NumField(); i++ {
            walk(val.Field(i).Interface(), fn)
        }
    case reflect.Slice:
        for i:=0; i<val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
    case reflect.String:
        fn(val.String())
    }
}
```

现在看起来代码要好不少！如果是一个struct(或slice)，我们就对每个字段(或每个索引)对应的值迭代调用`walk`。否则，如果是`reflect.String`的话，我们就直接调用`fn`。

对我来说，还可以重构得更好。对字段(或索引)对应值的迭代调用有重复，但是从概念上讲，它们是一样的。

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v6/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

如果`value`是一个`reflect.string`，那么我们就像之前一样直接调用`fn`。

否则，我们的`switch`将根据类型取出两样东西：

- 有多少个字段
- 如何取出值`Value`(通过`Field`或者`Index`函数)

一旦我们取得上述数据，我们就可以迭代`numberOfValues`次，每次调用`walk`，传入`getField`函数调用的结果值。

下面我们来考虑array情况，有了处理slice的经验，处理array应该不复杂。

## 先写测试

添加一个测试用例：

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v7/reflection_test.go)

```go
{
    "Arrays",
    [2]Profile {
        {33, "Shanghai"},
        {34, "Beijing"},
    },
    []string{"Shanghai", "Beijing"},
},
```

## 实现程序逻辑

Array的处理方式和slice类似，所以我们只需添加一个逗号分隔：

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v7/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

我们要处理的最后一个类型是`map`。

## 先写测试

```go
{
    "Maps",
    map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    },
    []string{"Bar", "Boz"},
},
```

## 写程序逻辑

再抽象思考一下，你会发现`map`和`struct`非常像，只是`map`的keys在编译时还未知。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

但是，你不能通过index从map中获取值，只能通过**key**，所以之前的抽象被破环了。

## 重构

你现在感觉如何？之前的抽象让我们感觉良好，但是现在代码的味道又不太好了。

**这很正常**，重构是一个过程，期间我们可能会犯错误。TDD的一大好处是，它给我们以试错的自由。

我们的每一步都有测试保护，所以我们完全可以回到之前的步骤。让我们回到重构之前。

[reflection.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v8/reflection.go)

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    walkValue := func(value reflect.Value) {
        walk(value.Interface(), fn)
    }

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        for i := 0; i< val.NumField(); i++ {
            walkValue(val.Field(i))
        }
    case reflect.Slice, reflect.Array:
        for i:= 0; i<val.Len(); i++ {
            walkValue(val.Index(i))
        }
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walkValue(val.MapIndex(key))
        }
    }
}
```

我们引入了一个`walkValue`匿名函数，这样我们不用在`switch`中直接调用`walk`，代码看起来抽象一致，更清晰。

### 最后一个问题

Go语言中的map并不保证顺序。所以你的测试有时会失败，因为我们的断言要求对`fn`是按特定顺序调用的。

为了修复这个问题，我们需要将对map用例的断言移出，做成一个单独的测试助手函数，这样我们就不必关心顺序问题。

[reflection_test.go](https://github.com/spring2go/learn-go-with-tests/blob/master/reflection/v8/reflection_test.go)

```go
t.Run("with maps", func(t *testing.T) {
    aMap := map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    }

    var got []string
    walk(aMap, func(input string) {
        got = append(got, input)
    })

    assertContains(t, got, "Bar")
    assertContains(t, got, "Boz")
})
```

这是`assertContains`的定义：

```go
func assertContains(t *testing.T, haystack []string, needle string)  {
    contains := false
    for _, x := range haystack {
        if x == needle {
            contains = true
        }
    }
    if !contains {
        t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
    }
}
```

## 总结

- 引入了`reflect`包中的一些概念。
- 使用递归遍历任意的数据结构
- 中间做了一次不成功的重构尝试，但是由于采用TDD和增量方法，我们很容易回到重构之前。
- 本章只是涉及了反射的很小一部分。[Go语言博客有一篇博文讲解关于反射的更多内容](https://blog.golang.org/laws-of-reflection)
- 虽然你已经学习了反射，但是在实践中，还是尽量避免使用反射。
