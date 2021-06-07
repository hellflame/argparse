// this is show case for multi parser in Action
//
// it act like sub command, but with less restriction
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
)

func main() {

	parser := argparse.NewParser("basic", "this is a basic program", nil)

	testParser := argparse.NewParser("test", "test command", nil)
	isOk := testParser.Flag("", "ok", nil)

	parser.Strings("", "test", &argparse.Option{Action: func(args []string) error {
		if e := testParser.Parse(args); e != nil {
			return e
		}
		fmt.Println("is it ok?", *isOk)
		return nil
	}})

	buildParser := argparse.NewParser("build", "build command", nil)
	count := buildParser.Int("c", "count", nil)

	parser.Strings("", "build", &argparse.Option{Action: func(args []string) error {
		if e := buildParser.Parse(args); e != nil {
			return e
		}
		fmt.Println("got an int", *count)
		return nil
	}})

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
}
