package main

import (
	"flag"
	"fmt"
	"github.com/boredomist/gob/parse"
	"os"
)

const GOB_VERSION = "0.0.0"

var (
	showVersion = flag.Bool("version", false, "Show version info")
	// TODO: other
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Gob v%s\n", GOB_VERSION)
	}

	for _, name := range flag.Args() {
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
