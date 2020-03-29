## 循环链表

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
```

其实看上去也没有什么不同，Ring 也是链表的一个节点，它也是一个双向链表，有前后两个指针。