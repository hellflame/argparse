package argparse

import (
	"fmt"
	"strings"
)

const fullPrefix = "--"
const shortPrefix = "-"

type arg struct {
	short  string
	full   string
	target interface{}
	Option
}

type Option struct {
	Meta       string
	Default    interface{}
	IsFlag     bool
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
	if a.Positional {
		if a.IsFlag {
			return fmt.Errorf("positional is a flag")
		}
	}
	if a.IsFlag {
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
	if a.IsFlag {
		return strings.Join(watchers, ", ")
	}
	var signedWatchers []string
	for _, w := range watchers {
		signedWatchers = append(signedWatchers, fmt.Sprintf("%s %s", w, metaName))
	}
	return strings.Join(signedWatchers, ", ")
}
