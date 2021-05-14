// this is show case for most simple use of argparse
//
// with a String Optional Argument created binding to variable 'name'
// and default help support
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
