// this is show case for creating argument groups
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"os"
)

func main() {
	p := argparse.NewParser("", "this is a show case about groups", &argparse.ParserConfig{DisableHelp: true, EpiLog: "try more"})
	p.Flag("n", "normal", nil)
	p.Float("f", "float", &argparse.Option{Positional: true})

	p.String("a", "aa", &argparse.Option{Group: "As", Help: "normal a"})
	p.String("aaa", "", &argparse.Option{Group: "As", Help: "triple a", Positional: true})

	p.Int("b", "bb", &argparse.Option{Group: "Bs", Help: "normal b"})
	p.Ints("", "bbb", &argparse.Option{Group: "Bs", Help: "triple b"})

	help := p.Flag("h", "", &argparse.Option{Group: "General", Help: "show help info"})
	if e := p.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
			os.Exit(1)
		default:
			fmt.Println(e.Error())
			return
		}
	}
	if *help {
		p.PrintHelp()
		return
	}
}
