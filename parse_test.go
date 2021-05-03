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

func TestParser_unrec(t *testing.T) {
	parser := NewParser("", "", nil)
	parser.String("a", "aa", nil)
	parser.String("b", "bb", &Option{Positional: true})
	if e := parser.Parse([]string{"x", "b"}); e != nil {
		if e.Error() != "unrecognized arguments: b" {
			t.Error("failed to un-recognize")
			return
		}
	}
	if e := parser.Parse([]string{"-a", "a", "b", "bx"}); e != nil {
		if e.Error() != "unrecognized arguments: bx" {
			t.Error("failed to un-recognize")
			return
		}
	}
}

func TestParser_Parse(t *testing.T) {
	parser := NewParser("", "", &ParserConfig{ContinueOnHelp: true})
	a := parser.String("a", "aa", nil)
	b := parser.Strings("b", "bb", nil)
	c := parser.String("", "cc", nil)
	f := parser.Flag("", "ff", nil)
	d := parser.String("", "dd", &Option{Positional: true})
	e := parser.String("", "ee", &Option{Positional: true})
	g := parser.Strings("", "gg", &Option{Positional: true})

	if e := parser.Parse([]string{"-a", "linux", "-b", "b1", "b2", "--cc", "x", "--ff"}); e != nil {
		t.Errorf(e.Error())
	}
	if *a != "linux" || len(*b) != 2 || (*b)[1] != "b2" || *c != "x" || !*f {
		t.Error("failed to parse string")
		return
	}
	if e := parser.Parse([]string{"-a"}); e == nil {
		t.Error("expect argument")
		return
	}
	if e := parser.Parse([]string{"linux", "ok"}); e != nil {
		t.Error("failed to parse")
		return
	}
	if *d != "linux" || *e != "ok" {
		t.Error("failed to parse position args")
	}
	if e := parser.Parse([]string{"linux", "ok", "g1", "g2"}); e != nil {
		t.Error(e.Error())
		return
	}
	if *e != "ok" || len(*g) != 2 || (*g)[0] != "g1" {
		t.Error("failed to parse")
		return
	}
}

func TestParser_types(t *testing.T) {
	parser := NewParser("", "", nil)
	a := parser.String("", "a", nil)
	b := parser.Strings("", "b", nil)
	c := parser.Int("", "c", nil)
	d := parser.Ints("", "d", nil)
	e := parser.Float("", "e", nil)
	f := parser.Floats("", "f", nil)
	g := parser.Flag("", "g", nil)
	if e := parser.Parse([]string{"--a", "a", "--b", "b1", "b2", "--c", "1",
		"--d", "1", "2", "--e", "3.14", "--f", "0.618", "2.7", "--g"}); e != nil {
		t.Error(e.Error())
		return
	}
	if *a != "a" || !*g || strings.Join(*b, ",") != "b1,b2" || *c != 1 ||
		len(*d) != 2 || (*d)[1] != 2 || *e != 3.14 ||
		len(*f) != 2 || (*f)[1] != 2.7 {
		t.Errorf("failed to apply values")
		return
	}
}
