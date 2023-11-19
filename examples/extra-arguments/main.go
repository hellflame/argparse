//go:build ignore

// this is show case for extra arguments after --
//
// if there are argument value happen to start with - or --, it may be regard as unknown arguments.
// in this scenario, you can input any input after a --, eg: program -- extra1 extra2
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hellflame/argparse"
)

func main() {
	// example:
	// go run main.go x y -- --x --y
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	names := parser.Strings("", "names", &argparse.Option{Positional: true})
	if e := parser.Parse(nil); e != nil {
		switch e {
		case argparse.BreakAfterHelpError:
			os.Exit(1)
		default:
			fmt.Println(e.Error())
		}
		return
	}
	fmt.Printf("hello %s\n", strings.Join(*names, ", "))
}
