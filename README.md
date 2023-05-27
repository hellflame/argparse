# argparse

[![GoDoc](https://godoc.org/github.com/hellflame/argparse?status.svg)](https://godoc.org/github.com/hellflame/argparse) [![Go Report Card](https://goreportcard.com/badge/github.com/hellflame/argparse)](https://goreportcard.com/report/github.com/hellflame/argparse) [![Coverage Status](https://coveralls.io/repos/github/hellflame/argparse/badge.svg?branch=master)](https://coveralls.io/github/hellflame/argparse?branch=master)

[中文文档](docs/cn.md)

Argparser is inspired by [python argparse](https://docs.python.org/3.9/library/argparse.html)

It's small but Powerful

Providing not only simple parsing args, but :

- [x] Sub Command
- [x] Argument Groups
- [x] Positional Arguments
- [x] Customizable Parse Formatter
- [x] Customizable Validate Checker
- [x] Argument Choice Support
- [x] Argument Action Support (infinite possible)
- [x] Shell Completion Support
- [x] Levenshtein Reminder
- [x] Output Color Schema Support
- [ ] ......

## Aim

The aim of the project is to help programers build better command line programs using `golang`.

## Installation

```bash
go get -u github.com/hellflame/argparse
```

> no third-party dependence required

## Usage

Just go:

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

Check output:

```bash
=> go run main.go
usage: basic [-h] [-n NAME]

this is a basic program

options:
  -h, --help            show this help message
  -n NAME, --name NAME

=> go run main.go -n hellflame
hello hellflame
```

A few points:

About the object __parser__ :

1. `NewParser`'s first argument is the name of your program, it's ok __to be empty__. When it's an empty string, the program's name will be `path.Base(os.Args[0])` . It can be handy when the release name is not decided yet.
2. `help` entry is usable by default, you can disable it with `&ParserConfig{DisableHelp: true}` when invoking `NewParser`, then you can define your own `help` .
3. When *help message* showed up, the program will default __exit with code 1 (version < v1.5)__ or __return error with type BreakAfterHelp (version >= 1.5)__ or the __error equals BreakAfterHelpError (version >= 1.10)__ , this is stoppable by setting `ParserConfig.ContinueOnHelp = true` .

About the action __parse__ :

1. The argument `name` is bond to user's input __after__  `parser.Parse` 
2. When give `parser.Parse` a `nil` as the argument, `os.Args[1:]` is used as parse source
3. The short name of your argument can be __more than one character__

---

Based on those points above, the code can be altered like this:

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

Check output:

```bash
=> go run main.go  # there will be no output

=> go run main.go -h
unrecognized arguments: -h
do you mean?: -n

# the real help entry is -help / --help-me
=> go run main.go -help
usage: main [-n NAME] [-help]

this is a basic program

options:
  -n NAME, --name NAME
  -help, --help-me

=> go run main.go --name hellflame
hello hellflame
```

A few more points:

1. `DisableHelp` only prevent  `-h/--help` flag from registering to parser, but the `help` function is still available using `PrintHelp` and `FormatHelp`.
2. When `DisableDefaultShowHelp` is false, and there is no user input, the `help` message will still show up as Default Action.
3. Display help message by invoking `parser.PrintHelp()` , `return` will put an end to the `main` function.
4. Notice the order of usage array, it's mostly the order of creating arguments, I tried to keep them this way. [example](example/change-help)

### Features

Some show cases

#### 1. Levenshtein error correct [ >= v1.2.0 ]

the `Parser` will try to give most posibile __options__ as recommend when there is no match

```go
parser := NewParser("", "", nil)
parser.String("a", "aa", nil)
if e := parser.Parse([]string{"--ax"}); e != nil {
  if e.Error() != "unrecognized arguments: --ax\ndo you mean?: --aa" {
    t.Error("failed to guess input")
    return
  }
}
// when user input '--ax', Parser will try to find best matches with nearrest levenshtein-distance
// here for example is --aa
```

Notice that if there are multiple `Positional Argument` , the `unrecognized arguments` will be regard as `Positional Argument` , and there will be no recommend. 

#### 2. Help info hint [ >= v1.6.0 ]

Help message can be generated with some hint info, like default value, choice range, required mark, or even any hint message. Like:

```bash
usage: sub-command test [--help] [--flag] [--other] [--float FLOAT] [--int INT] [--string STRING]

start a bug report

options:
  --help, -h                  show this help message
  --flag, -f                  from test parser
  --other, -o                 (optional => ∫)
  --float FLOAT               (options: [0.100000, 0.200000], required)
  --int INT, -i INT           this is int (default: 1)
  --string STRING, -s STRING  no hint message
```

Enable global hint by setting parser config `&argparse.ParserConfig{WithHint: true}` .

Disable one argument hint with `&argparse.Option{NoHint: true}`

Customize argument hint with `&argparse.Option{HintInfo: "customize info"}`

[example](examples/sub-command)

### Supported Arguments

#### 1. Flag

```go
parser.Flag(short, full, *Option)
```

`Flag` create a flag argument, return a `*bool` pointer to the parse result

Python version is like `add_argument("-s", "--full", action="store_true")`

Flag Argument can only be used as an __OptionalArguments__, more restrictions see [restrictions](#restriction-of-flags).

#### 2. String

```go
parser.String(short, full, *Option)
```

`String` create a string argument, return a `*string` pointer to the parse result

String Argument can be used as Optional or Positional Argument, default to be Optional. It's like `add_argument("-s", "--full")` in python

Set `Option.Positional = true` to use as Positional Argument, it's like `add_argument("s", "full")` in python

#### 3. String List

```go
parser.Strings(short, full, *Option)
```

`Strings` create a string list argument, return a `*[]string` pointer to the parse result

Options are mostly like `*Parser.String()`

Python version is like `add_argument("-s", "--full", nargs="*")` 

#### 4. Int

```go
parser.Int(short, full, *Option)
```

`Int` create an int argument, return a `*int` pointer to the parse result

Options are mostly like `*Parser.String()`, except the return type

Python version is like `add_argument("-s", "--full", type=int)`

#### 5. Int List

```go
parser.Ints(short, full, *Option)
```

`Ints` create an int list argument, return a `*[]int` pointer to the parse result

Options are mostly like `*Parser.Int()`

Python version is like `add_argument("-s", "--full", type=int, nargs="*")`

#### 6. Float

```go
parser.Float(short, full, *Option)
```

`Float` create a float argument, return a `*float64` pointer to the parse result

Options are mostly like `*Parser.String()`, except the return type

Python version is like `add_argument("-s", "--full", type=double)` 

#### 7. Float List

```go
parser.Floats(short, full, *Option)
```

`Floats` create a float list argument, return a `*[]float64` pointer to the parse result

Options are mostly like `*Parser.Float()`

Python version is like `add_argument("-s", "--full", type=double, nargs="*")` 

### Other Types

For complex types or even customized types, this library do __not directly support__ these feature , but it doesn't mean you can't do anything. Here are some cases:

#### 1. File type

You can check file's existence before reading it, and tells if it's a valid file, check modify time, etc. [example](examples/customzed-types/main.go)

Though the return type is still a `string` , but it's more garanteed to use them

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

The case above used `Validate` to do the trick, we'll talk about it later in more detail

Python code is like:

```python
def valid_type(arg):
  if not os.path.exist(arg):
    raise Exception("can't access {}".format(arg))
  return arg

parser.add_argument("-s", "--full", type=valid_type)
```

The difference is that, python can return any type from the type function `valid_type` , and you can just return a `file` type in there

There could arise some problems if a `*File` is returned in go. the `*File` might be used somewhere before, which makes it non-idempotent, and you need to `Close` the file somewhere, or the memory may leak, the resource management can be a problem. Instead of using `*File` with danger, you can manage the resouce in much safer way:

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

#### 2. Any Type

Checkout `Action` for example, then you can handle any type when parsing arguments !

### Advanced

#### 1. Argument Group

Argument group is useful to arrange arguments help info in to groups, only affects how the help info displays, use `Group` config to do so. [example](examples/yt-download/main.go)

```go
parser.Flag("", "version", &argparse.Option{
  Help: "Print program version and exit", 
  Group: "General Options",
})
```

#### 2. Display Meta

When the full name of the argument is too long to display, `Meta` can change how it displays in help info. More Control can be optimized by __MaxHeaderLength__. [example](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Meta: "NUMBER",
})
```

It will look like this in help message:

```bash
  --playlist-start NUMBER  Playlist video to start at (default is 1)
```

#### 3.Default Value

If the argument is not specified by user, default value will be applied. [example](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Default: "1",
})
```

Note that the Default value is a `String`, the value is used like an argument from `os.Args`, it has to get through `Validate` & `Formatter` & `parse` actions (if exist).

Also, the Default value can only be a `String` , and if you want an Array of arguments, you can only have one element Array as default value. You can apply your default array after *parse*.

#### 4. Required Argument

If the argument must be specified by the user, set `Required` to be `true`. [example](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Required: true,
})
```

Flag argument can not be `Required` (flag argument has more [restrictions](#restriction-of-flags), you will be noticed while choosing it)

#### 5. Positional Argument

If you want users to input arguments by position, set `Positional` to be true. [example](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Positional: true,
})
```

The position of the Positional Argument is quit flex, with not much restrictions, it's ok to be

1. in the middle of other arguments, `--play-list 2 xxxxxxxx --update`. If the argument before it is an Array argument,  `url` will be treated as one of the Array argument: `--user-ids id1 id2 url --update` 
2. after another single value Positional Argument, `--mode login username password` , the last `password` will be regard as another Positional Argument

So, use it carefully, it may __confuse__ (for the users), which is the same in argparse of Python Version.

#### 6. Argument Validation

Provide `Validate` function to check every one of its argument

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

`Validate` function is executed just after when `Default` value is set, which means, the default value has to go through `Validate` check.

> If the argument is an array type (`Ints`, `Strings`, `Floats` ...), every element will be validated by this function.

#### 7. Argument Formatter

Re-format input argument, but remember that, the return type of `Formatter` should be the same as your argument type

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

If `Validate` is set, `Formatter` is executed after `Validate`.

If error is raised in `Formatter`, it acts like `Validate`.

The return type of `interface{}` should be the same as your Argument Type, or Element Type of your Arguments, here,  is a `string` in the Example.

#### 8. Argument Choices

Restrict inputs to be within the given choices, use `Choices`

```go
parser.Ints("", "hours", &Option{
  Choices: []interface{}{1, 2, 3, 4},
})
```

The element type of the choice is the same as argument, or the element of the argument.

If `Formatter` is set, Choice check is right after `Formatter`

When it's single value, the value must be one of the `Choices`

When it's value array, each value must be one of of `Choices`

#### 9. Sub Commands

Create new parser scope, arguments won't interrupt other parsers(root parser and other sub parsers)

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

Output:

```bash
=> ./sub-command
usage: sub-command <cmd> [-h] [-f]

Go is a tool for managing Go source code.

commands:
  test        start a bug report

options:
  -h, --help  show this help message
  -f, --flag  from main parser

# when using sub command, it's a total different context
=> ./sub-command test
usage: sub-command test [-h] [-f] [-o] [-i INT]

start a bug report

options:
  -h, --help         show this help message
  -f, --flag         from test parser
  -o, --other
  -i INT, --int INT
```

The two `--flag` will parse seperately, so you can use `tFlag` & `t` to reference flag in `test` parser and `main` parser.

1. sub command has different context, so you can have two `--flag`, and different help message output
2. sub command show help message seperately, it's for user to understand your program *step by step*. While `Group Argument` helps user to understand your program *group by group*
2. Theoretically, you can create sub parser infinitely. Create a sub parser of a sub parser of a subparser ... Somehow, which would make the program harder to use, and easy to run out of users' patience.

**[v1.7.3]** Fix:

* Invoked Action

if subparser is invoked, the root parser will **NOT** be invoked.

* When Not Invoked

if subparser is **NOT** invoked, no argument's default value will be applied, and required argument will not be checked ... The subparser remains dead. (It's a bug if the subparser gets alive !)

#### 10. Argument Action

Argument Action allows you to do anything with the argument if there is any match, this enables infinite possibility when parsing arguments. [example](examples/any-type-action/main.go)

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

fmt.Println(sum)  // 6
```

A few points:

1. `Action` is a function with `args []string` as input，the `args` has two kinds of input
   * `nil` : which means it's a `Flag` argument
   * `[]string{"a1", "a2"}` : which means you have bond other type of argument, other than `Flag` argument
2. Errors can be returned if necessary, it can be normally captured in *parse*.
3. The return type of the argument is not of much importance, using `p.Strings` is the same as `p.Ints` , because `arg.Action` will be executed __before binding value__, which means, `Action` has __top priority__
4. If `Action` is executed,  `Validate`, `Formatter`, Choice check and value binding will be ignored.

#### 11. Default Parse Action [ >= v0.4 ]

Instead of showing help message as default, you can set your own default action when no user input is given, [example](examples/parse-action/main.go)

```go
parser := argparse.NewParser("basic", "this is a basic program", &argparse.ParserConfig{DefaultAction: func() {
  fmt.Println("hi ~\ntell me what to do?")
}})
parser.AddCommand("test", "testing", &argparse.ParserConfig{DefaultAction: func() {
  fmt.Println("ok, now you know we are testing")
}})
if e := parser.Parse(nil); e != nil {
  fmt.Println(e.Error())
  return
}
```

When `DefaultAction` is set, default action of showing help message will be ignored.

`DefaultAction` is effective on sub command, and if sub parser's `ParserConfig` is `nil`, `DefaultAction` from main parser will be inherited.

#### 12. Shell Completion Support [ >= v0.4 ]

Users can get terminal hint through tapping [tab]

Set `ParserConfig.AddShellCompletion` to `true` will register `--completion` to the parser. [example](examples/shell-completion/main.go)

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

Though, if you didn't set `ParserConfig.AddShellCompletion` to `true` , shell complete script is still available via `parser.FormatCompletionScript` , which will generate the script.

__Note__: 

1. the completion script only support `bash` & `zsh` for now
2. and it only generate simple complete code for basic use, it should be better than nothing.
3. sub command has __no__ completion entry
3. you will know if the user has triggered this input by checking the error returned from `Parse` function, it's a `BreakAfterShellScriptError`.

Save the output code (using `start --completion`) to `~/.bashrc` or `~/.zshrc` or `~/bash_profile` or some file at `/etc/bash_completion.d/` or `/usr/local/etc/bash_completion.d/` , then restart the shell or `source ~/.bashrc` will enable the completion. Or just save completion by appending this line in `~/.bashrc`:

```bash
eval `start --completion`
```

Completion the function is supported by shell, and the shell identify your program by its name, so you `MUST`  give your program a fix name.

#### 13. Hide Entry [ >= 1.3 ]

Sometimes, you want to hide en entry from users, because they should not see or are not necessary to know the entry, but you can still use the entry. Situations like:

1. the entry is to help generate completion candidates (which has mess or not much meaningful output)
2. secret back door that users should not know (you can use `os.Getenv` instead, but `argparse` can do more)

You only need to set `Option{HideEntry: true}` 

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

check ouput:

```bash
usage: basic [--help] [--name NAME]

this is a basic program

options:
  --help, -h               show this help message
  --name NAME, -n NAME
```

Which will have effect in `Shell Completion Script`

[example](examples/hide-help-entry/main.go)

#### 14. Invoked & InvokeAction [ >= 1.4 ]

When there is valid match for main parser or sub parser,  `Parser.Invoked` will be set true. If `Parser.InvokeAction` is set, it will be executed. 

```go
p := NewParser("", "", nil)
a := p.String("a", "", nil)
sub := p.AddCommand("sub", "", nil)
b := sub.String("b", "", nil)
p.InvokeAction = func(invoked bool) {
  // do things when main parser has any match
}
sub.InvokeAction = func(invoked bool) {
  // do things when sub parser has any match
}
subNo2 := p.AddCommand("sub2", "", nil)
subNo2.Int("a", "", nil)
subNo2.InvokeAction = func(invoked bool) {
  // do things when sub2 parser has any match
}

if e := p.Parse(nil); e != nil {
  t.Error(e.Error())
  return
}

// check parser Invoked

fmt.Println(p.Invoked, sub.Invoked, subNo2.Invoked)
```

#### 15. Limit args header length [ >= 1.7 ]

When argument is too long, you can set `ParserConfig.MaxHeaderLength` to a reasonable length.

Before setting `MaxHeaderLength` , the help info may display like (which is default to adjust to the longest argument length):

```bash
usage: long-args [--help] [--short SHORT] [--medium-size MEDIUM-SIZE] [--this-is-a-very-long-args THIS-IS-A-VERY-LONG-ARGS]
options:
  --help, -h                                                                        show this help message
  --short SHORT, -s SHORT                                                           this is a short args
  --medium-size MEDIUM-SIZE, -m MEDIUM-SIZE                                         this is a medium size args
  --this-is-a-very-long-args THIS-IS-A-VERY-LONG-ARGS, -l THIS-IS-A-VERY-LONG-ARGS  this is a very long args

```

After setting `ParserConfig.MaxHeaderLength = 20` (it's recommended to be around 20 ~ 30)，argument's help info will display on new line with 20 space indent, if its header is too long.

```bash
usage: long-args [--help] [--short SHORT] [--medium-size MEDIUM-SIZE] [--this-is-a-very-long-args THIS-IS-A-VERY-LONG-ARGS]
options:
  --help, -h        show this help message
  --short SHORT, -s SHORT
                    this is a short args
  --medium-size MEDIUM-SIZE, -m MEDIUM-SIZE
                    this is a medium size args
  --this-is-a-very-long-args THIS-IS-A-VERY-LONG-ARGS, -l THIS-IS-A-VERY-LONG-ARGS
                    this is a very long args
```

[example](examples/long-args/main.go)

#### 16. Inheriable argument [ >= 1.8 ]

When some argument represent the same thing among root parser and sub parsers, such as `debug`, `verbose` .... and you don't want to write duplicated arguments code for too many times, you have two recommended ways:

before *v1.8*, as all sub parsers are the same type as root parser, use a `for` loop with `Action` would do it.

```go
parser := argparse.NewParser("", "", nil)
sub := parser.AddCommand("sub", "", nil)
sub2 := parser.AddCommand("sub2", "", nil)

url := ""
for _, p := range []*argparse.Parser{parser, sub, sub2} {
    p.String("", "url", &argparse.Option{Action: func(args []string) error {
        url = args[0]
        return nil
    }})
}

if e := parser.Parse(nil); e != nil {
    return
}
print(url)
```

The code above might be less clean, and there comes `Inheritable` , which is an `argument option` . When argument sets __`&argparse.Option{Inheritable: true}`__, sub parsers added __after__ the argument will be able to inherit this argument or override it.

```go
parser := argparse.NewParser("", "", nil)
verbose := parser.Flag("v", "", &argparse.Option{Inheritable: true, Help: "show verbose info"}) // inheritable argument
local := parser.AddCommand("local", "", nil)
service := parser.AddCommand("service", "", nil)
version := service.Int("v", "version", &argparse.Option{Help: "choose version"})
```

As a result,  sub parser `local` will inhert `verbose` as a `Flag`, when pass user input `program local -v`, `*verbose` will be `true`, which means `*verbose` is shared among root parser and inherited parsers. However, as prefix `v` is also registered as `Int` by sub parser `service`, you can't use `-v` to `show verbose`, but `choose version` , it's `overrided`.

> Note `Inheritable` is valid for `Positional`, if their metaNames are the same, `argparse` will consider them the same.

[example](examples/inherit/main.go)

#### 17. Bind given parsers to an argument [ >= 1.9 ]

This is a very flexable way to create argument for multiple parsers. Create one argument using the main parser, and provide a series of parsers by setting the option when you create the argument:

 __`Option{BindParsers: []*Parser{a, b, ...}}`__.

```go
parser := argparse.NewParser("", "", nil)
a := parser.AddCommand("a", "", nil)
b := parser.AddCommand("b", "", nil)
c := parser.AddCommand("c", "", nil)

ab := parser.String("", "ab", &argparse.Option{
    BindParsers: []*argparse.Parser{a, b},
})
bc := parser.String("", "bc", &argparse.Option{
    BindParsers: []*argparse.Parser{b, c},
})
```

As a result, subparser `a` & `b` will bind a `String` argument with entry `--ab` and referred to by var `ab` , subparser `b` & `c` will bind a `String` argument with entry `--bc` and referred to by var `bc`. 

__Note__ that both arguments are detached from main `parser`, because you've set the `BindParsers`. If you still want to bind the argument to the main parser, append it to `BindParsers`, like `BindParsers: []*argparse.Parser{b, c, parser}`.

This is one way to share a single argument among different parsers. Be aware that creating argument like this still need to go through conflict check, if there's already a `--ab` exist in `a` or `b` parser, there will be a panic to notice the programer.

#### 18. Color Schema Support [ >= 1.11 ]

Set `ParserConfig.WithColor = true`, and the help message can be dyed with different colors, if the users' terminal support color.

Also Set `ParserConfig.EnsureColor = true`, and the help message will surely dye with colors. This is for some rare terminals without the environment variable `TERM`, the programer can check it for your own. Normally this is not necessary.

You can also set `ParserConfig.ColorSchema` to dye the help message with your own style. Take `DefaultColor` for reference, most part of the help message can be given a `Color` with *Code* and *Property*. A few knowledge of how color is presented in terminal is required. Normally you can try set Color *Code* within 30 and 49 to represent different text color and background color, and set Color *Property* within 10. Do some combinations, and you'll master them, just try!

[example](./examples/colorful/main.go)

![](./docs/images/colorful.png)

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

##### Restrictions Of Flag Argument {#restriction-of-flags}

1. Can't be Positional
2. Can't has Choices
3. Can't be Required
4. Can't set Formatter
5. Can't set Validate function

## Config

### 1. Parser Config

Relative struct: 

```go
type ParserConfig struct {
  Usage  string // manual usage display
  EpiLog string // message after help

  DisableHelp            bool // disable help entry register [-h/--help]
  ContinueOnHelp         bool // set true to: continue program after default help is printed
  DisableDefaultShowHelp bool // set false to: default show help when there is no args to parse (default action)

  DefaultAction      func() // set default action to replace default help action
  AddShellCompletion bool   // set true to register shell completion entry [--completion]
  WithHint           bool   // argument help message with argument default value hint
  MaxHeaderLength    int    // max argument header length in help menu, help info will start at new line if argument meta info is too long

  WithColor   bool         // enable colorful help message if the terminal has support for color
  EnsureColor bool         // use color code for sure, skip terminal env check
  ColorSchema *ColorSchema // use given color schema to draw help info
}
```

example:

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

options: # no [-h/--help] flag is registerd, which is affected by DisableHelp
  -n NAME, --name NAME
  -help, --help-me 

more detail please visit https://github.com/hellflame/argparse  # <=== EpiLog
```

Except the comment above, `ContinueOnHelp` is only affective on your program process, which gives you possibility to do something when `help` entry is invoked.

### 2. Argument Options

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
  Inheritable bool  // sub parsers after this argument can inherit it
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

## Error & Panic

The principle of returning `error` or just panic is that, __no panic for production use__ 

Cases where `argparse` will panic:

1. failed to add subcommand
2. failed to add argument entry, `Strings`, `Flag`, etc.

Those failures is not allowed, and you will notice when you test your program. The rest errors will be returned in `Parse`, which you should be able to tell users what to do.

## [Examples](examples)

there are some useful use cases to help you build your own command line program

feel free to add different use cases



1. [more type action](examples/any-type-action)
2. [parse action(don't help yet)](examples/parse-action)
3. [shell completion(tabs-able)](examples/shell-completion)
4. [hide help entry](examples/hide-help-entry)
5. [customized types(it's all your call)](examples/customzed-types)
6. [how to add sub command](examples/sub-command)
7. [deal with long args](examples/long-args)
8. [multiple parser in one program](examples/multi-parser)
9. [decide help's position](examples/change-help)
9. [argument groups](examples/argument-groups)
9. [batch create arguments](examples/batch-create-arguments)
9. [argument inherit](examples/inherit)
9. [color support](examples/colorful)

