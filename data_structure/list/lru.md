### LRU 缓存

设计和构建一个“最近最少使用”缓存，该缓存会删除最近最少使用的项目。缓存应该从键映射到值(允许你插入和检索特定键对应的值)，并在初始化时指定最大容量。当缓存被填满时，它应该删除最近最少使用的项目。

它应该支持以下操作： 获取数据 get 和 写入数据 put 。

- 获取数据 get(key) - 如果密钥 (key) 存在于缓存中，则获取密钥的值（总是正数），否则返回 -1。

- 写入数据 put(key, value) - 如果密钥不存在，则写入其数据值。当缓存容量达到上限时，它应该在写入新数据之前删除最近最少使用的数据值，从而为新的数据值留出空间。



基本的思路是，我们设计一个双向链表，一旦一个节点被访问，我们就把它挪动到链表的表头，也就是说，越往表头，这个节点越经常被访问，如果最少被访问，那么这个节点肯定在表尾。为了提高链表的查找效率，我们可以将链表对象同时用map做存储，克服链表查找效率低下的问题，同时也发挥了链表插入、删除效率高的特性。



我们首先可以设计一个结构体，用于存储缓存节点：

```go
type NodeList struct {
	val int
	key int
	pre *NodeList
	next *NodeList
}

func (c *NodeList) remove() {
	c.next.pre = c.pre
	c.pre.next = c.next
}
```

这里第一个了一个方法用于删除当前节点。

接着，我们设计 LRU 结构体：

```go
type LRUCache struct {
  // 链表表头
	head *NodeList
  // 链表表尾
	tail *NodeList
  // 容量
	capacity int
  // 缓存map
	cache map[int]*NodeList
}

func(c *LRUCache) setHeader(node *NodeList) {
	c.head.next.pre = node
	node.next = c.head.next
	node.pre = c.head
	c.head.next = node

}

// LRU 创建函数
func Constructor(capacity int) LRUCache {
	lru := LRUCache{head: new(NodeList),
					tail: new(NodeList),
					capacity: capacity,
					cache: make(map[int]*NodeList, capacity)}

	lru.tail.pre = lru.head
	lru.head.next = lru.tail
	return lru
}
```

这里的 setHeader 方法用于将当前节点放入链表表头。

这里的创建函数需要注意下，我们创建了一个带表头和表尾的双向链表，表头和表尾不存储元素，只是方便操作。

我们继续设计 LRU 的 Get 方法：

```go
func (this *LRUCache) Get(key int) int {
	if n, ok := this.cache[key]; !ok {
		return -1
	} else {
		n.remove()
		this.setHeader(n)
		return n.val
	}
}
```

如果查找到，就先将n从链表中当前位置移除并挪动到表头。

最后，我们设计 LRU 的 Put 方法：

```go
func (this *LRUCache) Put(key int, value int)  {
    if node, ok := this.cache[key]; ok {
        node.remove()
        delete(this.cache, key)
    } else if len(this.cache) >= this.capacity {
				toRemove := this.tail.pre
				this.tail.pre = nil
        toRemove.next = nil
        delete(this.cache, toRemove.key)
        this.tail = toRemove
	}
	newNode := &NodeList{val:value, key:key}
	this.setHeader(newNode)
	this.cache[key] = newNode
}
```

这里的做法比较巧妙，如果查找到 key 值在链表中已经存在，我们的处理方式是先将其从表中移除，且从map 中移除（防止出现 key 值重复）。接着如果容量已满，则删除表尾元素，将新节点创建并放入表头。

