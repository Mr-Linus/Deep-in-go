## 字符编码

在 Go 语言中，一个string类型的值既可以被拆分为一个包含多个字符的序列，也可以被拆分为一个包含多个字节的序列。

前者可以由一个以rune为元素类型的切片来表示，而后者则可以由一个以byte为元素类型的切片代表。

rune是 Go 语言特有的一个基本数据类型，它的一个值就代表一个字符，即：一个 Unicode 字符。

比如，'G'、'o'、'爱'、'好'、'者'代表的就都是一个 Unicode 字符。我们已经知道，UTF-8 编码方案会把一个 Unicode 字符编码为一个长度在[1, 4]范围内的字节序列。所以，一个rune类型的值也可以由一个或多个字节来代表。

```go
type rune = int32
```

根据rune类型的声明可知，它实际上就是int32类型的一个别名类型。也就是说，一个rune类型的值会由四个字节宽度的空间来存储。它的存储空间总是能够存下一个 UTF-8 编码值。

根据rune类型的声明可知，它实际上就是int32类型的一个别名类型。也就是说，一个rune类型的值会由四个字节宽度的空间来存储。它的存储空间总是能够存下一个 UTF-8 编码值。一个rune类型的值在底层其实就是一个 UTF-8 编码值。前者是（便于我们人类理解的）外部展现，后者是（便于计算机系统理解的）内在表达。

```go
str := "Go爱好者"
fmt.Printf("The string: %q\n", str)
fmt.Printf("  => runes(char): %q\n", []rune(str))
fmt.Printf("  => runes(hex): %x\n", []rune(str))
fmt.Printf("  => bytes(hex): [% x]\n", []byte(str))
```

字符串值"Go爱好者"如果被转换为[]rune类型的值的话，其中的每一个字符（不论是英文字符还是中文字符）就都会独立成为一个rune类型的元素值。因此，这段代码打印出的第二行内容就会如下所示：

```go
 => runes(char): ['G' 'o' '爱' '好' '者']
```

又由于，每个rune类型的值在底层都是由一个 UTF-8 编码值来表达的，所以我们可以换一种方式来展现这个字符序列：

```go
  => runes(hex): [47 6f 7231 597d 8005]
```

可以看到，五个十六进制数与五个字符相对应。很明显，前两个十六进制数47和6f代表的整数都比较小，它们分别表示字符'G'和'o'。

因为它们都是英文字符，所以对应的 UTF-8 编码值用一个字节表达就足够了。一个字节的编码值被转换为整数之后，不会大到哪里去。

而后三个十六进制数7231、597d和8005都相对较大，它们分别表示中文字符'爱'、'好'和'者'。

这些中文字符对应的 UTF-8 编码值，都需要使用三个字节来表达。所以，这三个数就是把对应的三个字节的编码值，转换为整数后得到的结果。

我们还可以进一步地拆分，把每个字符的 UTF-8 编码值都拆成相应的字节序列。上述代码中的第五行就是这么做的。它会得到如下的输出：

```go
  => bytes(hex): [47 6f e7 88 b1 e5 a5 bd e8 80 85]
```

这里得到的字节切片比前面的字符切片明显长了很多。这正是因为一个中文字符的 UTF-8 编码值需要用三个字节来表达。这个字节切片的前两个元素值与字符切片的前两个元素值是一致的，而在这之后，前者的每三个元素值才对应字符切片中的一个元素值。

注意，对于一个多字节的 UTF-8 编码值来说，我们可以把它当做一个整体转换为单一的整数，也可以先把它拆成字节序列，再把每个字节分别转换为一个整数，从而得到多个整数。这两种表示法展现出来的内容往往会很不一样。比如，对于中文字符'爱'来说，它的 UTF-8 编码值可以展现为单一的整数7231，也可以展现为三个整数，即：e7、88和b1。

![0d8dac40ccb2972dbceef33d03741085](https://static001.geekbang.org/resource/image/0d/85/0d8dac40ccb2972dbceef33d03741085.png)

一个string类型的值会由若干个 Unicode 字符组成，每个 Unicode 字符都可以由一个rune类型的值来承载。这些字符在底层都会被转换为 UTF-8 编码值，而这些 UTF-8 编码值又会以字节序列的形式表达和存储。因此，一个string类型的值在底层就是一个能够表达若干个 UTF-8 编码值的字节序列。

### 遍历字符串

带有range子句的for语句会先把被遍历的字符串值拆成一个字节序列，然后再试图找出这个字节序列中包含的每一个 UTF-8 编码值，或者说每一个 Unicode 字符。这样的for语句可以为两个迭代变量赋值。

如果存在两个迭代变量，那么赋给第一个变量的值，就将会是当前字节序列中的某个 UTF-8 编码值的第一个字节所对应的那个索引值。而赋给第二个变量的值，则是这个 UTF-8 编码值代表的那个 Unicode 字符，其类型会是rune。例如，有这么几行代码：

```go
str := "Go爱好者"
for i, c := range str {
 fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
}
```

这里被遍历的字符串值是"Go爱好者"。在每次迭代的时候，这段代码都会打印出两个迭代变量的值，以及第二个值的字节序列形式。完整的打印内容如下：

```go
0: 'G' [47]
1: 'o' [6f]
2: '爱' [e7 88 b1]
5: '好' [e5 a5 bd]
8: '者' [e8 80 85]
```

第一行内容中的关键信息有0、'G'和[47]。这是由于这个字符串值中的第一个 Unicode 字符是'G'。该字符是一个单字节字符，并且由相应的字节序列中的第一个字节表达。这个字节的十六进制表示为47。

第二行展示的内容与之类似，即：第二个 Unicode 字符是'o'，由字节序列中的第二个字节表达，其十六进制表示为6f。再往下看，第三行展示的是'爱'，也是第三个 Unicode 字符。

因为它是一个中文字符，所以由字节序列中的第三、四、五个字节共同表达，其十六进制表示也不再是单一的整数，而是e7、88和b1组成的序列。



正是因为'爱'是由三个字节共同表达的，所以第四个 Unicode 字符'好'对应的索引值并不是3，而是2加3后得到的5。

