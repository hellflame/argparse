# argparse

[![GoDoc](https://godoc.org/github.com/hellflame/argparse?status.svg)](https://godoc.org/github.com/hellflame/argparse) [![Go Report Card](https://goreportcard.com/badge/github.com/hellflame/argparse)](https://goreportcard.com/report/github.com/hellflame/argparse) [![Coverage Status](https://coveralls.io/repos/github/hellflame/argparse/badge.svg?branch=master)](https://coveralls.io/github/hellflame/argparse?branch=master)

Argparser 项目是受 [python argparse](https://docs.python.org/3.9/library/argparse.html) 启发所开发的 golang 命令行解析包，麻雀虽小五脏俱全。除了简单的解析命令行外，提供了如下特性：

- [x] 子命令
- [x] 命令行分组
- [x] 位置参数支持
- [x] 自定义解析函数
- [x] 自定义参数检查器
- [x] 参数范围限定支持
- [x] 参数行为支持 (释放无限可能性)
- [x] 命令行自动补全脚本支持
- [x] 根据编辑距离的纠错
- [ ] ......

## 目标

帮助程序员开发更好的命令行程序

## 安装

```bash
go get -u github.com/hellflame/argparse
```

> 无额外依赖

## 使用

栗子：

```go
package main

import (
    "fmt"
    "github.com/hellflame/argparse"
)

func main() {
    parser := argparse.NewParser("basic", "this is a basic program", nil)
  
    name := parser.String("n", "name", nil)
  
    if e := parser.Parse(nil); e != nil {
        fmt.Println(e.Error())
      	return
    }
    fmt.Printf("hello %s\n", *name)
}
```

[example](examples/basic)

检查输出:

```bash
=> go run main.go
usage: basic [-h] [-n NAME]

this is a basic program

optional arguments:
  -h, --help            show this help message
  -n NAME, --name NAME

=> go run main.go -n hellflame
hello hellflame
```

几点说明

关于 __parser__ 结构体 :

1. `NewParser` 第一个参数是程序的名字, 但__可以为空__ , 如果程序名为空, 则会使用 `os.Args[0]` 作为程序名 , 只是用 `go run` 时会有点怪怪的，但当发布为可执行文件，在运行时就会是可执行文件的名字(通过 `go build` 生成)，所以在发布名称不确定的时候会比较方便
2. `help` 方法会自动注入, 但也可以在 `NewParser` 时设定 `&ParserConfig{DisableHelp: true}` 来取消这个帮助入口，然后就可以用自己的帮助函数，甚至不给出帮助函数
3. 帮助信息显示后，程序会以状态码 1 退出程序。可以设置 `ParserConfig.ContinueOnHelp`  为 `true`, 阻止这种退出

关于 __parse__ 动作执行:

1. 参数 `name` 只在 `parser.Parse` 后可用
2. 若 `parser.Parse` 接受 `nil` 作为参数, `os.Args[1:]` 会作为解析来源
3. 参数缩写可以 __不止一个字符__

---

基于以上几点，可以写出下面的代码

```go
func main() {
    parser := argparse.NewParser("", "this is a basic program", &argparse.ParserConfig{
        DisableHelp:true,
        DisableDefaultShowHelp: true})
  
    name := parser.String("n", "name", nil)
    help := parser.Flag("help", "help-me", nil)
  
    if e := parser.Parse(os.Args[1:]); e != nil {
        fmt.Println(e.Error())
        return
    }
    if *help {
        parser.PrintHelp()
        return
    }
    if *name != "" {
        fmt.Printf("hello %s\n", *name)
    }
}
```

[example](examples/basic-bias)

检查输出:

```bash
=> go run main.go  # there will be no output

=> go run main.go -h
unrecognized arguments: -h
do you mean: -n

# the real help entry is -help / --help-me
=> go run main.go -help
usage: /var/folq1pddT/go-build42601/exe/main [-n NAME] [-help]

this is a basic program

optional arguments:
  -n NAME, --name NAME
  -help, --help-me

=> go run main.go --name hellflame
hello hellflame
```

几点说明:

1. `DisableHelp` 只是阻止了  `-h/--help` 注册到解析器, 但依然可以通过其他方式得到帮助信息
2. 如果保持 `DisableDefaultShowHelp` 为 `false`, 当没有用户输入时, 帮助信息会作为默认行为输出
3. 在手动执行 `parser.PrintHelp()` 后, `return` 会结束 `main` 方法
4. 注意使用信息输出的顺序, 基本和这些参数的创建顺序一致, 这是有意为之的

### 特点

用例展示

#### 1. 编辑距离纠错 [ >= v1.2.0 ]

`Parser` 会在没有匹配的情况下尽力匹配可选参数

```go
parser := NewParser("", "", nil)
parser.String("a", "aa", nil)
if e := parser.Parse([]string{"--ax"}); e != nil {
  if e.Error() != "unrecognized arguments: --ax\ndo you mean: --aa" {
    t.Error("failed to guess input")
    return
  }
}
// when user input '--ax', Parser will try to find best matches with smallest levenshtein-distance
// here for eg is: --aa
```

注意如果包含 `位置参数` 时 , 未知参数可能会被视为位置参数，就没有任何纠错提示了

### 支持的参数

#### 1. Flag

```go
parser.Flag("short", "full", nil)
```

`Flag` 会创建标记参数, 返回 `*bool` 指针保存结果

Python代码可能像这样： `add_argument("-s", "--full", action="store_true")`

标记参数只能是可选参数

#### 2. String

```go
parser.String("short", "full", nil)
```

`String` 会创建字符串参数, 返回 `*string` 指针保存结果

字符串参数可以作为可选或位置参数(默认为可选参数), Python代码如 `add_argument("-s", "--full")` 

设置 `Option.Positional = true` 变为位置参数, Python代码如 `add_argument("s", "full")`

#### 3. StringList

```go
parser.Strings("short", "full", nil)
```

`Strings` 创建字符串数组参数, 返回 `*[]string` 指针保存结果

和 `*Parser.String()` 差不多

Python代码如 `add_argument("-s", "--full", nargs="*")` 

#### 4. Int

```go
parser.Int("short", "full", nil)
```

`Int` 创建整数参数, 返回 `*int` 指针保存结果

除了返回类型外，和 `*Parser.String()` 差不多

Python代码如 `add_argument("-s", "--full", type=int)`

#### 5. IntList

```go
parser.Ints("short", "full", nil)
```

`Ints` 创建整数数组参数, 返回 `*[]int` 指针保存结果

和 `*Parser.Int()` 差不多

Python代码如 `add_argument("-s", "--full", type=int, nargs="*")`

#### 6. Float

```go
parser.Float("short", "full", nil)
```

`Float` 创建浮点数参数, 返回 `*float64` 指针保存结果

除了返回类型外，和 `*Parser.String()` 差不多

Python代码如 `add_argument("-s", "--full", type=double)` 

#### 7. FloatList

```go
parser.Floats("short", "full", nil)
```

`Floats` 创建浮点数数组参数, 返回 `*[]float64` 指针保存结果

和 `*Parser.Float()` 差不多

Python代码如 `add_argument("-s", "--full", type=double, nargs="*")` 

### 其他类型

这个项目对复杂类型甚至自定义类型是__没有直接支持__的，但这并不妨碍你在解析自己的类型前做点什么

#### 1. 文件类型

你可以在读文件前检查文件是否存在. [eg is here](examples/customzed-types/main.go)

虽然返回类型依然是字符串，但在用这个参数的时候会更有保障

```go
path := parser.String("f", "file", &argparse.Option{
  Validate: func(arg string) error {
    if _, e := os.Stat(arg); e != nil {
      return fmt.Errorf("unable to access '%s'", arg)
    }
    return nil
  },
})

if e := parser.Parse(nil); e != nil {
  fmt.Println(e.Error())
  return
}

if *path != "" {
  if read, e := ioutil.ReadFile(*path); e == nil {
    fmt.Println(string(read))
  }
}
```

主要是 `Validate` 在起作用, 之后会详细讨论它

Python代码如:

```python
function valid_type(arg) {
  if !os.path.exist(arg) {
    raise Exception("can't access {}".format(arg))
  }
  return arg
}

parser.add_argument("-s", "--full", type=valid_type)
```

和golang版本不一样的是，python可以用过`valid_type`返回任意类型, 比如文件类型

在go中返回`*File` 类型也会有一些问题. 文件句柄可能已经在前面使用过，导致它多次使用的结果不一致，并且你也需要管理文件句柄，比如关闭它，以免发生内存泄漏。所以除了使用危险的 `*File`, 你可以更安全的使用这个资源

```go
func dealFile(path) {
  f, e := os.Open("")
  if e != nil {
    fmt.Println(e.Error())
    return
  }
  defer f.Close()  // close file
  io.ReadAll(f)
}
```

#### 2. 任意类型

看看 `Action` 的栗子, 然后你就可以解析任意类型了

### 高级用法

#### 1. 命令行组

命令行组适合在帮助信息里讲参数分组。这只会影响帮助信息的显示，没别的作用。使用配置里的 `Group` 来实现[eg](examples/yt-download/main.go)

```go
parser.Flag("", "version", &argparse.Option{
  Help: "Print program version and exit", 
  Group: "General Options",
})
```

#### 2. 元信息

当参数的完整名称太长或者丑丑的，修改元信息可以改变帮助信息里的展示内容[eg](../examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Meta: "NUMBER",
})
```

帮助信息里看起来像这样：

```bash
  --playlist-start NUMBER  Playlist video to start at (default is 1)
```

#### 3.默认值

如果参数没有给，那么就会使用默认参数，[eg](../examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Default: "1",
})
```

注意默认参数的类型不是想要的 `Int` 类型，因为这个值是作为输入的命令行参数来使用的，它还必须通过 `Validate` & `Formatter` & `parse` 这些方法的处理,  `Validate` & `Formatter` 会在之后提到

并且默认值的类型只能是 `String`，如果想要给数组类型的参数赋默认值，只能得到只有一个参数的数组

#### 4. 必选参数

如果参数必须给出, 设 `Required` 为 `true` 即可, [eg](../examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Required: true,
})
```

标志类参数不能为 `Required` , 你应该知道为什么。当然，标志参数会有更多限制，在使用的过程中会发现的

#### 5. 位置参数

如果用户输入即为想要获取的参数, 设置 `Positional` 为 true 即可, [eg](../examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Positional: true,
})
```

位置参数的位置限制很少，以下情况皆可:

1. 在各种参数中间, `--play-list 2 xxxxxxxx --update`, 如果这个参数的前面是数组类型的参数，那么后面可选参数前的参数都会认为是该位置参数的值，如这里的 `url`: `--user-ids id1 id2 url --update` ，会被当作 `user-ids` 的参数之一
2. 在另一个单个位置参数之后, `--mode login username password` , 最后一个 `password` 会作为第二个位置参数的值

所以请小心使用，有时候会比较容易搞混，和python版本的命令行解析一样

#### 6. 参数检查

提供 `Validate` 方法来检查每一个输入参数

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links", 
  Validate: func(arg string) error {
    if !strings.HasPrefix(arg, "https://") {
      return fmt.Errorf("url should be start with 'https://'")
    }
    return nil
  },
})
```

`Validate` 有比较高的优先级, 会在默认值设置之后即执行, 这意味着默认值比如要通过 `Validate` 的检查

#### 7. 参数格式化

将你想要的参数进行格式化, 限制就是 `Formatter` 返回类型需要和参数类型保持一致

```go
parser.String("", "b", &Option{
  Formatter: func(arg string) (i interface{}, err error) {
    if arg == "False" {
      err = fmt.Errorf("no False")
      return
    }
    i = fmt.Sprintf("=> %s", arg)
    return
  },
})
```

如果给了 `Validate`, `Formatter` 将在 `Validate` 之后执行

如果 `Formatter` 返回错误类型, 就会表现的和 `Validate` 一样

返回类型 `interface{}` 应该和参数类型一致, 或与数组元素类型一致, 栗子里返回的是 `string` 类型

#### 8. 参数可选范围

限制输入参数是所给范围，设置 `Choices` 即可

```go
parser.Ints("", "hours", &Option{
  Choices: []interface{}{1, 2, 3, 4},
})
```

如果给了 `Formatter` , 可选范围在 `Formatter` 之后检查

如果参数仅接受单个值, 那么这个值必须在 `Choices` 范围内

如果参数接受数组, 那么每一个数组元素必须在 `Choices` 范围内

#### 9. 子命令

创建新的命令行解析域, 子命令的参数解析不会影响主命令解析

```go
func main() {
  parser := argparse.NewParser("sub-command", "Go is a tool for managing Go source code.", nil)
  t := parser.Flag("f", "flag", &argparse.Option{Help: "from main parser"})
  testCommand := parser.AddCommand("test", "start a bug report", nil)
  tFlag := testCommand.Flag("f", "flag", &argparse.Option{Help: "from test parser"})
  otherFlag := testCommand.Flag("o", "other", nil)
  defaultInt := testCommand.Int("i", "int", &argparse.Option{Default: "1"})
  if e := parser.Parse(nil); e != nil {
    fmt.Println(e.Error())
    return
  }
  println(*tFlag, *otherFlag, *t, *defaultInt)
}
```

输出:

```bash
=> ./sub-command
usage: sub-command <cmd> [-h] [-f]

