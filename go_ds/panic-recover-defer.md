## panic recover defer

- ### panic 

一种在我们意料之外的程序异常，它处理的不是错误，这种程序异常被叫做 panic，我把它翻译为运行时恐慌。其中的“恐慌”二字是由 panic 直译过来的，而之所以前面又加上了“运行时”三个字，是因为这种异常只会在程序运行的时候被抛出来。

比如说，一个 Go 程序里有一个切片，它的长度是 5，也就是说该切片中的元素值的索引分别为0、1、2、3、4，但是，我在程序里却想通过索引5访问其中的元素值，显而易见，这样的访问是不正确的。

Go 程序，确切地说是程序内嵌的 Go 语言运行时系统，会在执行到这行代码的时候抛出一个“index out of range”的 panic，用以提示你索引越界了。

这不仅仅是个提示。当 panic 被抛出之后，如果我们没有在程序里添加任何保护措施的话，程序（或者说代表它的那个进程）就会在打印出 panic 的详细情况（以下简称 panic 详情）之后，终止运行。

```go
panic: runtime error: index out of range

goroutine 1 [running]:
main.main()
 /root/main.go:5 +0x3d
exit status 2
```



这份详情的第一行是“panic: runtime error: index out of range”。其中的“runtime error”的含义是，这是一个runtime代码包中抛出的 panic。在这个 panic 中，包含了一个runtime.Error接口类型的值。runtime.Error接口内嵌了error接口，并做了一点点扩展，runtime包中有不少它的实现类型。



实际上，此详情中的“panic：”右边的内容，**正是这个 panic 包含的runtime.Error类型值的字符串表示形式**。

panic 详情中，一般还会包含与它的引发原因有关的 goroutine 的代码执行信息。正如前述详情中的“goroutine 1 [running]”，它表示有一个 ID 为1的 goroutine 在此 panic 被引发的时候正在运行。这里的 ID 其实并不重要，因为它只是 Go 语言运行时系统内部给予的一个 goroutine 编号，我们在程序中是无法获取和更改的。

“main.main()”表明了这个 goroutine 包装的go函数就是命令源码文件中的那个main函数，也就是说这里的 goroutine 正是主 goroutine。再下面的一行，指出的就是这个 goroutine 中的哪一行代码在此 panic 被引发时正在执行。

**这包含了此行代码在其所属的源码文件中的行数，以及这个源码文件的绝对路径**。这一行最后的**+0x3d代表的是：此行代码相对于其所属函数的入口程序计数偏移量。不过，一般情况下它的用处并不大**。

“exit status 2”表明我的这个程序是以退出状态码2结束运行的。在大多数操作系统中，只要退出状态码不是0，都意味着程序运行的非正常结束。**在 Go 语言中，因 panic 导致程序结束运行的退出状态码一般都会是2。**

从上边的这个 panic 详情可以看出，作为此 panic 的引发根源的代码处于文件中的第 5 行，同时被包含在main包（也就是命令源码文件所在的代码包）的main函数中。

### 引发 panic 的过程

大致的过程：某个函数中的某行代码有意或无意地引发了一个 panic。这时，初始的 panic 详情会被建立起来，并且**该程序的控制权会立即从此行代码转移至调用其所属函数的那行代码上，也就是调用栈中的上一级。**

这也意味着，**此行代码所属函数的执行随即终止。紧接着，控制权并不会在此有片刻的停留，它又会立即转移至再上一级的调用代码处**。**控制权如此一级一级地沿着调用栈的反方向传播至顶端，也就是我们编写的最外层函数那里。**

**最外层函数指的是go函数，对于主 goroutine 来说就是main函数**。但是**控制权也不会停留在那里，而是被 Go 语言运行时系统收回。**

随后，程序崩溃并终止运行，承载程序这次运行的进程也会随之死亡并消失。与此同时，在这个控制权传播的过程中，panic 详情会被逐渐地积累和完善，并会在程序终止之前被打印出来。

