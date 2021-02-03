package parser

import (
	"github.com/jakubDoka/gogen/dirs"
)

// File is a file... wow
type File struct {
	Raw     dirs.Paragraph
	Imports Imp
	Blocks  []BlockSlice
}

// ExtractBlocks extracts all blocks of given type, blocks will be removed
// as they should be parsed just once
func (f *File) ExtractBlocks(tp Block) (extracted []BlockSlice) {
	for _, b := range f.Blocks {
		if b.Type == tp {
			extracted = append(extracted, b)
		}
	}
	return
}

// BlockSlice stores info about a block time and interval containing its content
type BlockSlice struct {
	Type Block
	Raw  dirs.Paragraph
}
