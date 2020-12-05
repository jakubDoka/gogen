package parser

import (
	"gogen/dirs"
	"gogen/str"
	"path"
	"strings"
)

/*imp(
	gogen/gogentemps
)*/

/*gen(
	gogentemps.Set<string, SS>
)*/

// Pack is a go Package, it stores all information needed for code generation
type Pack struct {
	Name, Import, Path string

	Files   []File
	Content []string

	Defs      map[string]*Def
	Cont      Counter
	Generated map[string]*Rls
}

// NPack ...
func NPack(imp string, line *dirs.Line) (pack *Pack) {
	p := &Pack{
		Defs:      map[string]*Def{},
		Cont:      Counter{},
		Generated: map[string]*Rls{},
		Import:    imp,
	}

	p.Import = imp

	var ok bool
	p.Path, ok = dirs.PackPath(imp)
	if !ok {
		if line == nil {
			line = &dirs.Line{Path: &p.Path, Idx: -1, Content: "none"}
		}
		NError(*line, "package does not exist")
	}

	p.CollectFiles()

	p.CollectContent()

	p.LoadImports()

	p.ResolveDefBlocks()

	p.Generate()
	pack = p
	return
}

// Generate generates all code
func (p *Pack) Generate() (err error) {
	var (
		content string
		ignore  = SS{}

		def *Def
		ok  bool

		nonexistant = "template does not exist"
	)

	requests, imports := p.CollectGenRequests()
	for len(requests) != 0 {
		rls := requests[0]
		//fmt.Println(rls)
		u := rls.GetUniqueness()
		if _, ok := p.Generated[u]; ok {
			requests = requests[1:]
			continue
		} else {
			p.Generated[u] = rls
		}

		if rls.IsExternal() {
			pack, ok := AllPacks[rls.Pack]
			if !ok {
				NError(rls.Line, "nonexistant package")
			}

			def, ok = pack.Defs[rls.Name]
			if !ok {
				NError(rls.Line, nonexistant)
			}

			if !def.ImportSelf {
				ignore.Add(pack.Import)
			} else {
				ignore.Rem(pack.Import)
			}
		} else {
			def, ok = p.Defs[rls.Name]
			if !ok {
				NError(rls.Line, nonexistant)
			}

			ignore.Add(p.Import)
		}

		imports.Append(def.Imports)

		result, deps := def.Produce(rls, p.Cont, p.Generated)
		requests = append(deps, requests[1:]...)

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
func (p *Pack) CollectGenRequests() (req []*Rls, imports Imp) {
	imports = Imp{}
	for _, file := range p.Files {
		for _, block := range file.ExtractBlocks(Generators) {
			for _, line := range block.Raw {
				line.Content = str.RemInv(line.Content)
				if str.StartsWith(line.Content, "!") {
					imports.Add(line.Content[1:])
				} else {
					req = append(req, NRules(line, false))
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
					NError(line, "self import")
				}

				if _, ok := AllPacks[name]; ok {
					continue
				}

				AllPacks[name] = NPack(l, &line)
			}
		}
	}

	return
}

// ResolveDefBlocks ...
func (p *Pack) ResolveDefBlocks() (err error) {
	for _, f := range p.Files {
		for _, b := range f.ExtractBlocks(Definition) {
			df := NDef(
				p.Name,
				p.Content,
				b.Raw,
				f.Imports,
				p.Cont,
			)
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
