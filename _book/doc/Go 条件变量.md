## Go 条件变量

条件变量是基于互斥锁的，它必须有互斥锁的支撑才能发挥作用。

**条件变量并不是被用来保护临界区和共享资源的，它是用于协调想要访问共享资源的那些线程的**。**当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程。**

比如说，我们两个人在共同执行一项秘密任务，这需要在不直接联系和见面的前提下进行。我需要向一个信箱里放置情报，你需要从这个信箱中获取情报。这个信箱就相当于一个共享资源，而我们就分别是进行写操作的线程和进行读操作的线程。

如果我在放置的时候发现信箱里还有未被取走的情报，那就不再放置，而先返回。另一方面，如果你在获取的时候发现信箱里没有情报，那也只能先回去了。这就相当于写的线程或读的线程阻塞的情况。

虽然我们俩都有信箱的钥匙，但是同一时刻只能有一个人插入钥匙并打开信箱，这就是锁的作用了。更何况咱们俩是不能直接见面的，所以这个信箱本身就可以被视为一个临界区。

尽管没有协调好，咱们俩仍然要想方设法的完成任务。所以，如果信箱里有情报，而你却迟迟未取走，那我就需要每过一段时间带着新情报去检查一次，若发现信箱空了，我就需要及时地把新情报放到里面。

如果信箱里一直没有情报，那你也要每过一段时间去打开看看，一旦有了情报就及时地取走。这么做是可以的，但就是太危险了，很容易被敌人发现。

又想了一个计策，各自雇佣了一个不起眼的小孩儿。如果早上七点有一个戴红色帽子的小孩儿从你家楼下路过，那么就意味着信箱里有了新情报。另一边，如果上午九点有一个戴蓝色帽子的小孩儿从我家楼下路过，那就说明你已经从信箱中取走了情报。

这样一来，咱们执行任务的隐蔽性高多了，并且效率的提升非常显著。这两个戴不同颜色帽子的小孩儿就相当于条件变量，在共享资源的状态产生变化的时候，起到了通知的作用。

条件变量在这里的最大优势就是在效率方面的提升。当共享资源的状态不满足条件的时候，想操作它的线程再也不用循环往复地做检查了，只要等待通知就好了。

### 条件变量与互斥锁配合使用

条件变量提供的方法有三个：等待通知（wait）、单发通知（signal）和广播通知（broadcast）。

**在利用条件变量等待通知的时候，需要在它基于的那个互斥锁保护下进行**。而在进行单发通知或广播通知的时候，却是恰恰相反的，也就是说，需要在对应的互斥锁解锁之后再做这两种操作。

```go
var mailbox uint8
var lock sync.RWMutex
sendCond := sync.NewCond(&lock)
recvCond := sync.NewCond(lock.RLocker())
```

变量mailbox代表信箱，是uint8类型的。 若它的值为0则表示信箱中没有情报，而当它的值为1时则说明信箱中有情报。lock是一个类型为sync.RWMutex的变量，是一个读写锁，也可以被视为信箱上的那把锁。

基于这把锁，创建两个代表条件变量的变量，名字分别叫sendCond和recvCond。 它们都是*sync.Cond类型的，同时也都是由sync.NewCond函数来初始化的。

与sync.Mutex类型和sync.RWMutex类型不同，sync.Cond类型并不是开箱即用的。我们只能利用sync.NewCond函数创建它的指针值。这个函数需要一个sync.Locker类型的参数值。

条件变量是基于互斥锁的，它必须有互斥锁的支撑才能够起作用。因此，这里的参数值是不可或缺的，它会参与到条件变量的方法实现当中。

sync.Locker其实是一个接口，在它的声明中只包含了两个方法定义，即：Lock()和Unlock()。sync.Mutex类型和sync.RWMutex类型都拥有Lock方法和Unlock方法，只不过它们都是指针方法。因此，**这两个类型的指针类型才是sync.Locker接口的实现类型。**在为sendCond变量做初始化的时候，把基于lock变量的指针值传给了sync.NewCond函数。

lock变量的Lock方法和Unlock方法分别用于对其中写锁的锁定和解锁，它们与sendCond变量的含义是对应的。sendCond是专门为放置情报而准备的条件变量，向信箱里放置情报，可以被视为对共享资源的写操作。

