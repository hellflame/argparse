package argparse

import (
	"strings"
	"testing"
)

func Test_CreateParser(t *testing.T) {
	NewParser("test", "this is a desc", nil)
}

func TestParser_help(t *testing.T) {
	parser := NewParser("test", "this is a test program", &ParserConfig{EpiLog: "this is epi-log"})
	parser.Int("i", "int", nil)
	parser.Int("in", "another", &Option{Help: strings.Repeat("help ", 10)})
	parser.String("", "string", &Option{Help: "this is string", Positional: true, Group: "Ga options"})
	parser.String("", "s2", &Option{Help: "this is string", Group: "Ga options"})
	parser.String("", "s3", &Option{Help: "this is string", Positional: true})
	parser.Strings("", "p3", &Option{Help: "this is string", Positional: true})
	parser.Strings("", "p4", &Option{Help: "this is string", Positional: true, Required: true})
	parser.String("x", "test", &Option{Group: "test group", Meta: "t", Required: true})
	parser.String("", "test2", &Option{Group: "test group"})
	parser.String("", "test3", nil)
	parser.Strings("", "test4", nil)
	parser.Strings("", "t5", &Option{Required: true})
	parser.Flag("ok", "ok2", nil)
	sub := parser.AddCommand("sub", "this is a sub-command", nil)
	sub.Int("x", "no-show", nil)
	if len(parser.FormatHelp()) == 0 {
		t.Error("failed to format help")
		return
	}
	parser = NewParser("", "test program", &ParserConfig{Usage: "test [name]"})
	if len(parser.FormatHelp()) == 0 {
		t.Error("failed to format simple help")
		return
	}
}

func TestParser_AddCommand(t *testing.T) {

	parser := NewParser("", "test program", &ParserConfig{Usage: "test [name]"})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("failed to panic")
			}
		}()
		parser.AddCommand("", "", nil)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("failed to panic")
			}
		}()
		parser.AddCommand("a b", "", nil)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("failed to panic")
			}
		}()
		parser.AddCommand("ab", "", nil)
		parser.AddCommand("ab", "", nil)
	}()

}

func TestParser_Default(t *testing.T) {
	parser := NewParser("", "", &ParserConfig{ContinueOnHelp: true})
	parser.String("a", "aa", nil)
	if e := parser.Parse([]string{}); e != nil {
		t.Errorf(e.Error())
		return
	}
}

func TestParser_Parse(t *testing.T) {
	parser := NewParser("", "", &ParserConfig{ContinueOnHelp: true})
	a := parser.String("a", "aa", nil)
	b := parser.Strings("b", "bb", nil)
	c := parser.String("", "cc", nil)
	f := parser.Flag("", "ff", nil)

	if e := parser.Parse([]string{"-a", "linux", "-b", "b1", "b2", "--cc", "x", "--ff"}); e != nil {
		t.Errorf(e.Error())
	}
	if *a != "linux" || len(*b) != 2 || (*b)[1] != "b2" || *c != "x" {
		t.Error("failed to parse string")
		return
	}
	println(*f)
}
