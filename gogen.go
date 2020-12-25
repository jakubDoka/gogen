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

	parser.LoadConnections()

	pck := parser.NPack(pack, nil)

	if len(args) == 2 && args[1] == "r" { // regenerating all packages that depend on this package
	o:
		for p := range parser.Cons.Get(pck.Import) {
			for _, pck := range parser.AllPacks {
				if p == pck.Import {
					continue o
				}
			}
			_, ok := dirs.PackPath(p)
			if !ok {
				delete(parser.Cons.Get(pck.Import), p)
				continue
			}
			parser.NPack(p, nil)
		}
	}

	parser.Cons.Save()

	fmt.Println("Done")
}
