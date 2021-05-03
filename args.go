package argparse

import (
	"fmt"
	"strings"
)

const fullPrefix = "--"
const shortPrefix = "-"

type arg struct {
	short    string
	full     string
	target   interface{}
	assigned bool
	Option
}

type Option struct {
	Meta       string
	multi      bool
	Default    interface{}
	isFlag     bool
	Required   bool
	Positional bool
	Help       string
	Group      string
	Choices    []interface{}
	Validate   func(arg string) error
	Formatter  func(arg string) (interface{}, error)
}

// validate args setting before parsing args, right after adding to parser
// for conflict check & correction & restriction
func (a *arg) validate() error {
	if a.full == "" {
		return fmt.Errorf("arg name is empty")
	}
	if strings.Contains(a.full, " ") || strings.Contains(a.short, " ") {
		return fmt.Errorf("arg name with space")
	}
	if strings.HasPrefix(a.full, shortPrefix) || strings.HasPrefix(a.full, fullPrefix) {
		return fmt.Errorf("arg full name with extra prefix '%s'/'%s'", shortPrefix, fullPrefix)
	}
	if strings.HasPrefix(a.short, shortPrefix) {
		return fmt.Errorf("arg short name with extra prefix '%s'", shortPrefix)
	}
	if a.short == a.full {
		return fmt.Errorf("arg short is full")
	}
	if a.Positional {
		if a.isFlag {
			return fmt.Errorf("positional is a flag")
		}
	}
	if a.isFlag {
		if a.Meta != "" {
			return fmt.Errorf("flag with meta")
		}
		if len(a.Choices) != 0 {
			return fmt.Errorf("flag has choices")
		}
		if a.Required {
			return fmt.Errorf("flag with required")
		}
		if a.Formatter != nil {
			return fmt.Errorf("flag with formmater")
		}
		if a.Validate != nil {
			return fmt.Errorf("flag with validate")
		}
	}
	return nil
}

// argument watch list
// for parser use
func (a *arg) getWatchers() []string {
	if a.Positional {
		return []string{}
	}
	result := []string{fmt.Sprintf("%s%s", fullPrefix, a.full)}
	if a.short != "" && a.short != a.full {
		result = append([]string{fmt.Sprintf("%s%s", shortPrefix, a.short)}, result...)
	}
	return result
}

func (a *arg) getMetaName() string {
	if a.Meta != "" {
		return a.Meta
	}
	return strings.ToUpper(a.full)
}

func (a *arg) formatHelpHeader() string {
	metaName := a.getMetaName()
	if a.Positional {
		return metaName
	}
	watchers := a.getWatchers()
	if a.isFlag {
		return strings.Join(watchers, ", ")
	}
	var signedWatchers []string
	for _, w := range watchers {
		signedWatchers = append(signedWatchers, fmt.Sprintf("%s %s", w, metaName))
	}
	return strings.Join(signedWatchers, ", ")
}

func (a *arg) parseValue(values []string) error {
	a.assigned = true
	if a.isFlag {
		*a.target.(*bool) = true
		return nil
	}
	if a.Validate != nil {
		for _, v := range values {
			e := a.Validate(v)
			if e != nil {
				return e
			}
		}
	}
	var result []interface{}
	if a.Formatter != nil {
		for _, v := range values {
			f, e := a.Formatter(v)
			if e != nil {
				return e
			}
			result = append(result, f)
		}
	} else {
		switch a.target.(type) {
		case *string:
			result = append(result, values[0])
		case *[]string:
			for _, v := range values {
				result = append(result, v)
			}
		}
	}
	if len(result) == 0 {
		return fmt.Errorf("no value to parse")
	}
	if len(a.Choices) > 0 {
		for _, r := range result {
			found := false
			for _, c := range a.Choices {
				if c == r {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("args must one|some of %+v", a.Choices)
			}
		}
	}
	switch a.target.(type) {
	case *string:
		*a.target.(*string) = result[0].(string)
	case *int:
		*a.target.(*int) = result[0].(int)
	case *float64:
		*a.target.(*float64) = result[0].(float64)
	case *[]string:
		for _, r := range result {
			*a.target.(*[]string) = append(*a.target.(*[]string), r.(string))
		}
	case *[]int, *[]float64:
		a.target = result
	}
	return nil
}
