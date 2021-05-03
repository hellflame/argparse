package argparse

import (
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	name        string
	description string
	showHelp    bool
	config      *ParserConfig

	entries      []*arg
	entryMap     map[string]*arg
	positionArgs []*arg

	entryGroupOrder []string
	entryGroup      map[string][]*arg

	subParser    []*Parser
	subParserMap map[string]*Parser
}

type ParserConfig struct {
	Usage                  string
	EpiLog                 string
	DisableHelp            bool
	ContinueOnHelp         bool
	DisableDefaultShowHelp bool
}

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
		e := parser.registerArgument(&arg{short: "h", full: "help", target: parser.showHelp,
			Option: Option{isFlag: true, Help: "show this help message"}})
		if e != nil {
			panic(e.Error())
		}
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

func (p *Parser) PrintHelp() {
	fmt.Println(p.FormatHelp())
}

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

// args: set nil to use os.Args[1:] by default
func (p *Parser) Parse(args []string) error {
	if args == nil {
		args = os.Args[1:]
	}
	var subParser *Parser
	if len(args) == 0 {
		if !p.config.DisableDefaultShowHelp {
			p.showHelp = true
		}
	} else {
		if len(p.subParser) > 0 {
			if sub, ok := p.subParserMap[args[0]]; ok {
				subParser = sub
				e := sub.Parse(args[1:])
				if e != nil {
					return e
				}
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
						breakPoint := len(args) - 1
						for idx, a := range args[1:] {
							if _, isEntry := p.entryMap[a]; !isEntry {
								tillNext = append(tillNext, a)
							} else {
								breakPoint = idx
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
						} else {
							e := arg.parseValue(tillNext[0:1])
							if e != nil {
								return e
							}
						}
						args = args[breakPoint+1:]
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
	if targetParser.showHelp {
		targetParser.PrintHelp()
		if !targetParser.config.ContinueOnHelp {
			os.Exit(1)
		}
	}
	return nil
}

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
	if e := p.registerParser(parser); e != nil {
		panic(e.Error())
	}
	return parser
}

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
