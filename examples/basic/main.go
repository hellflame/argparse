package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"log"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	name := parser.String("n", "name", nil)
	if e := parser.Parse(nil); e != nil {
		log.Fatal(e.Error())
	}
	fmt.Printf("hello %s\n", *name)
}
