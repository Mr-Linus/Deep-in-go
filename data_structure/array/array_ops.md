### 常见的数组处理

#### 1. 最大上升子序列
给定一个无序的整数数组，找到其中最长上升子序列的长度。

示例:

输入: [10,9,2,5,3,7,101,18]
输出: 4 
解释: 最长的上升子序列是 [2,3,7,101]，它的长度是 4。

O(n^2)
```go
func lengthOfLIS(nums []int) int {
    var lens = len(nums)
    var store = make([]int,lens)
    var res int
    if lens == 0 {
        return 0
    }
    for i := range store{
        store[i]=1
    }
    for i := range nums{
        for j:=0;j<i;j++{
            if nums[i] > nums[j]{
                store[i] = Max(store[i],store[j]+1)
            }
        }
        res = Max(res,store[i])
    }
    return res
}

func Max (x,y int) int{
    if x < y{
        return y
    }
    return x
}
```

O(nlogn) 二分查找
```go
func lengthOfLIS(nums []int) int {
    var lens = len(nums)
    var store = make([]int,lens)
    if lens == 0 {
        return 0
    }
    lens = 0
    for _,x := range nums{
        i := BinaryFind(store,x,lens)
        store[i] = x
        if i == lens {
            lens++
        }
    }
    return lens
}

func BinaryFind(store []int,x int,lens int)int{
    low,high := 0,lens-1
    for low <= high{
        mid := (high+low)/2
        if x > store[mid] {
            low = mid +1
        }else if x < store[mid]{
            high = mid - 1
        }else{
            return mid 
        }
    }
    return low
}
```