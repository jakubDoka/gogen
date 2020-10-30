package commands

import (
	"fmt"
	"strings"
)

// ErrData remembers where something is in code
type ErrData struct {
	path string
	line int
	err  SyntaxErr
}

// Rules stores information about template or template request
type Rules struct {
	ErrData
	name string
	args []string
}

// NewRules decomposes raw template syntax
func NewRules(rules string, adr ErrData) (req Rules) {
	req.ErrData = adr
	parts := strings.Split(CleanTemplateKey(rules), "<")
	if len(parts) != 2 {
		req.err = IncorrectRules
		return
	}

	req.name = parts[0]
	req.args = strings.Split(parts[1][:len(parts[1])-1], ",")
	if len(req.args) == 0 || req.args[0] == "" {
		req.err = MissingTemplateArgs
		return
	}

	return
}

// Eq compares to Rules
func (r *Rules) Eq(o *Rules) bool {
	if r.name == o.name && len(r.args) == len(o.args) {
		for i, v := range r.args {
			if o.args[i] != v {
				return false
			}
		}
		return true
	}
	return false
}

// FilterDuplipcates ...
func FilterDuplipcates(slice []Rules) (seen []Rules) {

	for _, r := range slice {
		if contains(seen, r) {
			continue
		}
		seen = append(seen, r)
	}

	return
}

func contains(slice []Rules, r Rules) bool {
	for _, ru := range slice {
		if r.Eq(&ru) {
			return true
		}
	}

	return false
}

// Template stores all needed data for generating code
type Template struct {
	Rules
	content string
	imports map[string]bool
}

// NewTemplate ...
func NewTemplate(rules, content string, imports []string, adr ErrData) Template {
	temp := Template{content: content, imports: map[string]bool{}}
	temp.Rules = NewRules(rules, adr)
	if temp.err != -1 {
		return temp
	}

	for _, v := range imports {
		fmt.Println(v[1 : len(v)-1])
		idx := strings.Index(content, v[1:len(v)-1]+".")
		if idx == -1 {
			continue
		}

		if idx == 0 || strings.Contains(" \n\t\r", content[idx-1:idx]) {
			temp.imports[v] = true
		}
	}

	return temp
}

// Generate creates new variation of template that also can be impossible to compile
func (t *Template) Generate(subs []string) string {
	var result = t.content
	for i, a := range t.args {
		result = strings.Replace(result, a, subs[i], -1)
	}
	result = strings.Replace(result, Cf.PrefixSpecifier, subs[len(subs)-1], -1)
	return "\n" + result
}
