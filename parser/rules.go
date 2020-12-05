package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strconv"
	"strings"
)

// Rules are template rules, they can be part of a template definition, dependency or template request
type Rules struct {
	Args []string

	Idx                 int
	Name, Pack, OldName string

	Line dirs.Line
}

// NRules take line and parses a Rules, pas true if you need to parse definition rules
func NRules(line dirs.Line, isDef bool) (rules *Rules) {
	rules = &Rules{Line: line}

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
func (r *Rules) GetNameSub() string {
	return r.Args[len(r.Args)-1]
}

// UpdateNameSub updates name substitute so there is no name collizion
func (r *Rules) UpdateNameSub() {
	if r.Idx == 0 {
		return
	}
	r.Args[len(r.Args)-1] = r.GetNameSub() + strconv.Itoa(r.Idx)
}

// IsExternal returns whether definition is external
func (r *Rules) IsExternal() bool {
	return r.Pack != ""
}

// Summarize returns string that is used to determinate whether template is already generated
func (r *Rules) Summarize() (res string) {
	res = r.Name
	for _, a := range r.Args[:len(r.Args)-1] {
		res += a
	}

	return
}

// Copy copies Rules, this is necessary sa Rules contains slice
func (r *Rules) Copy() *Rules {
	nr := *r
	nr.Args = make([]string, len(r.Args))
	copy(nr.Args, r.Args)
	return &nr
}
