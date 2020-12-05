package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strconv"
	"strings"
)

// Def is a template definition, it has methods to return
// template results and also stores needed imports
type Def struct {
	*Rls

	Deps    []*Rls
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
		if str.StartsWith(line.Content, Rules) {
			_, line.Content = str.SplitToTwo(line.Content, ' ')
			if line.Content == "" {
				NError(line, "missing definition")
				return
			}

			ln = line
			def.Rls = NRules(line, true)
		} else {
			raw[j] = line
			j++
		}
	}

	raw = raw[:j]

	if def.Name == "" {
		NError(ln, "missing template rules")
	}

	args := make([]string, len(def.Args))
	copy(args, def.Args)

	for _, line := range raw {
		internal, external, dep := def.ParseLine(line, name, content, args, imports)
		if dep {
			line.Content = name + "." + internal
			rules := NRules(line, false)
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
func (d *Def) Produce(rules *Rls, cont Counter, done map[string]*Rls) (result string, deps []*Rls) {

	if len(rules.Args) != len(d.Args) {
		NError(rules.Line, "incorrect amount of arguments, expected: %d got: %d", len(d.Args), len(rules.Args))
	}

	if rules.IsExternal() {
		result = d.ExtCode
	} else {
		result = d.Code
	}

	deps = make([]*Rls, len(d.Deps))
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

	for _, dp := range deps {
		val, ok := done[dp.GetUniqueness()]
		if !ok {
			dp.Idx = cont.Process(dp.OldName)
			dp.UpdateNameSub()
		} else {
			val.OldName = dp.OldName
			dp = val
		}
		result = strings.ReplaceAll(result, Gibrich+dp.OldName, dp.GetNameSub())
	}

	return
}

// Rls are template rules, they can be part of a template definition or template request
type Rls struct {
	Args []string

	Idx                 int
	Name, Pack, OldName string

	Line dirs.Line
}

// NRules take line witch should contain template definition and parses it
func NRules(line dirs.Line, isDef bool) (rules *Rls) {
	rules = &Rls{Line: line}

	rw := str.RemInv(line.Content) // we don't want them

	rules.Name, rw = str.SplitToTwo(rw, '<')
	if rw == "" {
		NError(line, "missing parameters")
	}

	pck, name := str.SplitToTwo(rules.Name, '.')
	if name != "" {
		rules.Name, rules.Pack = name, pck
	}

	rw = rw[:len(rw)-1] // excluding ">"

	rules.Args = strings.Split(rw, ",")

	if isDef {
		if len(rules.Args) == 0 {
			NError(line, "template rules has less then 1 argument, that is considered redundant")
		}
		rules.Args = append(rules.Args, rules.Name)
	}

	rules.OldName = rules.GetNameSub()

	return
}

// GetNameSub return template name substitute
func (r *Rls) GetNameSub() string {
	return r.Args[len(r.Args)-1]
}

// UpdateNameSub updates name substitute so there is no name collizion
func (r *Rls) UpdateNameSub() {
	if r.Idx == 0 {
		return
	}
	r.Args[len(r.Args)-1] = r.GetNameSub() + strconv.Itoa(r.Idx)
}

// IsExternal returns whether definition is external
func (r *Rls) IsExternal() bool {
	return r.Pack != ""
}

func (r *Rls) GetUniqueness() (res string) {
	res = r.Name
	for _, a := range r.Args[:len(r.Args)-1] {
		res += a
	}

	return
}

func (r *Rls) Copy() *Rls {
	nr := *r
	nr.Args = make([]string, len(r.Args))
	copy(nr.Args, r.Args)
	return &nr
}
