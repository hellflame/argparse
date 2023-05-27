//go:build ignore

// this is show case for most simple use of argparse
//
// with a String Optional Argument created binding to variable 'name'
// and default help support
package main

import (
	"fmt"
	"os"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	name := parser.String("n", "name", nil)
	if e := parser.Parse(nil); e != nil {
		// // version < 1.10
		// switch e.(type) {
		// case argparse.BreakAfterHelp:
		switch e {
		case argparse.BreakAfterHelpError:
			os.Exit(1)
		default:
			fmt.Println(e.Error())
		}
		return
	}
	fmt.Printf("hello %s\n", *name)
}
