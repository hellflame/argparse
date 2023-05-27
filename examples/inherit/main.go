// inheritable arguments

package main

import (
	"fmt"
	"os"

	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("", "", nil)
	verbose := parser.Flag("v", "", &argparse.Option{Inheritable: true,
		Help: "show verbose info"})
	local := parser.AddCommand("local", "", nil)
	service := parser.AddCommand("service", "", nil)

	user := local.String("u", "user", nil)
	pwd := local.String("p", "password", nil)

	addr := service.String("", "address", nil)
	port := service.Int("", "port", nil)
	version := service.Int("v", "version", &argparse.Option{Help: "version choice"})

	if e := parser.Parse(nil); e != nil {
		switch e.(type) {
		case argparse.BreakAfterHelp:
			os.Exit(1)
		default:
			fmt.Println(e.Error())
		}
		return
	}
	if local.Invoked {
		fmt.Println("local run", *verbose, *user, *pwd)
	}
	if service.Invoked {
		fmt.Println("service run", *verbose, *addr, *port, *version)
	}
}
