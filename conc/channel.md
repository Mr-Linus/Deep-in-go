## Channel

在声明并初始化一个通道的时候，需要用到 Go 语言的内建函数make。就像用make初始化切片那样，我们传给这个函数的第一个参数应该是代表了通道的具体类型的类型字面量。

声明一个通道类型变量，首先要确定该**通道类型的元素类型**，这决定了我们可以**通过这个通道传递什么类型的数据**。

比如，类型字面量chan int，其中的chan是表示通道类型的关键字，而int则说明了该通道类型的元素类型。又比如，chan string代表了一个元素类型为string的通道类型。

在初始化通道的时候，make函数除了必须接收这样的类型字面量作为参数，还可以接收一个int类型的参数,用于表示该通道的容量。指通道最多可以缓存多少个元素值。由此，虽然这个参数是int类型的，但是它是不能小于0的。

一个通道相当于一个先进先出（FIFO）的队列。也就是说，通道中的各个元素值都是严格地按照发送的顺序排列的，先被发送通道的元素值一定会先被接收。元素值的发送和接收都需要用到操作符<-。我们也可以叫它接送操作符。一个左尖括号紧接着一个减号形象地代表了元素值的传输方向。

```go
package main

import "fmt"

func main() {
  ch1 := make(chan int, 3)
  ch1 <- 2
  ch1 <- 1
  ch1 <- 3
  elem1 := <-ch1
  fmt.Printf("The first element received from channel ch1: %v\n",
    elem1)
}
```

这里声明并初始化了一个元素类型为int、容量为3的通道ch1，并用三条语句，向该通道先后发送了三个元素值2、1和3。

依次敲入通道变量的名称（比如ch1）、接送操作符<-以及想要发送的元素值（比如2），并且这三者之间最好用空格进行分割。

当我们需要从通道接收元素值的时候，同样要用`接送操作符<-`，只不过，这时需要把它写在变量名的左边，用于表达“要从该通道接收一个元素值”的语义。

比如：<-ch1，这也可以被叫做接收表达式。在一般情况下，接收表达式的结果将会是通道中的一个元素值。

如果我们需要把如此得来的元素值存起来，那么在接收表达式的左边就需要依次添加赋值符号（=或:=）和用于存值的变量的名字。因此，语句`elem1 := <-ch1`会将最先进入ch1的元素2接收来并存入变量elem1。

### 1. 通道的发送和接收的基本特性

- **对于同一个通道，发送操作之间**是**互斥的**，**接收操作之间**也是**互斥的**。
- **发送操作和接收操作中对元素值**的处理都是**不可分割**的。
- **发送操作在完全完成之前会被阻塞**，**接收操作也是如此**。

在同一时刻，Go 语言的运行时系统只会**执行对同一个通道的任意个发送操作中的某一个**。**直到这个元素值被完全复制进该通道之后，其他针对该通道的发送操作才可能被执行。**类似的，**在同一时刻，运行时系统也只会执行，对同一个通道的任意个接收操作中的某一个。**

**直到这个元素值完全被移出该通道之后，其他针对该通道的接收操作才可能被执行**。**即使这些操作是并发执行的也是如此。**

这里所谓的并发执行，可以这样认为，多个代码块分别在不同的 goroutine 之中，并有机会在同一个时间段内被执行。

对于通道中的同一个元素值来说，发送操作和接收操作之间也是互斥的。例如，虽然会出现，正在被复制进通道但还未复制完成的元素值，但是这时它绝不会被想接收它的一方看到和取走。

元素值从外界进入通道时会被复制。更具体地说，**进入通道的并不是在接收操作符右边的那个元素值，而是它的副本。**

**元素值从通道进入外界时会被移动。**这个移动操作实际上包含了两步，**第一步是生成正在通道中的这个元素值的副本，并准备给到接收方**，**第二步是删除在通道中的这个元素值。**

特性中，“不可分割”的意思是，它们**处理元素值时都是一气呵成的，绝不会被打断。**例如，发送操作要么还没复制元素值，要么已经复制完毕，绝不会出现只复制了一部分的情况。又例如，接收操作在准备好元素值的副本之后，一定会删除掉通道中的原值，绝不会出现通道中仍有残留的情况。

这既是为了保证通道中元素值的**完整性**，也是为了**保证通道操作的唯一性**。对于通道中的同一个元素值来说，它只可能是某一个发送操作放入的，同时也只可能被某一个接收操作取出。

