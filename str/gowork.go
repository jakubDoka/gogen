package str

// GoDefNm get name of go definition, for example "type bar struct {}"
// returns "bar"
func GoDefNm(str string) string {
	_, str = SplitToTwo(str, ' ')
	for i := 0; i < len(str); i++ {
		if !IsIdent(str[i]) {
			return str[:i]
		}
	}

	return ""
}

// GoDefNms extract all names from bulk var or const definition
func GoDefNms(lines []string) (results []string) {
	var depth int
	for _, line := range lines {
		line = RemInv(line)
		var i int

		if depth == 0 {
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

// IsGoDef returns whether string is go definition
func IsGoDef(str string) (bool, string) {
	return StartsWithMany(str, "type", "var", "const", "func")
}
