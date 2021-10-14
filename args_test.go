package argparse

import "testing"

func TestArgs(t *testing.T) {
	if e := (&arg{}).validate(); e != nil {
		if e.Error() != "arg name is empty" {
			t.Error("arg name is empty")
			return
		}
	}
	if e := (&arg{full: "linux is"}).validate(); e != nil {
		if e.Error() != "arg name with space" {
			t.Error("arg name with space")
			return
		}
	}
	if e := (&arg{full: "-program"}).validate(); e != nil {
		if e.Error() != "arg full name with extra prefix '-'/'--'" {
			t.Error("arg full name with extra prefix '-'/'--'")
			return
		}
	}
	if e := (&arg{full: "program", short: "-p"}).validate(); e != nil {
		if e.Error() != "arg short name with extra prefix '-'" {
			t.Error("arg short name with extra prefix '-'")
			return
		}
	}
	if e := (&arg{full: "a", short: "a"}).validate(); e != nil {
		if e.Error() != "arg short is full" {
			t.Error("arg short is full")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{Positional: true, isFlag: true}}).validate(); e != nil {
		if e.Error() != "positional is a flag" {
			t.Error("positional is a flag")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{isFlag: true, Meta: "a"}}).validate(); e != nil {
		if e.Error() != "flag with meta" {
			t.Error("flag with meta")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{isFlag: true,
		Choices: []interface{}{"x"}}}).validate(); e != nil {
		if e.Error() != "flag has choices" {
			t.Error("flag has choices")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{isFlag: true, Required: true}}).validate(); e != nil {
		if e.Error() != "flag with required" {
			t.Error("flag with required")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{isFlag: true, Validate: func(arg string) error {
		return nil
	}}}).validate(); e != nil {
		if e.Error() != "flag with validate" {
			t.Error("flag with validate")
			return
		}
	}
	if e := (&arg{full: "a", Option: Option{isFlag: true, Formatter: func(arg string) (i interface{}, err error) {
		return nil, nil
	}}}).validate(); e != nil {
		if e.Error() != "flag with formatter" {
			t.Error("flag with formatter")
			return
		}
	}
}

func TestArgs_HideEntry(t *testing.T) {
	if (&arg{full: "a", Option: Option{HideEntry: true}}).formatUsage() != "" {
		t.Error("failed to hide usage entry")
		return
	}
}

func TestExtraInfo(t *testing.T) {
	if (&arg{full: "a", Option: Option{HintInfo: "value: 0 -> π"}}).formatHelpWithExtraInfo() != "(value: 0 -> π)" {
		t.Error("failed to use HintInfo")
		return
	}
	if (&arg{full: "a", Option: Option{Choices: []interface{}{1, 2, 3}}}).formatHelpWithExtraInfo() != "(options: [1, 2, 3])" {
		t.Error("failed to format extra choice")
		return
	}

	if (&arg{full: "a", Option: Option{Choices: []interface{}{"a", "b"}, Required: true}}).formatHelpWithExtraInfo() != "(options: [a, b], required)" {
		t.Error("failed to format Required args")
		return
	}
	if (&arg{full: "a", Option: Option{Choices: []interface{}{0.25, 0.5}}}).formatHelpWithExtraInfo() != "(options: [0.250000, 0.500000])" {
		t.Error("failed to generate choices for float")
		return
	}
}
