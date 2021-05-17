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
