## Go 数据结构 ——排序

### 1. 插入排序

- 直接插入排序 （内存6.3MB，用时 944ms）

插排的思想是排序的排到底i个的时候，保证前面下标是0~(i-1)的序列是有序的。

```go
func sortArray(nums []int) []int {
    for i,v := range nums {
        if i == 0{
            continue
        }
      	// 左子序列判断+挪位子
        for j := i-1 ; j>=0 ; j--{
            if v < nums[j]{
                nums[j],nums[j+1] = nums[j+1], nums[j]
            }else{
                break
            }
        }
    }
    return nums
}
```

> 时间复杂度 O(n^2) 空间复杂度O（1）

- 折半(二分)插入排序 （内存6.3MB，用时452ms）

直接插入排序是边查边挪，折半插入排序是先用二分查找找到插入点，然后统一挪位置

```go
func sortArray(nums []int) []int {
    var lens = len(nums)
    if lens <= 1 {
        return nums
    }
    for i,v := range nums{
        if i == 0 {
            continue
        }
        low := 0
        high := i-1
      	// 二分查找
        for low <= high {
            mid := (low+high)/2
            if (nums[mid]> v) {
                high = mid - 1
            }else {
                low = mid +1
            }
        }
      	// 往后挪
        for j := i-1 ; j>high ; j--{
            nums[j+1] = nums[j]
        }
        // 插入
        nums[high+1]=v
        // 重置高低下标
        low = 0
        high = i-1
    }
    return nums
}
```

> 比较次数降至 O(nlog2(n)）
>
> 时间复杂度 O(n^2) 空间复杂度O（1）

- 希尔排序 （内存6.3MB，用时28ms）

以步长为单位，将下表为i的与i+d的元素比较，若i元素大于i+d元素，则交换，每一趟之后再将步长缩小为一般再次进行比较。

```go
func sortArray(nums []int) []int {
    // 外层设置步长
    for d := len(nums)/2; d > 0; d /= 2 {
      	// 中层设置趟数
        for i := range nums{
            for j:=i; j-d >=0; j-=d{
                if nums[j] < nums[j-d]{
                    nums[j],nums[j-d] = nums[j-d],nums[j]
                }else{
                    break
                }
            }
        }
    }
    return nums
}
```

> 时间复杂度 O(n^2) 空间复杂度 O(1)

### 2. 交换排序

- 冒泡排序

从后往前或者两两比较交换次序，每一趟把最小的放在队首，或者将最大的队尾，每一趟之后将排序元素数减一进行多次交换 （小切片代码还行，大的直接超出时间限制）

```go
func sortArray(nums []int) []int {
	var lens = len(nums)
	for i := range nums{
		for j:=i+1; j<lens; j++ {
			if nums[j] < nums[i] {
				nums[j],nums[i] = nums[i],nums[j]
			}
		}
	}
	return nums
}
```

> 时间复杂度 O(n^2)  空间复杂度O(1)

- 快速排序

先找个枢纽元素，默认为下标为0的元素，然后数组头尾设置标记点和枢纽元素比，找到2个元素，前面的比枢纽大，后面的比枢纽小，就交换。每一趟下来，枢纽元素的左侧都是比枢纽元素小的，右侧都是比枢纽元素大的，当头尾标记移动到相同位置时，此位置即为枢纽的位置。然后再对左右两个子序列进行递归。

递归方法：（内存6.4MB，用时24ms）

```go
func sortArray(nums []int) []int {
    quickSort(nums)
    return nums
}

func quickSort(nums []int) {
    left, right := 0, len(nums) - 1
    for right > left {
        // 右边部分放大于
        if nums[right] > nums[0] {
            right--
            continue
        }
        // 左边部分放小于等于
        if nums[left] <= nums[0] {
            left++
            continue
        }
        nums[left], nums[right] = nums[right], nums[left]
    }
    nums[0], nums[right] = nums[right], nums[0]
    if len(nums[:right]) > 1 {
        sortArray(nums[:right])
    }
    if len(nums[right + 1:]) > 1 {
        sortArray(nums[right + 1:])
    }
}
```

> 时间复杂度 O(nlogn) 空间复杂度O（1）

### 3.选择排序

- 简单选择排序 每趟一次交换 (内存6.4MB，用时1456 ms)

每一趟都和第 (i+1) 个到第 n 个比较，在和把最小的元素和i交换。

```go
func sortArray(nums []int) []int {
	var (
		lens = len(nums)
		index int
		min int
	)
	for i,v := range nums{
		index = i
		min = v
		for j:=i+1; j<lens; j++ {
			if nums[j] < min {
				index = j 
				min = nums[j]
			}
		}
		if index != i {
			nums[i],nums[index] = nums[index],nums[i]
		}
	}
	return nums
}
```

- 堆排序 （内存6.4MB，用时28ms）

```go
func sortArray(nums []int) []int {
  	// 先构建最大堆
    for i:= len(nums)/2; i >=0; i--{
        MaxHead(nums,i)
    }
    var end = len(nums)-1
    for end > 0 {
       	// 每次把堆顶元素（切片下标1）和最后一个叶子节点互换
        nums[0],nums[end] = nums[end],nums[0]
        // 把最后一个叶子节点（上次排序的最大值）扣去，重新排序      
        MaxHead(nums[:end],0)
        end--
    }
    
    return nums
}

func MaxHead(nums []int,pos int){
  var (
      lens = len(nums) - 1 
      left = pos * 2 + 1
      right = left + 1
      step = left
  )
  if left > lens{
		return
	}
  // 右节点存在的处理
	if right <= lens {
        if  nums[right] > nums[left]{
            step = right
        }
	}
  // 子节点大就交换
	if nums[pos] < nums[step]{
		nums[pos],nums[step] = nums[step],nums[pos]
    // 递归子节点
		MaxHead(nums,step)
	}
  return
}

```

### 4.归并排序

- 归并排序 （执行用时 32 ms，内存消耗 7.4 MB）

```go
// 分治
func sortArray(nums []int) []int {
    var lens = len(nums)
    // 拆出来可能是单个的处理下
    if lens <= 1{
        return nums
    }
    // 比较+交换
    if lens == 2{
        if nums[0] > nums[1]{
            nums[1],nums[0] = nums[0],nums[1]
        }
        return nums
    }
    return Merge(sortArray(nums[:lens/2]),sortArray(nums[lens/2:]))
}
// 合并
func Merge(Pre,Post []int) []int{
    var lenPre,lenPost = len(Pre),len(Post)
    var list []int
    var i,j int 
    // 比较+合并
    for i < lenPre && j < lenPost{
        if Pre[i] <= Post[j]{
            list = append(list,Pre[i])
            i++
        }else{
            list = append(list,Post[j])
            j++
        }
    }
    // 多出来数的直接接到后面
    if i < lenPre{
        list = append(list,Pre[i:]...)
    }
    if j < lenPost{
        list = append(list,Post[j:]...)
    }
    return list
}
```

> 时间复杂度 O(nlogn) 
>
> 空间复杂度 O(n)
>
> Tips:
>
> 归并排序不是原地排序算法。

