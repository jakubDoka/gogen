package parser

import (
	"strings"

	"github.com/jakubDoka/gogen/dirs"
	"github.com/jakubDoka/gogen/str"
)

// Def is a template definition, it has methods to return
// template results and also stores needed imports
type Def struct {
	Rules

	Deps    []*Request
	Imports Imp

	ImportSelf bool

	Code, ExtCode string
}

// NDef ...
func NDef(pack string, content Content, raw dirs.Paragraph, imports Imp) (d *Def) {
	def := &Def{
		Imports: Imp{},
	}
	// Extracting rules
	for i, line := range raw {
		line.Content = str.RemInvStart(line.Content)
		if str.StartsWith(line.Content, RulesIdent) {
			_, line.Content = str.SplitToTwo(line.Content, ' ')
			if line.Content == "" {
				Exit(line, "missing definition")
				return
			}

			def.Rules = *NRules(line)
			raw = append(raw[:i], raw[i+1:]...)
			break
		}
	}

	if def.Name == "" {
		Exit(raw[0], "template is missing rules annotation")
	}

	def.Pack = pack

	args := append(def.Args, ConstructorPrefix+def.Name, def.Name)

	for _, line := range raw {
		internal, external, dep := def.ParseLine(line, content, args, imports)
		if dep {
			line.Content = internal
			dep := NRequest(pack, line, false)

			args = append(args, dep.SubName)
			args = append(args, ConstructorPrefix+dep.SubName)
			def.Deps = append(def.Deps, dep)
		} else {
			def.Code += internal + "\n"
			def.ExtCode += external + "\n"
		}

	}

	if def.ImportSelf {
		def.Imports[pack] = imports[pack]
	}

	d = def
	return
}

// ParseLine turns line of code to line that is usable for external generation and one for internal
func (d *Def) ParseLine(line dirs.Line, content Content, args []string, imports Imp) (code, exCode string, dep bool) {
	var (
		ln, i int
		lnl   = len(line.Content)
		cont  = line.Content
	)

	dep = str.StartsWith(cont, DependencyIdent)

	if dep {
		_, cont = str.SplitToTwo(cont, ' ')
		lnl = len(cont)
		if lnl == 0 {
			dep = false
			return
		}
	}
o:
	for ; i < lnl; i += ln {
		code += cont[i-ln : i]
		exCode += cont[i-ln : i]
		ln = 0

		var none bool
		for i+ln < lnl && !str.IsIdent(cont[i+ln]) {
			ln++
			none = true
		}
		if none {
			continue
		}

		for _, t := range args {
			ln = len(t)
			if !str.IsTheIdent(cont, t, i) {
				continue
			}

			code += Gibrich
			exCode += Gibrich
			continue o
		}

		for c := range content {

			ln = len(c)

			if !str.IsTheIdent(cont, c, i) {
				continue
			}
			d.ImportSelf = true
			exCode += d.Pack + "."
			continue o
		}

		for k, v := range imports {
			ln = len(k)

			if !str.IsTheImp(cont, k, i) {
				continue
			}

			d.Imports[k] = v
			continue o
		}

		ln = 1
		for i+ln < lnl && str.IsIdent(cont[i+ln]) {
			ln++
		}
	}
	code += cont[i-ln : i]
	exCode += cont[i-ln : i]
	return
}

// Produce forms a template
func (d *Def) Produce(pack string, r *Request, content Content, done map[string]*Request) (result string, deps []*Request) {
	if len(r.Args) != len(d.Args) {
		Exit(r.Line, "incorrect amount of arguments, expected: %d got: %d", len(d.Args), len(r.Args))
	}

	if pack != d.Pack {
		result = d.ExtCode
	} else {
		result = d.Code
	}

	// we need to take copy as Requests can get modified
	deps = make([]*Request, len(d.Deps))
	for i, d := range d.Deps {
		deps[i] = d.Copy()
	}

	// replacing identifier
	result = strings.ReplaceAll(result, Gibrich+d.Name, r.SubName)
	// replacing possible constructor
	result = strings.ReplaceAll(result, Gibrich+ConstructorPrefix+d.Name, ConstructorPrefix+r.SubName)
	// as this will be placed to the package, content needs to be updated
	content.NameFor(r.SubName)

	for i, a := range r.Args {
		// in case of nested argument we add it to dependencies and continue
		if !a.End {
			c := a.Copy()
			c.NestedSubstitute = d.Args[i]
			deps = append(deps, c)
			continue
		}

		ga := Gibrich + d.Args[i]

		result = strings.ReplaceAll(result, ga, a.Name) // simple argument can be replaced immediately

		// Substituting //dep annotation arguments with inputted template
		for _, dp := range deps {
			for i := range dp.Args {
				if dp.Args[i].Name == ga {
					dp.Args[i] = a
				}
			}
		}
	}

	for i := 0; i < len(deps); i++ {
		dp := deps[i]

		val, ok := done[dp.Summarize()]

		var sub string
		// if template with matching arguments is already generated reuse it and dump the current dependency
		if ok {
			sub = val.SubName
			deps = append(deps[:i], deps[i+1:]...)
			i--
		} else {
			dp.SubName = content.NameFor(dp.OriginalSub)
			sub = dp.SubName
		}

		// its getting little extrem but we also have to take dependant constructor into account
		result = strings.ReplaceAll(result, Gibrich+ConstructorPrefix+dp.NestedSubstitute, ConstructorPrefix+sub)
		// and fynally we are replacing the dependency identifier
		result = strings.ReplaceAll(result, Gibrich+dp.NestedSubstitute, sub)
	}

	return
}
