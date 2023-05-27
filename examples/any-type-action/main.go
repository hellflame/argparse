//go:build ignore

// this is show case for argument action
//
// argument action will be executed when user input has a match to the binding argument
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/hellflame/argparse"
)

// run like: go run main.go 2 3 4 --b main.go
func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	sum := 0
	content := ""
	parser.Strings("", "a", &argparse.Option{Positional: true, Action: func(args []string) error {
		for _, arg := range args {
			i, e := strconv.Atoi(arg)
			if e != nil {
				return e
			}
			sum += i
		}
		return nil
	}})

	parser.String("", "b", &argparse.Option{Action: func(args []string) error {
		raw, e := ioutil.ReadFile(args[0])
		if e != nil {
			return e
		}
		content = string(raw)
		return nil
	}})

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}

	fmt.Printf("totals: %d\n", sum)
	fmt.Println(content)
}