Go is a tool for managing Go source code.

available commands:
  test        start a bug report

optional arguments:
  -h, --help  show this help message
  -f, --flag  from main parser

# when using sub command, it's a total different context
=> ./sub-command test
usage: sub-command test [-h] [-f] [-o] [-i INT]

start a bug report

optional arguments:
  -h, --help         show this help message
  -f, --flag         from test parser
  -o, --other
  -i INT, --int INT
```

两个 `--flag` 会分开解析, 所以 `tFlag` & `t` 分别指向 `test` 和 `main` 中的两个标志参数

1. 子命令拥有不同的解析域, 所以你可以有两个 `--flag`, 以及不同的帮助输出
2. 子命令也会单独显示帮助信息, 可以让用户分布理解你的命令.  `Group Argument` 则会分组让用户理解你的命令

#### 10. 参数行为 √

参数行为在当出现匹配时允许你做任何操作, 这将开启无限的可能性, [eg](../examples/any-type-action/main.go)

```go
p := NewParser("action", "test action", nil)

sum := 0
p.Strings("", "number", &Option{Positional: true, Action: func(args []string) error {
  // here tries to sum every input number
  for _, a := range args {
    if i, e := strconv.Atoi(a); e != nil {
      return fmt.Errorf("I don't know this number: %s", a)
    } else {
      sum += i
    }
  }
  return nil
}})

