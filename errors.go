package argparse

// BreakAfterHelp will be thrown after help showed
type BreakAfterHelp struct {
}

func (b BreakAfterHelp) Error() string {
	return ""
}

// BreakAfterShellScript will be thrown after shell script showed
type BreakAfterShellScript struct {
}

func (b BreakAfterShellScript) Error() string {
	return ""
}
