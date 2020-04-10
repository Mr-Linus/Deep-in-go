## Go socket 与 IPC

socket，常被翻译为套接字，所谓 socket，是一种 IPC 方法。IPC 是 Inter-Process Communication 的缩写，可以被翻译为进程间通信。

顾名思义，IPC 这个概念（或者说规范）主要定义的是多个进程之间，相互通信的方法。这些方法主要包括：系统信号（signal）、管道（pipe）、套接字 （socket）、文件锁（file lock）、消息队列（message queue）、信号灯（semaphore，有的地方也称之为信号量）等。现存的主流操作系统大都对 IPC 提供了强有力的支持，尤其是 socket。

Go 语言对 IPC 也提供了一定的支持。比如，在os代码包和os/signal代码包中就有针对系统信号的 API。又比如，os.Pipe函数可以创建命名管道，而os/exec代码包则对另一类管道（匿名管道）提供了支持。对于 socket，Go 语言与之相应的程序实体都在其标准库的net代码包中。

在众多的 IPC 方法中，socket 是最为通用和灵活的一种。与其他的 IPC 方法不同，利用 socket 进行通信的进程，可以不局限在同一台计算机当中。

通信的双方无论存在于世界上的哪个角落，只要能够通过计算机的网卡端口以及网络进行互联，就可以使用 socket。

支持 socket 的操作系统一般都会对外提供一套 API。跑在它们之上的应用程序利用这套 API，就可以与互联网上的另一台计算机中的程序、同一台计算机中的其他程序，甚至同一个程序中的其他线程进行通信。

在 Linux 操作系统中，用于创建 socket 实例的 API，就是由一个名为socket的系统调用代表的。这个系统调用是 Linux 内核的一部分。

在 Go 语言标准库的syscall代码包中，有一个与这个socket系统调用相对应的函数。这两者的函数签名是基本一致的，它们都会接受三个int类型的参数，并会返回一个可以代表文件描述符的结果。

但不同的是，syscall包中的Socket函数本身是平台不相关的。在其底层，Go 语言为它支持的每个操作系统都做了适配，这才使得这个函数无论在哪个平台上，总是有效的。

Go 语言的net代码包中的很多程序实体，都会直接或间接地使用到syscall.Socket函数。

比如，我们在调用 net.Dial 函数的时候，会为它的两个参数设定值。其中的第一个参数名为network，它决定着 Go 程序在底层会创建什么样的 socket 实例，并使用什么样的协议与其他程序通信。

#### net.Dial函数

net.Dial函数会接受两个参数，分别名为network和address，都是string类型的。

参数network常用的可选值一共有 9 个。这些值分别代表了程序底层创建的 socket 实例可使用的不同通信协议，罗列如下。

- "tcp"：代表 TCP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
- "tcp4"：代表基于 IP 协议第四版的 TCP 协议。
- "tcp6"：代表基于 IP 协议第六版的 TCP 协议。
- "udp"：代表 UDP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
- "udp4"：代表基于 IP 协议第四版的 UDP 协议。
- "udp6"：代表基于 IP 协议第六版的 UDP 协议。
- "unix"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_STREAM 为 socket 类型。
- "unixgram"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_DGRAM 为 socket 类型。
- "unixpacket"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_SEQPACKET 为 socket 类型。

syscall.Socket函数接受的那三个参数，都是int类型的。这些参数所代表的分别是想要创建的 socket 实例通信域、类型以及使用的协议。

Socket 的通信域主要有这样几个可选项：IPv4 域、IPv6 域和 Unix 域。我想你应该能够猜出 IPv4 域、IPv6 域的含义，它们对应的分别是基于 IP 协议第四版的网络，和基于 IP 协议第六版的网络。

现在的计算机网络大都是基于 IP 协议第四版的，但是由于现有 IP 地址的逐渐枯竭，网络世界也在逐步地支持 IP 协议第六版。

