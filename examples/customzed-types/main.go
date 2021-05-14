// this is show case for some advanced use
//
// use Option.Validate to check if the user input is valid
//
// use Option.Formatter to pre-format the user input, before binding to the variable 'host'
//
// if there's more free-type argument need, check out another example 'any-type-action'
package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"io/ioutil"
	"net/url"
	"os"
)

func main() {
	parser := argparse.NewParser("basic", "this is a basic program", nil)
	path := parser.String("f", "file", &argparse.Option{
		Validate: func(arg string) error {
			if _, e := os.Stat(arg); e != nil {
				return fmt.Errorf("unable to access '%s'", arg)
			}
			return nil
		},
	})
	host := parser.String("u", "url", &argparse.Option{
		Help: "give me some url to parse",
		Formatter: func(arg string) (i interface{}, err error) {
			u, err := url.ParseRequestURI(arg)
			if err != nil {
				return
			}
			i = u.Host
			return
		},
	})

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}

	if *path != "" {
		if read, e := ioutil.ReadFile(*path); e == nil {
			fmt.Println(string(read))
		}
	}
	if *host != "" {
		fmt.Println(*host)
	}

}
