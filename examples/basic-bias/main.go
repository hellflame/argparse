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
