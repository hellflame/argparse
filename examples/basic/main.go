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
