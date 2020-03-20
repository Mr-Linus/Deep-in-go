## Go Context

在使用WaitGroup值的时候，我们最好用“先统一Add，再并发Done，最后Wait”的标准模式来构建协作流程。

如果在调用该值的Wait方法的同时，为了增大其计数器的值，而并发地调用该值的Add方法，那么就很可能会引发 panic。如果我们不能在一开始就确定执行子任务的 goroutine 的数量，那么使用WaitGroup值来协调它们和分发子任务的 goroutine，就是有一定风险的。一个解决方案是：**分批地启用执行子任务的 goroutine。**

WaitGroup值是可以被复用的，但需要保证其计数周期的完整性。尤其是涉及对其Wait方法调用的时候，它的下一个计数周期必须要等到，与当前计数周期对应的那个Wait方法调用完成之后，才能够开始。

前面提到的可能会引发 panic 的情况，就是由于没有遵循这条规则而导致的。

在严格遵循上述规则的前提下，分批地启用执行子任务的 goroutine，就肯定不会有问题。具体的实现方式有不少，其中最简单的方式就是使用for循环来作为辅助。这里的代码如下：

```go

func coordinateWithWaitGroup() {
 total := 12
 stride := 3
 var num int32
 fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
 var wg sync.WaitGroup
 for i := 1; i <= total; i = i + stride {
  wg.Add(stride)
  for j := 0; j < stride; j++ {
   go addNum(&num, i+j, wg.Done)
  }
  wg.Wait()
 }
 fmt.Println("End.")
}
```

经过改造后的coordinateWithWaitGroup函数，循环地使用了由变量wg代表的WaitGroup值。它运用的依然是“先统一Add，再并发Done，最后Wait”的这种模式，只不过它利用for语句，对此进行了复用。

用context包中的函数和Context类型作为实现工具。

```go
func coordinateWithContext() {
 total := 12
 var num int32
 fmt.Printf("The number: %d [with context.Context]\n", num)
 cxt, cancelFunc := context.WithCancel(context.Background())
 for i := 1; i <= total; i++ {
  go addNum(&num, i, func() {
   if atomic.LoadInt32(&num) == int32(total) {
    cancelFunc()
   }
  })
 }
 <-cxt.Done()
 fmt.Println("End.")
}
func addNum(numP *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		newNum := currNum + 1
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			break
		} else {
			//fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}
```

在这个函数体中，我先后调用了 context.Background 函数和 context.WithCancel 函数，并得到了一个可撤销的context.Context类型的值（由变量cxt代表），以及一个context.CancelFunc类型的撤销函数（由变量cancelFunc代表）。

在后面那条唯一的for语句中，我在每次迭代中都通过一条go语句，异步地调用addNum函数，调用的总次数只依据了total变量的值。

如果两个值相等，那么就调用cancelFunc函数。其含义是，**如果所有的addNum函数都执行完毕，那么就立即通知分发子任务的 goroutine。**

这里分发子任务的 goroutine，即为执行coordinateWithContext函数的 goroutine。它在执行完for语句后，会立即调用cxt变量的Done函数，并试图针对该函数返回的通道，进行接收操作。

由于一旦cancelFunc函数被调用，针对该通道的接收操作就会马上结束，所以，这样做就可以实现“等待所有的addNum函数都执行完毕”的功能。

Context类型之所以受到了标准库中众多代码包的积极支持，主要是因为它是一种非常通用的同步工具。它的值不但可以被任意地扩散，而且还可以被用来传递额外的信息和信号。

**Context类型可以提供一类代表上下文的值。此类值是并发安全的，也就是说它可以被传播给多个 goroutine。**

由于Context类型实际上是一个接口类型，而context包中实现该接口的所有私有类型，都是基于某个数据类型的指针类型，所以，如此传播并不会影响该类型值的功能和安全。

Context类型的值（以下简称Context值）是可以繁衍的，这意味着我们可以通过一个Context值产生出任意个子值。这些子值可以携带其父值的属性和数据，也可以响应我们通过其父值传达的信号。

正因为如此，所有的Context值共同构成了一颗代表了上下文全貌的树形结构。这棵树的树根（或者称上下文根节点）是一个已经在context包中预定义好的Context值，它是全局唯一的。通过调用context.Background函数，我们就可以获取到它。

这个上下文根节点仅仅是一个最基本的支点，它不提供任何额外的功能。也就是说，它既不可以被撤销（cancel），也不能携带任何数据。

context包中还包含了四个用于繁衍Context值的函数，即：WithCancel、WithDeadline、WithTimeout和WithValue。

这些函数的第一个参数的类型都是context.Context，而名称都为parent。顾名思义，这个位置上的参数对应的都是它们将会产生的Context值的父值。

WithCancel函数用于产生一个可撤销的parent的子值。在coordinateWithContext函数中，我通过调用该函数，获得了一个衍生自上下文根节点的Context值，和一个用于触发撤销信号的函数。

