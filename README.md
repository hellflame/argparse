# argparse

[![GoDoc](https://godoc.org/github.com/hellflame/argparse?status.svg)](https://godoc.org/github.com/hellflame/argparse) [![Go Report Card](https://goreportcard.com/badge/github.com/hellflame/argparse)](https://goreportcard.com/report/github.com/hellflame/argparse) [![Build Status](https://travis-ci.com/hellflame/argparse.svg?branch=master)](https://travis-ci.com/hellflame/argparse) [![Coverage Status](https://coveralls.io/repos/github/hellflame/argparse/badge.svg?branch=master)](https://coveralls.io/github/hellflame/argparse?branch=master)

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

a few points

about object __parser__ :

1. `NewParser` first argument is the name of your program, but it's ok __not to fill__ , when it's empty string, program name will be `os.Args[0]` , which can be wired when using `go run`, but it will be you executable file's name when you release the code. It can be convinient where the release name is uncertain
2. `help` function is auto injected, but you can disable it when `NewParser`, with `&ParserConfig{DisableHelp: true}`. then you can use any way to define the `help` function, or whether to `help` 
3. when `help` showed up, the program will default __exit with code 1__ , this is stoppable by setting `ParserConfig.ContinueOnHelp`  to by `true`, or just use your own help function instead

about __parse__ action:

1. the argument `name` is only usable __after__  `parser.Parse` , or there might be errors happening
2. when passing `parser.Parse` a `nil` as argument, `os.Args[1:]` is used as parse source, so notice to __only pass arguments__ to `parser.Parse` , or the program name may be `unrecognized`
3. the short name of your agument can be __more than one character__

---

based on these points, the code can be like this:

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
4. notice the order of usage array, it's mostly the order of creating arguments

### Advanced

#### 1. ArgumentGroup

argument group is useful to present argument help infos in group, only affects how the help info displays, using `Group` config to do so, [eg](examples/yt-download/main.go)

```go
parser.Flag("", "version", &argparse.Option{Help: "Print program version and exit", Group: "GeneralOptions"})
```

#### 2. DisplayMeta

when the full name of the argument is too long or seems ugly, `Meta` can change how it displays in help, [eg](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{Help: "Playlist video to start at (default is 1)", Meta: "NUMBER"})
```

looks like:

```bash
  --playlist-start NUMBER  Playlist video to start at (default is 1)
```

#### 3.DefaultValue

if the argument is not passed from arguments array (like `os.Args`), default value can be passed to continue, [eg](examples/yt-download/main.go)

```go
parser.Int("", "playlist-start", &argparse.Option{Help: "Playlist video to start at (default is 1)", Default: "1"})
```

> noted the Default value is not the type of `Int` , because the value is used like an argument from parse args (like `os.Args`), it's got to get through `Validate` & `Formatter` & `parse` actions (if these actions exist),  `Validate` & `Formatter` will be mentioned below
>
> also, the Default value can only be one `String` , if you want an Array arguments, you can only have one element Array as default value

#### 4. RequiredArgument

if the argument must be input, set `Required` to be `true`, [eg](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", Required: true})
```

> Flag argument can not be `Required` , you should know the reason, Flag argument has more restrictions, you will be noticed when using it

#### 5. PositionanArgument

if the input argument is the value you want, set `Positional` to be true, [eg](examples/yt-download/main.go)

```go
parser.Strings("", "url", &argparse.Option{Help: "youtube links, like 'https://www.youtube.com/watch?v=xxxxxxxx'", Positional: true})
```

> the position of the PositionalArgument is quit flex, with not much restrictions, it's ok to be
>
> 1. in the middle of arguments, `--play-list 2 xxxxxxxx --update`, if the argument before it is not an Array argument, won't parse `url` in this case: `--user-ids id1 id2 url --update` 
> 2. after another single value PositionalArgument, `--mode login username password` , the last `password` will be parsed as second PositionalArgument
>
> so, use it carefully

#### 6. ArgumentValidate

provide `Validate` function to check each passed-in argument

```go
parser.Strings("", "url", &argparse.Option{Help: "youtube links", 
     Validate: func(arg string) error {
				if !strings.HasPrefix(arg, "https://") {
					return fmt.Errorf("url should be start with 'https://'")
				}
				return nil
			}})
```

> `Validate` function has high priority, executed just after `Default` value is set, which means, the default value has to go through `Validate` check

#### 7. ArgumentFormatter

format input argument to most basic types you want

```go
parser.String("", "b", &Option{Formatter: func(arg string) (i interface{}, err error) {
		if arg == "False" {
			err = fmt.Errorf("no False")
			return
		}
		i = fmt.Sprintf("=> %s", arg)
		return
	}})
```

> if `Validate` is set, `Formatter` is right after `Validate`
>
> if raise errors in `Formatter`, it will partly act like `Validate` 
>
> the return type of `interface{}` should be the same as your Argument Type, or Element Type of your Arguments, to by `string` as Example shows

#### 8. ArgumentChoices

restrict inputs to be within the given choices, using `Choices`

```go
parser.Ints("", "hours", &Option{Choices: []interface{}{1, 2, 3, 4}})
```

> if `Formatter` is set, Choice check is right after `Formatter`
>
> when it's single value, the value must be one of the `Choices`
>
> when it's value array, each value must be one of of `Choices`

#### 9. SubCommands

create new parser realm, within the sub command parser, arguments won't interrupt each other

```go
func main() {
	parser := argparse.NewParser("", "Go is a tool for managing Go source code.", nil)
	testCommand := parser.AddCommand("test", "start a bug report", nil)
	tFlag := testCommand.Flag("f", "flag", nil)
	otherFlag := testCommand.Flag("o", "other", nil)
	t := parser.Flag("f", "flag", nil)
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	println(*tFlag, *otherFlag, *t)
}
```

output:

```bash
=> ./sub-command
usage: ./sub-command <cmd> [-h] [-f]

Go is a tool for managing Go source code.

available commands:
  test        start a bug report

optional arguments:
  -h, --help  show this help message
  -f, --flag
  
=> ./sub-command test
usage: test [-h] [-f] [-o]

start a bug report

optional arguments:
  -h, --help   show this help message
  -f, --flag
  -o, --other
```

the two `--flag` will parse seperately, so you can use `tFlag` & `t` to reference flag it `test` parser and `main` parser

##### Argument Process Map

```
                        ┌──────┐
 --date 20210102 --list │ arg1 │ arg2   arg3
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



## Config

### 1. ParserConfig

relative struct: 

```go
type ParserConfig struct {
	Usage                  string // manual usage display
	EpiLog                 string // message after help
	DisableHelp            bool   // disable help entry register [-h/--help]
	ContinueOnHelp         bool   // set true to: continue program after default help is printed
	DisableDefaultShowHelp bool   // set false to: default show help when there is no args to parse (default action)
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

output:

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

except the comment above, `ContinueOnHelp` is only affective on your program process, which give you possibility to do something when default `help` is shown

### 2. ArgumentOptions

related struct:

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

