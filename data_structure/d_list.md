## 双向链表

在 Go 语言的 SDK 中其实已经帮助我们实现了双向链表，该包位于 `container/list` 中。我们可以看看官方的大佬是怎么写双向链表的：

```go
type Element struct {

	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}

```

首先是定义链表中的节点， 他包含了指向前后两个元素的指针、存放的数据（Value）还有一个指向链表自身的指针。

我们再看看链表自身是个什么样子:

```go
type List struct {
	root Element
	len  int     
}

func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// 判断初始化函数，如果没有被初始化就立即初始化
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

func New() *List { return new(List).Init() }

func (l *List) Len() int { return l.len }
// 获取链表表头元素
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// 获取链表表尾元素
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}
```

我们可以发现 List 这个链表结构其实内部只有两个东西，一个 “前哨” （root）存放了对应的元素和整个链表的长度值。在链表调用初始化函数的时候，我们会发现，它把链表前哨元素的前后指针指向自己，作为一个初始化的流程。这里比较有意思的一点是，这个"前哨"的功能。

你会发现，前哨（root）在初始化的时候并没有给对应的 Value 赋值，开发者的本意其实是利用前哨来存储整个链表的表头和表尾（只用 root 的 prev 和 next 指针）。

前哨的前驱其实是链表的表头，前哨的后继其实是链表的表尾，root 的 `prev` 存放表尾， `next` 存放表头。



我们再回头看看之前的节点的两个方法实现：

```go
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
```

我们就明白他这里的处理是为什么了，还是要加层 `p != &e.list.root` 判断，因为List在初始化的时候会把前哨的前后指针指向自己。也就是说，如果当前元素的前后指针指向了 List 的前哨了，也就表示当前节点不存在前驱元素或者是后继元素。

- Push 操作

在此之前我们先看两个插入函数：

```go
// insert 代表把 e 插到 at 的后面
func (l *List) insert(e, at *Element) *Element {
	n := at.next // 保存at的后继
	at.next = e // 将 at 后继指向 e
	e.prev = at // 将 e 前驱指向 at
	e.next = n // 将 e 后继指向原理at的后继
	n.prev = e // 将 原来at的后继指向 e
	e.list = l // 将 e 的 list 指针指向 l
	l.len++ // 链表数量增加
	return e
}
// insertValue 其实是对 insert 方法的封装
func (l *List) insertValue(v interface{}, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}
```

两个 Push 操作分别表示在头部插入和在尾部插入：

```go
func (l *List) PushFront(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *List) PushBack(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}
```

比如在表头插入元素， root本身的后继就是链表头，所以直接对 root 操作就可以了。

这里还会使用 lazyInit 来判断链表有没有被初始化过，如果没有就立即初始化。

- Move 操作

其他设计的插入操作也大多殊途同归，都是对 `insertValue`  方法的封装，这里我们就简单看一个例子。

```go
func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

func (l *List) move(e, at *Element) *Element {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e

	return e
}
```

- 删除操作

```go
func (l *List) Remove(e *Element) interface{} {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}

func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next 
	e.next.prev = e.prev
	e.next = nil // 避免内存溢出
	e.prev = nil // 避免内存溢出
	e.list = nil
	l.len--
	return e
}
```



通过上面的几种方法，大致对官方写的 List 有了一定的了解，其实很多地方都能看出官方开发者的语言功底，比如 Front 方法和 Back 方法，一旦发现链表的长度为 0, 直接返回 nil 。

在用于删除元素、移动元素，以及一些用于插入元素的方法中，只要判断一下传入的元素中指向所属链表的指针，是否与当前链表的指针相等就可以了。这样做的可以快速判断当前操作节点是否在链表内部，且是否被初始化过（lazyInit）。

而且很多地方做了 lazyInit，即在必要时才初始化链表（申请内存资源），因为 Go 本身也不知道你存的数据到底是个什么东西，如果存放的值很大，那么一锅脑初始化可能会花费很多时间，所以通过 lazyInit去分散初始化的操作带来的计算量和存储空间消耗。

比如集中声明非常多的大容量切片， CPU 和内存空间的使用量肯定都会一个激增，并且只有设法让其中的切片及其底层数组被回收，内存使用量才会有所降低。所以把初始化分散再配合 GC 可以更好的利用内存资源同时减少CPU压力。

