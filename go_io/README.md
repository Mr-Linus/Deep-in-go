### io 包中的接口和工具

strings.Builder类型主要用于构建字符串，它的指针类型实现的接口有io.Writer、io.ByteWriter和fmt.Stringer。

另外，它其实还实现了一个io包的包级私有接口io.stringWriter（自 Go 1.12 起它会更名为io.StringWriter）。strings.Reader类型主要用于读取字符串，它的指针类型实现的接口比较多，包括：

- io.Reader；
- io.ReaderAt；
- io.ByteReader；
- io.RuneReader；
- io.Seeker；
- io.ByteScanner；
- io.RuneScanner；
- io.WriterTo；

共有 8 个，它们都是io包中的接口。其中，io.ByteScanner是io.ByteReader的扩展接口，而io.RuneScanner又是io.RuneReader的扩展接口。

bytes.Buffer是集读、写功能于一身的数据类型，它非常适合作为字节序列的缓冲区。 它的指针类型实现的接口就更多了。更具体地说，该指针类型实现的读取相关的接口有下面几个。

- io.Reader；
- io.ByteReader；
- io.RuneReader；
- io.ByteScanner；
- io.RuneScanner；
- io.WriterTo；

共有 6 个。而其实现的写入相关的接口则有这些。

- io.Writer；
- io.ByteWriter；
- io.stringWriter；
- io.ReaderFrom；

共 4 个。此外，它还实现了导出相关的接口fmt.Stringer。

在io包中，有这样几个用于拷贝数据的函数，它们是：

io.Copy；io.CopyBuffer；io.CopyN。

虽然这几个函数在功能上都略有差别，但是它们都首先会接受两个参数，即：用于代表数据目的地、io.Writer类型的参数dst，以及用于代表数据来源的、io.Reader类型的参数src。这些函数的功能大致上都是把数据从src拷贝到dst。

不论我们给予它们的第一个参数值是什么类型的，只要这个类型实现了io.Writer接口即可。

同样的，无论我们传给它们的第二个参数值的实际类型是什么，只要该类型实现了io.Reader接口就行。一旦我们满足了这两个条件，这些函数几乎就可以正常地执行了。

当然了，函数中还会对必要的参数值进行有效性的检查，如果检查不通过，它的执行也是不能够成功结束的。

```go
src := strings.NewReader(
 "CopyN copies n bytes (or until an error) from src to dst. " +
  "It returns the number of bytes copied and " +
  "the earliest error encountered while copying.")
dst := new(strings.Builder)
written, err := io.CopyN(dst, src, 58)
if err != nil {
 fmt.Printf("error: %v\n", err)
} else {
 fmt.Printf("Written(%d): %q\n", written, dst.String())
}
```

先使用strings.NewReader创建了一个字符串读取器，并把它赋给了变量src，然后我又new了一个字符串构建器，并将其赋予了变量dst。

之后，我在调用io.CopyN函数的时候，把这两个变量的值都传了进去，同时把给这个函数的第三个参数值设定为了58。也就是说，**我想从src中拷贝前58个字节到dst那里。**

虽然，变量src和dst的类型分别是strings.Reader和strings.Builder，但是当它们被传到io.CopyN函数的时候，就已经分别被包装成了io.Reader类型和io.Writer类型的值。io.CopyN函数也根本不会去在意，它们的实际类型到底是什么。

为了优化的目的，io.CopyN函数中的代码会对参数值进行再包装，也会检测这些参数值是否还实现了别的接口，甚至还会去探求某个参数值被包装后的实际类型，是否为某个特殊的类型。

从总体上来看，这些代码都是面向参数声明中的接口来做的。io.CopyN 函数的作者通过面向接口编程，极大地拓展了它的适用范围和应用场景。

换个角度看，正因为 strings.Reader 类型和 strings.Builder 类型都实现了不少接口，所以它们的值才能够被使用在更广阔的场景中。

Go 语言的各种库中，能够操作它们的函数和数据类型明显多了很多。这就是strings包和bytes包中的数据类型在实现了若干接口之后得到的最大好处。

这就是面向接口编程带来的最大优势。这些数据类型和函数的做法，也是非常值得我们在编程的过程中去效仿的。

可以看到，前文所述的几个类型实现的大都是io代码包中的接口。实际上，io包中的接口，对于 Go 语言的标准库和很多第三方库而言，都起着举足轻重的作用。它们非常基础也非常重要。

拿io.Reader和io.Writer这两个最核心的接口来说，它们是很多接口的扩展对象和设计源泉。同时，单从 Go 语言的标准库中统计，实现了它们的数据类型都（各自）有上百个，而引用它们的代码更是都（各自）有 400 多处。

