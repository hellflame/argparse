package argparse

import (
    "strings"
    "testing"
)

func TestFormatHelpRow(t *testing.T) {
    if strings.Count(formatHelpRow("this is header", strings.Repeat("C", 50), 30),
        "\n") > 0 {
        t.Error("should be only one line")
        return
    }
    if strings.Count(formatHelpRow("this is header", strings.Repeat("C", 51), 30),
        "\n") != 1 {
        t.Error("should be exactly one line")
        return
    }
}
