## 库源码文件

库源码文件是不能被直接运行的源码文件，它仅用于存放程序实体，这些程序实体可以被其他代码使用（只要遵从 Go 语言规范的话）。

### 1.关于程序实体

程序实体是变量、常量、函数、结构体和接口的统称。我们总是会先声明（或者说定义）程序实体，然后再去使用。程序实体的名字被统称为标识符。标识符可以是任何 Unicode 编码可以表示的字母字符、数字以及下划线“_”，但是其首字母不能是数字。从规则上说，我们可以用中文作为变量的名字。但是，我觉得这种命名方式非常不好，自己也会在开发团队中明令禁止这种做法。作为一名合格的程序员，我们应该向着编写国际水准的程序无限逼近。

将命令源码文件中的代码拆分到其他库源码文件时，源码文件声明的包名可以与其所在目录的名称不同，只要这些文件声明的包名一致就可以。

### 2.代码包声明的基本规则

- 同目录下的源码文件的代码包声明语句要一致。 它们要同属于一个代码包。这对于所有源码文件都是适用的。

  如果目录中有命令源码文件，那么其他种类的源码文件也应该声明属于 main 包。这也是我们能够成功构建和运行它们的前提。

- 源码文件声明的代码包的名称可以与其所在的目录的名称不同。在针对代码包进行构建时，生成的结果文件的主名称与其父目录的名称一致。

> 源码文件所在的目录相对于 src 目录的相对路径就是它的代码包导入路径，而实际使用其程序实体时给定的限定符要与它声明所属的代码包名称对应,为了不让该代码包的使用者产生困惑，我们总是应该让声明的包名与其父目录的名称一致。

### 3.程序实体的引用规则 & 访问权限规则

- 名称的首字母为大写的程序实体才可以被当前包外的代码引用，否则它就只能被当前包内的其他代码引用。通过名称，Go 语言自然地把程序实体的访问权限划分为了包级私有的和公开的。对于包级私有的程序实体，即使你导入了它所在的代码包也无法引用到它。
- 可以通过创建internal代码包让一些程序实体仅仅能被当前模块中的其他代码引用。这被称为 Go 程序实体的第三种访问权限：模块级私有。具体规则是，internal代码包中声明的公开程序实体仅能被该代码包的直接父包及其子包中的代码引用。当然，引用前需要先导入这个internal包。对于其他代码包，导入该internal包都是非法的，无法通过编译。

```go
// lib/internal/internal.go
package internal

import (
	"fmt"
	"io"
)

func Hello(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}
```