if e := p.Parse([]string{"1", "2", "3"}); e != nil {
  fmt.Println(e.Error())
  return
}

fmt.Println(sum)  // this is a 6 if everything goes on fine
```

有几点需要提一提:

1. `Action` 接受输入参数 `args []string` ， `args` 有两种可能
   * `nil` : 意味着这是一个标志参数
   * `[]string{"a1", "a2"}` : 意味着这是一个数组类型参数
2. 可以返回错误类型, 并且会被正常的捕捉到
3. 返回值的类型不重要, 使用 `p.Strings` 和 `p.Ints` 是一样的, 因为 `arg.Action` 会在 __绑定参数__ 前执行, 这意味着 `Action` 拥有 __最高的执行权限__

#### 11. 默认解析行为 [ >= v0.4 ]

如果不想默认显示帮助信息, 现在如果用户没有任何输入，你可以设置自己的默认行为, [eg](../examples/parse-action/main.go)

```go
parser := argparse.NewParser("basic", "this is a basic program", &argparse.ParserConfig{DefaultAction: func() {
  fmt.Println("hi ~\ntell me what to do?")
}})
parser.AddCommand("test", "testing", &argparse.ParserConfig{DefaultAction: func() {
  fmt.Println("ok, now you know you are testing")
}})
if e := parser.Parse(nil); e != nil {
  fmt.Println(e.Error())
  return
}
```

如果设置了 `DefaultAction`, 默认显示帮助信息会被忽略

`DefaultAction` 对子命令同样起作用, 并且如果子命令的 `ParserConfig` 为 `nil`, `DefaultAction` 会被继承

#### 12. 命令行补全支持 [ >= v0.4 ]

设置 `ParserConfig.AddShellCompletion` 为 `true` 将注册 `--completion` 参数, [eg](../examples/shell-completion/main.go)

```go
p := argparse.NewParser("start", "this is test", &argparse.ParserConfig{AddShellCompletion: true})
p.Strings("a", "aa", nil)
p.Int("", "bb", nil)
p.Float("c", "cc", &argparse.Option{Positional: true})
test := p.AddCommand("test", "", nil)
test.String("a", "aa", nil)
test.Int("", "bb", nil)
install := p.AddCommand("install", "", nil)
install.Strings("i", "in", nil)
if e := p.Parse(nil); e != nil {
  fmt.Println(e.Error())
  return
}
```

即使没有设置 `ParserConfig.AddShellCompletion` 为 `true` , 命令行补全脚本依然可以通过 `parser.FormatCompletionScript` 获取

__注意__: 

1. 命令行补全现在仅支持 `bash` & `zsh` 
2. 它只会生成简单的补全模式，总比没有好
3. 子命令不会注册该方法

保存输出脚本到 `~/.bashrc` or `~/.zshrc` or `~/bash_profile` or some file at `/etc/bash_completion.d/` or `/usr/local/etc/bash_completion.d/` , 然后重启脚本环境 或 `source ~/.bashrc` 会使脚本生效 

命令补全会将命令的名字作为注册入口注册到脚本环境，所以你最好给你的程序一个固定的名字

#### 13. 隐藏入口 [ >= 1.3.0 ]

有时候, 你可能想要对用户隐藏一些入口, 因为用户不应该知道这些入口或不需要知道，但你依然需要使用这些入口

比如:

1. 这是一个用来动态生成补全脚本候选的入口 (输出可能会很乱，没有意义)
2. 秘密的后门 (当然可以用 `os.Getenv` , 只是 `argparse` 可以帮你进行转换)

仅需要设置 `Option{HideEntry: true}` 即可

```go
func main() {
  parser := argparse.NewParser("basic", "this is a basic program", nil)
  name := parser.String("n", "name", nil)
  greet := parser.String("g", "greet", &argparse.Option{HideEntry: true})
  if e := parser.Parse(nil); e != nil {
    fmt.Println(e.Error())
    return
  }
  greetWord := "hello"
  if *greet != "" {
    greetWord = *greet
  }
  fmt.Printf("%s %s\n", greetWord, *name)
}
```

检查输出:

```bash
usage: basic [--help] [--name NAME]

