package main

import (
	"fmt"
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("", "Go is a tool for managing Go source code.", nil)
	t := parser.Flag("f", "flag", nil)
	testCommand := parser.AddCommand("test", "start a bug report", nil)
	tFlag := testCommand.Flag("f", "flag", nil)
	otherFlag := testCommand.Flag("o", "other", nil)
	defaultInt := testCommand.Int("i", "int", &argparse.Option{Default: "1"})
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	println(*tFlag, *otherFlag, *t, *defaultInt)
}
