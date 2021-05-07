package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"io/ioutil"
	"strconv"
)

// run like: go run main.go 2 3 4 --b main.go
func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	sum := 0
	content := ""
	parser.Strings("", "a", &argparse.Option{Positional: true, Action: func(args []string) error {
		for _, arg := range args {
			if i, e := strconv.Atoi(arg); e == nil {
				sum += i
			} else {
				return e
			}
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
