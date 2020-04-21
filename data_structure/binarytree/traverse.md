## 二叉树遍历

在二叉树的遍历过程中，除了其实递归，也可以通过迭代（栈、队列）实现。

这里我们首先看下用go 语言实现的栈长什么样子：

```go
type Stack struct {
	inter []*TreeNode
}

func (s *Stack)Push(x *TreeNode){
	s.inter = append(s.inter,x)
}

func (s *Stack)Len()int{
	return len(s.inter)
}

func (s *Stack)Pop() (*TreeNode,bool){
	if len(s.inter) == 0{
		return nil,false
	}
	key := s.inter[len(s.inter)-1]
	if len(s.inter) == 1{
		s.inter = []*TreeNode{}
	}else{
		s.inter = s.inter[:len(s.inter)-1]

	}
	return key,true
}
```





- 先序遍历

递归：

```go
func preorderTraversal(root *TreeNode) []int {
    var count []int
	DFS(root,&count)
	return count
}

func DFS(root *TreeNode,count *[]int){
    if root == nil {
        return
    }
    *count = append(*count,root.Val)
    DFS(root.Left,count)
    DFS(root.Right,count)
}
```

栈（迭代）：

```go
func preorderTraversal(root *TreeNode) []int {
    var res []int
    var cur = root
    var stack Stack
    for cur != nil || stack.Len() != 0 {
        for cur != nil {
            res = append(res,cur.Val)
            stack.Push(cur)
            cur = cur.Left
        }
        cur,_ = stack.Pop()
        cur = cur.Right
    }
    return res
}
```



中序遍历

递归：

```go
func inorderTraversal(root *TreeNode) []int {
    var local []int
    DFS(root,&local)
    return local
}

func DFS(root *TreeNode,count *[]int){
    if root == nil {
        return
    }
    DFS(root.Left,count)
    *count = append(*count,root.Val)
    DFS(root.Right,count)
}
```

栈（迭代）：

```go
type TreeNode struct {
	     Val int
	     Left *TreeNode
	     Right *TreeNode
}


func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var stack Stack
	var res []int
	cur := root
	for cur != nil || stack.Len() != 0{
		for cur != nil {
			stack.Push(root)
			cur = cur.Left
		}
		cur,_ = stack.Pop()
		res = append(res,cur.Val)
		cur = cur.Right
	}
	return res
}
```

后序遍历

```go
func postorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var stack Stack
	var res []int
	cur := root
	stack.Push(cur)
	for stack.Len() != 0{
		cur,_ = stack.Pop()
		res = append(res,cur.Val)
		if cur.Left != nil {
			stack.Push(cur.Left)
		}
		if cur.Right != nil {
			stack.Push(cur.Right)
		}
	}
    s, e := 0, len(res)-1
	for s < e {
		res[s], res[e] = res[e], res[s]
		s++
		e--
	}

	return res
}
```