this is a basic program

optional arguments:
  --help, -h               show this help message
  --name NAME, -n NAME
```

对 `Shell Completion Script` 同样起作用

[full eg](../examples/hide-help-entry/main.go)

#### 14. 匹配状态与匹配动作

当主解析或子命令解析被匹配到时，`Parser.Invoked` 会设置为 true 并且 `Parser.InvokeAction` 也会被执行

```go
p := NewParser("", "", nil)
a := p.String("a", "", nil)
sub := p.AddCommand("sub", "", nil)
b := sub.String("b", "", nil)
p.InvokeAction = func() {
  // do things when main parser has any match
}
sub.InvokeAction = func() {
  // do things when sub parser has any match
}
subNo2 := p.AddCommand("sub2", "", nil)
subNo2.Int("a", "", nil)
subNo2.InvokeAction = func() {
  // do things when sub2 parser has any match
}

if e := p.Parse(nil); e != nil {
  t.Error(e.Error())
  return
}

// check parser Invoked

fmt.Println(p.Invoked, sub.Invoked, subNo2.Invoked)
```

##### Argument Process Flow Map

```

                    ┌────►  BindAction？
                    │            │
                    │            │  Consume it's Arguments
                    │            ▼
                    │    ┌──────┐
  --date 20210102 --list │ arg1 │ arg2  arg3
                         └───┬──┘
                             │
                             │
                             ▼
                        ApplyDefault?
                             │
                             │
                       ┌─────▼──────┐
                       │  Validate  │
                       └─────┬──────┘
                             │
                       ┌─────▼──────┐
                       │  Formatter │
                       └─────┬──────┘
                             │
                       ┌─────▼───────┐
                       │ ChoiceCheck │
                       └─────────────┘
