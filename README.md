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

## install

```bash
go get github.com/hellflame/argparse
```

> no dependence needed

## usage

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

checkout output:

```bash
> go run main.go
usage: basic [-h] [-n NAME]

this is a basic program

optional arguments:
  -h, --help            show this help message
  -n NAME, --name NAME

```

a few point:

1. `NewParser` first argument is the name of your program, but it's ok __not to fill__ , when it's empty string, program name will be `os.Args[0]` , which can be wired when using `go run`, but it will be you executable file's name when you release the code. It can be convinient where the release name is uncertain
2. `help` function is auto injected, but you can disable it when `NewParser`, with `&ParserConfig{DisableHelp: true}`. then you can use any way to define the `help` function, or whether to `help` 

## config

## examples

