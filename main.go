package main

import (
	"fmt"
	opt "github.com/droundy/goopt"
	"github.com/erik/gob/emit"
	"github.com/erik/gob/parse"
	"os"
)

const GOB_VERSION = "0.0.0"

var (
	showVersion = opt.Flag([]string{"-v", "--version"}, []string{},
		"Show version info", "")
	parseOnly = opt.Flag([]string{"-p", "--parse-only"}, []string{},
		"Don't output anything, just parse", "")
	outFile = opt.Flag([]string{"-o"}, []string{}, "Specify output file", "")
)

func main() {
	opt.Parse(nil)

	if *showVersion {
		fmt.Printf("Gob v%s\n", GOB_VERSION)
		return
	}

	if len(opt.Args) < 1 {
		fmt.Println("Need to specify an input file")
		return
	}

	for _, name := range opt.Args {
		if len(opt.Args) > 1 {
			fmt.Printf("==== %s ====\n", name)
		}

		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parser := parse.NewParser(name, file)

		unit, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		}

		if err = unit.Verify(); err != nil {
			fmt.Println(err)
		}

		if *parseOnly {
			continue
		}

		var emit emit.CEmitter
		fmt.Println(emit.Emit(unit))
	}
}
