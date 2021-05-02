package argparse

import (
    "fmt"
    "strings"
)

const fullPrefix = "--"
const shortPrefix = "-"

type arg struct {
	short string
	full string
	target interface{}
    Option
}

type Option struct {
    Meta string
    Default interface{}
    IsFlag bool
    Required bool
    Help string
    Choices []interface{}
    Validate func(arg string) error
    Formatter func(arg string) (interface{}, error)
}


// validate args setting before parsing args, right after adding to parser
// for conflict check & correction & restrict
func (a *arg) validate() error {
    if a.full == "" {
        return fmt.Errorf("arg name is empty")
    }
    if strings.HasPrefix(a.full, shortPrefix) || strings.HasPrefix(a.full, fullPrefix) {
        return fmt.Errorf("arg full name with extra prefix '%s'/'%s'", shortPrefix, fullPrefix)
    }
    if strings.HasPrefix(a.short, shortPrefix) {
        return fmt.Errorf("arg short name with extra prefix '%s'", shortPrefix)
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
    result := []string{fmt.Sprintf("%s%s", fullPrefix, a.full)}
    if a.short != "" && a.short != a.full {
        result = append([]string{fmt.Sprintf("%s%s", shortPrefix, a.short)}, result...)
    }
    return result
}

