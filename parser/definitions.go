package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strings"
)

// Def is a template definition, it has methods to return
// template results and also stores needed imports
type Def struct {
	Args    []string
	Deps    dirs.Paragraph
	Imports Imp

	ImportSelf bool

	Name, Code, ExtCode string
}

// NDef ...
func NDef(name string, content []string, raw dirs.Paragraph, imports Imp) (d *Def, err error) {
	def := &Def{Imports: Imp{}}

	var ln dirs.Line
	j := 0
	for _, line := range raw {
		line.Content = str.RemInvStart(line.Content)
		if str.StartsWith(line.Content, Rules) {
			_, line.Content = str.SplitToTwo(line.Content, ' ')
			if line.Content == "" {
				err = NError(line, "missing definition")
				return
			}

			ln = line
			def.Args, def.Name, _, err = ParseRules(line)
			if err != nil {
				return
			}
		} else {
			raw[j] = line
			j++
		}
	}

	raw = raw[:j]

	if def.Name == "" {
		err = NError(ln, "missing template rules")
		return
	}

	for _, line := range raw {
		internal, external, dep := def.ParseLine(line, name, content, imports)
		if dep {
			line.Content = name + "." + internal
			def.Deps = append(def.Deps, line)
		} else {
			def.Code += internal + "\n"
			def.ExtCode += external + "\n"
		}

	}

	if def.ImportSelf {
		def.Imports[name] = imports[name]
	}

	d = def
	return
}

// ParseLine turns line of code to line that is usable for external generation and one for internal
func (d *Def) ParseLine(line dirs.Line, name string, content []string, imports Imp) (code, exCode string, dep bool) {
	var ln int
	var lnl = len(line.Content)
	var cont = line.Content
	var i int

	dep = str.StartsWith(cont, Dependency)

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

		for _, t := range d.Args {
			ln = len(t)

			if !str.IsTheIdent(cont, t, i) {
				continue
			}

			code += Gibrich
			exCode += Gibrich
			continue o
		}

		for _, c := range content {
			ln = len(c)

			if !str.IsTheIdent(cont, c, i) {
				continue
			}
			d.ImportSelf = true
			exCode += name + "."
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
func (d *Def) Produce(line dirs.Line, external bool, args ...string) (result string, err error) {
	if len(args) != len(d.Args) {
		err = NError(line, "incorrect amount of arguments, expected: %d got: %d", len(d.Args), len(args))
		return
	}

	if external {
		result = d.ExtCode
	} else {
		result = d.Code
	}

	for i, a := range args {
		result = strings.ReplaceAll(result, Gibrich+d.Args[i], a)
		for j := range d.Deps {
			d.Deps[j].Content = strings.ReplaceAll(d.Deps[j].Content, Gibrich+d.Args[i], a)
		}
	}

	return
}

// ParseRules take line witch should contain template definition and parses it
func ParseRules(line dirs.Line) (def []string, name, raw string, err error) {
	raw = str.RemInv(line.Content) // we don't want them

	name, raw = str.SplitToTwo(raw, '<')
	if raw == "" {
		err = NError(line, "missing parameters")
		return
	}

	raw = raw[:len(raw)-1] // excluding ">"

	def = strings.Split(raw, ",")
	return
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
