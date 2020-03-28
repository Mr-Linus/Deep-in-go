package code

import (
	"fmt"
)

// 节点数据结构
type Element struct {
	Value interface{}
	next *Element
}

func New(value interface{}) *Element{
	return &Element{
		Value: value,
		next:  nil,
	}
}

func (e *Element) Next() *Element{
	return e.next
}

// 单链表数据结构
type List struct {
	head,tail *Element
	len,cap int
}

func Init(cap int) *List {
	return &List{
		cap:cap,
		len:0,
		head:nil,
		tail:nil,
	}
}

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

func (l *List) Delete(value interface{}) bool{
	if l.len == 0 || l.cap == 0 {
		return false
	}
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

func (l *List) InsertAtTail(value interface{}) bool{
	var step = New(value)

	if l.cap > l.len{
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

func (l *List) InsertAtHead(value interface{}) bool{
	var step = New(value)
	if l.cap > l.len{
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

func main(){
	l := Init(9)
	fmt.Println(l.FindValue(1))
	fmt.Println(l.InsertAtTail(1))
	fmt.Println(l.Delete(1))

	fmt.Println(l.FindValue(1))
}