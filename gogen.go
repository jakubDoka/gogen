package main

import (
	"fmt"
	"os"

	"github.com/jakubDoka/gogen/dirs"
	"github.com/jakubDoka/gogen/parser"
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

	cons := parser.LoadConnections()

	pck := parser.NPack(pack, nil, cons, map[string]*parser.Pack{})

	if len(args) == 2 && args[1] == "r" { // regenerating all packages that depend on this package
	o:
		for p := range cons.Get(pck.Import) {
			for _, pck := range pck.Others {
				if p == pck.Import {
					continue o
				}
			}
			_, ok := dirs.PackPath(p)
			if !ok {
				delete(cons.Get(pck.Import), p)
				continue
			}
			parser.NPack(p, nil, cons, pck.Others)
		}
	}

	cons.Save()

	fmt.Println("Done")
}
