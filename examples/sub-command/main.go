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
	testCommand := parser.AddCommand("test", "start a bug report", &argparse.ParserConfig{WithHint: true})
	tFlag := testCommand.Flag("f", "flag", &argparse.Option{Help: "from test parser"})
	otherFlag := testCommand.Flag("o", "other", &argparse.Option{HintInfo: "optional => ∫"})
	floatWithChoice := testCommand.Float("", "float", &argparse.Option{Choices: []interface{}{0.1, 0.2}, Required: true})
	defaultInt := testCommand.Int("i", "int", &argparse.Option{Default: "1", Help: "this is int"})
	testCommand.String("s", "string", &argparse.Option{Default: "hello", Help: "no hint message", NoHint: true})
	if e := parser.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
		default:
			fmt.Println(e.Error())
		}
		return
	}
	println(*tFlag, *otherFlag, *t, *defaultInt, *floatWithChoice)
}
