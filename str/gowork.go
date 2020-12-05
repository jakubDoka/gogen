package str

import (
	"strings"
)

// ParseSimpleGoDef get name of go definition, for example "type bar struct {}"
// returns "bar"
func ParseSimpleGoDef(line string) []string {
	_, line = SplitToTwo(line, ' ')

	if IsRowGoDef(line) {
		return ParseRowGoDef(line)
	}

	for i := 0; i < len(line); i++ {
		if !IsIdent(line[i]) {
			return []string{line[:i]}
		}
	}

	return []string{}
}

// ParseMultilineGoDef extract all names from bulk var or const definition
// TODO this does not take multiple-online definition (var a, b, c int)
func ParseMultilineGoDef(lines []string) (results []string) {
	var depth int
	for _, line := range lines {
		line = RemInvStart(line)
		var i int

		if depth == 0 {
			if IsRowGoDef(line) {
				results = append(results, ParseRowGoDef(line)...)
				continue
			}
			for ; i < len(line); i++ {
				if !IsIdent(line[i]) {
					break
				}
			}
			if i != 0 && IsUpper(line[0]) {
				results = append(results, line[:i])
			}
		}

		for ; i < len(line); i++ {
			if line[i] == '(' || line[i] == '{' {
				depth++
			} else if line[i] == ')' || line[i] == '}' {
				depth--
			}
		}

		if depth == -1 {
			break
		}
	}
	return
}

// ParseRowGoDef extracts names from go bulk declaration
func ParseRowGoDef(line string) (results []string) {
	line, els := SplitToTwo(line, '=')

	rs := strings.Split(line, ",")
	i := len(rs) - 1

	rs[i] = RemInvEnd(RemInvStart(rs[i]))
	if els == "" {
		if first, second := SplitToTwo(rs[i], ' '); second != "" {
			if first != "" && IsUpper(first[0]) {
				results = append(results, first)
			}
		} else {
			i++
		}
	}

	for _, s := range rs[:i] {
		s = RemInv(s)
		if len(s) != 0 && IsUpper(s[0]) {
			results = append(results, s)
		}
	}
	return
}

// IsRowGoDef returns wether line is bulk go definition (var a, b = 6, 5 returns true)
func IsRowGoDef(line string) bool {
	a, b := strings.IndexByte(line, ','), strings.IndexByte(line, '=')
	return a != -1 && (b == -1 || b > a)
}

// IsGoDef returns whether string is go definition
func IsGoDef(line string) (bool, string) {
	return StartsWithMany(line, "type", "var", "const", "func")
}
