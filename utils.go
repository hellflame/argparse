package argparse

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const fallbackWidth = 80

func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin // this is important
	result, err := cmd.Output()
	if err == nil {
		parse := strings.Split(strings.TrimRight(string(result), "\n"), " ")
		w, e := strconv.Atoi(parse[1])
		if e == nil {
			return w
		}
	}
	return fallbackWidth
}

func formatHelpRow(head, content string, maxHeadLength int) string {
	terminalWidth := getTerminalWidth()
	content = strings.Replace(content, "\n", "", -1)
	result := fmt.Sprintf("  %s ", head)
	headLeftPadding := maxHeadLength - len(result)
	if headLeftPadding > 0 {
		result += strings.Repeat(" ", headLeftPadding)
	}
	contentPadding := strings.Repeat(" ", maxHeadLength)
	rows := []string{result + content}
	for len(rows[len(rows)-1]) > terminalWidth {
		lastIndex := len(rows) - 1
		lastOne := rows[lastIndex]
		if len(lastOne) < terminalWidth {
			break
		}
		rows[lastIndex] = rows[lastIndex][0:terminalWidth]
		rows = append(rows, contentPadding+lastOne[terminalWidth:])
	}
	return strings.Join(rows, "\n")
}
