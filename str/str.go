package str

// ImpNm returns name of import
func ImpNm(imp string) string {
	idx := LastByte(imp, '/')
	if idx == -1 {
		return imp
	}

	return imp[idx+1:]
}

// GoDefNm get name of go definition, for example "type bar struct {}"
// returns "bar"
func GoDefNm(str string) string {
	_, str = SplitToTwo(str, ' ')
	for i := 0; i < len(str); i++ {
		if !IsIdent(str[i]) {
			return str[:i]
		}
	}

	return str
}

// IsGoDef returns whether string is go definition
func IsGoDef(str string) (bool, string) {
	return StartsWithMany(str, "type", "var", "const", "func")
}

// StartsWithMany returns true and sequence with witch the string is starting
func StartsWithMany(str string, subs ...string) (bool, string) {
	for _, s := range subs {
		if StartsWith(str, s) {
			return true, s
		}
	}

	return false, ""
}

func StartsWith(str, sub string) bool {
	return len(str) >= len(sub) && str[:len(sub)] == sub
}

func EndsWith(str, sub string) bool {
	return len(str) >= len(sub) && str[len(str)-len(sub):] == sub
}

func SplitToTwo(str string, sep rune) (a, b string) {
	for i, r := range str {
		if r == sep {
			return str[:i], str[i+1:]
		}
	}
	return str, ""
}
