package parser

import (
	"fmt"
	"gogen/dirs"
	"gogen/str"
	"os"
)

// Block is pair of Start and end used whe parsing
type Block struct {
	Start, End string
}

// some dirty global state
var (
	// Blocks used when parsing
	Definition   = Block{"//def(", "//)"}
	Generators   = Block{"/*gen(", ")*/"}
	Imports      = Block{"/*imp(", ")*/"}
	Ignore       = Block{"//ign(", "//)"}
	MultiComment = Block{"/*", "*/"}

	Blocks = []Block{
		Definition,
		Generators,
		Imports,
		Ignore,
		MultiComment,
	}

	// def prefixes
	RulesIdent      = "//rules"
	DependencyIdent = "//dep"

	// its highly unlikely tha anyone will use 4 underscores in a row
	// so this is used to mark template arguments in a code
	Gibrich = "____"

	// Name of a output file
	OutputFile = "gogen-output.go"

	// Its better because i don't have to pass map everywhere
	AllPacks = map[string]*Pack{}
)

// IsBlockStart returns whether string is any block start
// and returns a blok of witch it is
func IsBlockStart(st string) (bool, Block) {
	for _, s := range Blocks {
		if str.StartsWith(st, s.Start) {
			return true, s
		}
	}

	return false, Block{}
}

// IsBlockEnd returns whether string is any block end
// and returns a blok of witch it is
func IsBlockEnd(st string) (bool, Block) {
	for _, s := range Blocks {
		if str.StartsWith(st, s.End) {
			return true, s
		}
	}

	return false, Block{}
}

// Exit prints an error and exits application, because this is just console app
// it is nice simplification that avoids tedious error handling
func Exit(line dirs.Line, message string, args ...interface{}) {
	fmt.Printf("file: %s\nline: %d\nerror: %s\n",
		*line.Path,
		line.Idx,
		fmt.Sprintf(message, args...),
	)
	os.Exit(2)
}
