package argparse

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func decideTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return width
}

func formatHelpRow(head, content string, bareHeadLength, maxHeadLength, terminalWidth int, withBreak bool) string {
	content = strings.Replace(content, "\n", "", -1)

	head = "  " + head + " " // head content
	bareHeadLength += 3      // head length without control chars

	// length of a single content row, constant
	contentRowLen := terminalWidth - maxHeadLength

	min := func(x, y int) int {
		if x < y {
			return x
		} else {
			return y
		}
	}

	var rows []string
	if withBreak && maxHeadLength < bareHeadLength {
		rows = append(rows, head)
	} else {
		// no break -> head is on the same row
		// as first content line
		var headRowPadding string
		if maxHeadLength > bareHeadLength {
			headRowPadding = strings.Repeat(" ", maxHeadLength-bareHeadLength)
		}
		rowLen := min(contentRowLen, len(content))
		rows = append(rows, head+headRowPadding+content[:rowLen])
		content = content[rowLen:]
	}

	rowPadding := strings.Repeat(" ", maxHeadLength)
	for content != "" {
		rowLen := min(contentRowLen, len(content))
		rows = append(rows, rowPadding+content[:rowLen])
		content = content[rowLen:]
	}

	return strings.Join(rows, "\n")
}
