//go:build ignore

// this is show case like npm install xxx xx
//
// with sub-command and its positional argument, you can run the code like "go run main.go install express vue"
package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("npm", "test npm install xxx", nil)
	install := parser.AddCommand("install", "install something", nil)
	pkgs := install.Strings("", "package", &argparse.Option{Positional: true})
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	for _, pkg := range *pkgs {
		fmt.Println(pkg)
	}
}
