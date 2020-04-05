### LFU 缓存

**最不经常使用算法（LFU）：**这个缓存算法使用一个计数器来记录条目被访问的频率。通过使用LFU缓存算法，最低访问数的条目首先被移除。这个方法并不经常使用，因为它无法对一个拥有最初高访问率之后长时间没有被访问的条目缓存负责。

它支持以下操作：get 和 put。

- get(key) - 如果键存在于缓存中，则获取键的值（总是正数），否则返回 -1。
- put(key, value) - 如果键不存在，请设置或插入值。当缓存达到其容量时，则应该在插入新项之前，使最不经常使用的项无效。在此问题中，当存在平局（即两个或更多个键具有相同使用频率）时，应该去除 最近 最少使用的键。

「项的使用次数」就是自插入该项以来对其调用 `get` 和 `put` 函数的次数之和。使用次数会在对应项被移除后置为 0 。

首先我们可以为每个缓存条目设计一个结构体，方便对访问频率做处理：

```go

type DoubleList struct {
	head, tail *Node
}

type Node struct {
	prev, next *Node
  // 存储 KV 还有访问次数
	key, val, freq int
}

func CreateDL() *DoubleList {
	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head
	return &DoubleList {
		head: head,
		tail: tail,
	}
}

func (this *DoubleList) AddFirst(node *Node) {
	node.next = this.head.next
	node.prev = this.head

	this.head.next.prev = node
	this.head.next = node
}

func (this *DoubleList) Remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev

	node.next = nil
	node.prev = nil
}

func (this *DoubleList) RemoveLast() *Node {
	if this.IsEmpty() {
		return nil
	}

	last := this.tail.prev
	this.Remove(last)

	return last
}

func (this *DoubleList) IsEmpty() bool {
	return this.head.next == this.tail
}
```

可以看出了，这里定义了一个双向链表，在初始化时定义了头尾两个节点，但是并不存储数据。

接着实现了头插、删除、删除末尾节点、判空等方法，用于处理缓存对象。

接着，我们定义一个缓存结构体：

```go
type LFUCache struct {
  // 缓存map
	cache map[int]*Node
  // 频率map
	freq map[int]*DoubleList
  // 记录容量、长度、最小频率
	ncap, size, minFreq int
}

func Constructor(capacity int) LFUCache {
	return LFUCache {
		cache: make(map[int]*Node),
		freq: make(map[int]*DoubleList),
		ncap: capacity,
	}
}
```

这里的缓存map记录的是所以 put 进 LFU 的节点，并且map的索引对应的就是 key。

频率map的索引对应的是各个节点的访问频率，相同频率的节点在同一个双向链表里存储。

然后我们定义 ncap, size, minFreq 分别记录 LFU的容量，长度和最小访问频率。

然后我们实现下Get 方法：

```go
func (this *LFUCache) Get(key int) int {
	if node, ok := this.cache[key]; ok {
		this.IncFreq(node)
		return node.val
	}
	return -1
}

func (this *LFUCache) IncFreq(node *Node) {
	_freq := node.freq
  // 将节点先从频率表中移除，因为当前链表可能不仅就这一个节点
	this.freq[_freq].Remove(node)
  // 如果最小频率和该节点频率相同并且缓存频率表为0
  // 即当前节点就是最小访问频率的节点
	if this.minFreq == _freq && this.freq[_freq].IsEmpty() {
    // 将最小缓存加1
		this.minFreq++
    // 哈希表移除该节点
		delete(this.freq, _freq)
	}
	// 当前节点访问次数加一
	node.freq++
  // 如果当前频率表里不存在，创建当前频率的频率链表
	if this.freq[node.freq] == nil {
		this.freq[node.freq] = CreateDL()
	}
  // 将当前节点，插入该频率的频率链表
	this.freq[node.freq].AddFirst(node)
}
```

首先，我们先通过缓存map去判断当前节点是否存在，如果存在，我们需要对频率map进行更新。

这里更新操作略微复杂，我们独立成一个 IncFreq 方法实现。由于我们不能确定当前节点所处的频率map里对应的链表是否只有它一个节点，所以我们先把当前节点从频率map中移除，然后如果当前节点是最小访问频率，并且该访问频率链表已经是空了（即，这个节点就是最小的，且只有它一个），那么将该频率链表删除，并且将当前频率加一同时将最小缓存频率加1，放入到新的频率链表中。



最后我们实现下Put 方法：

```go
func (this *LFUCache) Put(key, value int) {
  // 容量为0，直接返回
	if this.ncap == 0 {
		return
	}
  // 当前的key已经存在，则更新value
	if node, ok := this.cache[key]; ok {
		node.val = value
    // 刷新频率
		this.IncFreq(node)
	} else {
    // 缓存已满
		if this.size >= this.ncap {
      // 删除最小访问次数里的最后一个节点
			node := this.freq[this.minFreq].RemoveLast()
      // 删除该节点对应的数据
			delete(this.cache, node.key)
      // LFU 长度减1
			this.size--
		}
    // 创建新数据
		x := &Node{key: key, val: value, freq: 1}
		this.cache[key] = x
    // 如果访问次数1对应的双向链表不存在，则创建
		if this.freq[1] == nil {
			this.freq[1] = CreateDL()
		}
    // 将新节点插入频率表的链表头
		this.freq[1].AddFirst(x)
    // 将最小访问频率设置为1
		this.minFreq = 1
    // LFU缓存加1
		this.size++
	}
}
```

Put方法首先要判断，容量是否为0，接着判断节点舒服存在，如果存在，那么我们只需要更新值并且刷新一下节点访问频率就好了。

如果不存在，那么我们看看缓存是否已满，如果缓存已满，我们需要删除访问频率最小的链表里的最后一个元素（因为他是最少访问频率，而且是最久没有被访问的），然后我们还需要把它从缓存map中删除。接下来执行的增加节点的操作就和缓存未满一样了，也就是，我们首先创建数据，然后将他的访问频率设置为1，然后将其加入到缓存频率为1的链表中，然后我们将缓存最小频率设置为1，LFU长度加一。