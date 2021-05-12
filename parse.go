package argparse

import (
	"fmt"
	"os"
	"strings"
)

// Parser is the top level struct
// it's the only interface to interact with user input, parse & bind each `arg` value
type Parser struct {
	name        string
	description string
	config      *ParserConfig

	showHelp            *bool // flag to decide show help message
	showShellCompletion *bool // flag to  decide show shell completion

	entries      []*arg
	entryMap     map[string]*arg
	positionArgs []*arg

	entryGroupOrder []string
	entryGroup      map[string][]*arg

	subParser    []*Parser
	subParserMap map[string]*Parser
	parentList   []string
}

// ParserConfig is the only type to config `Parser`, programmers only need to use this type to control `Parser` action
type ParserConfig struct {
	Usage                  string // manual usage display
	EpiLog                 string // message after help
	DisableHelp            bool   // disable help entry register [-h/--help]
	ContinueOnHelp         bool   // set true to: continue program after default help is printed
	DisableDefaultShowHelp bool   // set false to: default show help when there is no args to parse (default action)
	AddShellCompletion     bool   // set true to register shell completion entry [--completion]
}

// NewParser create the parser object with optional name & description & ParserConfig
func NewParser(name string, description string, config *ParserConfig) *Parser {
	if config == nil {
		config = &ParserConfig{}
	}
	if name == "" && len(os.Args) > 0 {
		name = os.Args[0]
	}
	parser := &Parser{
		name:            name,
		description:     description,
		config:          config,
		entries:         []*arg{},
		entryMap:        make(map[string]*arg),
		entryGroup:      make(map[string][]*arg),
		entryGroupOrder: []string{},
		positionArgs:    []*arg{},
		subParser:       []*Parser{},
		subParserMap:    make(map[string]*Parser),
	}
	if !config.DisableHelp {
		parser.showHelp = parser.Flag("h", "help",
			&Option{Help: "show this help message"})
	}
	if config.AddShellCompletion {
		parser.showShellCompletion = parser.Flag("", "completion",
			&Option{Help: "show command completion script"})
	}
	return parser
}

func (p *Parser) registerArgument(a *arg) error {
	e := a.validate()
	if e != nil {
		return e
	}
	if a.Positional {
		p.positionArgs = append(p.positionArgs, a)
	}
	if a.Group != "" {
		if _, exist := p.entryGroup[a.Group]; !exist {
			p.entryGroupOrder = append(p.entryGroupOrder, a.Group)
		}
		p.entryGroup[a.Group] = append(p.entryGroup[a.Group], a)
	}
	for _, watcher := range a.getWatchers() {
		if match, exist := p.entryMap[watcher]; exist {
			return fmt.Errorf("conflict entry for '%s', say: '%s'", watcher, match.Help)
		}
		p.entryMap[watcher] = a
		p.entries = append(p.entries, a)
	}
	return nil
}

func (p *Parser) registerParser(parser *Parser) error {
	if match, exist := p.subParserMap[parser.name]; exist {
		return fmt.Errorf("conflict sub command for '%s', desc: '%s'",
			parser.name, match.description)
	}
	p.subParser = append(p.subParser, parser)
	p.subParserMap[parser.name] = parser
	return nil
}

// PrintHelp print help message to stdout
func (p *Parser) PrintHelp() {
	fmt.Println(p.FormatHelp())
}

// FormatHelp only format help message for manual use, for example: decide when to print help message
func (p *Parser) FormatHelp() string {
	result := p.formatUsage()
	if p.description != "" {
		result += "\n\n" + p.description + "\n"
	}
	headerLength := 10
	for _, parser := range p.subParser {
		l := len(parser.name)
		if l > headerLength {
			headerLength = l
		}
	}
	for _, arg := range p.positionArgs {
		l := len(arg.formatHelpHeader())
		if l > headerLength {
			headerLength = l
		}
	}
	for _, arg := range p.entries {
		l := len(arg.formatHelpHeader())
		if l > headerLength {
			headerLength = l
		}
	}
	headerLength += 4
	if len(p.subParser) > 0 {
		section := "\navailable commands:\n"
		for _, parser := range p.subParser {
			section += formatHelpRow(parser.name, parser.description, headerLength) + "\n"
		}
		result += section
	}
	if len(p.positionArgs) > 0 {
		section := "\npositional arguments:\n"
		for _, arg := range p.positionArgs {
			if arg.Group != "" {
				continue
			}
			section += formatHelpRow(arg.formatHelpHeader(), arg.Help, headerLength) + "\n"
		}
		result += section
	}
	if len(p.entries) > 0 {
		parsed := make(map[string]bool)
		section := "\noptional arguments:\n"
		for _, arg := range p.entries {
			if arg.Group != "" {
				continue
			}
			if _, exist := parsed[arg.full]; exist {
				continue
			}
			parsed[arg.full] = true
			section += formatHelpRow(arg.formatHelpHeader(), arg.Help, headerLength) + "\n"
		}
		result += section
	}
	for _, group := range p.entryGroupOrder {
		section := fmt.Sprintf("\n%s:\n", group)
		for _, arg := range p.entryGroup[group] {
			section += formatHelpRow(arg.formatHelpHeader(), arg.Help, headerLength) + "\n"
		}
		result += section
	}
	if p.config.EpiLog != "" {
		result += "\n" + p.config.EpiLog
	}

	return result
}

