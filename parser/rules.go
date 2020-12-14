package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strconv"
)

// Rules are template rules, they can be part of a template definition, dependency or template request
type Rules struct {
	Args []*Rules

	Idx int

	IsName bool

	Name, Pack, OriginalName, SubName, NestedName string

	Line dirs.Line
}

// NRules take line and parses a Rules, pas true if you need to parse definition rules
func NRules(line dirs.Line, isDef bool, sig ...struct{}) (r *Rules) {
	nr := len(sig) == 0
	r = &Rules{Line: line}

	rw := line.Content

	if nr {
		rw = str.RemInv(rw)
	}

	r.Name, rw = str.SplitToTwo(rw, '<')
	if rw == "" {
		if nr {
			Exit(line, "missing rules structure")
		}
		r.IsName = true
		return
	}

	pck, name := str.SplitToTwo(r.Name, '.')
	if name != "" {
		r.Name, r.Pack = name, pck
	}

	rw = rw[:len(rw)-1] // excluding ">"

	args := str.RevSplit(rw, ",", -1, str.Block{"<", ">"})
	if !isDef {
		r.SubName, args = args[len(args)-1], args[:len(args)-1]
	} else {
		r.SubName = r.Name
	}

	r.OriginalName = r.SubName
	r.NestedName = r.SubName

	for _, a := range args {
		line.Content = a
		n := NRules(line, true, struct{}{})
		if isDef && nr && !n.IsName {
			Exit(line, "nested definition not allowed")
		}
		r.Args = append(r.Args, n)
	}

	if isDef {
		if len(r.Args) == 0 {
			Exit(line, "template rules has less then 1 argument, template is redundant")
		}
	}

	return
}

// UpdateNameSub updates name substitute so there is no name collizion
func (r *Rules) UpdateNameSub() {
	if r.Idx == 0 {
		return
	}
	r.SubName = r.OriginalName + strconv.Itoa(r.Idx)
}

// IsExternal returns whether definition is external
func (r *Rules) IsExternal() bool {
	return r.Pack != ""
}

// Summarize returns string that is used to determinate whether template is already generated
func (r *Rules) Summarize() (res string) {
	res = r.Name + r.Pack
	for _, a := range r.Args {
		res += a.Name
	}

	return
}

// Copy copies Rules, this is necessary sa Rules contains slice
func (r *Rules) Copy() *Rules {
	nr := *r
	nr.Args = make([]*Rules, len(r.Args))
	for i := range nr.Args {
		nr.Args[i] = r.Args[i].Copy()
	}
	return &nr
}

// StringArgs returns args as slice of strings (only names of subRules)
func (r *Rules) StringArgs() []string {
	s := make([]string, len(r.Args)+1)
	for i, a := range r.Args {
		s[i] = a.Name
	}

	s[len(s)-1] = r.Name

	return s
}
