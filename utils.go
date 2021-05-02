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
	width, _ := strconv.Atoi(os.Getenv("COLUMNS"))
	if width <= 0 {
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
	return width
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
	cnt := 1
	for len(rows[len(rows)-1]) > terminalWidth {
		cnt += 1
		lastOne := rows[len(rows)-1]
		if len(lastOne) < terminalWidth {
			break
		}
		rows = append(rows, contentPadding+lastOne[terminalWidth:])
	}
	return strings.Join(rows, "\n")
}
