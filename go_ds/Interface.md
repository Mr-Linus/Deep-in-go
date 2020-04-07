## Interface

在 Go 语言的语境中，当我们在谈论“接口”的时候，一定指的是接口类型。因为接口类型与其他数据类型不同，它是没法被实例化的。既不能通过调用new函数或make函数创建出一个接口类型的值，也无法用字面量来表示一个接口类型的值。

如果没有任何数据类型可以作为它的实现，那么该接口的值就不可能存在。通过关键字type和interface，我们可以声明出接口类型。

接口类型声明中的这些方法所代表的就是该接口的方法集合。一个接口的方法集合就是它的全部特征。

对于任何数据类型，只要它的方法集合中完全包含了一个接口的全部特征（即全部的方法），那么它就一定是这个接口的实现类型。比如下面这样：

```go
type Pet interface {
  SetName(name string)
  Name() string
  Category() string
}
```

这里声明了一个接口类型Pet，它包含了 3 个方法定义，方法名称分别为SetName、Name和Category。这 3 个方法共同组成了接口类型Pet的方法集合。

只要一个数据类型的方法集合中有这 3 个方法，那么它就一定是Pet接口的实现类型。**这是一种无侵入式的接口实现方式。**这种方式还有一个专有名词，叫“Duck typing”，中文常译作“鸭子类型”。

#### 实现接口的两个条件

- 两个方法的签名需要完全一致
- 两个方法的名称要一模一样



#### Usecase

```go
package main

import "fmt"

type Pet interface {
	SetName(name string)
	Name() string
	Category() string
}

type Dog struct {
	name string // 名字。
}

func (dog *Dog) SetName(name string) {
	dog.name = name
}

func (dog Dog) Name() string {
	return dog.name
}

func (dog Dog) Category() string {
	return "dog"
}

func main() {
	// 示例1。
	dog := Dog{"little pig"}
	_, ok := interface{}(dog).(Pet)
	fmt.Printf("Dog implements interface Pet: %v\n", ok)
	_, ok = interface{}(&dog).(Pet)
	fmt.Printf("*Dog implements interface Pet: %v\n", ok)
	fmt.Println()

	// 示例2。
	var pet Pet = &dog
	fmt.Printf("This pet is a %s, the name is %q.\n",
		pet.Category(), pet.Name())
}
```

输出结果为：

```go
Dog implements interface Pet: false
*Dog implements interface Pet: true

This pet is a dog, the name is "little pig".
```

Dog类型本身的方法集合中只包含了 2 个方法，也就是所有的值方法。而它的指针类型*Dog方法集合却包含了 3 个方法。

 它拥有Dog类型附带的所有值方法和指针方法。又由于这 3 个方法恰恰分别是Pet接口中某个方法的实现，所以\*Dog类型就成为了Pet接口的实现类型。

可以声明并初始化一个Dog类型的变量dog，然后把它的指针值赋给类型为Pet的变量pet。

```go
dog := Dog{"little pig"}
var pet Pet = &dog
```

对于一个接口类型的变量来说，例如上面的变量pet，**我们赋给它的值可以被叫做它的实际值（也称动态值）**，而**该值的类型可以被叫做这个变量的实际类型（也称动态类型）**。

比如，把取址表达式&dog的结果值赋给了变量pet，这时**这个结果值就是变量pet的动态值**，而此**结果值的类型*Dog就是该变量的动态类型**。

动态类型这个叫法是相对于静态类型而言的。对于变量pet来讲，**它的静态类型就是Pet**，并且**永远是Pet**，但是**它的动态类型却会随着我们赋给它的动态值而变化**。

比如，只有我把一个\*Dog类型的值赋给变量pet之后，该变量的动态类型才会是\*Dog。如果还有一个Pet接口的实现类型\*Fish，并且我又把一个此类型的值赋给了pet，那么它的动态类型就会变为*Fish。

在我们给**一个接口类型的变量赋予实际的值之前，它的动态类型是不存在的**。

 #### Usecase

```go
dog := Dog{"little pig"}
var pet Pet = dog
dog.SetName("monster")
```

这里 pet变量的字段name的值依然是"little pig"

由于dog的SetName方法是指针方法，所以该方法持有的接收者就是指向dog的指针值的副本，因而其中对接收者的name字段的设置就是对变量dog的改动。那么当dog.SetName("monster")执行之后，dog的name字段的值就一定是"monster"。

这里有一条通用的规则需要你知晓：如果我们使用一个变量给另外一个变量赋值，那么真正赋给后者的，并不是前者持有的那个值，而是该值的一个副本。

接口类型本身是无法被值化的。在我们赋予它实际的值之前，它的值一定会是nil，这也是它的零值。

