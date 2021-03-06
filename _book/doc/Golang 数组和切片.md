## Golang 数组和切片

数组类型的值（以下简称数组）的长度是**固定**的，而切片类型的值（以下简称切片）是**可变长**的。 

**数组的长度在声明它的时候就必须给定，并且之后不会再改变。**可以说，**数组的长度是其类型的一部分。**比如，[1]string和[2]string就是两个不同的数组类型。而切片的类型字面量中只有元素的类型，而没有长度。切片的长度可以自动地随着其中元素数量的增长而增长，**但不会随着元素数量的减少而减小。**

![edb5acaf595673e083cdcf1ea7bb966c](https://static001.geekbang.org/resource/image/ed/6c/edb5acaf595673e083cdcf1ea7bb966c.png)

其实可以把切片看做是对数组的一层简单的封装，因为在**每个切片的底层数据结构中，一定会包含一个数组。** 数组可以被叫做切片的底层数组 ，而**切片也可以被看作是对数组的某个连续片段的引用。**

#### 关于传值和引用

> Go 语言的**切片类型属于引用类型**，**同属引用类型的还有字典类型、通道类型、函数类型等**；而 Go 语言的**数组类型则属于值类型**，同属值类型的有基础数据类型以及**结构体类型**。

Go 语言里不存在像 Java 等编程语言中令人困惑的“传值或传引用”问题。在 Go 语言中，我们判断所谓的“传值”或者“传引用”只要看被传递的值的类型就好了。

如果传递的值是引用类型的，那么就是“传引用”。如果传递的值是值类型的，那么就是“传值”。从**传递成本的角度讲，引用类型的值往往要比值类型的值低很多。**

数组和切片之上都可以应用索引表达式，得到的都会是某个元素。我们在它们之上也都可以应用切片表达式，也都会得到一个新的切片。

### 1.切片的长度和容量

```go
package main

import "fmt"

func main() {
  s1 := make([]int, 5)
  fmt.Printf("The length of s1: %d\n", len(s1)) 
  fmt.Printf("The capacity of s1: %d\n", cap(s1)) 
  fmt.Printf("The value of s1: %d\n", s1) 
  s2 := make([]int, 5, 8)
  fmt.Printf("The length of s2: %d\n", len(s2))
  fmt.Printf("The capacity of s2: %d\n", cap(s2))
  fmt.Printf("The value of s2: %d\n", s2)
}
```

内建函数make声明了一个[]int类型的变量s1。传给make函数的第二个参数是5，从而指明了该切片的长度。我用几乎同样的方式声明了切片s2，只不过多传入了一个参数8以指明该切片的容量。

切片s1和s2的容量分别是5和8，如果不指明其容量，那么它就会和长度一致。如果在初始化时指明了容量，那么切片的实际容量也就是它了。

可以把切片看做是对数组的一层简单的封装，因为在每个切片的底层数据结构中，一定会包含一个数组。数组可以被叫做切片的底层数组，而切片也可以被看作是对数组的某个连续片段的引用。

在这种情况下，**切片的容量实际上代表了它的底层数组的长度**，这里是8。**（注意，切片的底层数组等同于我们前面讲到的数组，其长度不可变。）**

```go
s3 := []int{1, 2, 3, 4, 5, 6, 7, 8}
s4 := s3[3:6]
fmt.Printf("The length of s4: %d\n", len(s4))
fmt.Printf("The capacity of s4: %d\n", cap(s4))
fmt.Printf("The value of s4: %d\n", s4)
```

[3:6]要表达的就是透过新窗口能看到的s3中元素的索引范围是从3到5（注意，不包括6）。

这里的3可被称为起始索引，6可被称为结束索引。那么s4的**长度**就是6减去3，即3。因此可以说，s4中的索引从0到2指向的元素对应的是s3及其底层数组中索引从3到5的那 3 个元素。

![96e2c7129793ee5e73a574ef8f3ad755](https://static001.geekbang.org/resource/image/96/55/96e2c7129793ee5e73a574ef8f3ad755.png)

切片的容量代表了它的底层数组的长度，但这仅限于使用make函数或者切片值字面量初始化切片的情况。

**一个切片的容量可以被看作是透过这个窗口最多可以看到的底层数组中元素的个数。**

由于s4是通过在s3上施加切片操作得来的，所以s3的底层数组就是s4的底层数组。又因为，在底层数组不变的情况下，**切片代表的窗口可以向右扩展**，直至其底层数组的末尾。所以，s4的容量就是其底层数组的长度8, 减去上述切片表达式中的那个起始索引3，即5。(因为切片只能向右扩展，所以左面的123看不见)

注意，切片代表的窗口是无法向左扩展的。也就是说，我们永远无法透过s4看到s3中最左边的那 3 个元素。

切片的窗口向右扩展到最大的方法：对于s4来说，切片表达式s4[0:cap(s4)]就可以做到。我想你应该能看懂。该表达式的结果值（即一个新的切片）会是[]int{4, 5, 6, 7, 8}，其长度和容量都是5。

### 2. 切片容量的增长

一旦一个切片无法容纳更多的元素，**Go 语言就会想办法扩容**。但**它并不会改变原来的切片**，而是会生成一个容量更大的切片，然后将把原有的元素和新元素一并拷贝到新切片中。在一般的情况下，你可以简单地认为新切片的容量（以下简称新容量）**将会是原切片容量（以下简称原容量）的 2 倍**。

但是，**当原切片的长度（以下简称原长度）大于或等于 1024 时，Go 语言将会以原容量的1.25倍作为新容量的基准（以下新容量基准）。**新容量基准会被调整（不断地与1.25相乘），**直到结果不小于原长度与要追加的元素数量之和（以下简称新长度）。**最终，**新容量往往会比新长度大一些，当然，相等也是可能的。**

**如果我们一次追加的元素过多，以至于使新长度比原容量的 2 倍还要大，那么新容量就会以新长度为基准。**注意，与前面那种情况一样，**最终的新容量在很多时候都要比新容量基准更大一些。**更多细节可参见runtime包中 slice.go 文件里的growslice及相关函数的具体实现。



### 3. 切片的底层数组的替换

一个切片的底层数组永远不会被替换。**虽然在扩容的时候 Go 语言一定会生成新的底层数组，但是它也同时生成了新的切片**。

它只是把**新的切片作为了新底层数组的窗口**，而**没有对原切片，及其底层数组做任何改动**。

在无需扩容时，**append函数返回的是指向原底层数组的新切片**，而在**需要扩容时，append函数返回的是指向新底层数组的新切片**。

所以，严格来讲，“扩容”这个词用在这里虽然形象但并不合适。不过鉴于这种称呼已经用得很广泛了，我们也没必要另找新词了。

只要新长度不会超过切片的原容量，那么使用append函数对其追加元素的时候就不会引起扩容。**这只会使紧邻切片窗口右边的（底层数组中的）元素被新的元素替换掉**。



### 4. 切片与数组的优劣

切片本身有着占用内存少和创建便捷等特点，但它的本质上还是数组。切片的一大好处是可以让我们通过窗口快速地定位并获取，或者修改底层数组中的元素。

但是，删除切片中的元素就很困难。元素复制一般是免不了的，就算只删除一个元素，有时也会造成大量元素的移动。这时还要注意空出的元素槽位的“清空”，否则很可能会造成内存泄漏。

另一方面，在切片被频繁“扩容”的情况下，新的底层数组会不断产生，这时内存分配的量以及元素复制的次数可能就很可观了，这肯定会对程序的性能产生负面的影响。

尤其是当我们没有一个合理、有效的”缩容“策略的时候，旧的底层数组无法被回收，新的底层数组中也会有大量无用的元素槽位。过度的内存浪费不但会降低程序的性能，还可能会使内存溢出并导致程序崩溃。

它们都是 Go 语言原生的数据结构，使用起来也都很方便. 不过，你的集合类工具箱中不应该只有它们。这就是我们使用**链表**的原因。