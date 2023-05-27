package argparse

import (
	"fmt"
	"os"
	"strings"
)

type Color struct {
	Code     int // color code for text color or background color, normally within 30 ~ 49
	Property int // property code, eg: 1 for bold, 4 for underline, etc.
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

// Just black & white
var NoColor = &ColorSchema{}

// Default color schema
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
