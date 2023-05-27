package argparse

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func Test_CreateParser(t *testing.T) {
	p := NewParser("test", "this is a desc", nil)
	if e := p.Parse(nil); e == nil {
		t.Error("not available for test env")
	}
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
	parser.Int("this-is-very-long-short", "this-is-very-long-full", &Option{Positional: true})
	sub := parser.AddCommand("sub", "this is a sub-command", nil)
	sub.Int("x", "no-show", nil)
	parser.AddCommand("this-is-long", "this is a sub-command", nil)
	if len(parser.FormatHelp()) == 0 {
		t.Error("failed to format help")
		return
	}
	parser = NewParser("", "test program", &ParserConfig{Usage: "test [name]"})
	if len(parser.FormatHelp()) == 0 {
		t.Error("failed to format simple help")
		return
	}

	p := NewParser("no-entry", "", nil)
	p.String("a", "", &Option{HideEntry: true, Help: "won't show"})
	p.Strings("b", "", &Option{HideEntry: true, Help: "won't show"})
	p.Int("c", "", &Option{HideEntry: true, Help: "won't show"})
	p.Ints("d", "", &Option{HideEntry: true, Help: "won't show"})
	p.Flag("e", "", &Option{HideEntry: true, Help: "won't show"})
	p.Strings("f", "", &Option{HideEntry: true, Help: "won't show", Group: "no-entry"})
	if len(strings.Split(p.FormatHelp(), "\n")) != 4 {
		t.Error("failed to hide help entry")
		return
	}
}

func TestParser_AddCommand(t *testing.T) {

	parser := NewParser("", "test program", &ParserConfig{Usage: "test [name]", ContinueOnHelp: true})
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

	sub := parser.AddCommand("test", "this is test", nil)
	sub.Flag("t", "test", nil)
	if e := parser.Parse([]string{"test", "-t"}); e != nil {
		t.Error(e.Error())
		return
	}
}

func TestParser_Default(t *testing.T) {
	parser := NewParser("", "", &ParserConfig{ContinueOnHelp: true})
	a := parser.String("", "a", &Option{Default: "default value"})
	b := parser.Int("", "b", &Option{Default: "1", Positional: true})
	sub := parser.AddCommand("test", "", nil)
	c := sub.Int("", "a", &Option{Default: "2"})
	if e := parser.Parse([]string{}); e != nil {
		t.Errorf(e.Error())
		return
	}
	if *a != "default value" {
		t.Error("failed to set default value")
		return
	}
	if *b != 1 {
		t.Error("failed to set default positional value")
		return
	}
	if *c == 2 {
		t.Error("sub command should not parse")
		return
	}
}

