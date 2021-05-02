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
	parser.String("", "s3", &Option{Help: "this is string", Positional: true})
	parser.String("", "s2", &Option{Help: "this is string", Group: "Ga options"})
	parser.String("x", "test", &Option{Group: "test group", Meta: "t", Required: true})
	parser.String("", "test2", &Option{Group: "test group"})
	parser.String("", "test3", nil)
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
