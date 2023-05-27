package argparse

import (
	"fmt"
	"os"
	"strings"
)

type Color struct {
	Code     int
	Property int
}

type ColorSchema struct {
	Usage       Color
	Description Color

	GroupTitle Color
	Command    Color

	Argument Color
	Meta     Color

	Epilog Color
}

var NoColor = &ColorSchema{}
var DefaultColor = &ColorSchema{
	Usage:      Color{37, 1},
	GroupTitle: Color{32, 1},
	Command:    Color{33, 1},
	Argument:   Color{36, 0},
	Epilog:     Color{37, 1},
}

func checkTerminalColorSupport() bool {
	return strings.Contains(os.Getenv("TERM"), "color")
}

func wrapperColor(content string, color Color) string {
	if color.Code == 0 && color.Property == 0 {
		return content
	}
	return fmt.Sprintf("\033[0%d;%dm%s\033[00m", color.Property, color.Code, content)
}
