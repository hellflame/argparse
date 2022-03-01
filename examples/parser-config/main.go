// this show case is for 'ParserConfig', showing how the config affect your parsing progress
//
// set ParserConfig.Usage will change your usage line, which can sometime be too complex for user to read
//
// set ParserConfig.EpiLog will append a message after help message, which is usually for contact info
//
// set ParserConfig.DisableHelp will prevent the default help entry injection, with ParserConfig.DisableHelp, you won't see '-h' or '--help' entry
//
// set ParserConfig.DisableDefaultShowHelp will farther prevent default help output when there is no user input
//
// set ParserConfig.ContinueOnHelp will keep program going on when the original help action is done
//
// the configs about 'Help' is not often used, most programmer may not care about it
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
	if _, e := parser.Parse(nil); e != nil {
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