一旦它被赋予了某个实现类型的值，它的值就不再是nil了。不过要注意，即使我们像前面那样把dog的值赋给了pet，pet的值与dog的值也是不同的。这不仅仅是副本与原值的那种不同。

当我们给一个接口变量赋值的时候，该变量的**动态类型**会与**它的动态值**一起被存储在一个专用的数据结构中。

严格来讲，这样一个变量的值其实是这个专用数据结构的一个实例，而不是我们赋给该变量的那个实际的值。所以我才说，pet的值与dog的值肯定是不同的，无论是从它们存储的内容，还是存储的结构上来看都是如此。不过，我们可以认为，这时pet的值中包含了dog值的副本。

在 Go 语言的runtime包中叫iface，是一种专用的数据结构。

iface的实例会包含两个指针，一个是指向类型信息的指针，另一个是指向动态值的指针。这里的类型信息是由另一个**专用数据结构的实例承载的**，其中包含了**动态值的类型，以及使它实现了接口的方法和调用它们的途径**，等等。接口变量被赋予动态值的时候，存储的是包含了这个动态值的副本的一个结构更加复杂的值。



### 接口变量的值为 nil

```go
var dog1 *Dog
fmt.Println("The first dog is nil. [wrap1]")
dog2 := dog1
fmt.Println("The second dog is nil. [wrap1]")
var pet Pet = dog2
if pet == nil {
  fmt.Println("The pet is nil. [wrap1]")
} else {
  fmt.Println("The pet is not nil. [wrap1]")
}
```

这里先声明了一个 *Dog 类型的变量 dog1，并且没有对它进行初始化。这时该变量的值是什么？显然是nil。然后我把该变量赋给了dog2，后者的值此时也必定是nil。把dog2赋给Pet类型的变量pet之后，变量pet的值不为 nil，虽然被包装的动态值是nil，但是pet的值却不会是nil，因为这个动态值只是pet值的一部分而已。这时的pet的动态类型就存在了，是\*Dog。我们可以通过fmt.Printf函数和占位符%T来验证这一点，另外reflect包的TypeOf函数也可以起到类似的作用。

在 Go 语言中，我们把由字面量 nil 表示的值叫做无类型的 nil。这是真正的 nil，因为它的类型也是 nil 的。虽然dog2的值是真正的 nil，但是当我们把这个变量赋给pet的时候，Go 语言会把它的类型和值放在一起考虑。

Go 语言会识别出赋予pet的值是一个*Dog类型的nil。然后，Go 语言就会用一个iface的实例包装它，包装后的产物肯定就不是nil了。

一个有类型的nil赋给接口变量，那么这个变量的值就一定不会是那个真正的nil。因此，当我们使用判等符号 == 判断pet是否与字面量nil相等的时候，答案一定会是false。

如果想让接口为 nil ，要么只声明它但不做初始化，要么直接把字面量nil赋给它。



### 接口之间的组合

接口类型间的嵌入也被称为接口的组合。接口类型间的嵌入要更简单一些，因为它不会涉及方法间的“屏蔽”。只要组合的接口之间有同名的方法就会产生冲突，从而无法通过编译，即使同名方法的签名彼此不同也会是如此。因此，接口的组合根本不可能导致“屏蔽”现象的出现。

与结构体类型间的嵌入很相似，我们只要把一个接口类型的名称直接写到另一个接口类型的成员列表中就可以了。

```go
type Animal interface {
  ScientificName() string
  Category() string
}

type Pet interface {
  Animal
  Name() string
}
```

接口类型Pet包含了两个成员，一个是代表了另一个接口类型的Animal，一个是方法Name的定义。它们都被包含在Pet的类型声明的花括号中，并且都各自独占一行。此时，Animal接口包含的所有方法也就成为了Pet接口的方法。

Go 语言团队鼓励我们声明体量较小的接口，并建议我们通过这种接口间的组合来扩展程序、增加程序的灵活性。这是因为相比于包含很多方法的大接口而言，小接口可以更加专注地表达某一种能力或某一类特征，同时也更容易被组合在一起。



Go 语言标准库代码包io中的`ReadWriteCloser`接口和`ReadWriter`接口就是这样的例子，它们都是由若干个小接口组合而成的。以`io.ReadWriteCloser`接口为例，它是由`io.Reader`、`io.Writer`和`io.Closer`这三个接口组成的。

这三个接口都只包含了一个方法，是典型的小接口。它们中的每一个都只代表了一种能力，分别是读出、写入和关闭。我们编写这几个小接口的实现类型通常都会很容易。

使我们只实现了io.Reader和io.Writer，那么也等同于实现了io.ReadWriter接口，因为后者就是前两个接口组成的。可以看到，这几个io包中的接口共同组成了一个接口矩阵。它们既相互关联又独立存在。

