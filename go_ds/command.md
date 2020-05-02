## 命令 & 程序入口

命令源码文件是程序运行的入口,是每个可独立运行的程序必须拥有的。我们可以通过build或install生成对应的可执行文件，这个可执行文件一般会与该命令源码文件的父目录同名。



如果一个源码文件声明属于一个 main 包，并且包含一个**无参数且无结果声明**的main函数，它就是命令源码文件。

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, world!")
}
```



> Tips:
>
> 对于一个独立的程序来说，命令源码文件永远只会也只能有一个。



### 命令源码文件接收参数

- #### flag 包

可以在 命名源码文件中加入 init 函数，并在 init 函数中添加 flag.StringVar：

```shell
flag.StringVar(&name, "name", "everyone", "The greeting object.")
```

其中 flag.StringVar 接收的四个参数：

- 第 1 个参数是用于存储该命令参数值的地址，具体到这里就是在前面声明的变量name的地址了，由表达式&name表示。
- 第 2 个参数是为了指定该命令参数的名称，这里是name。
- 第 3 个参数是为了指定在未追加该命令参数时的默认值，这里是everyone。
- 第 4 个函数参数，即是该命令参数的简短说明了，这在打印命令说明时会用到。

然后将 main 函数中加入 flag.Parse 用于真正解析命令参数，并把它们的值赋给相应的变量。

```go
package main

import (
  "flag"
  "fmt"
)

var name string

func init() {
  flag.StringVar(&name, "name", "everyone", "The greeting object.")
}

func main() {
  flag.Parse()
  fmt.Printf("Hello, %s!\n", name)
}

```

其中，flag.StringVar函数类似的函数，叫flag.String。这两个函数的区别是，后者会直接返回一个已经分配好的用于存储命令参数值的地址，如果想使用这个我们可以将代码改为：

```go
package main

import (
  "flag"
  "fmt"
)

var name = flag.String("name", "everyone", "The greeting object.")

func main() {
  flag.Parse()
  fmt.Printf("Hello, %s!\n", name)
}
```



运行该命令源码文件:

```shell
go run main.go -name="Robert"

Hello, Robert!
```



如果想查看该命令源码文件的参数说明:

```shell
Usage of /var/folders/ts/7lg_tl_x2gd_k1lm5g_48c7w0000gn/T/go-build155438482/b001/exe/main:
 -name string
    The greeting object. (default "everyone")
exit status 2
```

其中，`/var/folders/ts/7lg_tl_x2gd_k1lm5g_48c7w0000gn/T/go-build155438482/b001/exe/main` 是 `go run` 命令构建上述命令源码再运行生成的可执行文件。



### 自定义命令源码文件的参数使用说明

这有很多种方式，最简单的一种方式就是对变量 flag.Usage 重新赋值。flag.Usage 的类型是 func()，即一种无参数声明且无结果声明的函数类型。

- flag.Usage

flag.Usage 变量在声明时就已经被赋值了，所以我们才能够在运行命令 go run main.go —help 时看到正确的结果。

> Tips:
>
> 对flag.Usage的赋值必须在调用flag.Parse函数之前。



```go
package main

import (
  "flag"
  "fmt"
)

var name = flag.String("name", "everyone", "The greeting object.")

func main() {
  flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "question")
		flag.PrintDefaults()
	}
  flag.Parse()
  fmt.Printf("Hello, %s!\n", name)
}
```

运行：

```shell
go run main.go --help
Usage of question: 
	-name string The greeting object. (default "everyone")
exit status 2
```

事实上，在调用 StringVar、Parse 时，实际上在调用 flag.CommandLine 变量的对应方法。

- flag.CommandLine

flag.CommandLine相当于默认情况下的命令参数容器。

```go
package main

import (
  "flag"
  "fmt"
)

func init() {
  flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	flag.CommandLine.Usage = func() {
  	fmt.Fprintf(os.Stderr, "Usage of %s:\n", "question")
  	flag.PrintDefaults()
	}
}

func main() {
  flag.Parse()
  fmt.Printf("Hello, %s!\n", name)
}
```

运行：

```shell
go run main.go --help
Usage of question: 
	-name string The greeting object. (default "everyone")
exit status 2
```

还可以通过 将原有的 `flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)` 中的 `flag.ExitOnError` 修改为 `flag.PanicOnError` 重新执行 `go run main.go --help` 会产生另外一种效果。

>  Tips:
>
> 这里的 `flag.ExitOnError` 和  `flag.PanicOnError`  均为 flag 包中的常量。

- flag.ExitOnError：告诉命令参数容器，当命令后跟--help或者参数设置的不正确的时候，在打印命令参数使用说明后以状态码2结束当前程序。

  状态码2代表用户错误地使用了命令，而flag.PanicOnError与之的区别是在最后抛出“运行时恐慌（panic）”。

- flag.PanicOnError与之的区别是在最后抛出“运行时恐慌（panic）”。“运行时恐慌”是 Go 程序错误处理方面的概念。

##### 另一种方式

```go
package main

import (
  "flag"
  "fmt"
)

var name string

var cmdLine = flag.NewFlagSet("question", flag.ExitOnError)


func init() {
	cmdLine.StringVar(&name, "name", "everyone", "The greeting object.")
}

func main() {
  flag.Parse(os.Args[1:])
  fmt.Printf("Hello, %s!\n", name)
}

```

其中的os.Args[1:]指的就是我们给定的那些命令参数。这样做就完全脱离了flag.CommandLine。这样做的好处依然是更灵活地定制命令参数容器。但更重要的是，你的定制完全不会影响到那个全局变量 flag.CommandLine。



### 扩展阅读：

- Golang 编写命令行神器 Corba <https://github.com/spf13/cobra> Kubernetes 都在用！