而WithDeadline函数和WithTimeout函数则都可以被用来产生一个会定时撤销的parent的子值。至于WithValue函数，我们可以通过调用它，产生一个会携带额外数据的parent的子值。

正因为如此，在coordinateWithContext函数中，基于调用表达式cxt.Done()的接收操作，才能够起到感知撤销信号的作用。

### “可撤销的”在context 与 “撤销”一个Context值

Context类型这个接口中有两个方法与“撤销”息息相关。Done方法会返回一个元素类型为struct{}的接收通道。不过，这个接收通道的用途并不是传递元素值，而是让调用方去感知“撤销”当前Context值的那个信号。

除了让Context值的使用方感知到撤销信号，让它们得到“撤销”的具体原因，有时也是很有必要的。后者即是Context类型的Err方法的作用。该方法的结果是error类型的，并且其值只可能等于context.Canceled变量的值，或者context.DeadlineExceeded变量的值。

前者用于表示手动撤销，而后者则代表：由于我们给定的过期时间已到，而导致的撤销。

对于Context值来说，“撤销”这个词如果当名词讲，指的**其实就是被用来表达“撤销”状态的信号；如果当动词讲，指的就是对撤销信号的传达；而“可撤销的”指的则是具有传达这种撤销信号的能力。**

当我们通过调用context.WithCancel函数产生一个可撤销的Context值时，还会获得一个用于触发撤销信号的函数。通过调用这个函数，我们就可以触发针对这个Context值的撤销信号。一旦触发，撤销信号就会立即被传达给这个Context值，并由它的Done方法的结果值（一个接收通道）表达出来。

撤销函数只负责触发信号，而对应的可撤销的Context值也只负责传达信号，它们都不会去管后边具体的“撤销”操作。实际上，我们的代码可以在感知到撤销信号之后，进行任意的操作，Context值对此并没有任何的约束。

若再深究的话，这里的“撤销”最原始的含义其实就是，终止程序针对某种请求（比如 HTTP 请求）的响应，或者取消对某种指令（比如 SQL 指令）的处理。这也是 Go 语言团队在创建context代码包，和Context类型时的初衷。

我们可以去查看net包和database/sql包的 API 和源码，了解它们在这方面的典型应用。

### 撤销信号是在上下文树中的传播

context包中包含了四个用于繁衍Context值的函数。其中的WithCancel、WithDeadline和WithTimeout都是被用来基于给定的Context值产生可撤销的子值的。

context包的WithCancel函数在被调用后会产生两个结果值。**第一个结果值就是那个可撤销的Context值**，而**第二个结果值则是用于触发撤销信号的函数**。

在撤销函数被调用之后，对应的Context值会先关闭它内部的接收通道，也就是它的Done方法会返回的那个通道。

然后，**它会向它的所有子值（或者说子节点）传达撤销信号**。**这些子值会如法炮制，把撤销信号继续传播下去。最后，这个Context值会断开它与其父值之间的关联。**

![a801f8f2b5e89017ec2857bc1815fc9e](https://static001.geekbang.org/resource/image/a8/9e/a801f8f2b5e89017ec2857bc1815fc9e.png)



**我们通过调用context包的WithDeadline函数或者WithTimeout函数生成的Context值也是可撤销的**。**它们不但可以被手动撤销，还会依据在生成时被给定的过期时间，自动地进行定时撤销**。这里定时撤销的功能是借助它们内部的计时器来实现的。

**通过调用context.WithValue函数得到的Context值是不可撤销的。**撤销信号在被传播时，若遇到它们则会直接跨过，并试图将信号直接传给它们的子值。

### 通过Context值携带数据

WithValue函数在产生新的Context值（以下简称含数据的Context值）的时候需要三个参数，即：**父值、键和值。与“字典对于键的约束”类似，这里键的类型必须是可判等的。**

原因很简单，**当我们从中获取数据的时候，它需要根据给定的键来查找对应的值。不过，这种Context值并不是用字典来存储键和值的，后两者只是被简单地存储在前者的相应字段中而已。**

**Context类型的Value方法就是被用来获取数据的。**在我们**调用含数据的Context值的Value方法时，它会先判断给定的键，是否与当前值中存储的键相等，如果相等就把该值中存储的值直接返回，否则就到其父值中继续查找。**

如果其父值中仍然未存储相等的键，那么该方法就会沿着上下文根节点的方向一路查找下去。

除了含数据的Context值以外，其他几种Context值都是无法携带数据的。因此，Context值的Value方法在沿路查找的时候，会直接跨过那几种值。

**如果我们调用的Value方法的所属值本身就是不含数据的，那么实际调用的就将会是其父辈或祖辈的Value方法。**这是由于这几种Context值的实际类型，都属于结构体类型，并且**它们都是通过“将其父值嵌入到自身”，来表达父子关系的。**

**Context接口并没有提供改变数据的方法。因此，在通常情况下，我们只能通过在上下文树中添加含数据的Context值来存储新的数据，或者通过撤销此种值的父值丢弃掉相应的数据。**如果你存储在这里的数据可以从外部改变，那么必须自行保证安全。

