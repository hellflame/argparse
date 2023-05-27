package argparse

import (
	"fmt"
	"strconv"
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

// Option is the only type to config when creating argument
type Option struct {
	Meta        string                                // meta value for help/usage generate
	multi       bool                                  // take more than one argument
	Default     string                                // default argument value if not given
	isFlag      bool                                  // use as flag
	Required    bool                                  // require to be set
	Positional  bool                                  // is positional argument
	HideEntry   bool                                  // hide usage & help display
	Help        string                                // help message
	NoHint      bool                                  // disable hint info for the argument (usable when ParserConfig.WithHint = true)
	HintInfo    string                                // specified hint info suffixed after help message, when ParserConfig.WithHint = true
	Group       string                                // argument group info, default to be no group
	Inheritable bool                                  // sub parsers after this argument can inherit it
	Action      func(args []string) error             // bind actions when the match is found, 'args' can be nil to be a flag
	Choices     []interface{}                         // input argument must be one/some of the choice
	Validate    func(arg string) error                // customize function to check argument validation
	Formatter   func(arg string) (interface{}, error) // format input arguments by the given method
	BindParsers []*Parser                             // specify parsers to bind
}

// validate args setting before parsing args, right after adding to parser
// for conflict check & correction & restriction
func (a *arg) validate() error {
	if a.full == "" && a.short == "" { // argument must has a name
		return fmt.Errorf("arg name is empty")
	}
	if strings.Contains(a.full, " ") || strings.Contains(a.short, " ") { // space will interrupt
		return fmt.Errorf("arg name with space")
	}
	if strings.HasPrefix(a.full, shortPrefix) || strings.HasPrefix(a.full, fullPrefix) { // argument sign with be auto prefixed
		return fmt.Errorf("arg full name with extra prefix '%s'/'%s'", shortPrefix, fullPrefix)
	}
	if strings.HasPrefix(a.short, shortPrefix) { // argument will be auto prefixed
		return fmt.Errorf("arg short name with extra prefix '%s'", shortPrefix)
	}
	if a.short == a.full { // this will cause register conflict
		return fmt.Errorf("arg short is full")
	}
	if a.isFlag {
		if a.Positional { // positional argument can't be a flag, use flag instead
			return fmt.Errorf("positional is a flag")
		}
		if a.Meta != "" { // flag has no meta info to show
			return fmt.Errorf("flag with meta")
		}
		if len(a.Choices) != 0 { // if argument has a flag, it only has true as a choice
			return fmt.Errorf("flag has choices")
		}
		if a.Required { // if flag is a must, the result must be true
			return fmt.Errorf("flag with required")
		}
		if a.Formatter != nil { // flag has no need to reformat
			return fmt.Errorf("flag with formatter")
		}
		if a.Validate != nil { // flag has no need to be validated
			return fmt.Errorf("flag with validate")
		}
	}
	return nil
}

// get argument watch list for parser use
func (a *arg) getWatchers() []string {
	if a.Positional { // positional argument has nothing to watch, only positions
		return []string{}
	}
	var result []string
	if a.full != "" {
		result = append(result, fmt.Sprintf("%s%s", fullPrefix, a.full))
	}
	if a.short != "" {
		result = append(result, fmt.Sprintf("%s%s", shortPrefix, a.short))
	}
	return result
}

func (a *arg) getMetaName() string {
	if a.Meta != "" {
		return a.Meta // Meta variable given by programmer
	}
	return strings.ToUpper(a.getIdentifier())
}

func (a *arg) formatUsage() string {
	if a.HideEntry {
		return ""
	}

	meta := a.getMetaName()
	if a.Positional { // positional usage
		if a.Required {
			if a.multi {
				return fmt.Sprintf("%s [%s ...] ", meta, meta)
			}
			return fmt.Sprintf("%s ", meta)
		}
		if a.multi {
			return fmt.Sprintf("[%s [%s ...]] ", meta, meta)
		}
		return fmt.Sprintf("[%s] ", meta)
	}

	// other optional usage
	sign := a.getWatchers()[0]
	if a.isFlag {
		return fmt.Sprintf("[%s] ", sign)
	}
	u := fmt.Sprintf("%s %s", sign, meta)
	if a.Required {
		if a.multi {
			return fmt.Sprintf("%s [%s ...] ", u, meta)
		}
		return fmt.Sprintf("%s ", u)
	}
	if a.multi {
		return fmt.Sprintf("[%s [%s ...]] ", u, meta)
	}
	return fmt.Sprintf("[%s] ", u)
}

func (a *arg) getIdentifier() string {
	if a.full != "" {
		return a.full
	}
	return a.short
}

func (a *arg) formatHelpHeader(argument, meta Color) (size int, content string) {
	metaName := a.getMetaName()
	if a.Positional {
		size = len(metaName)
		content = wrapperColor(metaName, argument)
		return
	}

	wrapped := []string{}
	watchers := a.getWatchers()
	for _, w := range watchers {
		if a.isFlag {
			wrapped = append(wrapped, wrapperColor(w, argument))
			size += len(w)
		} else {
			wrapped = append(wrapped,
				fmt.Sprintf("%s %s", wrapperColor(w, argument),
					wrapperColor(metaName, meta)))
			size += len(w) + len(metaName) + 1
		}
	}
	size += (len(wrapped) - 1) * 2
	content = strings.Join(wrapped, ", ")
	return
}

func (a *arg) formatHelpWithExtraInfo() string {
	help := a.Help
	if help != "" {
		help += " " // append a space after help info
	}
	if a.HintInfo != "" {
		return fmt.Sprintf("%s(%s)", help, a.HintInfo)
	}
	var extraInfo []string
	if a.Default != "" {
		extraInfo = append(extraInfo, fmt.Sprintf("default: %s", a.Default))
	}
	if len(a.Choices) > 0 {
		extraInfo = append(extraInfo, fmt.Sprintf("options: [%s]", a.dumpChoices()))
	}
	if a.Required {
		extraInfo = append(extraInfo, "required")
	}
	if extra := strings.Join(extraInfo, ", "); extra != "" {
		return fmt.Sprintf("%s(%s)", help, extra)
	}
	return help
}

func (a *arg) dumpChoices() string {
	var choices []string
	e := a.Choices[0]
	switch e.(type) {
	case string:
		for _, e := range a.Choices {
			choices = append(choices, e.(string))
		}
	case int:
		for _, e := range a.Choices {
			choices = append(choices, fmt.Sprintf("%d", e.(int)))
		}
	case float64:
		for _, e := range a.Choices {
			choices = append(choices, fmt.Sprintf("%f", e.(float64)))
		}
	}
	return strings.Join(choices, ", ")
}

// parse input & bind (default) value to target
func (a *arg) parseValue(values []string) error {
	a.assigned = true
	if a.Action != nil {
		return a.Action(values)
	}
	if a.isFlag {
		*a.target.(*bool) = true
		return nil
	}
	if len(values) == 0 && a.Default != "" {
		values = append(values, a.Default) // add default value in the parse flow
	}
	if a.Validate != nil { // execute user given Validate function for each input
		for _, v := range values {
			e := a.Validate(v)
			if e != nil {
				return e
			}
		}
	}
	var result []interface{}
	if a.Formatter != nil { // format each input
		for _, v := range values {
			f, e := a.Formatter(v)
			if e != nil {
				return e
			}
			result = append(result, f)
		}
	} else {
		switch a.target.(type) {
		case *string, *[]string:
			for _, v := range values {
				result = append(result, v)
			}
		case *int, *[]int:
			for _, raw := range values {
				v, e := strconv.Atoi(raw)
				if e != nil {
					return fmt.Errorf("invalid int value: %s", raw)
				}
				result = append(result, v)
			}
		case *float64, *[]float64:
			for _, raw := range values {
				v, e := strconv.ParseFloat(raw, 64)
				if e != nil {
					return fmt.Errorf("invalid float value: %s", raw)
				}
				result = append(result, v)
			}
		}
	}
	//if len(result) == 0 {
	//	return fmt.Errorf("no value to parse") // normally you can't reach this area
	//}
	if len(a.Choices) > 0 { // check if user input is among given Choices
		for _, r := range result {
			found := false
			for _, c := range a.Choices {
				if c == r {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("args must be one|some of %+v", a.Choices)
			}
		}
	}
	switch a.target.(type) { // bind different types
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
	case *[]int:
		for _, r := range result {
			*a.target.(*[]int) = append(*a.target.(*[]int), r.(int))
		}
	case *[]float64:
		for _, r := range result {
			*a.target.(*[]float64) = append(*a.target.(*[]float64), r.(float64))
		}
	}
	return nil
}
