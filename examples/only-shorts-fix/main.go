package main

import (
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("", "", nil)
	v := parser.Flag("v", "", nil)
	d := parser.Flag("d", "", nil)
	f := parser.String("f", "", &argparse.Option{Positional: true})
	if e := parser.Parse(nil); e != nil {
		print(e.Error())
		return
	}
	println("v =>", *v)
	println("d =>", *d)
	println("f =>", *f)
}
