package main

import (
	"fmt"
	"github.com/hellflame/argparse"
)

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
