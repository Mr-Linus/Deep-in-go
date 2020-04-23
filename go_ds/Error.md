## Error

error 类型其实是一个接口类型，也是一个 Go 语言的内建类型。在这个接口类型的声明中只包含了一个方法 Error。Error 方法不接受任何参数，但是会返回一个 string 类型的结果。它的作用是返回错误信息的字符串表示形式。

使用 error 类型的方式通常是，在函数声明的结果列表的最后，声明一个该类型的结果，同时在调用这个函数之后，先判断它返回的最后一个结果值是否“不为nil”。

如果这个值“不为nil”，那么就进入错误处理流程，否则就继续进行正常的流程。

#### Usecase

```go
package main

import (
  "errors"
  "fmt"
)

func echo(request string) (response string, err error) {
  if request == "" {
    err = errors.New("empty request")
    return
  }
  response = fmt.Sprintf("echo: %s", request)
  return
}

func main() {
  for _, req := range []string{"", "hello!"} {
    fmt.Printf("request: %s\n", req)
    resp, err := echo(req)
    if err != nil {
      fmt.Printf("error: %s\n", err)
      continue
    }
    fmt.Printf("response: %s\n", resp)
  }
}
```

先看echo函数的声明。echo 函数接受一个 string 类型的参数 request，并会返回两个结果。

这两个结果都是有名称的，第一个结果 response 也是 string 类型的，它代表了这个函数正常执行后的结果值。

第二个结果 err 就是 error 类型的，它代表了函数执行出错时的结果值，同时也包含了具体的错误信息。

当echo函数被调用时，它会先检查参数 request 的值。如果该值为空字符串，那么它就会通过调用 errors.New函数，为结果err赋值，然后忽略掉后边的操作并直接返回。此时，结果response的值也会是一个空字符串。如果request的值并不是空字符串，那么它就为结果response赋一个适当的值，然后返回，此时结果err的值会是nil。

在生成error类型值的时候，用到了errors.New函数。这是一种最基本的生成错误值的方式。我们调用它的时候传入一个由字符串代表的错误信息，它会给返回给我们一个包含了这个错误信息的error类型值。该值的静态类型当然是error，而动态类型则是一个在errors包中的，包级私有的类型*errorString。

errorString类型拥有的一个指针方法实现了error接口中的Error方法。这个方法在被调用后，会原封不动地返回我们之前传入的错误信息。实际上，error类型值的Error方法就相当于其他类型值的String方法。

对于其他类型的值来说，只要我们能为这个类型编写一个String方法，就可以自定义它的字符串表示形式。而对于error类型值，它的字符串表示形式则取决于它的Error方法。



### 错误类型判断

- 对于类型在已知范围内的一系列错误值，一般使用类型断言表达式或类型switch语句来判断；
- 对于已有相应变量且类型相同的一系列错误值，一般直接使用判等操作来判断；
- 对于没有相应变量且类型未知的一系列错误值，只能使用其错误信息的字符串表示形式来做判断。

类型在已知范围内的错误值其实是最容易分辨的。就拿os包中的几个代表错误的类型os.PathError、os.LinkError、os.SyscallError和os/exec.Error来说，它们的指针类型都是error接口的实现类型，同时它们也都包含了一个名叫Err，类型为error接口类型的代表潜在错误的字段。

如果我们得到一个error类型值，并且知道该值的实际类型肯定是它们中的某一个，那么就可以用类型switch语句去做判断。例如：

```go
func underlyingError(err error) error {
  switch err := err.(type) {
  case *os.PathError:
    return err.Err
  case *os.LinkError:
    return err.Err
  case *os.SyscallError:
    return err.Err
  case *exec.Error:
    return err.Err
  }
  return err
}
```

函数underlyingError的作用是：

获取和返回已知的操作系统相关错误的潜在错误值。其中的类型switch语句中有若干个case子句，分别对应了上述几个错误类型。当它们被选中时，都会把函数参数err的Err字段作为结果值返回。如果它们都未被选中，那么该函数就会直接把参数值作为结果返回，即放弃获取潜在错误值。

只要类型不同，我们就可以如此分辨。但是在错误值类型相同的情况下，这些手段就无能为力了。在 Go 语言的标准库中也有不少以相同方式创建的同类型的错误值。

我们还拿os包来说，其中不少的错误值都是通过调用errors.New函数来初始化的，比如：os.ErrClosed、os.ErrInvalid 以及 os.ErrPermission，等等。

注意，与前面讲到的那些错误类型不同，这几个都是已经定义好的、确切的错误值。os包中的代码有时候会把它们当做潜在错误值，封装进前面那些错误类型的值中。

```go
printError := func(i int, err error) {
  if err == nil {
    fmt.Println("nil error")
    return
  }
  err = underlyingError(err)
  switch err {
  case os.ErrClosed:
    fmt.Printf("error(closed)[%d]: %s\n", i, err)
  case os.ErrInvalid:
    fmt.Printf("error(invalid)[%d]: %s\n", i, err)
  case os.ErrPermission:
    fmt.Printf("error(permission)[%d]: %s\n", i, err)
  }
}
```

