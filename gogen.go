package main

import (
	"fmt"
	"gogen/dirs"
	"gogen/parser"
	"os"
)

func main() {
	args := os.Args[1:]
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("problem with accessing working directory")
		os.Exit(2)
	}

	var pack string
	if len(args) == 0 {
		pack = dirs.PackImport(wd)
	} else {
		pack = args[0]
	}

	parser.NPack(pack, nil)

	fmt.Println("Done")
}