一般情况下，发送操作包括了**“复制元素值”和“放置副本到通道内部”**这两个步骤。

这两个步骤完全完成之前，发起这个发送操作的那句代码会一直阻塞在那里。也就是说，在它之后的代码不会有执行的机会，直到这句代码的阻塞解除。在通道完成发送操作之后，运行时系统会通知这句代码所在的 goroutine，以使它去争取继续运行代码的机会。

接收操作通常包含了**“复制通道内的元素值”、“放置副本到接收方”、“删掉原值”**三个步骤。

在所有这些步骤完全完成之前，发起该操作的代码也会一直阻塞，直到该代码所在的 goroutine 收到了运行时系统的通知并重新获得运行机会为止。如此阻塞代码其实就是为了实现操作的互斥和元素值的完整。

### 2. 信道被长时间阻塞

- 针对缓冲通道

针对缓冲通道的情况。如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道中有元素值被接收走。

这时，通道会优先通知最早因此而等待的、那个发送操作所在的 goroutine，后者会再次执行发送操作。

由于发送操作在这种情况下被阻塞后，它们所在的 goroutine 会顺序地进入通道内部的发送等待队列，所以通知的顺序总是公平的。

相对的，如果通道已空，那么对它的所有接收操作都会被阻塞，直到通道中有新的元素值出现。这时，通道会通知最早等待的那个接收操作所在的 goroutine，并使它再次执行接收操作。因此而等待的、所有接收操作所在的 goroutine，都会按照先后顺序被放入通道内部的接收等待队列。

- 针对非缓冲通道

对于非缓冲通道，情况要简单一些。无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。由此可见，非缓冲通道是在用同步的方式传递数据。也就是说，只有收发双方对接上了，数据才会被传递。

并且，数据是直接从发送方复制到接收方的，中间并不会用非缓冲通道做中转。相比之下，**缓冲通道则在用异步的方式传递数据。**

在大多数情况下，**缓冲通道会作为收发双方的中间件**。但是，当发送操作在执行的时候发现空的通道中，正好有等待的接收操作，那么它会直接把元素值复制给接收方。

#### 错误使用通道而造成的阻塞

对于值为nil的通道，不论它的具体类型是什么，对它的发送操作和接收操作都会永久地处于阻塞状态。它们所属的 goroutine 中的任何代码，都不再会被执行。

> 注意，由于通道类型是引用类型，所以它的零值就是nil。换句话说，当我们只声明该类型的变量但没有用make函数对它进行初始化时，该变量的值就会是nil。我们一定不要忘记初始化通道！

#### 错误使用通道而造成的Panic

对于一个已初始化，但并未关闭的通道来说，收发操作一定不会引发 panic。但是通道一旦关闭，再对它进行发送操作，就会引发 panic。

如果我们试图关闭一个已经关闭了的通道，也会引发 panic。注意，接收操作是可以感知到通道的关闭的，并能够安全退出。

更具体地说，**当我们把接收表达式的结果同时赋给两个变量时**，**第二个变量的类型就是一定bool类型**。**它的值如果为false就说明通道已经关闭，并且再没有元素值可取了。**

> 注意，**如果通道关闭时，里面还有元素值未被取出，那么接收表达式的第一个结果，仍会是通道中的某一个元素值，而第二个结果值一定会是true。**

通过接收表达式的第二个结果值，来判断通道是否关闭是可能有延时的。

除非有特殊的保障措施，我们**千万不要让接收方关闭通道，而应当让发送方做这件事。**

### 3. 单向通道

单向通道就是，只能发不能收，或者只能收不能发的通道。一个通道是双向的，还是单向的是由它的类型字面量体现的。

```go
var uselessChan = make(chan<- int, 1)
```

声明并初始化了一个名叫uselessChan的变量。这个变量的类型是chan<- int，容量是1。

紧挨在关键字chan右边的那个<-，这表示了这个通道是单向的，并且只能发而不能收。

如果这个操作符紧挨在chan的左边，那么就说明该通道只能收不能发。所以，前者可以被简称为发送通道，后者可以被简称为接收通道。

> 注意，与发送操作和接收操作对应，这里的“发”和“收”都是站在操作通道的代码的角度上说的。

从上述变量的名字上你也能猜到，这样的通道是没用的。通道就是为了传递数据而存在的，声明一个只有一端（发送端或者接收端）能用的通道没有任何意义。

