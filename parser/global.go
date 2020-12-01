package parser

import (
	"gogen/str"
)

// Block is pair of Start and end used whe parsing
type Block struct {
	Start, End string
}

// Blocks used when parsing
var (
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

	Rules      = "//rules"
	Dependency = "//dependency"
	Gibrich    = "____"
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