很多数据类型实现了io.Reader接口，是因为它们提供了从某处读取数据的功能。类似的，许多能够把数据写入某处的数据类型，也都会去实现io.Writer接口。

有不少类型的设计初衷都是：实现这两个核心接口的某个，或某些扩展接口，以提供比单纯的字节序列读取或写入，更加丰富的功能，就像前面讲到的那几个strings包和bytes包中的数据类型那样。

在 Go 语言中，对接口的扩展是通过接口类型之间的嵌入来实现的，这也常被叫做接口的组合。

Go 语言提倡使用小接口加接口组合的方式，来扩展程序的行为以及增加程序的灵活性。io代码包恰恰就可以作为这样的一个标杆，它可以成为我们运用这种技巧时的一个参考标准。

io.Reader接口为对象提出一个与接口扩展和实现有关的问题。如果你研究过这个核心接口以及相关的数据类型的话，这个问题回答起来就并不困难。

### io.Reader的扩展接口和实现类型

- io.ReadWriter：此接口既是io.Reader的扩展接口，也是io.Writer的扩展接口。换句话说，该接口定义了一组行为，包含且**仅包含了基本的字节序列读取方法Read，和字节序列写入方法Write**。

- io.ReadCloser：此接口除了包含基本的字节序列读取方法之外，**还拥有一个基本的关闭方法Close。后者一般用于关闭数据读写的通路。这个接口其实是io.Reader接口和io.Closer接口的组合**。

- io.ReadWriteCloser：很明显，此接口是**io.Reader、io.Writer和io.Closer这三个接口的组合。**

- io.ReadSeeker：此接口的特点是拥有一个**用于寻找读写位置的基本方法Seek**。更具体地说，该方法可以根据给定的偏移量基于数据的起始位置、末尾位置，或者当前读写位置去寻找新的读写位置。这个新的读写位置用于表明下一次读或写时的起始索引。**Seek是io.Seeker接口唯一拥有的方法**。

再来说说io包中的io.Reader接口的实现类型，它们包括下面几项内容。

- \*io.LimitedReader：此类型的**基本类型会包装io.Reader类型的值，并提供一个额外的受限读取的功能**。所谓的受限读取指的是，**此类型的读取方法Read返回的总数据量会受到限制，无论该方法被调用多少次。这个限制由该类型的字段N指明，单位是字节。**
- \*io.SectionReader：此类型的基本类型可以包装 io.ReaderAt类型的值，并**且会限制它的Read方法，只能够读取原始数据中的某一个部分（或者说某一段）**。这个**数据段的起始位置和末尾位置，需要在它被初始化的时候就指明，并且之后无法变更。该类型值的行为与切片有些类似，它只会对外暴露在其窗口之中的那些数据**。
- \*io.teeReader：此类型是一个包级私有的数据类型，也是io.TeeReader函数结果值的实际类型。这个**函数接受两个参数r和w，类型分别是io.Reader和io.Writer**。其结果值的**Read方法会把r中的数据经过作为方法参数的字节切片p写入到w**。可以说，这个值就是r和w之间的数据桥梁，而那个参数p就是这座桥上的数据搬运者。
- *io.multiReader：此类型也是一个包级私有的数据类型。类似的，**io包中有一个名为 MultiReader 的函数，它可以接受若干个io.Reader类型的参数值，并返回一个实际类型为io.multiReader的结果值**。当这个结果值的Read方法被调用时，它会顺序地从前面那些io.Reader类型的参数值中读取数据。因此，我们也可以称之为多对象读取器。

- \*io.pipe：此类型为一个包级私有的数据类型，它比上述类型都要复杂得多。它不但实现了io.Reader接口，而且还实现了io.Writer接口。实际上，**io.PipeReader类型和io.PipeWriter类型拥有的所有指针方法都是以它为基础的**。这些方法都**只是代理了io.pipe类型值所拥有的某一个方法而已**。又因为**io.Pipe函数会返回这两个类型的指针值并分别把它们作为其生成的同步内存管道的两端**，所以可以说，*io.pipe类型就是io包提供的同步内存管道的核心实现。

- *io.PipeReader：此类型可以被**视为io.pipe类型的代理类型**。它代理了后者的一部分功能，并**基于后者实现了io.ReadCloser接口**。同时，它还定义了同步内存管道的读取端。

### io包中的核心接口

io包中的核心接口只有 3 个，它们是：**io.Reader、io.Writer和io.Closer。**

还可以把io包中的简单接口分为四大类。这四大类接口分别针对于四种操作，即：**读取、写入、关闭和读写位置设定**。前三种操作属于基本的 I/O 操作。

关于读取操作，我们在前面已经重点讨论过核心接口io.Reader。它在io包中有 5 个扩展接口，并有 6 个实现类型。除了它，这个包中针对读取操作的接口还有不少。我们下面就来梳理一下。