func (p *Parser) formatUsage() string {
	usage := "usage: "
	if p.config.Usage != "" {
		return usage + p.config.Usage
	}
	for _, parent := range p.parentList {
		usage += parent + " "
	}
	usage += p.name + " "
	if len(p.subParser) > 0 {
		usage += "<cmd> "
	}
	parsed := make(map[string]bool)
	for _, arg := range p.entries {
		if _, exist := parsed[arg.full]; exist {
			continue
		}
		parsed[arg.full] = true
		sign := arg.getWatchers()[0]
		if arg.isFlag {
			usage += fmt.Sprintf("[%s] ", sign)
		} else {
			meta := arg.getMetaName()
			u := fmt.Sprintf("%s %s", sign, meta)
			if arg.Required {
				usage += u + " "
				if arg.multi {
					usage += fmt.Sprintf("[%s ...] ", meta)
				}
			} else {
				if arg.multi {
					usage += fmt.Sprintf("[%s [%s ...]] ", u, meta)
				} else {
					usage += fmt.Sprintf("[%s] ", u)
				}
			}
		}
	}
	for _, arg := range p.positionArgs {
		meta := arg.getMetaName()
		if arg.Required {
			usage += meta + " "
			if arg.multi {
				usage += fmt.Sprintf("[%s ...] ", meta)
			}
		} else {
			if arg.multi {
				usage += fmt.Sprintf("[%s [%s ...]] ", meta, meta)
			} else {
				usage += fmt.Sprintf("[%s] ", meta)
			}
		}
	}
	return usage
}

// Parse will parse given args to bind to any registered arguments
// args: set nil to use os.Args[1:] by default
func (p *Parser) Parse(args []string) error {
	if args == nil {
		args = os.Args[1:]
	}
	matchSub := false
	if len(p.subParser) > 0 && len(args) > 0 {
		_, matchSub = p.subParserMap[args[0]]
	}
	var subParser *Parser
	if len(args) == 0 {
		if !p.config.DisableDefaultShowHelp {
			help := true
			p.showHelp = &help
		}
	} else {
		if matchSub {
			subParser = p.subParserMap[args[0]]
			e := subParser.Parse(args[1:])
			if e != nil {
				return e
			}
		} else {
			lastPositionArgIndex := 0
			registeredPositionsLength := len(p.positionArgs)
			for len(args) > 0 {
				sign := args[0]
				if arg, ok := p.entryMap[sign]; ok {
					if arg.isFlag {
						_ = arg.parseValue(nil)
						args = args[1:]
					} else {
						var tillNext []string
						for _, a := range args[1:] {
							if _, isEntry := p.entryMap[a]; !isEntry {
								tillNext = append(tillNext, a)
							} else {
								break
							}
						}
						if len(tillNext) == 0 {
							return fmt.Errorf("argument %s expect argument",
								strings.Join(arg.getWatchers(), "/"))
						}
						if arg.multi {
							e := arg.parseValue(tillNext)
							if e != nil {
								return e
							}
							args = args[len(tillNext)+1:]
						} else {
							e := arg.parseValue(tillNext[0:1])
							if e != nil {
								return e
							}
							args = args[2:]
						}
					}
				} else {
					if registeredPositionsLength > lastPositionArgIndex {
						arg := p.positionArgs[lastPositionArgIndex]
						lastPositionArgIndex += 1
						var tillNext []string
						for _, a := range args {
							if _, isEntry := p.entryMap[a]; !isEntry {
								tillNext = append(tillNext, a)
							} else {
								break
							}
						}
						if arg.multi {
							e := arg.parseValue(tillNext)
							if e != nil {
								return e
							}
							args = args[len(tillNext):]
						} else {
							e := arg.parseValue(tillNext[0:1])
							if e != nil {
								return e
							}
							args = args[1:]
						}
					} else {
						return fmt.Errorf("unrecognized arguments: %s", sign)
					}
				}
			}
		}
	}
	targetParser := p
	if subParser != nil {
		targetParser = subParser
	}
	if targetParser.showHelp != nil && *targetParser.showHelp {
		targetParser.PrintHelp()
		if !targetParser.config.ContinueOnHelp {
			os.Exit(1)
		}
	}
	entries := append(p.entries, p.positionArgs...)
	for _, _p := range p.subParser {
		entries = append(entries, append(_p.entries, _p.positionArgs...)...)
	}
	for _, arg := range entries {
		if !arg.assigned && arg.Default != "" {
			if e := arg.parseValue(nil); e != nil {
				return e
			}
		}
		if arg.Required && !arg.assigned {
			return fmt.Errorf("%s is required", arg.getMetaName())
		}
	}
	return nil
}

