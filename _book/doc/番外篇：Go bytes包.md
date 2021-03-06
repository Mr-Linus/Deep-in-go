## 番外篇：Go bytes包

strings包和bytes包可以说是一对孪生兄弟，它们在 API 方面非常的相似。单从它们提供的函数的数量和功能上讲，差别可以说是微乎其微。

**只不过，strings包主要面向的是 Unicode 字符和经过 UTF-8 编码的字符串，而bytes包面对的则主要是字节和字节切片。**

### bytes.Buffer

bytes.Buffer类型的用途主要是**作为字节序列的缓冲区**。与strings.Builder类型一样，bytes.Buffer也是开箱即用的。

strings.Builder 只能拼接和导出字符串，而 **bytes.Buffer 不但可以拼接、截断其中的字节序列**，以各种形式导出其中的内容，还可以顺序地读取其中的子序列。

可以说，bytes.Buffer是集读、写功能于一身的数据类型。当然了，这些也基本上都是作为一个缓冲区应该拥有的功能。

bytes.Buffer类型同样是**使用字节切片作为内容容器的**。并且，与strings.Reader类型类似，bytes.Buffer有一个int类型的字段，用于代表已读字节的计数，可以简称为已读计数。

```go
var buffer1 bytes.Buffer
contents := "Simple byte buffer for marshaling data."
fmt.Printf("Writing contents %q ...\n", contents)
buffer1.WriteString(contents)
fmt.Printf("The length of buffer: %d\n", buffer1.Len())
fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
```



我先声明了一个 bytes.Buffer 类型的变量 buffer1，并写入了一个字符串。然后，我想打印出这个 bytes.Buffer 类型的值（以下简称Buffer值）的长度和容量。在运行这段代码之后，我们将会看到如下的输出：

```go
Writing contents "Simple byte buffer for marshaling data." ...
The length of buffer: 39
The capacity of buffer: 64
```

乍一看这没什么问题。长度39和容量64的含义看起来与我们已知的概念是一致的。我向缓冲区中写入了一个长度为39的字符串，所以buffer1的长度就是39。

根据切片的自动扩容策略，64 这个数字也是合理的。另外，可以想象，这时的已读计数的值应该是 0，这是因为我还没有调用任何用于读取其中内容的方法。

可实际上，与 strings.Reader 类型的 Len 方法一样，**buffer1 的 Len 方法返回的也是内容容器中未被读取部分的长度，而不是其中已存内容的总长度（以下简称内容长度）**。示例如下：

```go
p1 := make([]byte, 7)
n, _ := buffer1.Read(p1)
fmt.Printf("%d bytes were read. (call Read)\n", n)
fmt.Printf("The length of buffer: %d\n", buffer1.Len())
fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
```

当我从 buffer1 中读取一部分内容，并用它们填满长度为7的字节切片 p1 之后，buffer1 的 Len 方法返回的结果值也会随即发生变化。如果运行这段代码，我们会发现，这个缓冲区的长度已经变为了32。

另外，因为我们并没有再向该缓冲区中写入任何内容，所以它的容量会保持不变，仍是64。

**Buffer 值的长度是未读内容的长度，而不是已存内容的总长度**。 它与在当前值之上的读操作和写操作都有关系，并会随着这两种操作的进行而改变，它可能会变得更小，也可能会变得更大。

而 Buffer 值的容量指的是它的**内容容器（也就是那个字节切片）的容量**，它**只与在当前值之上的写操作有关，并会随着内容的写入而不断增长**。

由于 strings.Reader 还有一个 **Size 方法可以给出内容长度的值，所以我们用内容长度减去未读部分的长度，就可以很方便地得到它的已读计数**。

然而，bytes.Buffer类型却没有这样一个方法，**它只有 Cap 方法。可是 Cap 方法提供的是内容容器的容量，也不是内容长度。**

这里的内容容器容量在很多时候都与内容长度不相同。**因此，没有了现成的计算公式，只要遇到稍微复杂些的情况，我们就很难估算出Buffer值的已读计数。**

**一旦理解了已读计数这个概念，并且能够在读写的过程中**，**实时地获得已读计数和内容长度的值，我们就可以很直观地了解到当前Buffer值各种方法的行为了**。



