// this show case is for sub command
//
// sub command is created by AddCommand, which returns a *Parser for programmer to bind arguments
//
// sub command has different parse context from main parser (created by NewParser)
//
// mainly use sub command to help user understand your program step by step
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("sub-command", "Go is a tool for managing Go source code.", nil)
	t := parser.Flag("f", "flag", &argparse.Option{Help: "from main parser"})
	testCommand := parser.AddCommand("test", "start a bug report", nil)
	tFlag := testCommand.Flag("f", "flag", &argparse.Option{Help: "from test parser"})
	otherFlag := testCommand.Flag("o", "other", nil)
	defaultInt := testCommand.Int("i", "int", &argparse.Option{Default: "1"})
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	println(*tFlag, *otherFlag, *t, *defaultInt)
}