func TestParser_unrec(t *testing.T) {
	parser := NewParser("", "", nil)
	parser.String("a", "aa", nil)
	parser.String("b", "bb", &Option{Positional: true})
	parser.Strings("", "ab", &Option{Help: "this is abcc"})
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
	if e := parser.Parse([]string{"cover-bb", "--ax"}); e != nil {
		err := e.Error()
		if !strings.Contains(err, "unrecognized arguments: --ax") ||
			!strings.Contains(err, "--ab (this is abcc)") ||
			!strings.Contains(err, "--aa") {
			t.Error("failed to guess input")
			return
		}
	}
	if e := parser.Parse([]string{"cover-bb", "--abc"}); e != nil {
		if e.Error() != "unrecognized arguments: --abc\ndo you mean?: --ab (this is abcc)" {
			t.Error("failed to output help msg")
			return
		}
	}
	if e := parser.Parse([]string{"cover-bb", "cd"}); e != nil {
		if e.Error() != "unrecognized arguments: cd" {
			t.Error("failed to tell")
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

func TestParser_Choices(t *testing.T) {
	parser := NewParser("", "", nil)
	parser.String("", "a",
		&Option{Choices: []interface{}{"x"}})
	parser.Ints("", "b",
		&Option{Choices: []interface{}{1, 2}})
	if e := parser.Parse([]string{"--a", "y"}); e != nil {
		if e.Error() != "args must be one|some of [x]" {
			t.Error("failed to make a choice")
			return
		}
	}
	if e := parser.Parse([]string{"--a", "x"}); e != nil {
		t.Error("error choice")
		return
	}
	if e := parser.Parse([]string{"--b", "3"}); e != nil {
		if e.Error() != "args must be one|some of [1 2]" {
			t.Error("failed to make choices")
			return
		}
	}

}

func TestParser_Validate(t *testing.T) {
	parser := NewParser("", "", nil)
	a := parser.String("", "a", &Option{Validate: func(arg string) error {
		if arg == "ok" {
			return fmt.Errorf("not ok")
		}
		return nil
	}})
	if e := parser.Parse([]string{"--a", "ok"}); e != nil {
		if e.Error() != "not ok" {
			t.Error("not ok")
			return
		}
	}
	if *a != "" {
		t.Error("this is invalid value")
		return
	}
	if e := parser.Parse([]string{"--a", "this is ok"}); e != nil {
		t.Error("this is not ok")
		return
	}
	if *a != "this is ok" {
		t.Error("this should be ok")
		return
	}
}

func TestParser_Formatter(t *testing.T) {
	parser := NewParser("", "", nil)
	a := parser.Ints("", "a", &Option{Formatter: func(arg string) (i interface{}, err error) {
		v, err := strconv.Atoi(arg)
		if err != nil {
			return
		}
		i = v + 1
		return
	}})
	b := parser.String("", "b", &Option{Formatter: func(arg string) (i interface{}, err error) {
		if arg == "False" {
			err = fmt.Errorf("no False")
			return
		}
		i = fmt.Sprintf("=> %s", arg)
		return
	}})
	if e := parser.Parse([]string{"--a", "1"}); e != nil {
		t.Error(e.Error())
		return
	}
	if (*a)[0] != 2 {
		t.Error("failed to format value")
		return
	}
	if e := parser.Parse([]string{"--b", "False"}); e != nil {
		if e.Error() != "no False" {
			t.Error("formatter should filtered this")
			return
		}
	}
	if e := parser.Parse([]string{"--b", "b"}); e != nil {
		t.Error(e.Error())
		return
	}
	if *b != "=> b" {
		t.Error("formatter is not functioned")
		return
	}
}

func TestParser_Required(t *testing.T) {
	parser := NewParser("", "", nil)
	a := parser.String("", "a", &Option{Required: true})
	parser.String("", "b", nil)
	c := parser.Int("", "c", &Option{Required: true, Positional: true})
	if e := parser.Parse([]string{"--b", "linux"}); e != nil {
		if e.Error() != "A is required" {
			t.Error("A is required but unknown")
			return
		}
	}
	// parser should be new, strictly speaking
	if e := parser.Parse([]string{"--a", "x", "3"}); e != nil {
		t.Error("failed to parse required value")
		return
	}
	if *a != "x" || *c != 3 {
		t.Error("failed to parse")
		return
	}

	p := NewParser("", "", nil)
	p.String("", "a", &Option{Required: true, Positional: true})
	p.String("", "b", nil)
	if e := p.Parse([]string{"--b", "x"}); e != nil {
		if e.Error() != "A is required" {
			t.Error("A is required")
			return
		}
	}
}

func TestParse_Actions(t *testing.T) {
	p := NewParser("", "", nil)
	sum := 0
	p.Strings("", "a", &Option{Positional: true, Action: func(args []string) error {
		for _, a := range args {
			i, e := strconv.Atoi(a)
			if e != nil {
				return fmt.Errorf("invalid int value %s", a)
			}
			sum += i
		}
		return nil
	}})
	if e := p.Parse([]string{"1", "2", "3"}); e != nil {
		t.Error(e.Error())
		return
	}
	if sum != 6 {
		t.Error("fail to sum up")
		return
	}

	p = NewParser("", "", nil)

	buttonColor := "gray"
	p.Flag("", "a", &Option{Action: func(args []string) error {
		buttonColor = "RED"
		return nil
	}})
	if e := p.Parse([]string{"--a"}); e != nil {
		t.Error(e.Error())
		return
	}
	if buttonColor != "RED" {
		t.Error("failed to do the action")
		return
	}
}

func TestParser_Fail(t *testing.T) {
	p := NewParser("", "", nil)
	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Errorf("failed to panic")
			}
		}()
		p.String("", "", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p.String("a", "aa", nil)
		p.Int("", "aa", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p := NewParser("", "", nil)
		p.String("a", "aa", nil)
		p.Ints("", "aa", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p := NewParser("", "", nil)
		p.String("a", "aa", nil)
		p.Float("", "aa", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p := NewParser("", "", nil)
		p.String("a", "aa", nil)
		p.Floats("", "aa", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p := NewParser("", "", nil)
		p.String("a", "aa", nil)
		p.Strings("", "aa", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("failed to panic")
			}
		}()
		p := NewParser("", "", nil)
		p.String("a", "aa", nil)
		p.Flag("", "aa", nil)
	}()

	p = NewParser("", "", nil)
	p.Int("", "a", nil)
	p.Float("", "b", nil)
	p.Floats("", "c", &Option{Positional: true})
	if e := p.Parse([]string{"--a", "x"}); e == nil || e.Error() != "invalid int value: x" {
		t.Error("failed to check invalid int")
		return
	}
	if e := p.Parse([]string{"--b", "x"}); e == nil || e.Error() != "invalid float value: x" {
		t.Error("failed to check invalid float")
		return
	}
	if e := p.Parse([]string{"a", "b", "c"}); e == nil || e.Error() != "invalid float value: a" {
		t.Error("failed to check invalid value in multi positional")
		return
	}
	p = NewParser("", "", nil)
	p.Int("", "a", &Option{Positional: true})
	if e := p.Parse([]string{"a"}); e == nil || e.Error() != "invalid int value: a" {
		t.Error("failed to check invalid value in positional")
		return
	}
	test := p.AddCommand("test", "", nil)
	test.Int("", "a", nil)
	test.Float("", "b", &Option{Default: "x"})
	if e := p.Parse([]string{"test", "--a", "x"}); e == nil || e.Error() != "invalid int value: x" {
		t.Error("failed to check invalid value in sub command")
		return
	}

	p = NewParser("", "", nil)
	p.Float("", "a", &Option{Default: "x"})
	p.Int("", "b", nil)
	if e := p.Parse([]string{"--b", "1"}); e == nil || e.Error() != "invalid float value: x" {
		t.Error("failed to check invalid value with default")
		return
	}
}

func TestParser_WithSubCommand(t *testing.T) {
	p := NewParser("", "", nil)
	f := p.Flag("", "a", nil)
	sub := p.AddCommand("test", "", nil)
	x := sub.Int("", "a", &Option{Default: "1"})
	if e := p.Parse([]string{"--a"}); e != nil {
		t.Error(e.Error())
		return
	}
	if !*f {
		t.Error("failed to parse ")
		return
	}
	if *x == 1 {
		t.Error("sub command is not parsed, default value should not be bond")
		return
	}
}

func TestParser_usage(t *testing.T) {
	p := NewParser("a", "", &ParserConfig{ContinueOnHelp: true})
	test := p.AddCommand("test", "", nil)
	if !strings.HasPrefix(test.formatUsage(), "usage: a test") {
		t.Error("failed to generate sub usage")
	}
}

func TestDefaultAction(t *testing.T) {
	d := false
	p := NewParser("", "", &ParserConfig{DefaultAction: func() {
		d = true
	}})
	if e := p.Parse([]string{}); e != nil {
		t.Error(e.Error())
		return
	}
	if !d {
		t.Error("failed to do default action")
		return
	}

	test := false
	p = NewParser("", "", &ParserConfig{ContinueOnHelp: true, DisableDefaultShowHelp: true})
	p.AddCommand("test", "", &ParserConfig{DefaultAction: func() {
		test = true
	}})
	if e := p.Parse([]string{}); e != nil {
		t.Error(e.Error())
		return
	}
	if test {
		t.Error("test should not run default action when it's not hit")
		return
	}

	test = false
	p = NewParser("", "", &ParserConfig{ContinueOnHelp: true})
	p.AddCommand("test", "", &ParserConfig{DefaultAction: func() {
		test = true
	}})
	if e := p.Parse([]string{"test"}); e != nil {
		t.Error(e.Error())
		return
	}
	if !test {
		t.Error("default action should be ran")
		return
	}

}

func TestParse_Completion(t *testing.T) {
	p := NewParser("cp.sh", "this is test", nil)
	p.Strings("a", "aa", nil)
	p.Int("", "bb", nil)
	p.Float("c", "cc", &Option{Positional: true})
	p.String("d", "", &Option{HideEntry: true})
	test := p.AddCommand("test", "", nil)
	test.String("a", "aa", nil)
	test.Int("", "bb", nil)
	install := p.AddCommand("install", "", nil)
	install.Strings("i", "in", nil)
	install.Strings("a", "", &Option{HideEntry: true})
	install.Flag("b", "", &Option{HideEntry: true})
	if p.FormatCompletionScript() == "" {
		t.Error("failed to generate completion script")
		return
	}
}

func TestParse_AddCompletion(t *testing.T) {
	p := NewParser("", "", &ParserConfig{AddShellCompletion: true})
	p.Int("", "a", nil)
}

func TestParse_AllowShort(t *testing.T) {
	p := NewParser("", "", nil)
	x := p.Int("x", "", nil)
	if e := p.Parse([]string{"-x", "1"}); e != nil {
		t.Error(e.Error())
		return
	}
	if *x != 1 {
		t.Error("failed to parse only short arg")
	}
	if p.FormatHelp() == "" {
		t.Error("failed to generate help")
		return
	}
}

func TestParser_Invoke(t *testing.T) {
	p := NewParser("", "", nil)
	a := p.String("a", "", &Option{Default: "no!"})
	sub := p.AddCommand("sub", "", nil)
	b := sub.String("b", "", nil)
	mainParsed := false
	subParsed := false
	No2Parsed := false
	p.InvokeAction = func(invoked bool) {
		mainParsed = invoked
		if *a == "no!" {
			t.Error("not parsed parser should not bind default")
		}
	}
	sub.InvokeAction = func(invoked bool) {
		subParsed = invoked
		if *b != "linux" {
			t.Error("failed to bind")
			return
		}
	}
	subNo2 := p.AddCommand("sub2", "", nil)
	subNo2.Int("a", "", nil)
	subNo2.InvokeAction = func(invoked bool) {
		No2Parsed = invoked
	}

	if e := p.Parse([]string{"sub", "-b", "linux"}); e != nil {
		t.Error(e.Error())
		return
	}
	if p.Invoked || mainParsed {
		t.Error("root parse should not be parsed")
		return
	}
	if !sub.Invoked || !subParsed {
		t.Error("sub parse state error")
		return
	}
	if subNo2.Invoked || No2Parsed {
		t.Error("irrelevant parse state error")
		return
	}
}

func TestParser_Error(t *testing.T) {
	p := NewParser("", "", nil)
	p.Int("x", "", nil)
	if e := p.Parse([]string{}); e != nil {
		switch e.(type) {
		case BreakAfterHelp:
			if e.Error() != "" {
				t.Error("error not empty")
			}
		default:
			t.Error("not normal error")
		}
	}

	p = NewParser("", "", &ParserConfig{AddShellCompletion: true})
	p.Int("x", "", nil)
	if e := p.Parse([]string{"--completion"}); e != nil {
		switch e.(type) {
		case BreakAfterShellScript:
			if e.Error() != "" {
				t.Error("error not empty")
			}
		default:
			t.Error("completion is not normal error")
		}
	}
}

func TestWithHint(t *testing.T) {
	p := NewParser("", "", &ParserConfig{WithHint: true})
	p.Int("a", "", &Option{Default: "34", Help: "this is an int"})
	p.String("b", "", &Option{Default: "a", Positional: true})
	if !strings.Contains(p.FormatHelp(), "(default: 34)") {
		t.Error("default value missing")
		return
	}
	if !strings.Contains(p.FormatHelp(), "(default: a)") {
		t.Error("default value missing for positional")
		return
	}
}

func TestMaxHeaderLength(t *testing.T) {
	p := NewParser("", "", &ParserConfig{MaxHeaderLength: 20})
	p.String("s", "short", &Option{Help: "this is short help message"})
	p.String("l", "this-is-a-very-long-args", &Option{Help: "this is long"})
	p.FormatHelp()
}

func TestParseOptionalAfterPositional(t *testing.T) {
	p := NewParser("", "", nil)
	a := p.String("a", "", &Option{Positional: true})
	bs := p.Strings("b", "", nil)
	if e := p.Parse([]string{"position", "-b", "1", "2"}); e != nil {
		t.Error(e.Error())
	}
	if *a != "position" || (*bs)[0] != "1" || (*bs)[1] != "2" {
		t.Error("failed to set values")
	}
}

func TestConflictPositional(t *testing.T) {
	p := NewParser("", "", nil)
	p.String("a", "", &Option{Positional: true})
	closure := func() {
		defer func() {
			if e := recover(); e == nil {
				t.Error("failed to panic for conflict positional!")
			}
		}()
		p.Int("", "a", &Option{Positional: true})
	}
	closure()
}

func TestInheritable(t *testing.T) {
	p := NewParser("", "", nil)
	f1 := p.Flag("f", "flag", &Option{Inheritable: true})
	p.Flag("f2", "flag2", &Option{Inheritable: true})
	pos1 := p.String("s", "str", &Option{Positional: true, Inheritable: true, Help: "position str"})
	p.Int("p2", "pos2", &Option{Positional: true, Inheritable: true, Help: "positional 2"})

	sub := p.AddCommand("sub1", "", nil)
	f2 := sub.Float("f2", "", nil)                                                          // override flag2
	pos2 := sub.Float("", "pos2", &Option{Positional: true, Help: "override positional 2"}) // override pos2

	p.String("f3", "", &Option{Inheritable: true})

	if e := p.Parse([]string{"sub1", "-f", "pos"}); e != nil {
		t.Error("f1 should be inherited by sub1")
		return
	}
	if !*f1 || *pos1 != "pos" {
		t.Error("failed to inherit from root parser")
		return
	}

	if e := p.Parse([]string{"sub1", "-f3", "linux"}); e == nil {
		t.Error("f3 should not be inherited")
		return
	}
	// sub.PrintHelp()
	if e := p.Parse([]string{"sub1", "-f2", "0.2", "fill inherit string", "9.9"}); e != nil {
		t.Error("override failed!")
		return
	}
	if *f2 > 1 || *pos2 < 1 {
		println(*f2, *pos2)
		t.Error("mess parse")
		return
	}
}

func TestBindParser(t *testing.T) {
	parser := NewParser("", "", &ParserConfig{DisableHelp: true})
	a := parser.AddCommand("a", "", nil)
	b := parser.AddCommand("b", "", nil)
	c := parser.AddCommand("c", "", nil)

	ab := parser.String("", "ab", &Option{
		BindParsers: []*Parser{a, b},
	})
	bc := parser.String("", "bc", &Option{
		BindParsers: []*Parser{b, c},
	})
	if len(parser.entries) != 0 || len(a.entries) != 1 || len(b.entries) != 2 || len(c.entries) != 1 {
		t.Error("error binding")
		return
	}
	if e := parser.Parse([]string{"b", "--ab", "x", "--bc", "y"}); e != nil {
		t.Error(e)
		return
	}
	if *ab != "x" || *bc != "y" {
		t.Error("failed to parse")
	}
}

func TestBindParserConflict(t *testing.T) {
	defer func() {
		if e := recover(); e == nil || e.(string) != "conflict entry for '-c', say: ''" {
			t.Error("failed to confilict")
		}
	}()
	parser := NewParser("", "", nil)
	a := parser.AddCommand("a", "", nil)

	parser.String("c", "", &Option{BindParsers: []*Parser{
		a, a,
	}})
}

func TestColorful(t *testing.T) {
	os.Setenv("TERM", "")
	p := NewParser("", "", &ParserConfig{WithColor: true})
	p.String("l", "", nil)
	if strings.Contains(p.FormatHelp(), "\033[00m") {
		t.Error("should be no color")
		return
	}

	p = NewParser("", "", &ParserConfig{WithHint: true,
		WithColor: true, EnsureColor: true, ColorSchema: DefaultColor})
	p.String("l", "", &Option{HintInfo: "", Group: "base"})
	if !strings.Contains(p.FormatHelp(), "\033[00m") {
		t.Error("should be colorful")
		return
	}
}
