##  循环链表

同样的，循环链表本身在 Go 的 SDK 里已经实现了，位于 `container/ring` 包中。

我们可以观摩下官方实现的循环链表：

```go
type Ring struct {
	next, prev *Ring
	Value      interface{} // for use by client; untouched by this library
}

func (r *Ring) init() *Ring {
	r.next = r
	r.prev = r
	return r
}

// Next returns the next ring element. r must not be empty.
func (r *Ring) Next() *Ring {
	if r.next == nil {
		return r.init()
	}
	return r.next
}
// Prev returns the previous ring element. r must not be empty.
func (r *Ring) Prev() *Ring {
	if r.next == nil {
		return r.init()
	}
	return r.prev
}
```

其实看上去也没有什么不同，Ring 也是链表的一个节点，它既是一个循环链表，也是一个双向链表，有前后两个指针。

和官方 SDK 定义的双向链表相比，其实非常相似，Ring 这里就可以看成是双向链表里面的 Element，也就是链表中的节点，但是这个 Ring 节点并没有指向一个独立的链表结构体，也就是说，循环链表中并没有定义像双向链表中的 List 的结构体。



- 循环链表中的的创建：

```go
// New creates a ring of n elements.
func New(n int) *Ring {
	if n <= 0 {
		return nil
	}
	r := new(Ring)
	p := r
	for i := 1; i < n; i++ {
		p.next = &Ring{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}
```

创建和双向链表相比也显得十分的简洁，其中参数 n 代表链表元素的个数，最后返回的是创建好的链表的链表首元素指针（这里也是比较罕见的用 new 去创建一个结构体，赶紧拿笔记下来）。



- 获取链表长度

```go
// Len computes the number of elements in ring r.
// It executes in time proportional to the number of elements.
//
func (r *Ring) Len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}
```

这里其实就是对链表的遍历啦，通过创建一个 n 作为计数器，遍历节点最终计算出链表的长度。那么可以看出来，这个计算链表的时间复杂度是 O（n），如果链表中元素较多的时候，这个操作还是比较耗时的，建议在使用的时候最好在外部把链表的长度先存起来，做修改的时候直接操作外部的值，以免造成额外的开销。



- 链表中节点的获取（移动）

```go
// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.
//
func (r *Ring) Move(n int) *Ring {
	if r.next == nil {
		return r.init()
	}
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}
```

Move 方法通过传入参数获取到指针移动的次数和参数，这个比较好玩的是，他实现了往前移动和往后移动，取负值就是往后移动。



- 链表的连接和解连接

```go
func (r *Ring) Link(s *Ring) *Ring {
	n := r.Next()
	if s != nil {
		p := s.Prev()
		// Note: Cannot use multiple assignment because
		// evaluation order of LHS is not specified.
		r.next = s
		s.prev = r
		n.prev = p
		p.next = n
	}
	return n
}

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed subring. r must not be empty.
//
func (r *Ring) Unlink(n int) *Ring {
	if n <= 0 {
		return nil
	}
	return r.Link(r.Move(n + 1))
}
```

这个 Link 方法其实就是把 s 链表接到 r 链表的后面。

在解连接的时候，传入的参数是 n，其实是删除 n % r.Len() 个元素，从 r 的下一个点开始，到 r 往后的第 n 个节点，这里非常巧妙通过 move 指向 r 往后的第 n + 1 个节点，然后将 r 节点和后面的第 n + 1 个节点连接，完成整个删除的操作。

不过，需要注意的是，r 本身不能为空，并且传入参数的 n 如果和链表本身长度相同，那么这个解连接会失效（相当于什么都没做，因为这个时候 n+1 就是自己的下一个节点）。



- Do 方法（遍历）

```GO
func (r *Ring) Do(f func(interface{})) {
	if r != nil {
		f(r.Value)
		for p := r.Next(); p != r; p = p.next {
			f(p.Value)
		}
	}
}
```

其实就是实现对循环链表遍历的操作，你可以通过传入 f 函数实现一些定制化的功能。



### 总结

通过看 ring 的源码，与双向链表相比，我们还是可以发现 ring 更加精简，实现的复杂性明显降低的很多，实现特别简单，但是在计算链表长度时可能会比较耗时，建议在使用循环链表的时候最好在外部把链表的长度先存起来，做修改的时候直接操作外部的值，以免造成额外的开销。

