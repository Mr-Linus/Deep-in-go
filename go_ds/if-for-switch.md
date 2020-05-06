## 循环与判断

### range & for


```go
numbers1 := []int{1, 2, 3, 4, 5, 6}
for i := range numbers1 {
  if i == 3 {
    numbers1[i] |= i
  }
}
fmt.Println(numbers1)
```

打印的内容会是[1 2 3 7 5 6]。

在这条for语句中，只有一个迭代变量i。我在每次迭代时，都会先去判断i的值是否等于3，如果结果为true，那么就让numbers1的第i个元素值与i本身做按位或的操作，再把操作结果作为numbers1的新的第i个元素值。最后我会打印出numbers1的值。



当for语句被执行的时候，在range关键字右边的numbers1会先被求值。这个位置上的代码被称为range表达式。range表达式的结果值可以是数组、数组的指针、切片、字符串、字典或者允许接收操作的通道中的某一个，并且结果值只能有一个。

对于不同种类的range表达式结果值，for语句的迭代变量的数量可以有所不同。

当只有一个迭代变量的时候，数组、数组的指针、切片和字符串的元素值都是无处安放的，我们只能拿到按照从小到大顺序给出的一个个索引值。



```go
numbers2 := [...]int{1, 2, 3, 4, 5, 6}
maxIndex2 := len(numbers2) - 1
for i, e := range numbers2 {
  if i == maxIndex2 {
    numbers2[0] += e
  } else {
    numbers2[i+1] += e
  }
}
fmt.Println(numbers2)
```

打印的内容会是[7 3 5 7 9 11],当for语句被执行的时候，在range关键字右边的numbers2会先被求值,需要注意：

- **range表达式只会在for语句开始执行时被求值一次，无论后边会有多少次迭代；**
- **range表达式的求值结果会被复制**，也就是说，被迭代的对象是range表达式结果值的副本而不是原值。

基于这两个规则，我们接着往下看。在第一次迭代时，我改变的是numbers2的第二个元素的值，新值为3，也就是1和2之和。但是，被迭代的对象的第二个元素却没有任何改变，毕竟它与numbers2已经是毫不相关的两个数组了。因此，在第二次迭代时，我会把numbers2的第三个元素的值修改为5，即被迭代对象的第二个元素值2和第三个元素值3的和。



```go
numbers3 := []int{1, 2, 3, 4, 5, 6}
	maxIndex3 := len(numbers3) - 1
	for i, e := range numbers3 {
		if i == maxIndex3 {
			numbers3[0] += e
		} else {
			numbers3[i+1] += e
		}
	}
	fmt.Println(numbers3)
```

打印的内容会是 [22 3 6 10 15 21], **对于这些引用类型的值来说，即使有复制也只会复制一些指针而已，底层数据结构是不会被赋值的**。因此这里的 e对应的其实是底层数组对应的值。

### Switch

#### Usecase

```go
value1 := [...]int8{0, 1, 2, 3, 4, 5, 6}
switch 1 + 3 {
case value1[0], value1[1]:
  fmt.Println("0 or 1")
case value1[2], value1[3]:
  fmt.Println("2 or 3")
case value1[4], value1[5], value1[6]:
  fmt.Println("4 or 5 or 6")
}
```