panic 可能是我们在无意间（或者说一不小心）引发的，如前文所述的索引越界。这类 panic 是真正的、在我们意料之外的程序异常。不过，除此之外，我们还是可以有意地引发 panic。Go 语言的内建函数panic是专门用于引发 panic 的。panic函数使程序开发者可以在程序运行期间报告异常。

注意，这与从函数返回错误值的意义是完全不同的。当我们的函数返回一个非nil的错误值时，函数的调用方有权选择不处理，并且不处理的后果往往是不致命的。意思是，不至于使程序无法提供任何功能（也可以说僵死）或者直接崩溃并终止运行（也就是真死）。

panic 详情会在控制权传播的过程中，被逐渐地积累和完善，并且，控制权会一级一级地沿着调用栈的反方向传播至顶端。

在针对某个 goroutine 的代码执行信息中，调用栈底端的信息会先出现，然后是上一级调用的信息，以此类推，最后才是此调用栈顶端的信息。

#### Usecase

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter function main.")
	caller1()
	fmt.Println("Exit function main.")
}

func caller1() {
	fmt.Println("Enter function caller1.")
	caller2()
	fmt.Println("Exit function caller1.")
}

func caller2() {
	fmt.Println("Enter function caller2.")
	s1 := []int{0, 1, 2, 3, 4}
	e5 := s1[5]
	_ = e5
	fmt.Println("Exit function caller2.")
}
```

```shell
Enter function main.
Enter function caller1.
Enter function caller2.
panic: runtime error: index out of range [5] with length 5

goroutine 1 [running]:
main.caller2()
        /Users/funky/Projects/Test/main.go:22 +0x82
main.caller1()
        /Users/funky/Projects/Test/main.go:15 +0x7e
main.main()
        /Users/funky/Projects/Test/main.go:9 +0x7e
exit status 2

```

main函数调用了caller1函数，而caller1函数又调用了caller2函数，那么caller2函数中代码的执行信息会先出现，然后是caller1函数中代码的执行信息，最后才是main函数的信息。

![606ff433a6b58510f215e57792822bd7](https://static001.geekbang.org/resource/image/60/d7/606ff433a6b58510f215e57792822bd7.png)



如果一个 panic 是我们在无意间引发的，那么其中的值只能由 Go 语言运行时系统给定。但是，当我们使用panic函数有意地引发一个 panic 的时候，却可以自行指定其包含的值。



在调用panic函数时，把某个值作为参数传给该函数就可以了。由于panic函数的唯一一个参数是空接口（也就是interface{}）类型的，所以从语法上讲，它可以接受任何类型的值。但是，我们最好传入error类型的错误值，或者其他的可以被有效序列化的值。这里的“有效序列化”指的是，可以更易读地去表示形式转换。

对于fmt包下的各种打印函数来说，error类型值的Error方法与其他类型值的String方法是等价的，它们的唯一结果都是string类型的。

一旦程序异常了，我们就一定要把异常的相关信息记录下来，这通常都是记到程序日志里。

在为程序排查错误的时候，首先要做的就是查看和解读程序日志；而最常用也是最方便的日志记录方式，就是记下相关值的字符串表示形式。

如果你觉得某个值有可能会被记到日志里，那么就应该为它关联String方法。如果这个值是error类型的，那么让它的Error方法返回你为它定制的字符串表示形式就可以了。

你可能会想到 fmt.Sprintf，以及 fmt.Fprintf 这类可以格式化并输出参数的函数。

它们本身就可以被用来输出值的某种表示形式。不过，它们在功能上，肯定远不如我们自己定义的 Error 方法或者String 方法。因此，为不同的数据类型分别编写这两种方法总是首选。

传给panic函数的参数值，至少在程序崩溃的时候，panic 包含的那个值字符串表示形式会被打印出来。我们还可以施加某种保护措施，避免程序的崩溃。这个时候，panic 包含的值会被取出，而在取出之后，它一般都会被打印出来或者记录到日志里。

- ### recover & defer

Go 语言的内建函数recover专用于恢复 panic，或者说平息运行时恐慌。recover函数无需任何参数，并且会返回一个空接口类型的值。

如果用法正确，这个值实际上就是即将恢复的 panic 包含的值。并且，如果这个 panic 是因我们调用panic函数而引发的，那么该值同时也会是我们此次调用 panic 函数时，传入的参数值副本。请注意，这里强调用法的正确。

#### 错误示范

```go
package main

