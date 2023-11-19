package argparse

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func decideTerminalWidth() int {
	// decide terminal width
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin // this is important
	result, e := cmd.Output()
	if e != nil {
		result = []byte("0 80")
	}
	parse := strings.Split(strings.TrimRight(string(result), "\n"), " ")
	if w, e := strconv.Atoi(parse[1]); e == nil {
		return w
	}
	return 80
}

func formatHelpRow(head, content string, bareHeadLength, maxHeadLength, terminalWidth int, withBreak bool) string {
	content = strings.Replace(content, "\n", "", -1)
	result := fmt.Sprintf("  %s ", head)
	headLeftPadding := maxHeadLength - bareHeadLength - 3
	if headLeftPadding > 0 { // fill left padding
		result += strings.Repeat(" ", headLeftPadding)
	}
	contentPadding := strings.Repeat(" ", maxHeadLength)
	var rows []string
	if withBreak && headLeftPadding < 0 {
		rows = append(rows, result, contentPadding+content)
	} else {
		rows = append(rows, result+content)
	}
	for len(rows[len(rows)-1]) > terminalWidth { // break into lines
		lastIndex := len(rows) - 1
		lastOne := rows[lastIndex]
		rows[lastIndex] = rows[lastIndex][0:terminalWidth]
		rows = append(rows, contentPadding+lastOne[terminalWidth:])
	}
	return strings.Join(rows, "\n")
}
