# 字典Maps

**[本章代码](https://github.com/spring2go/learn-go-with-tests/tree/master/maps)**

在[数组和切片](arrays-and-slices.md)章节，我们学习了如何顺序存储数据。现在，我们来学习如何通过键`key`来存储数据，然后快速查找数据。

Map这种数据类型和字典类似，它支持以key/value对方式存取数据。你可以把`key`看成是字典里头的字(或词)，`value`可以看成是字典里头的对字(或词)的定义。下面我们将实际动手创建Map这种数据结构，来进一步学习它。

首先，假设我们的字典里头已经存在一些字和对应的定义，如果我们按字搜索，字典就会返回对应的定义。

## 先写测试

`dictionary_test.go`

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    expected := "this is just a test"

    if got != expected {
        t.Errorf("got %q expected %q given, %q", got, expected, "test")
    }
}
```

声明字典的方式和声明数组有点类似，只不过字典是以`map`关键字开始声明，并且需要两个类型。第一个是键的类型，这个键key写在方括号`[]`中。第二个是值的类型，写在方括号之后。

键的类型比较特殊，它只能是可以比较的类型，显然，如果无法比较两个key是否相等，我们就无法确保获得正确的值。可比较类型在[语言规范](https://golang.org/ref/spec#Comparison_operators)中有详细解释。

而值value可以是任意类型，甚至可以是另一个map。

测试中的其它部分你应该已经熟悉了。

## 写程序逻辑

`dictionary.go`

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

从字典中获取值的语法`map[key]`。

## 重构

```go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    expected := "this is just a test"

    assertStrings(t, got, expected)
}

func assertStrings(t *testing.T, got, expected string) {
    t.Helper()

    if got != expected {
        t.Errorf("got %q expected %q", got, expected)
    }
}
```

我把`assertStrings`助手函数抽取出来，让测试更清晰。

### 使用一个定制类型

We can improve our dictionary's usage by creating a new type around map and making `Search` a method.
我们可以用类型别名改进代码，在map基础上创建一个新类型，然后在新类型上添加`Search`方法。

在[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v1/dictionary_test.go)文件中:

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    expected := "this is just a test"

    assertStrings(t, got, expected)
}
```

上面的测试中我们用了`Dictionary`类型，然后在这个类型的实例`dictionary`上调用了`Search`方法。`assertStrings`无需变化。

下面我们来定义`Dictionary`类型。

在[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v1/dictionary.go)文件中:

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

我们创建了一个`Dictionary`类型，它实际上是`map`的一个封装类型。有了定制类型以后，我们就可以创建`Search`方法。

## 先写测试

对字典的基本查找很容易实现，但是如果查找的字在字典中不存在会怎样？我们应该什么也拿不到。这是OK的，程序可以继续运行，但是还有一个更好的做法 ～ 函数可以明确报告该字在字典中不存在，这样，用户不至于疑惑。

我们先写测试:


```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        expected := "this is just a test"

        assertStrings(t, got, expected)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        expected := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), expected)
    })
}
```

Go语言中处理这种场景的方式，就是返回第二个类型为`Error`的返回值。

通过调用`Error`实例的`.Error()`方法，`Error`可以被转换成一个字符串，我们在断言中就是这样转的。我们也对`assertStrings`加了一个`if`判断作为保护，确保我们不会在`nil`上调用`.Error()`。

## 写程序逻辑

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

为了让测试通过，我们使用了map的一种特别的查找语法，它可以返回2个值。第二个值是一个布尔值，表明对应的键是否存在。这样，我们就可以区分某个键是不存在，还是没有对应的定义。

## 重构

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v2/dictionary.go)

```go
var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

通过把error抽取为一个常量，我们的测试代码会变得更清晰。

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v2/dictionary_test.go)

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    assertError(t, got, ErrNotFound)
})
}

func assertError(t *testing.T, got, expected error) {
    t.Helper()

    if got == nil {
        t.Fatal("expected to get an error.")
    }
    
    if got != expected {
        t.Errorf("got error %q expected %q", got, expected)
    }
}
```

再重构下测试代码，把`assertError`抽取出来，这样可以简化测试。通过重用`ErrNotFound`变量，我们的测试代码可维护性增强了(后续修改错误消息只需集中修改一个地方)。

## 先写测试

我们已经可以搜索字典了，但我们还需要支持向字典添加新字。

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    expected := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if expected != got {
        t.Errorf("got %q expected %q", got, expected)
    }
}
```

先添加新字和定义，然后查找，再断言。

## 编写代码逻辑

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v3/dictionary.go)

```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

Adding to a map is also similar to an array. You just need to specify a key and set it equal to a value.

### 引用类型

