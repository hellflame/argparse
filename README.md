# argparse

argparser inspired by [python argparse](https://docs.python.org/3.9/library/argparse.html)

provide not just simple parse args, but :

- [x] sub command
- [x] argument groups
- [x] positional arguments
- [x] optimizable parse method
- [x] optimizable validate checker
- [x] argument choice support
- [ ] ......

## Installation

```bash
go get github.com/hellflame/argparse
```

> no dependence needed

## Usage

just go:

```go
package main

import (
    "fmt"
    "github.com/hellflame/argparse"
    "log"
)

func main() {
    parser := argparse.NewParser("basic", "this is a basic program", nil)
    name := parser.String("n", "name", nil)
    if e := parser.Parse(nil); e != nil {
        log.Fatal(e.Error())
    }
    fmt.Printf("hello %s\n", *name)
}
```

[example](examples/basic)

checkout output:

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

a few point:

1. `NewParser` first argument is the name of your program, but it's ok __not to fill__ , when it's empty string, program name will be `os.Args[0]` , which can be wired when using `go run`, but it will be you executable file's name when you release the code. It can be convinient where the release name is uncertain
2. `help` function is auto injected, but you can disable it when `NewParser`, with `&ParserConfig{DisableHelp: true}`. then you can use any way to define the `help` function, or whether to `help` 
3. the argument `name` is only usable __after__  `parser.Parse` , or there might be errors happening
4. when passing `parser.Parse` a `nil` as argument, `os.Args[1:]` is used as parse source, so notice to __only pass arguments__ to `parser.Parse` , or the program name may be `unrecognized`
5. the short name of your agument can be __more than one character__
6. when `help` showed up, the program will default __exit with code 1__ , this is stoppable by setting `ParserConfig.ContinueOnHelp`  to by `true`, or just use your own help function instead

based on these points, the code can be like this:

```go
package main

import (
    "fmt"
    "github.com/hellflame/argparse"
    "os"
)

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

check output:

```bash
=> go run main.go

=> go run main.go -h
unrecognized arguments: -h

# the real help entry is -help / --help-me
=> go run main.go -help
usage: /var/folq1pddT/go-build42601/exe/main [-n NAME] [-help]

this is a basic program

optional arguments:
  -n NAME, --name NAME
  -help, --help-me

# still functional
=> go run main.go --name hellflame
hello hellflame
```

a few points:

1. `DisableHelp` only avoid `-h/--help` flag to register to parser, but the `help` is still fully functional
2. if keep `DisableDefaultShowHelp` to be false, where there is no argument, the `help` function will still show up
3. after the manual call of `parser.PrintHelp()` , program goes on
4. notice the order of usage array, it's mostly the order of your arguments

## Config

## Examples

