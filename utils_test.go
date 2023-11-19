package argparse

import (
	"strings"
	"testing"
)

func TestFormatHelpRow(t *testing.T) {
	width := decideTerminalWidth()
	header := "this is header"
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 50), len(header), 30, width, false),
		"\n") > 0 {
		t.Error("should be only one line")
		return
	}
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 51), len(header), 30, width, false),
		"\n") != 1 {
		t.Error("should be exactly one line")
		return
	}
	header = "short_e"
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 51), len(header), 10, width, true), "\n") > 0 {
		t.Error("should be exactly one line")
		return
	}
	header = "short_f"
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 91), len(header), 10, width, true), "\n") != 1 {
		t.Error("should be exactly two line")
		return
	}
	header = "the_long"
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 51), len(header), 10, width, true), "\n") != 1 {
		t.Error("two line")
		return
	}
	header = "this is very long too"
	if strings.Count(formatHelpRow(header, strings.Repeat("C", 91), len(header), 10, width, true), "\n") != 2 {
		t.Error("line break error")
		return
	}
}
