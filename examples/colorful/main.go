// this is show case for using color in argparse
//
// demo for ColorSchema
package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", &argparse.ParserConfig{
		WithColor: true,
		WithHint:  true,

		EpiLog: "more info please visit https://github.com/hellflame/argparse",
	})
	sub := parser.AddCommand("run", "run your program", nil)
	parser.AddCommand("test", "test for your program", nil)

	sub.Flag("d", "dir", &argparse.Option{Help: "give me a directory"})
	parser.String("n", "name", &argparse.Option{Default: "flame", Help: "your name"})
	parser.Ints("t", "times", &argparse.Option{HintInfo: "run times", Group: "base", Help: "how many times"})
	parser.Float("s", "size", &argparse.Option{Help: "give me a size", Group: "base", Required: true})
	parser.String("u", "url", &argparse.Option{Positional: true, Help: "target url"})
	parser.String("l", "", nil)
	if e := parser.Parse(nil); e != nil {
		switch e {
		case argparse.BreakAfterHelpError:
			return
		default:
			fmt.Println(e.Error())
		}
		return
	}
}