// AddCommand will add sub command entry parser
// Return a new pointer to sub command parser
func (p *Parser) AddCommand(name string, description string, config *ParserConfig) *Parser {
	if config == nil {
		config = p.config
	}
	if name == "" {
		panic("sub command name is empty")
	}
	if strings.Contains(name, " ") {
		panic("sub command name has space")
	}
	parser := NewParser(name, description, config)
	parser.parentList = append(p.parentList, p.name)
	if e := p.registerParser(parser); e != nil {
		panic(e.Error())
	}
	return parser
}

// Flag create flag argument, Return a `*bool` point to the parse result
// python version is like `add_argument("-s", "--full", action="store_true")`
// Flag Argument can only be used as an OptionalArguments
func (p *Parser) Flag(short, full string, opts *Option) *bool {
	var result bool
	if opts == nil {
		opts = &Option{}
	}
	opts.isFlag = true
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// String create string argument, return a `*string` point to the parse result
// String Argument can be used as Optional or Positional Arguments, default to be Optional, then it's like `add_argument("-s", "--full")` in python
// set `Option.Positional = true` to use as Positional Argument, then it's like `add_argument("s", "full")` in python
func (p *Parser) String(short, full string, opts *Option) *string {
	var result string
	if opts == nil {
		opts = &Option{}
	}
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// Strings create string list argument, return a `*[]string` point to the parse result
// mostly like `*Parser.String()`
// python version is like `add_argument("-s", "--full", nargs="*")` or `add_argument("s", "full", nargs="*")`
func (p *Parser) Strings(short, full string, opts *Option) *[]string {
	var result []string
	if opts == nil {
		opts = &Option{}
	}
	opts.multi = true
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// Int create int argument, return a `*int` point to the parse result
// mostly like `*Parser.String()`, except the return type
// python version is like `add_argument("s", "full", type=int)` or `add_argument("-s", "--full", type=int)`
func (p *Parser) Int(short, full string, opts *Option) *int {
	var result int
	if opts == nil {
		opts = &Option{}
	}
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// Ints create int list argument, return a `*[]int` point to the parse result
// mostly like `*Parser.Int()`
// python version is like `add_argument("s", "full", type=int, nargs="*")` or `add_argument("-s", "--full", type=int, nargs="*")`
func (p *Parser) Ints(short, full string, opts *Option) *[]int {
	var result []int
	if opts == nil {
		opts = &Option{}
	}
	opts.multi = true
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// Float create float argument, return a `*float64` point to the parse result
// mostly like `*Parser.String()`, except the return type
// python version is like `add_argument("-s", "--full", type=double)` or `add_argument("s", "full", type=double)`
func (p *Parser) Float(short, full string, opts *Option) *float64 {
	var result float64
	if opts == nil {
		opts = &Option{}
	}
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}

// Floats create float list argument, return a `*[]float64` point to the parse result
// mostly like `*Parser.Float()`
// python version is like `add_argument("-s", "--full", type=double, nargs="*")` or `add_argument("s", "full", type=double, nargs="*")`
func (p *Parser) Floats(short, full string, opts *Option) *[]float64 {
	var result []float64
	if opts == nil {
		opts = &Option{}
	}
	opts.multi = true
	if e := p.registerArgument(&arg{
		short:  short,
		full:   full,
		target: &result,
		Option: *opts,
	}); e != nil {
		panic(e.Error())
	}
	return &result
}
