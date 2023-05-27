//go:build ignore

// this is show case for Hide entry
//
// you won't see argument 'greet', but you can still use the entry
package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program",
		&argparse.ParserConfig{AddShellCompletion: true})
	name := parser.String("n", "name", nil)
	greet := parser.String("g", "greet", &argparse.Option{HideEntry: true})
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	greetWord := "hello"
	if *greet != "" {
		greetWord = *greet
	}
	fmt.Printf("%s %s\n", greetWord, *name)

}
