package argparse

import (
	"strings"
	"testing"
)

func TestFormatHelpRow(t *testing.T) {
	if strings.Count(formatHelpRow("this is header", strings.Repeat("C", 50), 30, false),
		"\n") > 0 {
		t.Error("should be only one line")
		return
	}
	if strings.Count(formatHelpRow("this is header", strings.Repeat("C", 51), 30, false),
		"\n") != 1 {
		t.Error("should be exactly one line")
		return
	}
	if strings.Count(formatHelpRow("short_e", strings.Repeat("C", 51), 10, true), "\n") > 0 {
		t.Error("should be exactly one line")
		return
	}
	if strings.Count(formatHelpRow("short_f", strings.Repeat("C", 91), 10, true), "\n") != 1 {
		t.Error("should be exactly two line")
		return
	}
	if strings.Count(formatHelpRow("the_long", strings.Repeat("C", 51), 10, true), "\n") != 1 {
		t.Error("two line")
		return
	}
	if strings.Count(formatHelpRow("this is very long too", strings.Repeat("C", 91), 10, true), "\n") != 2 {
		t.Error("line break error")
		return
	}
}