虽然我不知道这些错误值的类型的范围，但却知道它们或它们的潜在错误值一定是某个已经在os包中定义的值。

所以，我先用underlyingError函数得到它们的潜在错误值，当然也可能只得到原错误值而已。然后，我用switch语句对错误值进行判等操作，三个case子句分别对应我刚刚提到的那三个已存在于os包中的错误值。如此一来，我就能分辨出具体错误了。

对于上面这两种情况，我们都有明确的方式去解决。但是，如果我们对一个错误值可能代表的含义知之甚少，那么就只能通过它拥有的错误信息去做判断了。



### 构建错误值体系

构建错误值体系的基本方式有两种，即：创建立体的错误类型体系和创建扁平的错误值列表。

先说错误类型体系。由于在 Go 语言中实现接口是非侵入式的，所以我们可以做得很灵活。比如，在标准库的net代码包中，有一个名为Error的接口类型。它算是内建接口类型error的一个扩展接口，因为error是net.Error的嵌入接口。

net.Error接口除了拥有error接口的Error方法之外，还有两个自己声明的方法：Timeout和Temporary。

net包中有很多错误类型都实现了net.Error接口，比如：\*net.OpError；*net.AddrError；net.UnknownNetworkError等等。

当net包的使用者拿到一个错误值的时候，可以先判断它是否是net.Error类型的，也就是说该值是否代表了一个网络相关的错误。如果是，那么我们还可以再进一步判断它的类型是哪一个更具体的错误类型，这样就能知道这个网络相关的错误具体是由于操作不当引起的，还是因为网络地址问题引起的，又或是由于网络协议不正确引起的。

当我们细看net包中的这些具体错误类型的实现时，还会发现，与os包中的一些错误类型类似，它们也都有一个名为Err、类型为error接口类型的字段，代表的也是当前错误的潜在错误。

所以说，这些错误类型的值之间还可以有另外一种关系，即：链式关系。比如说，使用者调用net.DialTCP之类的函数时，net包中的代码可能会返回给他一个*net.OpError类型的错误值，以表示由于他的操作不当造成了一个错误。

同时，这些代码还可能会把一个*net.AddrError或net.UnknownNetworkError类型的值赋给该错误值的Err字段，以表明导致这个错误的潜在原因。如果，此处的潜在错误值的Err字段也有非nil的值，那么将会指明更深层次的错误原因。如此一级又一级就像链条一样最终会指向问题的根源。

用类型建立起树形结构的错误体系，用统一字段建立起可追根溯源的链式错误关联。这是 Go 语言标准库给予我们的优秀范本，非常有借鉴意义。

注意，如果你不想让包外代码改动你返回的错误值的话，一定要**小写其中字段的名称首字母**。你可以通过暴露某些方法让包外代码有进一步获取错误信息的权限，比如编写一个可以返回包级私有的err字段值的公开方法Err。



### 设计错误值列表

相比于立体的错误类型体系，扁平的错误值列表就要简单得多了。当我们只是想预先创建一些代表已知错误的错误值时候，用这种扁平化的方式就很恰当了。

由于error是接口类型，所以通过errors.New函数生成的错误值只能被赋给变量，而不能赋给常量，又由于这些代表错误的变量需要给包外代码使用，所以其访问权限只能是公开的。

这就带来了一个问题，如果有恶意代码改变了这些公开变量的值，那么程序的功能就必然会受到影响。因为在这种情况下我们往往会通过判等操作来判断拿到的错误值具体是哪一个错误，如果这些公开变量的值被改变了，那么相应的判等操作的结果也会随之改变。

有两个解决方案。第一个方案（不推荐）是，先私有化此类变量，也就是说，让它们的名称首字母变成小写，然后编写公开的用于获取错误值以及用于判等错误值的函数。比如，对于错误值os.ErrClosed，先改写它的名称，让其变成os.errClosed，然后再编写ErrClosed函数和IsErrClosed函数。这不是说让你去改动标准库中已有的代码，这样做的危害会很大，甚至是致命的。我只能说，对于你可控的代码，最好还是要尽量收紧访问权限。



第二个方案，此方案存在于syscall包中。该包中有一个类型叫做Errno，该类型代表了系统调用时可能发生的底层错误。这个错误类型是error接口的实现类型，同时也是对内建类型uintptr的再定义类型。

由于uintptr可以作为常量的类型，所以syscall.Errno自然也可以。syscall包中声明有大量的Errno类型的常量，每个常量都对应一种系统调用错误。syscall包外的代码可以拿到这些代表错误的常量，但却无法改变它们。可以仿照这种声明方式来构建我们自己的错误值列表，这样就可以保证错误值的只读特性了。

