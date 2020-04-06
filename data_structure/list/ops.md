### 链表的常见操作

#### 1.链表的反转

一般也叫链表的逆序，比如

```
输入: 1->2->3->4->5->NULL
输出: 5->4->3->2->1->NULL
```

算法：

```go
func reverseList(head *ListNode) *ListNode {
  	// 设置前驱节点
    var prev = head
  	// 处理链表为空或只有一个元素的情况
    if head == nil {
        return head
    }
    if head.Next == nil{
        return head
    }
  	// 设定当前指针
    cur := head.Next
    for  cur.Next != nil{
      	// 设置后驱指针
        post := cur.Next
      	// 将当前指针的下一跳指向前驱节点
        cur.Next = prev
      	// 将前驱设置为当前节点
        prev = cur
      	// 将当前节点为后驱（后移）
        cur = post
      	// 后驱指针后移
        post = post.Next
    }
  	// 处理最后一个节点
    cur.Next = prev
  	// 处理原来的投节点（即当前的尾节点）
    head.Next = nil
    return cur
}
```



#### 2. 合并两个链表

将两个升序链表合并为一个新的升序链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。 

**示例：**

```
输入：1->2->4, 1->3->4
输出：1->1->2->3->4->4
```

算法：

```go
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    if l1 == nil {
        return l2
    }
    if l2 == nil {
        return l1
    }
    var res *ListNode
    // 当l1节点的值大于l2节点的值，那么res指向l2的节点，从l2开始遍历，反之从l1开始
    if l1.Val >= l2.Val {
        res = l2
        res.Next = mergeTwoLists(l1, l2.Next)
    } else {
        res = l1
        res.Next = mergeTwoLists(l1.Next, l2)
    }
    return res
}
```



#### 3. 链表判环

给定一个链表，判断链表中是否有环。

算法：

```go
func hasCycle(head *ListNode) bool {
    fast,slow := head,head
    for fast != nil && slow != nil && fast.Next != nil {
        fast = fast.Next.Next
        slow = slow.Next
        if fast==slow {
            return true
        }
    }
    return false
}
```

快慢指针的思想：

- 每次快指针走2步，慢指针走一步

- 如果链表里有环，那么快指针和慢指针一定对相遇（指向相同的节点）

这里需要考虑下边界条件，如果快指针、慢指针或者快指针的下一跳为空则表示当前链表无环。



#### 4. 链表的中间节点

给定一个带有头结点 `head` 的非空单链表，返回链表的中间结点。

如果有两个中间结点，则返回第二个中间结点。

算法：

```go
func middleNode(head *ListNode) *ListNode {
    var step *ListNode
  	// 获取链表长度
    var lens int
    step = head
    for step != nil {
        lens++
        step = step.Next
    }
  	// 折半遍历
    step = head
    for i:=1;i<=lens/2;i++{
        step = step.Next
    }
    return step
}
```



#### 5. 删除中间节点

实现一种算法，删除单向链表中间的某个节点（除了第一个和最后一个节点，不一定是中间节点），假定你只能访问该节点。

思想：

不能获取链表头，那我们可以通过修改值实现中间节点的删除。

算法：

```go
func deleteNode(node *ListNode) {
    var head = node
    var lens int
    // 后面的元素值赋给自己
    for head.Next != nil {
        lens ++
        head.Val = head.Next.Val
        head = head.Next
    }
    head = node
    // 删除最后一个节点
    for i:=0; i<lens-1; i++{
        head = head.Next
    }
    head.Next = nil
}
```



#### 6. 两个链表的第一个公共节点

输入两个链表，找出它们的第一个公共节点。

如下面的两个链表**：**

[![img](https://assets.leetcode-cn.com/aliyun-lc-upload/uploads/2018/12/14/160_statement.png)](https://assets.leetcode-cn.com/aliyun-lc-upload/uploads/2018/12/14/160_statement.png)

在节点 c1 开始相交。

思想：

我们可以先计算两个链表的长度，然后让长的链表多走几步，走到两个链表后部分长度一样时，一个一个比较，直到找到相交的第一个元素。

算法：

```go
func getIntersectionNode(headA, headB *ListNode) *ListNode {
    var (
        l1,l2,lose int
        ha = headA
        hb = headB
    )
    for h1:= headA; h1!=nil; h1=h1.Next{
        l1++
    }
    for h2:= headB; h2!=nil; h2=h2.Next{
        l2++
    }
    if l1 > l2 {
        lose = l1 - l2
        for i:=0;i< lose;i++ {
            ha = ha.Next
        }
    }else {
        lose = l2 - l1
        for i:=0;i< lose;i++ {
            hb = hb.Next
        }
    }
    for ha != nil && hb != nil{
        if ha == hb {
            return ha
        }
        ha = ha.Next
        hb = hb.Next
    }
    return nil
}
```