Unix 域，指的是一种类 Unix 操作系统中特有的通信域。在装有此类操作系统的同一台计算机中，应用程序可以基于此域建立 socket 连接。以上三种通信域分别可以由syscall代码包中的常量AF_INET、AF_INET6和AF_UNIX表示。

Socket 的类型一共有 4 种，分别是：SOCK_DGRAM、SOCK_STREAM、SOCK_SEQPACKET以及SOCK_RAW。syscall代码包中也都有同名的常量与之对应。

前两者更加常用一些。**SOCK_DGRAM中的“DGRAM”代表的是 datagram，即数据报文**。它是一种**有消息边界，但没有逻辑连接的非可靠 socket 类型，我们熟知的基于 UDP 协议的网络通信就属于此类**。

有消息边界的意思是，与 **socket 相关的操作系统内核中的程序（以下简称内核程序）在发送或接收数据的时候是以消息为单位的。**

你可以把消息理解为带有固定边界的一段数据。内核程序可以自动地识别和维护这种边界，并在必要的时候，把数据切割成一个一个的消息，或者把多个消息串接成连续的数据。

如此一来，应用程序只需要面向消息进行处理就可以了。所谓的有逻辑连接是指，通信双方在收发数据之前必须先建立网络连接。待连接建立好之后，双方就可以一对一地进行数据传输了。

显然，基于 UDP 协议的网络通信并不需要这样，它是没有逻辑连接的。只要应用程序指定好对方的网络地址，内核程序就可以立即把数据报文发送出去。这有优势，也有劣势。

优势是发送速度快，不长期占用网络资源，并且每次发送都可以指定不同的网络地址。当然了，最后一个优势有时候也是劣势，因为这会使数据报文更长一些。其他的劣势有，无法保证传输的可靠性，不能实现数据的有序性，以及数据只能单向进行传输。

**而 SOCK_STREAM 这个 socket 类型，恰恰与 SOCK_DGRAM 相反**。**它没有消息边界，但有逻辑连接，能够保证传输的可靠性和数据的有序性，同时还可以实现数据的双向传输**。**众所周知的基于 TCP 协议的网络通信就属于此类。**

这样的网络通信传输数据的形式是字节流，而不是数据报文。

字节流是以字节为单位的。内核程序无法感知一段字节流中包含了多少个消息，以及这些消息是否完整，这完全需要应用程序自己去把控。

不过，此类网络通信中的一端，总是会忠实地按照另一端发送数据时的字节排列顺序，接收和缓存它们。所以，应用程序需要根据双方的约定去数据中查找消息边界，并按照边界切割数据，仅此而已。

syscall.Socket函数的第三个参数用于表示 socket 实例所使用的协议。

通常，只要明确指定了前两个参数的值，我们就无需再去确定第三个参数值了，一般把它置为0就可以了。这时，内核程序会自行选择最合适的协议。

比如，当前两个参数值分别为syscall.AF_INET和syscall.SOCK_DGRAM的时候，内核程序会选择 UDP 作为协议。又比如，在前两个参数值分别为syscall.AF_INET6和syscall.SOCK_STREAM时，内核程序可能会选择 TCP 作为协议。

