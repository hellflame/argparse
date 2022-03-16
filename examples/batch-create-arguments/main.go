// this is show case for creating arguments by batch
package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

func main() {
	p := argparse.NewParser("", "example to batch create arguments", nil)

	nameArgs := make(map[string]*string)
	for _, n := range []string{"firstName", "midName", "lastName"} {
		nameArgs[n] = p.String("", n, &argparse.Option{Help: "this is part of your name", Group: "Names"})
	}

	statusArgs := make(map[string]*bool)
	for _, s := range []string{"employed", "graduated", "merried", "divorced", "underage", "crazy", "happy"} {
		statusArgs[s] = p.Flag("", s, &argparse.Option{Help: fmt.Sprintf("are you %s?", s)})
	}

	if e := p.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
			return
		default:
			fmt.Println(e.Error())
			return
		}
	}

	for arg, bind := range nameArgs {
		fmt.Println(arg, *bind)
	}
	for arg, bind := range statusArgs {
		fmt.Println(arg, *bind)
	}
}