### bytes.Buffer 类型的值记录的已读计数的作用

- 读取内容时，相应方法会依据已读计数找到未读部分，并在读取后更新计数。
- 写入内容时，如需扩容，相应方法会根据已读计数实现扩容策略。
- 截断内容时，相应方法截掉的是已读计数代表索引之后的未读部分。
- 读回退时，相应方法需要用已读计数记录回退点。
- 重置内容时，相应方法会把已读计数置为0。
- 导出内容时，相应方法只会导出已读计数代表的索引之后的未读部分。
- 获取长度时，相应方法会依据已读计数和内容容器的长度，计算未读部分的长度并返回。

在读取内容的时候，相应方法会先根据已读计数，判断一下内容容器中是否还有未读的内容。如果有，那么它就会从已读计数代表的索引处开始读取。

在读取完成后，它还会及时地更新已读计数。也就是说，它会记录一下又有多少个字节被读取了。这里所说的相应方法包括了所有名称以Read开头的方法，以及Next方法和WriteTo方法。

在写入内容的时候，绝大多数的相应方法都会先检查当前的内容容器，是否有足够的容量容纳新的内容。如果没有，那么它们就会对内容容器进行扩容。然后，方法将会把已读计数的值置为0，以表示下一次读取需要从内容容器的第一个字节开始。用于写入内容的相应方法，包括了所有名称以Write开头的方法，以及ReadFrom方法。

用于截断内容的方法Truncate，会让很多对bytes.Buffer不太了解的程序开发者迷惑。 它会接受一个int类型的参数，这个参数的值代表了：在截断时需要保留头部的多少个字节。

需要注意的是，这里说的头部指的并不是内容容器的头部，而是其中的未读部分的头部。头部的起始索引正是由已读计数的值表示的。因此，在这种情况下，已读计数的值再加上参数值后得到的和，就是内容容器新的总长度。

在bytes.Buffer中，用于读回退的方法有UnreadByte和UnreadRune。 这两个方法分别用于回退一个字节和回退一个 Unicode 字符。调用它们一般都是为了退回在上一次被读取内容末尾的那个分隔符，或者为重新读取前一个字节或字符做准备。

不过，退回的前提是，在调用它们之前的那一个操作必须是“读取”，并且是成功的读取，否则这些方法就只能忽略后续操作并返回一个非nil的错误值。UnreadByte方法的做法比较简单，把已读计数的值减1就好了。而UnreadRune方法需要从已读计数中减去的，是上一次被读取的 Unicode 字符所占用的字节数。

这个字节数由bytes.Buffer的另一个字段负责存储，它在这里的有效取值范围是[1, 4]。只有ReadRune方法才会把这个字段的值设定在此范围之内。由此可见，只有紧接在调用ReadRune方法之后，对UnreadRune方法的调用才能够成功完成。该方法明显比UnreadByte方法的适用面更窄。

### bytes.Buffer 的扩容策略

Buffer 值既可以被手动扩容，也可以进行自动扩容。并且，这两种扩容方式的策略是基本一致的。所以，除非我们完全确定后续内容所需的字节数，否则让 Buffer 值自动去扩容就好了。

在扩容的时候，Buffer值中相应的代码（以下简称扩容代码）会先判断内容容器的剩余容量，是否可以满足调用方的要求，或者是否足够容纳新的内容。如果可以，那么扩容代码会在当前的内容容器之上，进行长度扩充。更具体地说，如果内容容器的容量与其长度的差，大于或等于另需的字节数，那么**扩容代码就会通过切片操作对原有的内容容器的长度进行扩充，就像下面这样：**

```go
b.buf = b.buf[:length+need]
```

如果内容容器的剩余容量不够了，那么扩容代码可能就会用新的内容容器去替代原有的内容容器，从而实现扩容。

如果当前内容容器的容量的一半，仍然大于或等于其现有长度再加上另需的字节数的和:

```go
cap(b.buf)/2 >= len(b.buf)+need
```

扩容代码就会复用现有的内容容器，并把容器中的未读内容拷贝到它的头部位置。这也意味着其中的已读内容，将会全部被未读内容和之后的新内容覆盖掉。