```

The return value of last process will be the input of next process, if shows it in code, it's like

```python
with MatchFound:
  if MatchFound.BindAction:
    return MatchFound.BindAction(*args)
  else:
    for arg in args:
      if Validate(arg):
        yield ChoiceCheck(Formatter(arg))
```

## Config

### 1. ParserConfig

Relative struct: 

```go
type ParserConfig struct {
  Usage                  string // manual usage display
  EpiLog                 string // message after help
  DisableHelp            bool   // disable help entry register [-h/--help]
  ContinueOnHelp         bool   // set true to: continue program after default help is printed
  DisableDefaultShowHelp bool   // set false to: default show help when there is no args to parse (default action)
  DefaultAction          func() // set default action to replace default help action
  AddShellCompletion     bool   // set true to register shell completion entry [--completion]
}
```

eg:

```go
func main() {
  parser := argparse.NewParser("basic", "this is a basic program",
		&argparse.ParserConfig{
      Usage:                  "basic xxx",
      EpiLog:                 "more detail please visit https://github.com/hellflame/argparse",
      DisableHelp:            true,
      ContinueOnHelp:         true,
      DisableDefaultShowHelp: true,
    })

  name := parser.String("n", "name", nil)
  help := parser.Flag("help", "help-me", nil)

  if e := parser.Parse(nil); e != nil {
    fmt.Println(e.Error())
    return
  }
  if *help {
    parser.PrintHelp()
    return
  }
  if *name != "" {
    fmt.Printf("hello %s\n", *name)
  }
}
```

[example](examples/parser-config)

Output:

```bash
=> go run main.go
# there will be no help message
# affected by DisableDefaultShowHelp