import (
 "fmt"
 "errors"
)

func main() {
 fmt.Println("Enter function main.")
 // 引发panic。
 panic(errors.New("something wrong"))
 p := recover()
 fmt.Printf("panic: %s\n", p)
 fmt.Println("Exit function main.")
}
```

程序依然会崩溃，这个recover函数调用并不会起到任何作用，甚至都没有机会执行。

顾名思义，defer语句就是被用来延迟执行代码的。延迟到该语句所在的函数即将执行结束的那一刻，无论结束执行的原因是什么。

与go语句有些类似，一个defer语句总是由一个defer关键字和一个调用表达式组成。

这里被调用的函数可以是有名称的，也可以是匿名的。我们可以把这里的函数叫做defer函数或者延迟函数。注意，被延迟执行的是defer函数，而不是defer语句。

无论函数结束执行的原因是什么，其中的defer函数调用都会在它即将结束执行的那一刻执行。即使导致它执行结束的原因是一个 panic 也会是这样。正因为如此，我们需要联用defer语句和recover函数调用，才能够恢复一个已经发生的 panic。

#### usecase

```go
package main

import (
 "fmt"
 "errors"
)

func main() {
 fmt.Println("Enter function main.")
 defer func(){
  fmt.Println("Enter defer function.")
  if p := recover(); p != nil {
   fmt.Printf("panic: %s\n", p)
  }
  fmt.Println("Exit defer function.")
 }()
 // 引发panic。
 panic(errors.New("something wrong"))
 fmt.Println("Exit function main.")
}
```

这个main函数中，我先编写了一条defer语句，并在defer函数中调用了recover函数。仅当调用的结果值不为nil时，也就是说只有 panic 确实已发生时，我才会打印一行以“panic:”为前缀的内容。

紧接着，我调用了 panic 函数，并传入了一个error类型值。这里一定要注意，我们要尽量把 defer 语句写在函数体的开始处，**因为在引发 panic 的语句之后的所有语句，都不会有任何执行机会**。因此上述**defer 函数中的 recover 函数调用才会拦截，并恢复 defer 语句所属的函数，及其调用的代码中发生的所有 panic。**

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Enter function main.")
	caller()
	fmt.Println("Exit function main.")
}

func caller() {
	fmt.Println("Enter function caller.")
	panic(errors.New("something wrong")) // 正例。
	panic(fmt.Println)                   // 反例。
	fmt.Println("Exit function caller.")
}
```



### 多条defer语句调用顺序

在同一个函数中，defer 函数调用的执行顺序与它们分别所属的 defer 语句的出现顺序（更严谨地说，是执行顺序）完全相反。

当一个函数即将结束执行时，**其中的写在最下边的defer函数调用会最先执行**，其次是写在它上边、与它的距离最近的那个defer函数调用，以此类推，**最上边的defer函数调用会最后一个执行。**

如果函数中有一条for语句，并且这条for语句中包含了一条defer语句，那么，**显然这条defer语句的执行次数，就取决于for语句的迭代次数。**

同一条defer语句每被执行一次，**其中的defer函数调用就会产生一次，而且，这些函数调用同样不会被立即执行。**

在defer语句每次执行的时候，**Go 语言会把它携带的defer函数及其参数值另行存储到一个队列中。**

这个队列与该defer语句所属的函数是对应的，并且，它是先进后出（FILO）的，相当于一个栈。

在需要执行某个函数中的defer函数调用的时候，Go 语言会先拿到对应的队列，然后从该队列中一个一个地取出defer函数及其参数值，并逐个执行调用。

```go
package main

import "fmt"

func main() {
	defer fmt.Println("first defer")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer in for [%d]\n", i)
	}
	defer fmt.Println("last defer")
}
```

输出：

```go
last defer
defer in for [2]
defer in for [1]
defer in for [0]
first defer
```