##### 单向通道的用途

```go
func SendInt(ch chan<- int) {
  ch <- rand.Intn(1000)
}
```

func关键字声明了一个叫做SendInt的函数。这个函数只接受一个chan<- int类型的参数。在这个函数中的代码只能向**参数ch发送元素值，而不能从它那里接收元素值**。这就起到了约束函数行为的作用。

在实际场景中，这种约束一般会出现在接口类型声明中的某个方法定义上。请看这个叫Notifier的接口类型声明：

```go
type Notifier interface {
  SendInt(ch chan<- int)
}
```

在接口类型声明的花括号中，每一行都代表着一个方法的定义。接口中的方法定义与函数声明很类似，但是只包含了方法名称、参数列表和结果列表。

一个类型如果想成为一个接口类型的实现类型，那么就必须实现这个接口中定义的所有方法。因此，如果我们在**某个方法的定义中使用了单向通道类型，那么就相当于在对它的所有实现做出约束。**

Notifier接口中的SendInt方法只会接受一个发送通道作为参数，所以，**在该接口的所有实现类型中的SendInt方法都会受到限制**。这种约束方式还是很有用的，尤其是在我们编写模板代码或者可扩展的程序库的时候。

在调用SendInt函数的时候，只需要**把一个元素类型匹配的双向通道传给它就行了**，没必要用发送通道，**因为 Go 语言在这种情况下会自动地把双向通道转换为函数所需的单向通道。**

```go
intChan1 := make(chan int, 3)
SendInt(intChan1)
```

还可以在函数声明的结果列表中使用单向通道。

```go
func getIntChan() <-chan int {
  num := 5
  ch := make(chan int, num)
  for i := 0; i < num; i++ {
    ch <- i
  }
  close(ch)
  return ch
}
```

函数getIntChan会返回一个<-chan int类型的通道，这就意味着得到该通道的程序，只能从通道中接收元素值。**这实际上就是对函数调用方的一种约束了。**

在 Go 语言中还可以声明函数类型，**如果我们在函数类型中使用了单向通道，那么就相等于在约束所有实现了这个函数类型的函数。**

```go
intChan2 := getIntChan()
for elem := range intChan2 {
  fmt.Printf("The element in intChan2: %v\n", elem)
}
```

调用getIntChan得到的结果值赋给了变量intChan2，然后用for语句循环地取出了该通道中的所有元素值，并打印出来。

上述 for 语句会不断地尝试从通道 intChan2 中取出元素值。**即使 intChan2 已经被关闭了，它也会在取出所有剩余的元素值之后再结束执行。**

**倘若通道intChan2的值为nil，那么这条for语句就会被永远地阻塞在有for关键字的那一行。**

### 4. select语句与通道的联用

select语句只能与通道联用，它一般由若干个分支组成。每次执行这种语句的时候，一般**只有一个分支中的代码会被运行。**

select语句的分支分为两种，一种叫做**候选分支**，另一种叫做**默认分支**。候选分支总是以关键字case开头，后跟一个case表达式和一个冒号，然后我们可以从下一行开始写入当分支被选中时需要执行的语句。

默认分支其实就是 default case，因为，**当且仅当没有候选分支被选中时它才会被执行**，所以它以关键字default开头并直接后跟一个冒号。同样的，我们可以在default:的下一行写入要执行的语句。

select语句是专为通道而设计的，所以每个case表达式中都只能包含操作通道的表达式，比如接收表达式。

如果我们需要把接收表达式的结果赋给变量的话，还可以把这里写成赋值语句或者短变量声明。下面展示一个简单的例子。

```go
// 准备好几个通道。
intChannels := [3]chan int{
  make(chan int, 1),
  make(chan int, 1),
  make(chan int, 1),
}
// 随机选择一个通道，并向它发送元素值。
index := rand.Intn(3)
fmt.Printf("The index: %d\n", index)
intChannels[index] <- index
// 哪一个通道中有可取的元素值，哪个对应的分支就会被执行。
select {
case <-intChannels[0]:
  fmt.Println("The first candidate case is selected.")
case <-intChannels[1]:
  fmt.Println("The second candidate case is selected.")
case elem := <-intChannels[2]:
  fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
default:
  fmt.Println("No candidate case is selected!")
}
```

这里先准备好了三个类型为chan int、容量为1的通道，并把它们存入了一个叫做intChannels的数组。

