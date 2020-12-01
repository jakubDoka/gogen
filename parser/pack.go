package parser

import (
	"gogen/dirs"
)

// Pack is a go Package, it stores all information needed for code generation
type Pack struct {
	Name, Import, Path string

	Files   []File
	Content []string
	Defs    []Def
}

// NPack ...
func NPack(imp string) (pack *Pack, err error) {
	p := &Pack{}

	p.Import = imp
	p.Path, err = dirs.PackPath(imp)
	if err != nil {
		return
	}

	err = p.CollectFiles()
	if err != nil {
		return
	}

	p.CollectContent()

	err = p.ResolveDefBlocks()
	if err != nil {
		return
	}

	pack = p
	return
}

// ResolveDefBlocks ...
func (p *Pack) ResolveDefBlocks() (err error) {
	for _, f := range p.Files {
		for _, b := range f.ExtractBlocks(Definition) {
			df, err := NDef(
				p.Name,
				p.Content,
				f.Raw[b.Start:b.End],
				f.Imports,
			)
			if err != nil {
				return err
			}
			p.Defs = append(p.Defs, *df)
		}
	}

	return
}

// CollectContent collects all names of definitions in packages and all
// blocks important for gogen
func (p *Pack) CollectContent() {
	for i, f := range p.Files {
		var content []string
		content, p.Files[i].Blocks = CollectContent(f.Raw)
		p.Content = append(p.Content, content...)
	}
}

// CollectFiles stores all files as lines
func (p *Pack) CollectFiles() (err error) {
	fl, err := dirs.ListFilePaths(p.Path)
	if err != nil {
		return
	}

	p.Files = make([]File, len(fl))

	var last int
	for i, f := range fl {
		fl := File{}

		fl.Raw, p.Name, err = dirs.FileAsLines(f)
		if err != nil {
			return
		}

		fl.Imports, last = ExtractImps(fl.Raw)
		fl.Imports[p.Name] = p.Import

		fl.Raw = fl.Raw[last:]

		p.Files[i] = fl
	}

	return
}

// File is a file... wow
type File struct {
	Raw     []string
	Imports Imp
	Blocks  []BlockSlice
}

// ExtractBlocks extracts all blocks of given type, blocks will be removed
// as they should be parsed just once
func (f *File) ExtractBlocks(tp Block) (extracted []BlockSlice) {
	var i int
	for _, b := range f.Blocks {
		if b.Type == tp {
			extracted = append(extracted, b)
		} else {
			f.Blocks[i] = b
			i++
		}
	}
	f.Blocks = f.Blocks[:i]

	return
}
