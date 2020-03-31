## 单链表

单链表在 Go 语言本身并没有在标准库里有相应的实现，我们可以自己设计单链表并实现相关的方法。

我们首先定义单链表中每个节点的数结构：

```go
// 节点数据结构
type Element struct {
	Value interface{}
	next *Element
}
// 创建节点
func New(value interface{}) *Element{
	return &Element{
		Value: value,
		next:  nil,
	}
}
// 获取节点的下一跳 
func (e *Element) Next() *Element{
	return e.next
}
```

这里我定义的链表元素只存放两个属性，当前元素的元素值和后驱指针。

这里其实也可以参考 Go SDK 实现的双向链表，再定义个指向链表结构体的指针，为了简单实现，这里我就不具体定义了。

接着定义链表的数据结构：

```go
// 单链表数据结构
type List struct {
	head,tail *Element
	len,cap int
}
// 初始化单链表
func Init(cap int) *List {
	return &List{
		cap:cap,
		len:0,
		head:nil,
		tail:nil,
	}
}
```

我们有了结构体，那么下面我们可以根据结构体定义一些相应的方法为整个链表实现增减操作。

- 设计查找的方法：

```go
// 根据索引号查找对应元素
func (l *List) FindIndex(index int) interface{}{
	var step = l.head
	if index+1 > l.len{
		return nil
	}
	for i:=0; i<=index; i++{
		step = step.next
	}
	return step.Value
}

// 根据对应值查找元素是否存在
func (l *List) FindValue(value interface{}) bool{
	if l.len == 0 || l.cap == 0 {
		return false
	}
	var step = l.head
	for step != nil{
		if step.Value == value {
			return true
		}
		step = step.next
	}
	return false
}
```

- 设计插入方法

我这里插入分2种情况，第一种，从链表头部插入，第二种从链表尾部插入：

```go
// 表尾插入
func (l *List) InsertAtTail(value interface{}) bool{
	var step = New(value)
  // 溢出判断
	if l.cap > l.len{
    // 空链表判断
		if l.len == 0 {
			l.head = step
			l.tail = step
			l.len++
			return true
		}
		l.tail.next = step
		l.tail = step
		l.len ++
		return true
	}
	return false
}
// 表头插入
func (l *List) InsertAtHead(value interface{}) bool{
	var step = New(value)
  // 溢出判断
	if l.cap > l.len{
    // 空链表判断
		if l.len == 0 {
			l.head = step
			l.tail = step
			l.len ++
			return true
		}
		step.next = l.head
		l.head = step
		l.len ++
		return true
	}
	return false
}
```

- 设计删除方法

```go
func (l *List) Delete(value interface{}) bool{
  // 容量判断
	if l.len == 0 || l.cap == 0 {
		return false
	}
  // 单元素链表处理
	if l.len == 1 {
		if l.head.Value == value{
			l.head = nil
			l.tail = nil
			l.len--
			return true
		}else {
			return false
		}
	}
  // 2个以上元素的单链表处理
	var (
		prev = l.head
		post = l.head.next
	)
	if prev.Value == value{
		l.head = post
		l.len--
		return true
	}
	for post != nil {
		if post.Value == value{
      
			if post.next != nil{
				prev.next = post.next
			}else {
        // 删除的是最后一个节点的处理
				prev.next = nil
				l.tail = prev
			}
			l.len--
			return true
		}
		prev = post
		post = post.next
	}
	return  false
}
```

