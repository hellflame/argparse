//go:build ignore

// this is show case for default parse action when no user input is given, also has effect on sub command
//
// run code like "go run main.go" or "go run main.go test"
package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", &argparse.ParserConfig{DefaultAction: func() {
		fmt.Println("hi ~\ntell me what to do?")
	}})
	parser.AddCommand("test", "testing", &argparse.ParserConfig{DefaultAction: func() {
		fmt.Println("ok, now you know you are testing")
	}})
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
}
