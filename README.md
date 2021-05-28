# argparse

[![GoDoc](https://godoc.org/github.com/hellflame/argparse?status.svg)](https://godoc.org/github.com/hellflame/argparse) [![Go Report Card](https://goreportcard.com/badge/github.com/hellflame/argparse)](https://goreportcard.com/report/github.com/hellflame/argparse) [![Build Status](https://travis-ci.com/hellflame/argparse.svg?branch=master)](https://travis-ci.com/hellflame/argparse) [![Coverage Status](https://coveralls.io/repos/github/hellflame/argparse/badge.svg?branch=master)](https://coveralls.io/github/hellflame/argparse?branch=master)

Argparser is inspired by [python argparse](https://docs.python.org/3.9/library/argparse.html)

It's small (about 700 rows of code) but fully Functional & Powerful

Provide not just simple parse args, but :

- [x] Sub Command
- [x] Argument Groups
- [x] Positional Arguments
- [x] Customize Parse Formatter
- [x] Customize Validate checker
- [x] Argument Choice support
- [x] Argument Action support (infinite possible)
- [x] Shell Completion support
- [x] Levenshtein Error Correction
- [ ] ......

## Installation

```bash
go get github.com/hellflame/argparse
```

> no dependence needed

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

Checkout output:

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

A few points

About object __parser__ :

1. `NewParser` first argument is the name of your program, but it's ok __not to fill__ , when it's empty string, program name will be `os.Args[0]` , which can be wired when using `go run`, but it will be the executable file's name when you release the code(using `go build`). It can be convinient where the release name is uncertain
2. `help` function is auto injected, but you can disable it when using `NewParser`, with `&ParserConfig{DisableHelp: true}`. then you can use any way to define the `help` function, or whether to `help` 
3. When `help` showed up, the program will default __exit with code 1__ , this is stoppable by setting `ParserConfig.ContinueOnHelp`  to be `true`, you can use your own help function instead

About __parse__ action:

1. The argument `name` is only usable __after__  `parser.Parse` 
2. When passing `parser.Parse` a `nil` as argument, `os.Args[1:]` is used as parse source
3. The short name of your argument can be __more than one character__

---

Based on those points above, the code can be like this:

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

A few points:

1. `DisableHelp` only prevent  `-h/--help` flag to register to parser, but the `help` is still available
2. If keep `DisableDefaultShowHelp` to be false, where there is no argument, the `help` message will still show up as Default Action
3. After the manually call of `parser.PrintHelp()` , `return` will put an end to `main`
4. Notice the order of usage array, it's mostly the order of creating arguments, I tried to keep them this way

### Features

some show case

#### 1. levenshtein error correct ( >= v1.2.0)

the `Parser` will try to match __flag arguments__ when there is no match

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

### Supported Arguments

#### 1. Flag

```go
parser.Flag("short", "full", nil)
```

`Flag` create flag argument, Return a `*bool` point to the parse result

Python version is like `add_argument("-s", "--full", action="store_true")`

Flag Argument can only be used as an OptionalArguments

#### 2. String

```go
parser.String("short", "full", nil)
```

`String` create string argument, return a `*string` point to the parse result

String Argument can be used as Optional or Positional Arguments, default to be Optional, then it's like `add_argument("-s", "--full")` in python

Set `Option.Positional = true` to use as Positional Argument, then it's like `add_argument("s", "full")` in python

#### 3. StringList

```go
parser.Strings("short", "full", nil)
```

`Strings` create string list argument, return a `*[]string` point to the parse result

Mostly like `*Parser.String()`

Python version is like `add_argument("-s", "--full", nargs="*")` 

#### 4. Int

```go
parser.Int("short", "full", nil)
```

`Int` create int argument, return a `*int` point to the parse result

Mostly like `*Parser.String()`, except the return type

Python version is like `add_argument("-s", "--full", type=int)`

#### 5. IntList

```go
parser.Ints("short", "full", nil)
```

`Ints` create int list argument, return a `*[]int` point to the parse result

Mostly like `*Parser.Int()`

Python version is like `add_argument("-s", "--full", type=int, nargs="*")`

#### 6. Float

```go
parser.Float("short", "full", nil)
```

`Float` create float argument, return a `*float64` point to the parse result

Mostly like `*Parser.String()`, except the return type

Python version is like `add_argument("-s", "--full", type=double)` 

#### 7. FloatList

```go
parser.Floats("short", "full", nil)
```

`Floats` create float list argument, return a `*[]float64` point to the parse result

Mostly like `*Parser.Float()`

Python version is like `add_argument("-s", "--full", type=double, nargs="*")` 

### Other Types

For complex type or even customized types are __not directly supported__ , but it doesn't mean you can't do anything before parsing to your own type, here shows some cases:

#### 1. File type

You can check file's existence before read it, and tell if it's a valid file, etc. [eg is here](examples/customzed-types/main.go)

Though the return type is still a `string` , but it's more garanteed to use the argument as what you wanted

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

It used `Validate` to do the magic, we'll talk about it later in more detail

Python code is like:

```python
function valid_type(arg) {
  if !os.path.exist(arg) {
    raise Exception("can't access {}".format(arg))
  }
  return arg
}

parser.add_argument("-s", "--full", type=valid_type)
```

The difference is that, python can return any type from the type function `valid_type` , and you can just return a `File` type in there

There is a little problem if Argument return a `*File` in go. the `*File` might be used somewhere before, which makes it non-idempotent, and you need to `Close` the file somewhere, or the memory may leak. Instead of a `*File` to use with danger, you can manage the resouce much safer:

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

#### 1. ArgumentGroup

Argument group is useful to present argument help infos in group, only affects how the help info displays, using `Group` config to do so, [eg](examples/yt-download/main.go)

```go
parser.Flag("", "version", &argparse.Option{
  Help: "Print program version and exit", 
  Group: "General Options",
})
```

#### 2. DisplayMeta

When the full name of the argument is too long or seems ugly, `Meta` can change how it displays in help, [eg](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Meta: "NUMBER",
})
```

It will looks like this in help message:

```bash
  --playlist-start NUMBER  Playlist video to start at (default is 1)
```

#### 3.DefaultValue

If the argument is not passed from arguments array (like `os.Args`), default value can be passed to continue, [eg](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{
  Help: "Playlist video to start at (default is 1)", 
  Default: "1",
})
```

Noted the Default value is not the type of `Int` , because the value is used like an argument from parse args (like `os.Args`), it's got to get through `Validate` & `Formatter` & `parse` actions (if these actions exist),  `Validate` & `Formatter` will be mentioned below

Also, the Default value can only be one `String` , if you want an Array arguments, you can only have one element Array as default value

#### 4. RequiredArgument

If the argument must be input, set `Required` to be `true`, [eg](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Required: true,
})
```

Flag argument can not be `Required` , you should know the reason, Flag argument has more restrictions, you will be noticed when using it

#### 5. PositionalArgument

If the input argument is the value you want, set `Positional` to be true, [eg](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{
  Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", 
  Positional: true,
})
```

The position of the PositionalArgument is quit flex, with not much restrictions, it's ok to be

1. in the middle of arguments, `--play-list 2 xxxxxxxx --update`, if the argument before it is not an Array argument, won't parse `url` in this case: `--user-ids id1 id2 url --update` 
2. after another single value PositionalArgument, `--mode login username password` , the last `password` will be parsed as second PositionalArgument

So, use it carefully, it cause fusion sometime, which is the same as Python Version of argparse

#### 6. ArgumentValidate

Provide `Validate` function to check each passed-in argument

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

`Validate` function has high priority, executed just after `Default` value is set, which means, the default value has to go through `Validate` check

#### 7. ArgumentFormatter

Format input argument to most basic types you want, the limitation is that, the return type of `Formatter` should be the same as your argument type

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

If `Validate` is set, `Formatter` is right after `Validate`

If raise errors in `Formatter`, it will do some job like `Validate` 

The return type of `interface{}` should be the same as your Argument Type, or Element Type of your Arguments, here,  to be `string` as Example shows

#### 8. ArgumentChoices

Restrict inputs to be within the given choices, using `Choices`

```go
parser.Ints("", "hours", &Option{
  Choices: []interface{}{1, 2, 3, 4},
})
```

If `Formatter` is set, Choice check is right after `Formatter`

When it's single value, the value must be one of the `Choices`

When it's value array, each value must be one of of `Choices`

#### 9. SubCommands

Create new parser scope, within the sub command parser, arguments won't interrupt main parser

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

The two `--flag` will parse seperately, so you can use `tFlag` & `t` to reference flag in `test` parser and `main` parser.

As you can see, though main parser & test parser has different context, but they do parse user input at the same time, it's quite like `Argument Group` , except:

1. sub command has different context, so you can have two `--flag`, and different help message output
2. sub command show help message seperately, it's for user to understand your program step by step. While `Group Argument` helps user to understand your program group by group

#### 10. Argument Action √

Argument Action allows you to do anything with the argument when there is any match, this enables infinite possibility when parsing arguments, [eg](examples/any-type-action/main.go)

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

A few points to be noted:

1. `Action` takes function with `args []string` as input，the `args` has two kind of input
   * `nil` : which means it's a `Flag` argument
   * `[]string{"a1", "a2"}` : which means you have bind other type of argument, other than `Flag` argument
2. Errors can be returned if necessary, it can be normally captured

#### 11. Default Parse Action [ >= v0.4 ]

Instead of showing help message as default, now you can set your own default action when no user input is given, [eg](examples/parse-action/main.go)

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

When `DefaultAction` is set, default show help message will be ignored.

`DefaultAction` is effective on sub-command, and if sub parser's `ParserConfig` is `nil`, `DefaultAction` from main parser will be inherited.

#### 12. Shell Completion Support [ >= v0.4 ]

Set `ParserConfig.AddShellCompletion` to `true` will register `--completion` to the parser, [eg](examples/shell-completion/main.go)

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
3. sub command has no completion entry

Save the output code to `~/.bashrc` or `~/.zshrc` or `~/bash_profile` or some file at `/etc/bash_completion.d/` or `/usr/local/etc/bash_completion.d/` , then restart the shell or `source ~/.bashrc` will enable the completion. 

Completion will register to your shell by your program name, so, it's best to give your program a fix name

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

## [Examples](examples)

there are some useful use cases to help program build there own command

feel free to add different use cases