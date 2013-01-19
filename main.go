package main

import (
	"fmt"
	"github.com/boredomist/gob/parse"
	opt "github.com/droundy/goopt"
	"os"
)

const GOB_VERSION = "0.0.0"

var (
	showVersion = opt.Flag([]string{"-v", "--version"}, []string{},
		"Show version info", "")
	// TODO: other
)

func main() {
	opt.Parse(nil)

	if *showVersion {
		fmt.Printf("Gob v%s\n", GOB_VERSION)
	}

	for _, name := range opt.Args {
		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parser := parse.NewParser(name, file)

		unit, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(unit)
		}
	}

}
