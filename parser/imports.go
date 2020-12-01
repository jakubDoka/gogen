package parser

import "gogen/str"

// Imp symbolizes imports object
type Imp = map[string]string

// ExtractImps collects all imports in a file and saves to a map
func ExtractImps(raw []string) (Imp, int) {
	imports := Imp{}
	var inside bool
	var last int
	for i, l := range raw {
		l = str.RemInv(l)
		ln := str.LastByte(l, '"')

		if inside {
			if ln != -1 {
				st := l[1:ln]
				imports[str.ImpNm(st)] = st
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
				st := l[len("import")+1 : ln]
				imports[str.ImpNm(st)] = st
				last = i
			}
		} else if ok, _ := str.IsGoDef(l); ok {
			break
		}
	}

	return imports, last
}

// BuildImps turns Imp to valid go import syntax
func BuildImps(imports Imp, ignore string) string {
	result := "import (\n"
	for _, v := range imports {
		if v == ignore {
			continue
		}
		result += "\t\"" + v + "\"\n"
	}
	return result + ")\n"
}

// CollectContent collects all package content that can be imported
func CollectContent(raw []string) (content []string, blocks []BlockSlice) {
	var inBlock bool
	var current BlockSlice
	for i, l := range raw {
		l = str.RemInvStart(l)
		if inBlock {
			if ok, _ := IsBlockEnd(l); ok {
				current.End = i
				blocks = append(blocks, current)

				inBlock = false
			}
			continue
		} else {
			if ok, tp := IsBlockStart(l); ok {
				current.Type = tp
				current.Start = i + 1

				inBlock = true
				continue
			}
		}

		if ok, _ := str.IsGoDef(l); ok {
			name := str.GoDefNm(l)
			if name == "" { // its a struct method
				continue
			}
			content = append(content, name)
		}
	}

	return
}

// BlockSlice stores info about a block time and interval containing its content
type BlockSlice struct {
	Type       Block
	Start, End int
}
