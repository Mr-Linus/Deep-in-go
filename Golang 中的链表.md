## Golang 中的链表

Go 语言的链表实现在标准库的 [container/list](https://golang.google.cn/pkg/container/list/) 代码包中。这个代码包中有两个公开的程序实体——List和Element，List 实现了一个双向链表（以下简称链表），而 Element 则代表了链表中元素的结构。

```go
// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}
// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}
```

在List包含的方法中，用于插入新元素的那些方法都只接受interface{}类型的值。这些方法在内部会使用Element值，包装接收到的新元素。

这样做正是为了避免直接使用我们自己生成的元素，**主要原因是避免链表的内部关联，遭到外界破坏，这对于链表本身以及我们这些使用者来说都是有益的。**

- List 包里的方法

```go
func (l *List) MoveBefore(e, mark *Element)
func (l *List) MoveAfter(e, mark *Element)

func (l *List) MoveToFront(e *Element)
func (l *List) MoveToBack(e *Element)
func (l *List) Front() *Element
func (l *List) Back() *Element
func (l *List) InsertBefore(v interface{}, mark *Element) *Element
func (l *List) InsertAfter(v interface{}, mark *Element) *Element
func (l *List) PushFront(v interface{}) *Element
func (l *List) PushBack(v interface{}) *Element
```



MoveBefore方法和MoveAfter方法，它们分别用于把给定的元素移动到另一个元素的前面和后面。

MoveToFront方法和MoveToBack方法，分别用于把给定的元素移动到链表的最前端和最后端。

Front和Back方法分别用于获取链表中最前端和最后端的元素。

InsertBefore和InsertAfter方法分别用于在指定的元素之前和之后插入新元素。

PushFront和PushBack方法则分别用于在链表的最前端和最后端插入新元素。

这些方法都会把一个Element值的指针作为结果返回，**它们就是链表留给我们的安全“接口”**。拿到这些内部元素的指针，**我们就可以去调用前面提到的用于移动元素的方法了**。

### 1. 链表的开箱即用

List 和 Element 都是结构体类型。**结构体类型有一个特点，那就是它们的零值都会是拥有特定结构，但是没有任何定制化内容的值，相当于一个空壳。**值中的**字段也都会被分别赋予各自类型的零值。**

> 所谓的**零值就是只做了声明，但还未做初始化的变量被给予的缺省值。**每个类型的**零值都会依据该类型的特性而被设定**。比如，**经过语句var a [2]int声明的变量a的值，将会是一个包含了两个0的整数数组**。又比如，**经过语句var s []int声明的变量s的值将会是一个[]int类型的、值为nil的切片。**



var l list.List声明的变量 l 的值, 这个零值将会是一个长度为0的链表。这个链表持有的根元素也将会是一个空壳，其中只会包含缺省的内容。这样的链表可以被“开箱即用”。

Go 语言标准库中很多结构体类型的程序实体都做到了开箱即用。这也是在编写可供别人使用的代码包（或者说程序库）时，我们推荐遵循的最佳实践之一。

var l list.List声明的链表l可以直接使用，关键在于它的“延迟初始化”机制。

### 2. “延迟初始化”机制

延迟初始化可以理解为把初始化操作延后，仅在实际需要的时候才进行。延迟初始化的优点在于“延后”，它可以分散初始化操作带来的计算量和存储空间消耗。

如果我们需要集中声明非常多的大容量切片的话，那么那时的 CPU 和内存空间的使用量肯定都会一个激增，并且只有设法让其中的切片及其底层数组被回收，内存使用量才会有所降低。

如果数组是可以被延迟初始化的，那么计算量和存储空间的压力就可以被分散到实际使用它们的时候。**这些数组被实际使用的时间越分散，延迟初始化带来的优势就会越明显。**

> 延迟初始化的缺点恰恰也在于“延后”。你可以想象一下，**如果我在调用链表的每个方法的时候，它们都需要先去判断链表是否已经被初始化，那这也会是一个计算量上的浪费**。在这些方法被**非常频繁地调用的情况下，这种浪费的影响就开始显现了，程序的性能将会降低。**

在这里的链表实现中，一些方法是无需对是否初始化做判断的。比如Front方法和Back方法，一旦发现链表的长度为0, 直接返回nil就好了。比如，在用于删除元素、移动元素，以及一些用于插入元素的方法中，只要判断一下传入的元素中指向所属链表的指针，是否与当前链表的指针相等就可以了。

如果不相等，就一定说明传入的元素不是这个链表中的，后续的操作就不用做了。反之，就一定说明这个链表已经被初始化了。

原因在于，链表的PushFront方法、PushBack方法、PushBackList方法以及PushFrontList方法总会先判断链表的状态，并在必要时进行初始化，这就是延迟初始化。

我们在向一个空的链表中添加新元素的时候，肯定会调用这四个方法中的一个，这时新元素中指向所属链表的指针，一定会被设定为当前链表的指针。所以，指针相等是链表已经初始化的充分必要条件。？List利用了自身以及Element在结构上的特点，巧妙地平衡了延迟初始化的优缺点，使得链表可以开箱即用，并且在性能上可以达到最优。

### 3. Ring与List的区别

container/ring包中的Ring类型实现的是一个循环链表，也就是我们俗称的环。其实List在内部就是一个循环链表。它的**根元素永远不会持有任何实际的元素值**，而**该元素的存在就是为了连接这个循环链表的首尾两端**。也可以说，**List的零值是一个只包含了根元素，但不包含任何实际元素值的空链表。**

```go
// A Ring is an element of a circular list, or ring.
// Rings do not have a beginning or end; a pointer to any ring element
// serves as reference to the entire ring. Empty rings are represented
// as nil Ring pointers. The zero value for a Ring is a one-element
// ring with a nil Value.
//
type Ring struct {
	next, prev *Ring
	Value      interface{} // for use by client; untouched by this library
}

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}
// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}
```



那么区别是：

- **Ring类型的数据结构仅由它自身即可代表**，而**List类型则需要由它以及Element类型联合表示。**这是表示方式上的不同，也是结构复杂度上的不同。(个人感觉 Ring 更像 C 语言里的链表)
- 一个Ring类型的值只代表了**其所属的循环链表中的一个元素**，而一个**List类型的值则代表了一个完整的链表**。这是表示维度上的不同。
- **创建并初始化一个Ring值的时，可以指定它包含的元素的数量**，但是对于**一个List值来说却不能这样做**。**循环链表（Ring）一旦被创建，其长度是不可变的。**这是两个代码包中的New函数在功能上的不同，也是两个类型在初始化值方面的第一个不同。
- 仅通过**var r ring.Ring**语句声明的r将会是**一个长度为1的循环链表**，而**List类型的零值则是一个长度为0的链表。**别忘了**List中的根元素不会持有实际元素值**，因此计算长度时不会包含它。这是两个类型在初始化值方面的第二个不同。
- **Ring值的Len方法的算法复杂度是 O(N) 的，而List值的Len方法的算法复杂度则是 O(1) 的。**这是两者在性能方面最显而易见的差别。(List 的len就在List结构体里面)

由于内部结构的不同，也导致各个方法的实现也是不同的。

### 4. 链表与数组的对比

一个链表所占用的内存空间，往往要比包含相同元素的数组所占内存大得多。这是由于链表的元素并不是连续存储的，所以相邻的元素之间需要互相保存对方的指针。不但如此，每个元素还要存有它所属链表的指针。

有了这些关联，链表的结构反倒更简单了。它只持有头部元素（或称为根元素）基本上就可以了。当然了，为了防止不必要的遍历和计算，链表的长度记录在内也是必须的。