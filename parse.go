package argparse

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Parser is the top level struct. Don't use it directly, use NewParser to create one
//
// it's the only struct to interact with user input, parse & bind each `arg` value
type Parser struct {
	name        string
	description string
	config      *ParserConfig

	Invoked      bool       // whether the parser is invoked
	InvokeAction func(bool) // execute after parse

	showHelp            *bool // flag to decide show help message
	showShellCompletion *bool // flag to decide show shell completion

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
	DefaultAction          func() // set default action to replace default help action
	AddShellCompletion     bool   // set true to register shell completion entry [--completion]
	WithHint               bool   // argument help message with argument default value hint
	MaxHeaderLength        int    // max argument header length in help menu, help info will start at new line if argument meta info is too long
}

// NewParser create the parser object with optional name & description & ParserConfig
func NewParser(name string, description string, config *ParserConfig) *Parser {
	if config == nil {
		config = &ParserConfig{}
	}
	if name == "" && len(os.Args) > 0 {
		name = strings.ReplaceAll(path.Base(os.Args[0]), " ", "") // avoid space for shell complete code generate
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
	for _, watcher := range a.getWatchers() { // register optional arguments to 'entryMap'
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
	headerLength := 10 // here set minimum header length, the code after will find the max length of headers
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
	headerLength += 4 // 2 space padding around header, before & after
	helpBreak := false
	if p.config.MaxHeaderLength > 0 {
		headerLength = p.config.MaxHeaderLength
		helpBreak = true
	}
	if len(p.subParser) > 0 {
		section := "\navailable commands:\n"
		for _, parser := range p.subParser {
			section += formatHelpRow(parser.name, parser.description, headerLength, helpBreak) + "\n"
		}
		result += section
	}
	withHint := p.config.WithHint
	if len(p.positionArgs) > 0 { // dealing positional arguments present
		section := "\npositional arguments:\n"
		for _, arg := range p.positionArgs {
			if arg.Group != "" || arg.HideEntry {
				continue
			}
			help := arg.Help
			if withHint && !arg.NoHint {
				help = arg.formatHelpWithExtraInfo()
			}
			section += formatHelpRow(arg.formatHelpHeader(), help, headerLength, helpBreak) + "\n"
		}
		result += section
	}
	if len(p.entries) > 0 { // dealing optional arguments present
		parsed := make(map[string]bool)
		section := "\noptions:\n"
		for _, arg := range p.entries {
			if arg.Group != "" {
				continue
			}
			if _, exist := parsed[arg.full]; exist {
				continue
			}
			parsed[arg.full] = true
			if arg.HideEntry {
				continue
			}
			help := arg.Help
			if withHint && !arg.NoHint {
				help = arg.formatHelpWithExtraInfo()
			}
			section += formatHelpRow(arg.formatHelpHeader(), help, headerLength, helpBreak) + "\n"
		}
		result += section
	}
	for _, group := range p.entryGroupOrder { // dealing arguments group present
		section := fmt.Sprintf("\n%s:\n", group)
		content := ""
		for _, arg := range p.entryGroup[group] {
			if arg.HideEntry {
				continue
			}
			content += formatHelpRow(arg.formatHelpHeader(), arg.Help, headerLength, helpBreak) + "\n"
		}
		if content != "" {
			result += section + content
		}
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
	for _, parent := range p.parentList { // sub command usage
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
		argUsage := arg.formatUsage()
		if argUsage != "" {
			usage += argUsage
		}
	}
	for _, arg := range p.positionArgs {
		argUsage := arg.formatUsage()
		if argUsage != "" {
			usage += argUsage
		}
	}
	return usage
}

// formatBashCompletionScript will generate bash shell script
func (p *Parser) formatBashCompletionScript() string {
	completionName := fmt.Sprintf("_%s_completion", p.name)
	var topLevel []string
	subLevelMap := make(map[string]string)
	for entry, arg := range p.entryMap {
		if arg.HideEntry {
			continue
		}
		topLevel = append(topLevel, entry)
	}
	for entry, subParser := range p.subParserMap {
		topLevel = append(topLevel, entry)
		var subOptions []string
		for subOption, arg := range subParser.entryMap {
			if arg.HideEntry {
				continue
			}
			subOptions = append(subOptions, subOption)
		}
		subLevelMap[entry] = strings.Join(subOptions, " ")
	}
	var subCompletions []string
	for entry, candidates := range subLevelMap {
		subCompletions = append(subCompletions,
			fmt.Sprintf("    %s) COMPREPLY=( $(compgen -W \"%s\" -- $cur ) ) ;;", entry, candidates))
	}
	subCompletionsScript := ""
	if len(subCompletions) > 0 {
		subCompletionsScript = fmt.Sprintf(`
    case "$cmd" in
  %s
    esac`, strings.Join(subCompletions, "\n"))
	}

	return fmt.Sprintf(`
  %s() {
    local i=1 cur="${COMP_WORDS[COMP_CWORD]}" cmd

    while [[ "$i" -lt "$COMP_CWORD" ]]
    do
      local s="${COMP_WORDS[i]}"
      case "$s" in
        %s*)
          cmd="$s"
          ;;
        *)
          cmd="$s"
          break
          ;;
      esac
      (( i++ ))
    done

    if [[ "$i" -eq "$COMP_CWORD" ]]
    then
      COMPREPLY=($(compgen -W "%s" -- $cur))
      return
    fi
  %s
  }

  complete -o bashdefault -o default -F %s %s
`, completionName, shortPrefix, strings.Join(topLevel, " "), subCompletionsScript, completionName, p.name)
}

// formatZshCompletionScript will generate zsh shell script
func (p *Parser) formatZshCompletionScript() string {
	completionName := fmt.Sprintf("_%s_completion", p.name)
	var positional []string
	var positionalFirstSection []string
	for entry, arg := range p.entryMap {
		if arg.HideEntry {
			continue
		}
		positional = append(positional, entry)
		positionalFirstSection = append(positionalFirstSection,
			fmt.Sprintf("\"%s\"", entry))
	}

	subLevelPosition := ""
	subLevelMap := make(map[string]string)
	for entry, subParser := range p.subParserMap {
		var subOptions []string
		for subOption, arg := range subParser.entryMap {
			if arg.HideEntry {
				continue
			}
			subOptions = append(subOptions, fmt.Sprintf("\"%s\"", subOption))
		}
		subLevelPosition += entry + " "
		subLevelMap[entry] = strings.Join(subOptions, " ")
	}
	var subCompletions []string
	for entry, candidates := range subLevelMap {
		subCompletions = append(subCompletions,
			fmt.Sprintf("    %s) _arguments %s ;;", entry, candidates))
	}
	subCompletionsScript := ""
	if len(subCompletions) > 0 {
		subCompletionsScript = fmt.Sprintf(`
    case $line[1] in
  %s
    esac`, strings.Join(subCompletions, "\n"))
	}

	return fmt.Sprintf(`
  function %s {
    local line
    _arguments -C %s "1: :(%s %s)" "*::arg:->args"
    %s
  }
  compdef %s %s
`, completionName, strings.Join(positionalFirstSection, " "), subLevelPosition, strings.Join(positional, " "), subCompletionsScript, completionName, p.name)
}

// FormatCompletionScript generate simple shell complete script, which support bash & zsh for completion
func (p *Parser) FormatCompletionScript() string {
	return fmt.Sprintf(`
###-begin-completion-###
# save the output to ~/.bashrc (or ~/.zshrc)
# or save file to your completion path like /usr/local/etc/bash_completion.d/ or /etc/bash_completion.d/
if type complete &>/dev/null; then
%s
elif type compctl &>/dev/null; then
%s
fi
###-end-completion-###
`, p.formatBashCompletionScript(), p.formatZshCompletionScript())
}

// Parse will parse given args to bind to any registered arguments
//
// args: set nil to use os.Args[1:] by default
func (p *Parser) Parse(args []string) (*Parser, error) {
	if args == nil {
		args = os.Args[1:]
	}
	matchSub := false
	if len(p.subParser) > 0 && len(args) > 0 {
		_, matchSub = p.subParserMap[args[0]]
	}
	var subParser *Parser
	if len(args) == 0 {
		if p.config.DefaultAction != nil {
			p.config.DefaultAction()
		} else if !p.config.DisableDefaultShowHelp {
			help := true
			p.showHelp = &help
		}
	} else {
		p.Invoked = true // when there is any match, it's invoked, or the default action will be called
		if matchSub {
			subParser = p.subParserMap[args[0]]
			var e error
			subParser, e = subParser.Parse(args[1:])
			if e != nil {
				return subParser, e
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
						var tillNext []string // find user inputs before next registered optional argument
						for _, a := range args[1:] {
							if _, isEntry := p.entryMap[a]; !isEntry {
								tillNext = append(tillNext, a)
							} else {
								break
							}
						}
						if len(tillNext) == 0 { // argument takes at least one input as argument, but there is 0
							return p, fmt.Errorf("argument %s expect argument",
								strings.Join(arg.getWatchers(), "/"))
						}
						if arg.multi { // if argument takes more than one arguments, it will take all user input before next registered argument, and proceed 'args' parsing to next registered argument
							e := arg.parseValue(tillNext)
							if e != nil {
								return p, e
							}
							args = args[len(tillNext)+1:]
						} else { // then the argument takes only one argument, and proceed the left arguments for positional argument parsing
							e := arg.parseValue(tillNext[0:1])
							if e != nil {
								return p, e
							}
							args = args[2:]
						}
					}
				} else {
					if registeredPositionsLength > lastPositionArgIndex { // while there is unparsed positional argument
						arg := p.positionArgs[lastPositionArgIndex]
						lastPositionArgIndex += 1
						var tillNext []string // find user inputs before next registered optional argument
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
								return p, e
							}
							args = args[len(tillNext):]
						} else {
							e := arg.parseValue(tillNext[0:1])
							if e != nil {
								return p, e
							}
							args = args[1:]
						}
					} else {
						if strings.HasPrefix(sign, shortPrefix) {
							var candidates []string
							for k := range p.entryMap {
								candidates = append(candidates, k)
							}
							var tips []string
							for _, m := range decideMatch(sign, candidates) {
								helpInfo := p.entryMap[m].Help
								if helpInfo != "" {
									helpInfo = fmt.Sprintf(" (%s)", helpInfo)
								}
								tips = append(tips, fmt.Sprintf("%s%s", m, helpInfo))
							}
							match := strings.Join(tips, "\nor ")
							if match != "" {
								return p, fmt.Errorf("unrecognized arguments: %s\ndo you mean?: %s", sign, match)
							}
						}
						return p, fmt.Errorf("unrecognized arguments: %s", sign)
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
			return targetParser, BreakAfterHelp{}
		}
	}
	if p.showShellCompletion != nil && *p.showShellCompletion {
		fmt.Println(p.FormatCompletionScript())
		return targetParser, BreakAfterShellScript{}
	}
	entries := append(p.entries, p.positionArgs...) // ready for Required checking & Default parsing
	for _, _p := range p.subParser {
		entries = append(entries, append(_p.entries, _p.positionArgs...)...)
	}
	for _, arg := range entries {
		if !arg.assigned && arg.Default != "" {
			if e := arg.parseValue(nil); e != nil {
				return targetParser, e
			}
		}
		if arg.Required && !arg.assigned {
			return targetParser, fmt.Errorf("%s is required", arg.getMetaName())
		}
	}
	if p.InvokeAction != nil {
		p.InvokeAction(p.Invoked)
	}
	return targetParser, nil
}

// AddCommand will add sub command entry parser
//
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
	config.AddShellCompletion = false // disable sub command completion
	parser := NewParser(name, description, config)
	parser.parentList = append(p.parentList, p.name)
	if e := p.registerParser(parser); e != nil {
		panic(e.Error())
	}
	return parser
}

// Flag create flag argument, Return a "*bool" point to the parse result
//
// python version is like add_argument("-s", "--full", action="store_true")
//
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

// String create string argument, return a "*string" point to the parse result
//
// String Argument can be used as Optional or Positional Arguments, default to be Optional, then it's like add_argument("-s", "--full") in python
//
// set Option.Positional = true to use as Positional Argument, then it's like add_argument("s", "full") in python
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

// Strings create string list argument, return a "*[]string" point to the parse result
//
// mostly like "*Parser.String()"
//
// python version is like add_argument("-s", "--full", nargs="*") or add_argument("s", "full", nargs="*")
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

// Int create int argument, return a *int point to the parse result
//
// mostly like *Parser.String(), except the return type
//
// python version is like add_argument("s", "full", type=int) or add_argument("-s", "--full", type=int)
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

// Ints create int list argument, return a *[]int point to the parse result
//
// mostly like *Parser.Int()
//
// python version is like add_argument("s", "full", type=int, nargs="*") or add_argument("-s", "--full", type=int, nargs="*")
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

// Float create float argument, return a *float64 point to the parse result
//
// mostly like *Parser.String(), except the return type
//
// python version is like add_argument("-s", "--full", type=double) or add_argument("s", "full", type=double)
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

// Floats create float list argument, return a *[]float64 point to the parse result
//
// mostly like *Parser.Float()
//
// python version is like add_argument("-s", "--full", type=double, nargs="*") or add_argument("s", "full", type=double, nargs="*")
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
