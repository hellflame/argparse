// this is show case for expand * for positional arguments
// run the code like "go run main.go ~" or "go run main.go ~/*"
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	name := parser.Strings("n", "name", &argparse.Option{Positional: true})
	files := parser.Strings("f", "file", nil)
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	for i, n := range *name {
		fmt.Println("name:", i, n)
	}
	for i, f := range *files {
		fmt.Println("file:", i, f)
	}
}
