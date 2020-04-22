### Slice 切片



#### 1.切片数据结构

Go 语言中切片数据结构在源码包 src 的 runtime/slice.go

```go
type slice struct {
	array unsafe.Pointer  // 数据部分
	len   int          	// 长度
	cap   int             // 容量
}
```



#### 2.切片扩容

当 Go 中切片 append 当容量超过了现有容量，就需要进行扩容

1. 确定扩容的大小

```go
func growslice(et *_type, old slice, cap int) slice {
	// 得到旧容量大小
	newcap := old.cap
	// 扩容容量 = 旧容量 * 2
	doublecap := newcap + newcap
	// cap是新申请容量
    // 如果申请容量(cap)溢出（大于doublecap），就直接申请
	if cap > doublecap {
		newcap = cap
    // 不足，就再判断旧切片长度大于等于1024，则最终容量（newcap）从旧容量（old.cap）开始循环增加原来的1/4，即（newcap=old.cap,for {newcap += newcap/4}）直到最终容量（newcap）大于等于新申请的容量(cap)，即（newcap >= cap）
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
}
```

2. 根据扩容大小和切片类型，确定不同的内存分配大小，同时保证内存的对齐。因此，申请的内存可能会大于实际的 `et.size * newcap`

```go
   	var overflow bool
   	var lenmem, newlenmem, capmem uintptr
   	// Specialize for common values of et.size.
   	// For 1 we don't need any division/multiplication.
   	// For sys.PtrSize, compiler will optimize division/multiplication into a shift by a constant.
   	// For powers of 2, use a variable shift.
   	switch {
       // 根据et.size大小执行不同内存大小计算逻辑，这里不赘述，知道默认如何分配即可
   	case et.size == 1:
   		...
   	case et.size == sys.PtrSize:
   		...
   	case isPowerOfTwo(et.size):
   		...
   	default:
   		lenmem = uintptr(old.len) * et.size
   		newlenmem = uintptr(cap) * et.size
   		capmem, overflow = math.MulUintptr(et.size, uintptr(newcap))
   		capmem = roundupsize(capmem)
   		newcap = int(capmem / et.size)
	}

```


3. 最后核心是申请内存。要注意的是，新的切片不一定意味着新的地址。

```go

   	if overflow || capmem > maxAlloc {
   		panic(errorString("growslice: cap out of range"))
   	}
   	// 针对指针
   	var p unsafe.Pointer
   	if et.ptrdata == 0 {
   		p = mallocgc(capmem, nil, false)
   		// The append() that calls growslice is going to overwrite from old.len to cap (which will be the new length).
   		// Only clear the part that will not be overwritten.
   		memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
   	} else {
   		// Note: can't use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.
   		p = mallocgc(capmem, et, true)
   		if lenmem > 0 && writeBarrier.enabled {
   			// Only shade the pointers in old.array since we know the destination slice p
   			// only contains nil pointers because it has been cleared during alloc.
   			bulkBarrierPreWriteSrcOnly(uintptr(p), uintptr(old.array), lenmem)
   		}
   	}
   	// 高效拷贝字节指令，from old.array to p
   	memmove(p, old.array, lenmem)
   
	return slice{p, old.len, newcap}

```


#### 切片复制

- 浅拷贝 

  1. 先将 data 的成员数据拷贝到寄存器，然后从寄存器拷贝到 shallowCopy 的对象中。

  2. 注意到只是拷贝了指针而已, 所以是浅拷贝

```go

  shallowCopy := data[:1] // 默认的浅拷贝方式，Go直接汇编操作没有源码

```

- 深拷贝

  1. 检查切片长度与元素大小

  2. 静态分析和内存扫描

  3. 遍历内存，去复制每个元素

```go

  func slicecopy(to, fm slice, width uintptr) int {  if fm.len == 0 || to.len == 0 {    return 0
    }
  
    n := fm.len
    if to.len < n {
      n = to.len
    }
    // 元素大小为0，则直接返回
    if width == 0 {
        return n
    }
    // 静态分析和内存扫描
    // 直接内存拷贝
    size := uintptr(n) * width
    if size == 1 { // common case worth about 2x to do here
      *(*byte)(to.array) = *(*byte)(fm.array) // known to be a byte pointer
    } else {
      memmove(to.array, fm.array, size)
    }  return n
  }
  
  // 针对字符串slice的拷贝
  func slicestringcopy(to []byte, fm string) int {  
      if len(fm) == 0 || len(to) == 0 {
          return 0
    	}
  
  	n := len(fm)  if len(to) < n {
      	n = len(to)
    	}    
      // 静态分析和内存扫描
      // ...
  
   	memmove(unsafe.Pointer(&to[0]), stringStructOf(&fm).str, uintptr(n))  
      return n
  }
```