随机选择一个范围在 [0, 2] 的整数，把它作为索引在上述数组中选择一个通道，并向其中发送一个元素值。

我用一个包含了三个候选分支的 select 语句，分别尝试从上述三个通道中接收元素值，哪一个通道中有值，哪一个对应的候选分支就会被执行。后面还有一个默认分支，不过在这里它是不可能被选中的。

#### Select 语句的注意事项

- 如果加入了默认分支，那么无论涉及通道操作的表达式是否有阻塞，select语句都不会被阻塞。**如果那几个表达式都阻塞了，或者说都没有满足求值的条件，那么默认分支就会被选中并执行。**
- 如果没有加入默认分支，那么一旦所有的case表达式都没有满足求值条件，那么select语句就会被阻塞。直到至少有一个case表达式满足条件为止。
- 可以通过接收表达式的第二个结果值来判断通道是否已经关闭。一旦发现某个通道关闭了，我们就应该及时地屏蔽掉对应的分支或者采取其他措施。这对于程序逻辑和程序性能都是有好处的。
- select 语句只能对其中的每一个 case 表达式各求值一次。所以，如果我们想连续或定时地操作其中的通道的话，就往往需要通过在for语句中嵌入select语句的方式实现。但是注意，**简单地在select语句的分支中使用break语句，只能结束当前的select语句的执行，而并不会对外层的for语句产生作用。这种错误的用法可能会让这个for语句无休止地运行下去。**

```go
intChan := make(chan int, 1)
	// 一秒后关闭通道。
time.AfterFunc(time.Second, func() {
  close(intChan)
})
select {
case _, ok := <-intChan:
  if !ok {
    fmt.Println("The candidate case is closed.")
    break
  }
  fmt.Println("The candidate case is selected.")
}
```

这里，先声明并初始化了一个叫做intChan的通道，然后通过 time 包中的 AfterFunc 函数约定在一秒钟之后关闭该通道。

后面的 select 语句只有一个候选分支，我在其中利用接收表达式的第二个结果值对 intChan 通道是否已关闭做了判断，并在得到肯定结果后，通过 break 语句立即结束当前 select 语句的执行。

#### Select 分支选择规则

- **比如，如果 case 表达式是包含了接收表达式的短变量声明时，那么在赋值符号左边的就可以是一个或两个表达式，不过此处的表达式的结果必须是可以被赋值的。**当**这样的case表达式被求值时，它包含的多个表达式总会以从左到右的顺序被求值。**

- **select 语句包含的候选分支中的 case 表达式都会在该语句执行开始时先被求值，并且求值的顺序是依从代码编写的顺序从上到下的。**结合上一条规则，在select语句开始执行时，排在最上边的候选分支中最左边的表达式会最先被求值，然后是它右边的表达式。仅当最上边的候选分支中的所有表达式都被求值完毕后，从上边数第二个候选分支中的表达式才会被求值，顺序同样是从左到右，然后是第三个候选分支、第四个候选分支，以此类推。

- 对于每一个case表达式，**如果其中的发送表达式或者接收表达式在被求值时，相应的操作正处于阻塞状态，那么对该 case 表达式的求值就是不成功的**。在这种情况下，我们可以说，这个case表达式所在的候选分支是不满足选择条件的。

- **仅当select语句中的所有case表达式都被求值完毕后**，**它才会开始选择候选分支。这时候，它只会挑选满足选择条件的候选分支执行。**如果所有的候选分支都不满足选择条件，那么默认分支就会被执行。如果这时没有默认分支，那么select语句就会立即进入阻塞状态，直到至少有一个候选分支满足选择条件为止。一旦有一个候选分支满足选择条件，select语句（或者说它所在的 goroutine）就会被唤醒，这个候选分支就会被执行。

- **如果 select 语句发现同时有多个候选分支满足选择条件，那么它就会用一种伪随机的算法在这些分支中选择一个并执行。注意，即使select语句是在被唤醒时发现的这种情况，也会这样做。**

- **一条 select 语句中只能够有一个默认分支**。并且，默认分支只在无候选分支可选时才会被执行，这与它的编写位置无关。

- **select 语句的每次执行，包括 case 表达式求值和分支选择，都是独立的**。不过，至于它的执行是否是并发安全的，就要看其中的case表达式以及分支中，是否包含并发不安全的代码了。