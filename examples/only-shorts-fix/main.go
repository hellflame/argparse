//go:build ignore

// show case for testing multi arguments with only short tags
package main

import (
	"os"
	"strings"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("", "", &argparse.ParserConfig{WithHint: true})
	v := parser.Flag("v", "", &argparse.Option{Help: "version info"})
	d := parser.Flag("d", "", &argparse.Option{Help: "some d"})
	f := parser.Strings("f", "", &argparse.Option{Positional: true, Help: "file path", Required: true})
	if e := parser.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
			os.Exit(1)
		case argparse.BreakAfterShellScript:
		default:
			println(e.Error())
		}
		return
	}
	println("v =>", *v)
	println("d =>", *d)
	println("f =>", strings.Join(*f, ", "))
}
