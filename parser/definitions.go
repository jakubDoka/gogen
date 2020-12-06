package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strings"
)

// Def is a template definition, it has methods to return
// template results and also stores needed imports
type Def struct {
	*Rules

	Deps    []*Rules
	Imports Imp

	ImportSelf bool

	Code, ExtCode string
}

// NDef ...
func NDef(name string, content []string, raw dirs.Paragraph, imports Imp, cont Counter) (d *Def) {
	def := &Def{Imports: Imp{}}

	var ln dirs.Line
	j := 0
	for _, line := range raw {
		line.Content = str.RemInvStart(line.Content)
		if str.StartsWith(line.Content, RulesIdent) {
			_, line.Content = str.SplitToTwo(line.Content, ' ')
			if line.Content == "" {
				Exit(line, "missing definition")
				return
			}

			ln = line
			def.Rules = NRules(line, true)
		} else {
			raw[j] = line
			j++
		}
	}

	raw = raw[:j]

	if def.Name == "" {
		Exit(ln, "missing template rules")
	}

	args := make([]string, len(def.Args))
	copy(args, def.Args)

	for _, line := range raw {
		internal, external, dep := def.ParseLine(line, name, content, args, imports)
		if dep {
			line.Content = internal
			rules := NRules(line, false)
			if !rules.IsExternal() {
				rules.Pack = name
			}
			args = append(args, rules.GetNameSub())
			def.Deps = append(def.Deps, rules)
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
func (d *Def) ParseLine(line dirs.Line, name string, content, args []string, imports Imp) (code, exCode string, dep bool) {
	var (
		ln, i int
		lnl   = len(line.Content)
		cont  = line.Content
	)

	args = append(args, ConstructorPrefix+args[len(args)-1])
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
func (d *Def) Produce(rules *Rules, cont Counter, done map[string]*Rules) (result string, deps []*Rules) {

	if len(rules.Args) != len(d.Args) {
		Exit(rules.Line, "incorrect amount of arguments, expected: %d got: %d", len(d.Args), len(rules.Args))
	}

	if rules.IsExternal() {
		result = d.ExtCode
	} else {
		result = d.Code
	}

	deps = make([]*Rules, len(d.Deps))
	for i, d := range d.Deps {
		deps[i] = d.Copy()
	}

	for i, a := range rules.Args {
		ga := Gibrich + d.Args[i]
		result = strings.ReplaceAll(result, ga, a)
		for _, dp := range deps {
			for i := range dp.Args {
				if dp.Args[i] == ga {
					dp.Args[i] = a
				}
			}
		}
	}

	result = strings.ReplaceAll(result, Gibrich+ConstructorPrefix+d.Name, ConstructorPrefix+rules.GetNameSub())

	for _, dp := range deps {
		val, ok := done[dp.Summarize()]
		if !ok {
			dp.Idx = cont.Increment(dp.OldName)
			dp.UpdateNameSub()
		} else {
			val.OldName = dp.OldName
			dp = val
		}
		result = strings.ReplaceAll(result, Gibrich+dp.OldName, dp.GetNameSub())
	}

	return
}