若这一步优化未能达成，也就是说，当前内容容器的容量小于新长度的二倍。

那么，扩容代码就只能再创建一个新的内容容器，并把原有容器中的未读内容拷贝进去，最后再用新的容器替换掉原有的容器。这个新容器的容量将会等于原有容量的二倍再加上另需字节数的和。

**新容器的容量 =2* 原有容量 + 所需字节数**

通过上面这些步骤，对内容容器的扩充基本上就完成了。不过，为了内部数据的一致性，以及避免原有的已读内容可能造成的数据混乱，扩容代码还会把已读计数置为0，并再对内容容器做一下切片操作，以掩盖掉原有的已读内容。

顺便说一下，对于处在零值状态的 Buffer 值来说，如果第一次扩容时的另需字节数不大于64，那么该值就会基于一个预先定义好的、长度为64的字节数组来创建内容容器。

在这种情况下，这个内容容器的容量就是64。这样做的目的是为了让Buffer值在刚被真正使用的时候就可以快速地做好准备。

### bytes.Buffer中内容泄露

内容泄露是指，使用Buffer值的一方通过某种非标准的（或者说不正式的）方式，得到了本不该得到的内容。

比如说，我通过调用 Buffer 值的某个用于读取内容的方法，得到了一部分未读内容。我应该，也只应该通过这个方法的结果值，拿到在那一时刻 Buffer 值中的未读内容。

但是，在这个 Buffer 值又有了一些新内容之后，我却可以通过当时得到的结果值，直接获得新的内容，而不需要再次调用相应的方法。

这就是典型的非标准读取方式。这种读取方式是不应该存在的，即使存在，我们也不应该使用。因为它是在无意中（或者说一不小心）暴露出来的，其行为很可能是不稳定的。

在 bytes.Buffer 中，Bytes方法和 Next 方法都可能会造成内容的泄露。原因在于，它们都把基于内容容器的切片直接返回给了方法的调用方。

通过切片，我们可以直接访问和操纵它的底层数组。不论这个切片是基于某个数组得来的，还是通过对另一个切片做切片操作获得的，都是如此。

Bytes 方法和 Next 方法返回的字节切片，都是通过对内容容器做切片操作得到的。也就是说，它们与内容容器共用了同一个底层数组，起码在一段时期之内是这样的。

以Bytes方法为例。它会返回在调用那一刻其所属值中的所有未读内容。示例代码如下：

```go
contents := "ab"
buffer1 := bytes.NewBufferString(contents)
fmt.Printf("The capacity of new buffer with contents %q: %d\n",
 contents, buffer1.Cap()) // 内容容器的容量为：8。
unreadBytes := buffer1.Bytes()
fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // 未读内容为：[97 98]。
```

用字符串值"ab"初始化了一个Buffer值，由变量buffer1代表，并打印了当时该值的一些状态。

又向该值写入了字符串值"cdefg"，此时，其容量仍然是8。我在前面通过调用buffer1的Bytes方法得到的结果值unreadBytes，包含了在那时其中的所有未读内容。

但是，由于这个结果值与buffer1的内容容器在此时还共用着同一个底层数组，所以，我只需通过简单的再切片操作，就可以利用这个结果值拿到buffer1在此时的所有未读内容。如此一来，buffer1的新内容就被泄露出来了。

```go
buffer1.WriteString("cdefg")
fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 内容容器的容量仍为：8。
unreadBytes = unreadBytes[:cap(unreadBytes)]
fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // 基于前面获取到的结果值可得，未读内容为：[97 98 99 100 101 102 103 0]。
```

如果我当时把unreadBytes的值传到了外界，那么外界就可以通过该值操纵buffer1的内容了，就像下面这样：

```go
unreadBytes[len(unreadBytes)-2] = byte('X') // 'X'的ASCII编码为88。
fmt.Printf("The unread bytes of the buffer: %v\n", buffer1.Bytes()) // 未读内容变为了：[97 98 99 100 101 102 88]。
```

不过，如果经过扩容，Buffer值的内容容器或者它的底层数组被重新设定了，那么之前的内容泄露问题就无法再进一步发展了。