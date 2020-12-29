package parser

import (
	"gogen/dirs"
	"gogen/str"
	"strconv"
)

// ExtractImps collects all imports in a file and saves them to a map
func ExtractImps(raw dirs.Paragraph) (Imp, int) {
	imports := Imp{}
	var inside bool
	var last int
	for i, line := range raw {
		l := str.RemInv(line.Content)
		ln := str.LastByte(l, '"')

		if inside {
			if ln != -1 {
				imports.Add(l[1:ln])
			}

			if str.EndsWith(l, ")") {
				inside = false
				last = i
			}

			continue
		}

		if str.StartsWith(l, "import") {
			switch l[len("import")] {
			case '(':
				inside = true
			case '"':
				imports.Add(l[len("import")+1 : ln])
				last = i
			}
		} else if ok, _ := str.IsGoDef(l); ok {
			break
		}
	}

	return imports, last
}

// Imp symbolizes imports object
type Imp map[string]string

// Append appens argument to caller
func (i Imp) Append(imp Imp) {
	for k, v := range imp {
		i[k] = v
	}
}

// Add adds import to Imp
func (i Imp) Add(imp string) {
	i[str.ImpNm(imp)] = imp
}

// Build turns Imp to valid go import syntax
func (i Imp) Build(ignore SS) string {
	result := "import (\n"
	for _, v := range i {
		if ignore[v] {
			continue
		}
		result += "\t\"" + v + "\"\n"
	}

	if result == "import (\n" {
		return ""
	}

	return result + ")\n"
}

// CollectContent collects all package content that can be imported from other package.
// This is important for external generation.
func CollectContent(raw dirs.Paragraph) (content []string, blocks []BlockSlice) {
	var inBlock bool
	var current BlockSlice
	for i, line := range raw {
		l := str.RemInvStart(line.Content)
		if inBlock {
			if ok, _ := IsBlockEnd(l); ok {
				blocks = append(blocks, current)
				current = BlockSlice{}

				inBlock = false
			} else {
				current.Raw = append(current.Raw, line)
			}
			continue
		} else {
			if ok, tp := IsBlockStart(l); ok {
				current.Type = tp

				inBlock = true
				continue
			}
		}

		if ok, _ := str.IsGoDef(l); ok {
			defs := str.ParseSimpleGoDef(l)
			if len(defs) == 0 {
				if str.EndsWith(str.RemInv(line.Content), "(") { // multiline def
					for _, d := range str.ParseMultilineGoDef(raw.GetContent()[i+1:]) {
						content = append(content, d)
					}
				}
			} else {
				for _, d := range defs {
					content = append(content, d)
				}
			}

		}
	}

	return
}

// Content stores go definitions names and returns unique names for shadows within one package
type Content map[string]int

// NameFor returns unique name
func (c Content) NameFor(name string) string {
	val := c[name]

	c[name]++

	if val == 0 {
		return name
	}

	return c.NameFor(name + strconv.Itoa(c[name]))
}