=> go run main.go --help-me
usage: basic xxx   # <=== Usage

this is a basic program

optional arguments: # no [-h/--help] flag is registerd, which is affected by DisableHelp
  -n NAME, --name NAME
  -help, --help-me 

more detail please visit https://github.com/hellflame/argparse  # <=== EpiLog
```

Except the comment above, `ContinueOnHelp` is only affective on your program process, which give you possibility to do something when default `help` is shown

### 2. ArgumentOptions

Related struct:

```go
type Option struct {
  Meta       string // meta value for help/usage generate
  multi      bool   // take more than one argument
  Default    string // default argument value if not given
  isFlag     bool   // use as flag
  Required   bool   // require to be set
  Positional bool   // is positional argument
  HideEntry  bool   // hide usage & help display
  Help       string // help message
  Group      string // argument group info, default to be no group
  Action     func(args []string) error // bind actions when the match is found, 'args' can be nil to be a flag
  Choices    []interface{}  // input argument must be one/some of the choice
  Validate   func(arg string) error  // customize function to check argument validation
  Formatter  func(arg string) (interface{}, error) // format input arguments by the given method
}
```

## How it works

```
  ┌──────────────────────┐ ┌──────────────────────┐
  │                      │ │                      │
  │     OptionArgsMap    │ │  PositionalArgsList  │
  │                      │ │                      │
  │      -h ───► helpArg │ │                      │
  │                      │ │[  posArg1  posArg2  ]│
  │      -n ──┐          │ │                      │
  │           │► nameArg │ │                      │
  │  --name ──┘          │ │                      │
  │                      │ │                      │
  └──────────────────────┘ └──────────────────────┘
             ▲ yes                  no ▲
             │                         │
             │ match?──────────────────┘
             │
             │
           ┌─┴──┐                   match helpArg:
    args:  │ -h │-n  hellflame           ▼
           └────┘                  ┌──isflag?───┐
                                   ▼            ▼
                                  done ┌──MultiValue?───┐
                                       ▼                ▼
                                   ┌──parse     consume untill
                                   ▼            NextOptionArgs
                                 done
```

## 错误 & 宕机

原则是, __生产环境不产生宕机__ 

以下场景 `argparse` 将宕机:

1. 添加子命令失败
2. 添加参数失败

这些失误是不允许的, 你也会在开发过程中发现这些问题. 剩下的被 `Parse` 返回的错误将指导你如何提示用户正确的输入

## [栗子](../examples)

这里有一些有用的栗子来帮助你搭建自己的命令行，可以帮忙添加一些特别的栗子