#### 读接口

io.ByteReader和io.RuneReader这两个简单接口。它们分别定义了一个读取方法，即：**ReadByte 和ReadRune**。

但与io.Reader接口中Read方法不同的是，这两个读取方法分别**只能够读取下一个单一的字节和 Unicode 字符**。

之前讲过的数据类型 strings.Reader 和 bytes.Buffer 都是 io.ByteReader 和 io.RuneReader 的实现类型。

不仅如此，这两个类型还都实现了io.ByteScanner接口和io.RuneScanner接口。

io.ByteScanner接口内嵌了**简单接口io.ByteReader，并定义了额外的UnreadByte方法**。

如此一来，它就抽象出了一个能够**读取和读回退单个字节的功能集**。与之类似，io.RuneScanner**内嵌了简单接口io.RuneReader**，并定义了额外的UnreadRune方法。它抽象的是可以读取和读回退单个 Unicode 字符的功能集。

io.ReaderAt接口也是一个简单接口，其中只定义了一个方法ReadAt。与我们在前面说过的读取方法都不同，ReadAt是一个纯粹的只读方法。**只去读取其所属值中包含的字节**，而**不对这个值进行任何的改动**，比如，它绝对**不能去修改已读计数的值。这也是io.ReaderAt接口与其实现类型之间最重要的一个约定**。因此，**如果仅仅并发地调用某一个值的ReadAt方法，那么安全性应该是可以得到保障的**。

io.WriterTo 接口定义了一个名为WriteTo的方法。WriteTo方法**其实是一个读取方法**。它会**接受一个io.Writer类型的参数值**，并会**把其所属值中的数据读出并写入到这个参数值中**。与之相对应的是 **io.ReaderFrom 接口**。它**定义了一个名叫ReadFrom的写入方法**。该方法会**接受一个io.Reader类型的参数值，并会从该参数值中读出数据, 并写入到其所属值中**。

前面用到过的 io.CopyN 函数，在复制数据的时候会先检测其参数 src 的值，**是否实现了io.WriterTo接口**。**如果是，那么它就直接利用该值的WriteTo方法，把其中的数据拷贝给参数dst代表的值。**

类似的，这个函数还会检测dst的值是否实现了io.ReaderFrom接口。**如果是，那么它就会利用这个值的ReadFrom方法，直接从src那里把数据拷贝进该值**。

实际上，对于io.Copy函数和io.CopyBuffer函数来说也是如此，因为它们在内部做数据复制的时候用的都是同一套代码。

**io.ReaderFrom接口与io.WriterTo接口对应得很规整**。实际上，在io包中，与写入操作有关的接口都与读取操作的相关接口有着一定的对应关系。下面，我们就来说说写入操作相关的接口。

#### 写接口

首先当然是核心接口io.Writer。基于它的扩展接口除了有我们已知的io.ReadWriter、io.ReadWriteCloser和io.ReadWriteSeeker之外，还有io.WriteCloser和io.WriteSeeker。

之前提及的*io.pipe就是io.ReadWriter接口的实现类型。然而，在io包中并没有io.ReadWriteCloser接口的实现，它的实现类型主要集中在net包中。

除此之外，写入操作相关的简单接口还有io.ByteWriter和io.WriterAt。可惜，io包中也没有它们的实现类型。不过，有一个数据类型值得在这里提一句，那就是*os.File。这个类型不但是io.WriterAt接口的实现类型，还同时实现了io.ReadWriteCloser接口和io.ReadWriteSeeker接口。也就是说，该类型支持的 I/O 操作非常的丰富。io.Seeker接口作为一个读写位置设定相关的简单接口，也仅仅定义了一个方法，名叫Seek。该方法主要用于寻找并设定下一次读取或写入时的起始索引位置。io包中有几个基于io.Seeker的扩展接口，包括前面讲过的io.ReadSeeker和io.ReadWriteSeeker，以及还未曾提过的io.WriteSeeker。io.WriteSeeker是基于io.Writer和io.Seeker的扩展接口。之前多次提到的两个指针类型strings.Reader和io.SectionReader都实现了io.Seeker接口。顺便说一句，这两个类型也都是io.ReaderAt接口的实现类型。

最后，关闭操作相关的接口io.Closer非常通用，它的扩展接口和实现类型都不少。我们单从名称上就能够一眼看出io包中的哪些接口是它的扩展接口。至于它的实现类型，io包中只有io.PipeReader和io.PipeWriter。

![e5b4af00105769cdc9f0ab729bb3b30b](https://static001.geekbang.org/resource/image/e5/0b/e5b4af00105769cdc9f0ab729bb3b30b.png)



