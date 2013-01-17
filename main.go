package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	showVersion = flag.Bool("version", false, "Show version info")
	// TODO: other
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("Gob 0.0.0")
	}

	for _, name := range flag.Args() {
		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parser := NewParser(name, file)

		node, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(*node)
		}
	}

}