**recvCond 变量代表的是专门为获取情报而准备的条件变量**。 虽然获取情报也会涉及对信箱状态的改变，但是好在做这件事的人只会有你一个，而且我们也需要借此了解一下，条件变量与读写锁中的读锁的联用方式。所以，在这里，我们暂且把获取情报看做是对共享资源的读操作。

为了初始化recvCond这个条件变量，我们需要的是lock变量中的读锁，并且还需要是sync.Locker类型的。可是，lock变量中用于对读锁进行锁定和解锁的方法却是RLock和RUnlock，它们与sync.Locker接口中定义的方法并不匹配。

好在sync.RWMutex类型的RLocker方法可以实现这一需求。我们只要在调用sync.NewCond函数时，传入调用表达式lock.RLocker()的结果值，就可以使该函数返回符合要求的条件变量了。

现在有四个变量。一个是代表信箱的mailbox，一个是代表信箱上的锁的lock。还有两个是，代表了蓝帽子小孩儿的sendCond，以及代表了红帽子小孩儿的recvCond。

![3619456ade9d45a4d9c0fbd22bb6fd5d](https://static001.geekbang.org/resource/image/36/5d/3619456ade9d45a4d9c0fbd22bb6fd5d.png)

我，现在是一个 goroutine（携带的go函数），想要适时地向信箱里放置情报并通知你：

```go
lock.Lock()
for mailbox == 1 {
 sendCond.Wait()
}
mailbox = 1
lock.Unlock()
recvCond.Signal()
```

先调用lock变量的Lock方法。检查mailbox变量的值是否等于1，也就是说，要看看信箱里是不是还存有情报。如果还有情报，那么我就回家去等蓝帽子小孩儿了。

这就是那条for语句以及其中的调用表达式sendCond.Wait()所表示的含义了。

如果信箱里没有情报，那么我就把新情报放进去，关上信箱、锁上锁，然后离开。用代码表达出来就是mailbox = 1和lock.Unlock()。那就是让红帽子小孩儿准时去你家楼下路过。也就是说，我会及时地通知你“信箱里已经有新情报了”，我们调用recvCond的Signal方法就可以实现这一步骤。

另一方面，你现在是另一个 goroutine，想要适时地从信箱中获取情报，然后通知我。

```go
lock.RLock()
for mailbox == 0 {
  recvCond.Wait() // Wait执行会先进入通知队列，然后RUnLock，然后等待通知，被唤醒后，再RLock()，往下跑
}
mailbox = 0
lock.RUnlock()
sendCond.Signal()
```

你跟我做的事情在流程上其实基本一致，只不过每一步操作的对象是不同的。你需要调用的是lock变量的RLock方法。因为你要进行的是读操作，并且会使用recvCond变量作为辅助。recvCond与lock变量的读锁是对应的。

在打开信箱后，你要关注的是信箱里是不是没有情报，也就是检查 mailbox 变量的值是否等于0。如果它确实等于0，那么你就需要回家去等红帽子小孩儿，也就是调用recvCond的Wait方法。这里使用的依然是for语句。

如果信箱里有情报，那么你就应该取走情报，关上信箱、锁上锁，然后离开。对应的代码是mailbox = 0和lock.RUnlock()。之后，你还需要让蓝帽子小孩儿准时去我家楼下路过。这样我就知道信箱中的情报已经被你获取了。

只要条件不满足，我就会通过调用sendCond变量的Wait方法，去等待你的通知，只有在收到通知之后我才会再次检查信箱。

另外，当我需要通知你的时候，我会调用recvCond变量的Signal方法。你使用这两个条件变量的方式正好与我相反。你可能也看出来了，利用条件变量可以实现单向的通知，而双向的通知则需要两个条件变量。这也是条件变量的基本使用规则。

条件变量是基于互斥锁的一种同步工具，它必须有互斥锁的支撑才能发挥作用。  条件变量可以协调那些想要访问共享资源的线程。当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程。

### 条件变量的Wait方法

条件变量的Wait方法主要做了四件事:

- 把调用它的 goroutine（也就是当前的 goroutine）加入到当前条件变量的通知队列中。
- 解锁当前的条件变量基于的那个互斥锁。
- 让当前的 goroutine 处于等待状态，等到通知到来时再决定是否唤醒它。此时，这个 goroutine 就会阻塞在调用这个Wait方法的那行代码上。
- 如果通知到来并且决定唤醒这个 goroutine，那么就在唤醒它之后重新锁定当前条件变量基于的互斥锁。自此之后，当前的 goroutine 就会继续执行后面的代码了。

所以为什么先要锁定条件变量基于的互斥锁，才能调用它的Wait方法，因为条件变量的Wait方法在阻塞当前的 goroutine 之前，**会解锁它基于的互斥锁**，所以在调用该Wait方法之前，我们必须先锁定那个互斥锁，否则在调用这个Wait方法时，就会引发一个不可恢复的 panic。如果 Wait 方法在互斥锁已经锁定的情况下，阻塞了当前的 goroutine，那么该 goroutine一直处于阻塞状态，无法恢复。成对的锁定和解锁，就算别的 goroutine 可以来解锁，那万一解锁重复了怎么办？由此引发的 panic 可是无法恢复的。

如果当前的 goroutine 无法解锁，别的 goroutine 也都不来解锁，那么又由谁来进入临界区，并改变共享资源的状态呢？只要共享资源的状态不变，即使当前的 goroutine 因收到通知而被唤醒，也依然会再次执行这个Wait方法，并再次被阻塞。

所以说，如果条件变量的Wait方法不先解锁互斥锁的话，那么就只会造成两种后果：**不是当前的程序因 panic 而崩溃，就是相关的 goroutine 全面阻塞。**

为了保险起见，一般把wait函数包裹在for语句里。如果一个 goroutine 因收到通知而被唤醒，但却发现共享资源的状态，依然不符合它的要求，那么就应该再次调用条件变量的Wait方法，并继续等待下次通知的到来。

这种情况是很有可能发生的，具体如下面所示。

**有多个 goroutine 在等待共享资源的同一种状态**。比如，**它们都在等mailbox变量的值不为0的时候再把它的值变为0，这就相当于有多个人在等着我向信箱里放置情报。**虽然等待的 goroutine 有多个，**但每次成功的 goroutine 却只可能有一个。**别忘了，**条件变量的Wait方法会在当前的 goroutine 醒来后先重新锁定那个互斥锁。在成功的 goroutine 最终解锁互斥锁之后**，其他的 goroutine 会先后进入临界区，但它们会发现共享资源的状态依然不是它们想要的。这个时候，for循环就很有必要了。

共享资源可能有的状态不是两个，而是更多。比如，mailbox变量的可能值不只有0和1，还有2、3、4。这种情况下，由于状态在每次改变后的结果只可能有一个，所以，在设计合理的前提下，单一的结果一定不可能满足所有 goroutine 的条件。那些未被满足的 goroutine 显然还需要继续等待和检查。

有一种可能，共享资源的状态只有两个，并且每种状态都只有一个 goroutine 在关注，就像我们在主问题当中实现的那个例子那样。不过，即使是这样，使用for语句仍然是有必要的。原因是，在一些多 CPU 核心的计算机系统中，**即使没有收到条件变量的通知，调用其Wait方法的 goroutine 也是有可能被唤醒的**。这是由计算机硬件层面决定的，即使是操作系统（比如 Linux）本身提供的条件变量也会如此。

### 条件变量的Signal方法和Broadcast方法

条件变量的Signal方法和Broadcast方法都是被用来发送通知的，不同的是，前者的通知只会唤醒一个因此而等待的 goroutine，而后者的通知却会唤醒所有为此等待的 goroutine。

条件变量的Wait方法总会把当前的 goroutine 添加到通知队列的队尾，而它的Signal方法总会从通知队列的队首开始，查找可被唤醒的 goroutine。所以，因Signal方法的通知，**而被唤醒的 goroutine 一般都是最早等待的那一个**。

这两个方法的行为决定了它们的适用场景。如果你确定只有一个 goroutine 在等待通知，或者只需唤醒任意一个 goroutine 就可以满足要求，那么使用条件变量的Signal方法就好了。

否则，使用Broadcast方法总没错，只要你设置好各个 goroutine 所期望的共享资源状态就可以了。

强调一下，与Wait方法不同，**条件变量的Signal方法和Broadcast方法并不需要在互斥锁的保护下执行**。恰恰相反，我们**最好在解锁条件变量基于的那个互斥锁之后，再去调用它的这两个方法。这更有利于程序的运行效率。**

请注意，**条件变量的通知具有即时性**。也就是说，**如果发送通知的时候没有 goroutine 为此等待，那么该通知就会被直接丢弃。在这之后才开始等待的 goroutine 只可能被后面的通知唤醒。**

