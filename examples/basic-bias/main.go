// this is show case for a little bit complex than 'basic'
// which disabled the default help menu entry '-h' or '--help',
// instead, it use '-help' and '--help-me' as help menu entry and handle help print manually
// also, it pass 'os.Args[1:]' manually to 'Parse' method, which is the same as pass 'nil'
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"os"
)

func main() {
	parser := argparse.NewParser("", "this is a basic program", &argparse.ParserConfig{
		DisableHelp:            true,
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
