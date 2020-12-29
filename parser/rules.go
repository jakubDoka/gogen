package parser

import (
	"gogen/dirs"
	"gogen/str"
)

// Rules are template rules, they can be part of a template definition, dependency or template request
type Rules struct {
	Args []string

	Name, Pack string

	Line dirs.Line
}

// NRules take line and parses a Rules, pas true if you need to parse definition rules
func NRules(line dirs.Line) (r *Rules) {
	r = &Rules{Line: line}

	raw := str.RemInv(line.Content)

	r.Name, raw = str.SplitToTwo(raw, '<')
	if raw == "" {
		Exit(line, "missing rules structure")
	}

	raw = raw[:len(raw)-1] // excluding ">"

	r.Args = str.RevSplit(raw, ",", -1)

	return
}

// Request is code generation request. It supports nesting.
type Request struct {
	Args []*Request

	End bool

	Name, Pack, SubName, NestedSubstitute, OriginalSub string

	Line dirs.Line
}

// NRequest creates new request, this can fail and exit program
func NRequest(pack string, line dirs.Line, recursive bool) (r *Request) {
	r = &Request{Line: line}

	raw := line.Content

	if !recursive {
		raw = str.RemInv(raw)
	}

	r.Name, raw = str.SplitToTwo(raw, '<')
	if raw == "" {
		if !recursive {
			Exit(line, "missing rules structure")
		}

		r.End = true

		return
	}

	pck, name := str.SplitToTwo(r.Name, '.')
	if name != "" {
		r.Name, r.Pack = name, pck
	} else {
		r.Pack = pack
	}

	raw = raw[:len(raw)-1] // excluding ">"

	args := str.RevSplit(raw, ",", -1, str.Block{"<", ">"})
	l := len(args) - 1

	if !recursive {
		r.SubName, args = args[l], args[:l]
	} else {
		r.SubName = r.Name
	}

	r.OriginalSub = r.SubName
	r.NestedSubstitute = r.SubName

	for _, a := range args {
		line.Content = a
		n := NRequest(pack, line, true)
		r.Args = append(r.Args, n)
	}

	return
}

// Summarize returns string that is used to determinate whether template is already generated
func (r *Request) Summarize() (res string) {
	res = r.Name + r.Pack
	for _, a := range r.Args {
		res += a.Name
	}

	return
}

// Copy copies Rules, this is necessary sa Rules contains slice
func (r *Request) Copy() *Request {
	nr := *r
	nr.Args = make([]*Request, len(r.Args))
	for i := range nr.Args {
		nr.Args[i] = r.Args[i].Copy()
	}
	return &nr
}

// StringArgs returns args as slice of strings (only names of subRequest)
func (r *Request) StringArgs() []string {
	s := make([]string, len(r.Args)+1)
	for i, a := range r.Args {
		s[i] = a.Name
	}

	s[len(s)-1] = r.Name

	return s
}
