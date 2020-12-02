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
	var dir string
	if len(args) == 0 {
		dir = dirs.PackImport(wd)
	} else {
		dir = args[0]
	}

	_, err = parser.NPack(dir, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Done")
	}
}

//def(
//rules Max<int, Ident>
func Ident(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//)

/*gen(
	Max<float64, MaxF64>
	Max<float32, MaxF32>
	Max<byte, MaxB>
)*/