先声明了一个数组类型的变量value1，该变量的元素类型是int8。在后边的switch语句中，被夹在switch关键字和左花括号{之间的是1 + 3，这个位置上的代码被称为switch表达式。这个switch语句还包含了三个case子句，而每个case子句又各包含了一个case表达式和一条打印语句。

所谓的case表达式一般由case关键字和一个表达式列表组成，表达式列表中的多个表达式之间需要有英文逗号,分割，比如，上面代码中的case value1[0], value1[1]就是一个case表达式，其中的两个子表达式都是由索引表达式表示的。

另外的两个case表达式分别是case value1[2], value1[3]和case value1[4], value1[5], value1[6]。

在这里的每个case子句中的那些打印语句，会分别打印出不同的内容，这些内容用于表示case子句被选中的原因，比如，打印内容0 or 1表示当前case子句被选中是因为switch表达式的结果值等于0或1中的某一个。另外两条打印语句会分别打印出2 or 3和4 or 5 or 6。

**一旦某个case子句被选中，其中的附带在case表达式后边的那些语句就会被执行。与此同时，其他的所有case子句都会被忽略。**

**当然了，如果被选中的case子句附带的语句列表中包含了fallthrough语句，那么紧挨在它下边的那个case子句附带的语句也会被执行。**

正因为存在上述判断相等的操作（以下简称判等操作），switch语句对switch表达式的结果类型，以及各个case表达式中子表达式的结果类型都是有要求的。毕竟，**在 Go 语言中，只有类型相同的值之间才有可能被允许进行判等操作。**

如果switch表达式的结果值是无类型的常量，比如1 + 3的求值结果就是无类型的常量4，那么这个常量会被自动地转换为此种常量的默认类型的值，比如整数4的默认类型是int，又比如浮点数3.14的默认类型是float64。

**因此，由于上述代码中的switch表达式的结果类型是int，而那些case表达式中子表达式的结果类型却是int8，它们的类型并不相同，所以这条switch语句是无法通过编译的。**


```go
value2 := [...]int8{0, 1, 2, 3, 4, 5, 6}
switch value2[4] {
case 0, 1:
  fmt.Println("0 or 1")
case 2, 3:
  fmt.Println("2 or 3")
case 4, 5, 6:
  fmt.Println("4 or 5 or 6")
}
```

其中的变量value2与value1的值是完全相同的。但不同的是，我把switch表达式换成了value2[4]，并把下边那三个case表达式分别换为了case 0, 1、case 2, 3和case 4, 5, 6。

switch表达式的结果值是int8类型的，而那些case表达式中子表达式的结果值却是无类型的常量了。因为，**如果case表达式中子表达式的结果值是无类型的常量，那么它的类型会被自动地转换为switch表达式的结果类型**，**又由于上述那几个整数都可以被转换为int8类型的值，所以对这些表达式的结果值进行判等操作是没有问题的。**

![91add0a66b9956f81086285aabc20c1c](https://static001.geekbang.org/resource/image/91/1c/91add0a66b9956f81086285aabc20c1c.png)

**switch语句会进行有限的类型转换**，**但肯定不能保证这种转换可以统一它们的类型**。还要注意，**如果这些表达式的结果类型有某个接口类型，那么一定要小心检查它们的动态值是否都具有可比性（或者说是否允许判等操作）。**否则，虽然不会造成编译错误，但是后果会更加严重：引发 panic（也就是运行时恐慌）。

#### usecase

```go

value3 := [...]int8{0, 1, 2, 3, 4, 5, 6}
switch value3[4] {
case 0, 1, 2:
  fmt.Println("0 or 1 or 2")
case 2, 3, 4:
  fmt.Println("2 or 3 or 4")
case 4, 5, 6:
  fmt.Println("4 or 5 or 6")
}
```

switch语句在case子句的选择上是具有唯一性的，由于在这三个case表达式中存在结果值相等的子表达式，所以这个switch语句无法通过编译。

子表达式1+1和2不能同时出现，1+3和4也不能同时出现。有了这个约束的约束，我们就可以想办法绕过这个对子表达式的限制了。


```go
value5 := [...]int8{0, 1, 2, 3, 4, 5, 6}
switch value5[4] {
case value5[0], value5[1], value5[2]:
  fmt.Println("0 or 1 or 2")
case value5[2], value5[3], value5[4]:
  fmt.Println("2 or 3 or 4")
case value5[4], value5[5], value5[6]:
  fmt.Println("4 or 5 or 6")
}
```

把case表达式中的常量都换成了诸如value5[0]这样的索引表达式。虽然第一个case表达式和第二个case表达式都包含了value5[2]，并且第二个case表达式和第三个case表达式都包含了value5[4]，但这已经不是问题了。这条switch语句可以成功通过编译。

不过，这种绕过方式对用于类型判断的switch语句（以下简称为类型switch语句）就无效了。因为类型switch语句中的case表达式的子表达式，都必须直接由类型字面量表示，而无法通过间接的方式表示。


```go
value6 := interface{}(byte(127))
switch t := value6.(type) {
case uint8, uint16:
  fmt.Println("uint8 or uint16")
case byte:
  fmt.Printf("byte")
default:
  fmt.Printf("unsupported type: %T", t)
}
```

byte类型是uint8类型的别名类型。它们两个本质上是同一个类型，只是类型名称不同罢了。在这种情况下，这个类型switch语句是无法通过编译的，因为子表达式byte和uint8重复了。