![99f8a0405a98ea16495364be352fe969](https://static001.geekbang.org/resource/image/99/69/99f8a0405a98ea16495364be352fe969.png)



#### net.DialTimeout函数

net.DialTimeout函数时给定的超时时间，代表着函数为网络连接建立完成而等待的最长时间。这是一个相对的时间。它会由这个函数的参数timeout的值表示。

**开始的时间点几乎是我们调用net.DialTimeout函数的那一刻**。在这之后，**时间会主要花费在“解析参数network和address的值”，以及“创建 socket 实例并建立网络连接”这两件事情上**。

不论执行到哪一步，只要在绝对的超时时间达到的那一刻，网络连接还没有建立完成，该函数就会返回一个代表了 I/O 操作超时的错误值。

值得注意的是，**在解析address的值的时候，函数会确定网络服务的 IP 地址、端口号等必要信息，并在需要时访问 DNS 服务。**

另外，**如果解析出的 IP 地址有多个，那么函数会串行或并发地尝试建立连接。但无论用什么样的方式尝试，函数总会以最先建立成功的那个连接为准。**同时，**它还会根据超时前的剩余时间，去设定针对每次连接尝试的超时时间，以便让它们都有适当的时间执行。**

在 net 包中还有一个名为 Dialer 的结构体类型。该类型有一个名叫 Timeout 的字段，它与上述的 timeout 参数的含义是完全一致的。实际上，net.DialTimeout 函数正是利用了这个类型的值才得以实现功能的。net.Dialer类型值得你好好学习一下，尤其是它的每个字段的功用以及它的DialContext方法。

#### net/http 包

HTTP 协议是基于 TCP/IP 协议栈的，并且它也是一个面向普通文本的协议。

原则上，我们使用任何一个文本编辑器，都可以轻易地写出一个完整的 HTTP 请求报文。只要你搞清楚了请求报文的头部（header）和主体（body）应该包含的内容，这样做就会很容易。所以，在这种情况下，即便直接使用net.Dial函数，你应该也不会感觉到困难。不过，不困难并不意味着很方便。如果我们只是访问基于 HTTP 协议的网络服务的话，那么使用net/http代码包中的程序实体来做，显然会更加便捷。其中，最便捷的是使用http.Get函数。我们在调用它的时候只需要传给它一个 URL 就可以了，比如像下面这样：

```go
url1 := "http://google.cn"
fmt.Printf("Send request to %q with method GET ...\n", url1)
resp1, err := http.Get(url1)
if err != nil {
  fmt.Printf("request sending error: %v\n", err)
}
defer resp1.Body.Close()
line1 := resp1.Proto + " " + resp1.Status
fmt.Printf("The first line of response:\n%s\n", line1)
```

http.Get函数会返回两个结果值。第一个结果值的类型是\*http.Response，它是网络服务给我们传回来的响应内容的结构化表示。

第二个结果值是error类型的，它代表了在创建和发送 HTTP 请求，以及接收和解析 HTTP 响应的过程中可能发生的错误。http.Get函数会在内部使用缺省的 HTTP 客户端，并且调用它的Get方法以完成功能。这个缺省的 HTTP 客户端是由net/http包中的公开变量DefaultClient代表的，其类型是*http.Client。它的基本类型也是可以被拿来使用的，甚至它还是开箱即用的。下面的这两行代码：

```go
var httpClient1 http.Client
resp2, err := httpClient1.Get(url1)
```

与前面的这一行代码

```go
resp1, err := http.Get(url1)
```

是等价的。

http.Client是一个结构体类型，并且它包含的字段都是公开的。**之所以该类型的零值仍然可用，是因为它的这些字段要么存在着相应的缺省值，要么其零值直接就可以使用，且代表着特定的含义。**

#### http.Client Transport

http.Client类型中的Transport字段代表着：**向网络服务发送 HTTP 请求，并从网络服务接收 HTTP 响应的操作过程**。也就是说，该字段的方法**RoundTrip应该实现单次 HTTP 事务（或者说基于 HTTP 协议的单次交互）需要的所有步骤。**

**这个字段是 http.RoundTripper 接口类型的，它有一个由 http.DefaultTransport 变量代表的缺省值（以下简称DefaultTransport）**。当我们在初始化一个http.Client类型的值（以下简称Client值）的时候，如果**没有显式地为该字段赋值，那么这个Client值就会直接使用DefaultTransport**。

http.Client类型的 Timeout字段，代表的正是前面所说的单次 HTTP 事务的超时时间，它是time.Duration类型的。它的零值是可用的，用于表示没有设置超时时间。

DefaultTransport的实际类型是*http.Transport，后者即为http.RoundTripper接口的默认实现。这个类型是可以被复用的，也推荐被复用，同时，它也是并发安全的。正因为如此，http.Client类型也拥有着同样的特质。

http.Transport类型，会在内部使用一个net.Dialer类型的值（以下简称Dialer值），并且，它会把该值的Timeout字段的值，设定为30秒。

也就是说，这个Dialer值如果**在 30 秒内还没有建立好网络连接，那么就会被判定为操作超时**。在DefaultTransport的值被初始化的时候，这样的Dialer值的DialContext方法会被赋给前者的DialContext字段。

http.Transport类型还包含了很多其他的字段，其中有一些字段是关于操作超时的:

- IdleConnTimeout：含义是空闲的连接在多久之后就应该被关闭。
- DefaultTransport会把该字段的值设定为90秒。如果该值为0，那么就表示不关闭空闲的连接。注意，这样很可能会造成资源的泄露。
- ResponseHeaderTimeout：含义是，从客户端把请求完全递交给操作系统到从操作系统那里接收到响应报文头的最大时长。DefaultTransport并没有设定该字段的值。
- ExpectContinueTimeout：含义是，在**客户端递交了请求报文头之后，等待接收第一个响应报文头的最长时间**。在客户端想要使用 HTTP 的“POST”方法把一个很大的报文体发送给服务端的时候，它可以先通过发送一个包含了“Expect: 100-continue”的请求报文头，来询问服务端是否愿意接收这个大报文体。这个字段就是用于设定在这种情况下的超时时间的。注意，**如果该字段的值不大于0，那么无论多大的请求报文体都将会被立即发送出去。这样可能会造成网络资源的浪费**。DefaultTransport把该字段的值设定为了1秒。
- TLSHandshakeTimeout：TLS 是 Transport Layer Security 的缩写，可以被翻译为传输层安全。这个字段代表了基于 TLS 协议的连接在被建立时的握手阶段的超时时间。**若该值为0，则表示对这个时间不设限。DefaultTransport把该字段的值设定为了10秒。**

此外，还有一些与IdleConnTimeout相关的字段值得我们关注，即：MaxIdleConns、MaxIdleConnsPerHost以及MaxConnsPerHost。

无论当前的http.Transport类型的值（以下简称Transport值）访问了多少个网络服务，**MaxIdleConns字段都只会对空闲连接的总数做出限定**。而MaxIdleConnsPerHost字段限定的则是，该Transport值访问的每一个网络服务的最大空闲连接数。

每一个网络服务都会有自己的网络地址，可能会使用不同的网络协议，对于一些 HTTP 请求也可能会用到代理。Transport值正是通过这三个方面的具体情况，来鉴别不同的网络服务的。

MaxIdleConnsPerHost字段的缺省值，由http.DefaultMaxIdleConnsPerHost变量代表，值为2。也就是说，在默认情况下，对于某一个Transport值访问的每一个网络服务，它的空闲连接数都最多只能有两个。

与MaxIdleConnsPerHost字段的含义相似的，是MaxConnsPerHost字段。不过，后者限制的是，针对某一个Transport值访问的每一个网络服务的最大连接数，不论这些连接是否是空闲的。并且，该字段没有相应的缺省值，它的零值表示不对此设限。

DefaultTransport并没有显式地为MaxIdleConnsPerHost和MaxConnsPerHost这两个字段赋值，但是它却把MaxIdleConns字段的值设定为了100。在默认情况下，空闲连接的总数最大为100，而针对每个网络服务的最大空闲连接数为2。

HTTP 协议有一个请求报文头叫做“Connection”。在 HTTP 协议的 1.1 版本中，这个报文头的值默认是“keep-alive”。在这种情况下的网络连接都是持久连接，它们会在当前的 HTTP 事务完成后仍然保持着连通性，因此是可以被复用的。既然连接可以被复用，那么就会有两种可能。一种可能是，**针对于同一个网络服务，有新的 HTTP 请求被递交，该连接被再次使用**。另一种可能是，**不再有对该网络服务的 HTTP 请求，该连接被闲置**。

显然，后一种可能就产生了空闲的连接。另外，如果分配给某一个网络服务的连接过多的话，也可能会导致空闲连接的产生，因为每一个新递交的 HTTP 请求，都只会征用一个空闲的连接。所以，为空闲连接设定限制，在大多数情况下都是很有必要的，也是需要斟酌的。

如果我们想彻底地杜绝空闲连接的产生，那么可以在初始化Transport值的时候把**它的DisableKeepAlives字段的值设定为true**。这时，**HTTP 请求的“Connection”报文头的值就会被设置为“close”。这会告诉网络服务，这个网络连接不必保持，当前的 HTTP 事务完成后就可以断开它了**。如此一来，每当一个 HTTP 请求被递交时，就都会产生一个新的网络连接。这样做**会明显地加重网络服务以及客户端的负载**，并会让每个 HTTP 事务都耗费更多的时间。所以，**在一般情况下，我们都不要去设置这个DisableKeepAlives字段。**

在net.Dialer类型中，也有一个看起来很相似的字段KeepAlive。不过，它与前面所说的 HTTP 持久连接并不是一个概念，KeepAlive是直接作用在底层的 socket 上的。它的背后是一种针对网络连接（更确切地说，是 TCP 连接）的存活探测机制。它的值用于表示每间隔多长时间发送一次探测包。当该值不大于0时，则表示不开启这种机制。DefaultTransport会把这个字段的值设定为30秒。

#### http.Server ListenAndServe

http.Server类型与http.Client是相对应的。http.Server代表的是基于 HTTP 协议的服务端，或者说网络服务。

http.Server类型的ListenAndServe方法的功能是：**监听一个基于 TCP 协议的网络地址，并对接收到的 HTTP 请求进行处理**。这个方法会默认开启针对网络连接的存活探测机制，以保证连接是持久的。同时，该方法会一直执行，直到有严重的错误发生或者被外界关掉。当被外界关掉时，它**会返回一个由http.ErrServerClosed变量代表的错误值**。

这个ListenAndServe方法主要会做下面这几件事情。

- 检查当前的http.Server类型的值（以下简称当前值）的Addr字段。该字段的值代表了当前的网络服务需要使用的网络地址，即：IP 地址和端口号. 如果这个字段的值为空字符串，那么就用":http"代替。也就是说，使用任何可以代表本机的域名和 IP 地址，并且端口号为80。
- 通过调用net.Listen函数在已确定的网络地址上启动基于 TCP 协议的监听。
- 检查net.Listen函数返回的错误值。如果该错误值不为nil，那么就直接返回该值。否则，通过调用当前值的Serve方法准备接受和处理将要到来的 HTTP 请求。

net.Listen函数会先解析参数值中包含的网络地址隐含的 IP 地址和端口号，根据给定的网络协议，确定监听的方法，并开始进行监听。

http.Server类型的Serve方法，在一个for循环中，网络监听器的Accept方法会被不断地调用，该方法会返回两个结果值；第一个结果值是net.Conn类型的，它会代表包含了新到来的 HTTP 请求的网络连接；第二个结果值是代表了可能发生的错误的error类型值。

如果这个错误值不为nil，除非它代表了一个暂时性的错误，否则循环都会被终止。如果是暂时性的错误，那么循环的下一次迭代将会在一段时间之后开始执行。

如果这里的Accept方法没有返回非nil的错误值，那么这里的程序将会先把它的第一个结果值包装成一个*http.conn类型的值（以下简称conn值），然后通过在新的 goroutine 中调用这个conn值的serve方法，来对当前的 HTTP 请求进行处理。

