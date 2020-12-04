package parser

import (
	"fmt"
	"gogen/dirs"
	"gogen/gogentemps"
	"gogen/str"
	"path"
	"strings"
)

/*imp(
	gogen/gogentemps
)*/

/*gen(
	gogentemps.Set<string, SS>
	gogentemps.DoubleSet<int, string, DS>
)*/

// Pack is a go Package, it stores all information needed for code generation
type Pack struct {
	Name, Import, Path string

	Files   []File
	Content []string
	Defs    map[string]*Def

	Generated gogentemps.Set
}

// NPack ...
func NPack(imp string, line *dirs.Line) (pack *Pack, err error) {
	p := &Pack{Defs: map[string]*Def{}, Generated: gogentemps.Set{}}

	p.Import = imp
	var ok bool
	p.Path, ok = dirs.PackPath(imp)
	if !ok {
		if line == nil {
			line = &dirs.Line{Path: &p.Path, Idx: -1, Content: "none"}
		}
		err = NError(*line, "package does not exist")
		return
	}

	err = p.CollectFiles()
	if err != nil {
		return
	}

	p.CollectContent()

	err = p.LoadImports()
	if err != nil {
		return
	}

	err = p.ResolveDefBlocks()
	if err != nil {
		return
	}

	err = p.Generate()
	if err != nil {
		return
	}

	pack = p
	return
}

// Generate generates all code
func (p *Pack) Generate() (err error) {
	var (
		content string
		ignore  string
		result  string

		def *Def
		ok  bool

		nonexistant = "template does not exist"
	)

	requests, imports := p.CollectGenRequests()

	for i := 0; i < len(requests); i++ {
		line := requests[i]

		args, name, raw, err := ParseRules(line)
		if err != nil {
			return err
		}

		if p.Generated[raw] {
			continue
		} else {
			p.Generated.Add(raw)
		}

		if strings.Contains(name, ".") {
			packName, name := str.SplitToTwo(name, '.')

			pack, ok := AllPacks[packName]
			if !ok {
				return NError(line, "nonexistant package")
			}

			def, ok = pack.Defs[name]
			if !ok {
				return NError(line, nonexistant)
			}

			imports.Append(def.Imports)

			result, err = def.Produce(line, true, args...)

			if !def.ImportSelf {
				ignore = pack.Import
			}
		} else {
			def, ok = p.Defs[name]
			if !ok {
				fmt.Println(p.Defs)
				return NError(line, nonexistant)
			}

			imports.Append(def.Imports)

			result, err = def.Produce(line, false, args...)

			ignore = p.Import

		}
		requests = append(requests, def.Deps...)
		if err != nil {
			return err
		}
		content += result
	}

	path := path.Join(p.Path, OutputFile)

	if content == "" {
		return dirs.DeleteIfPresent(path)
	}

	return dirs.CreateFile(
		path,
		"package "+p.Name+"\n\n"+
			imports.Build(ignore)+
			content,
	)
}

// CollectGenRequests ...
func (p *Pack) CollectGenRequests() (req dirs.Paragraph, imports Imp) {
	imports = Imp{}
	for _, file := range p.Files {
		for _, block := range file.ExtractBlocks(Generators) {
			for _, line := range block.Raw {
				line.Content = str.RemInv(line.Content)
				if str.StartsWith(line.Content, "!") {
					imports.Add(line.Content[1:])
				} else {
					req = append(req, line)
				}
			}
		}
	}
	return
}

// LoadImports loads all dependant packages
func (p *Pack) LoadImports() (err error) {
	for _, file := range p.Files {
		for _, block := range file.ExtractBlocks(Imports) {
			for _, line := range block.Raw {
				l := str.RemInv(line.Content)
				name := str.ImpNm(l)
				if name == p.Name {
					return NError(line, "self import")
				}
				if _, ok := AllPacks[name]; ok {
					continue
				}
				var pack *Pack
				pack, err = NPack(l, &line)
				if err != nil {
					return
				}
				AllPacks[name] = pack
			}
		}
	}

	return
}

// ResolveDefBlocks ...
func (p *Pack) ResolveDefBlocks() (err error) {
	for _, f := range p.Files {
		for _, b := range f.ExtractBlocks(Definition) {
			df, err := NDef(
				p.Name,
				p.Content,
				b.Raw,
				f.Imports,
			)
			if err != nil {
				return err
			}
			p.Defs[df.Name] = df
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
	fl, err := dirs.ListFilePaths(p.Path, ".go")
	if err != nil {
		return
	}

	p.Files = make([]File, len(fl))

	var last int
	for i, f := range fl {
		if strings.Contains(f, OutputFile) {
			continue
		}

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