字典的一个特性是你可以直接修改它们，而无需传递指针。**因为`map`是引用类型** ～ 它对底层数据结构有一个引用，非常像一个指针。底层数据结构是一个哈希表，关于哈希表，可以参考[这里](https://en.wikipedia.org/wiki/Hash_table)。

字典属于引用类型非常有用，因为不管字典长多大，它始终只有一份拷贝。

关于引用类型要注意的一点是，字典可能为`nil`。当试图读取的时候，`nil`字典的行为和空字典是一样的，但是如果试图写入一个`nil`字典，那么程序会抛**runtime panic**。关于字典的更多信息，可以参考[这里](https://blog.golang.org/go-maps-in-action)。

因此，你不应该以如下方式初始化一个空字典变量:

```go
var m map[string]string
```
而是应该用下面的方式，或者使用`make`关键字初始化空字典:

```go
var dictionary = map[string]string{}

// OR

var dictionary = make(map[string]string)
```

上面两种方法都可以创建空字典(和指向空字典的指针)，这两种初始化方法可以确保不会产生**runtime panic**。

## 重构

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v3/dictionary_test.go)

代码无需重构，但是测试可以再简化一下。

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    word := "test"
    definition := "this is just a test"

    dictionary.Add(word, definition)

    assertDefinition(t, dictionary, word, definition)
}

func assertDefinition(t *testing.T, dictionary Dictionary, word, definition string) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if definition != got {
        t.Errorf("got %q expected %q", got, definition)
    }
}
```

我们为`word`和`definition`创建了变量，并且把对definition的断言移到了助手函数中。

我们的`Add`方法看起来可以了。但是，我们还没有考虑试图添加已经存在的键的情况！

如果键已经存在，向字典添加重复键不会抛错，它只会用新值覆盖现有的值。实践中这一行为是蛮方便的，但让我们的函数名变得意义不明确 ～ `Add`不应该修改现有的值，它应该只向字典添加新的键值对。

## 先写测试

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v4/dictionary_test.go)

```go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}
        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
```

For this test, we modified `Add` to return an error, which we are validating against a new error variable, `ErrWordExists`. We also modified the previous test to check for a `nil` error.

为了让这个测试通过，我们需要修改`Add`方法返回一个错误，然后在测试中，将错误和一个新的错误变量`ErrWordExists`进行比对。之前的测试我们也修改了一下，检查err是`nil`。

## 写代码逻辑

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

这里我们用了`switch`语句来匹配错误，如果`Search`返回一个除`ErrNotFound`以外的错误，`switch`语句提供了额外的检查和返回，这样更简洁安全。

## 重构

没有太多需要重构，但是因为用了几个error，我们可以做些小修改。

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v4/dictionary.go)

```go
const (
    ErrNotFound   = DictionaryErr("could not find the word you were looking for")
    ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

我们把errors改成了常量，这要求我们创建定制的`DictionaryErr`类型，这个类型要实现`error`接口。关于这种用法的细节，[Dave Cheney](https://dave.cheney.net/2016/04/07/constant-errors)写了一篇很不错的文章。简单讲，它让错误变得可重用，并且是不可变的(immutable)。

下一步，我们来创建一个`Update`方法，可以更新字典中字的定义。

## 先写测试

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v5/dictionary_test.go)

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

`Update` is very closely related to `Add` and will be our next implementation.

`Update`和`Add`类似，我们马上来实现。

## 写程序逻辑

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v5/dictionary.go)

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

代码很少，但是我们有一个和之前`Add`类似的问题 ～ 如果我们传入一个新字，`Update`也会把它添加到字典中.

## 先写测试

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v6/dictionary_test.go)

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    newDefinition := "new definition"
    dictionary := Dictionary{word: definition}

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

我们需要新加一个错误类型`ErrWordDoesNotExist`，如果更新时键key不存在，就返回这个错误。

## 写程序逻辑

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v6/dictionary.go)

```go

const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

这个函数和`Add`很像，只是字典更新和错误返回逻辑有调整。

### 为更新声明一个新错误类型

我们可以重用`ErrNotFound`，但最好再创建一个新的错误类型，这样在更新失败时可以获得更明确错误提示。

在出错时，明确的错误会给你更多提示信息。例如在一个web应用中:

> 如果碰到一个`ErrNotFound`错误，你可以将用户重定向，而当碰到一个`ErrWordDoesNotExist`，你可以显示一个明确错误消息。

下面，我们来为字典创建一个`Delete`功能。

## 先写测试

[`dictionary_test.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v6/dictionary_test.go)

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected %q to be deleted", word)
    }
}
```

先创建一个`Dictionary`，初始化一个字，然后删除这个字，最后检查这个字确实被删除。

## 写程序逻辑

[`dictionary.go`](https://github.com/spring2go/learn-go-with-tests/blob/master/maps/v7/dictionary.go)

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Go内置支持`delete`函数，它可以应用于字典。它接收两个参数，第一个是字典(map)，第二个是要删除的键(key)。

`delete`函数没有返回，所以我们的`Delete`方法也没有返回。因为删除一个不存在的键是没有效果的，所以我们没必要像`Update`和`Delete`那样再写switch判断逻辑。

## 总结

本章我们讲了很多东西，为我们自己定义的字典开发了完整的增删改查(CRUD，Create/Read/Update/Delete)API，通过这个过程我们学到:

* 创建字典
* 在字典中查找项
* 给字典添加新项
* 更新字典中的项
* 从字典中删除项
* 学习了更多错误处理技术
  * 如何创建常量型错我
  * 编写错误封装
