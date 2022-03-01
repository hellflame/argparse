// this is show case for most simple use of argparse
//
// with a String Optional Argument created binding to variable 'name'
// and default help support
//
// this show case shows how to handle after help or shell script showed,
// but it's ok not to handle them separately, because their error msg is empty,
// only one more blank line will be added in the tail output
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"os"
)

func main() {
	p := argparse.NewParser("start", "this is test",
		&argparse.ParserConfig{AddShellCompletion: true})
	p.Strings("a", "aa", nil)
	p.Int("", "bb", nil)
	p.Float("c", "cc", &argparse.Option{Positional: true})
	test := p.AddCommand("test", "", nil)
	test.String("a", "aa", nil)
	test.Int("", "bb", nil)
	install := p.AddCommand("install", "", nil)
	install.Strings("i", "in", nil)
	if _, e := p.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
			os.Exit(1)
		case argparse.BreakAfterShellScript:
		default:
			fmt.Printf(e.Error())
		}
		return
	}
}
