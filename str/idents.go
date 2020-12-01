package str

import (
	"strings"
	"unicode/utf8"
)

// ContainsIdent returns whether string contains given identifier
func ContainsIdent(str, ident string) bool {
	idx := strings.Index(str, ident)
	if idx == -1 {
		return false
	}

	return IsSliceIdent(str, idx, len(ident))
}

// IsTheIdent whether there is ident on given index
func IsTheIdent(str, ident string, idx int) bool {
	l := len(ident)
	return len(str) > idx+l && str[idx:idx+l] == ident && IsSliceIdent(str, idx, l)
}

// IsSliceIdent returns whether slice of string is identifier
// assuming content of slice is valid identifier
func IsSliceIdent(str string, idx, ln int) bool {
	if idx != 0 {
		b := str[idx-1]
		if IsIdent(b) {
			return false
		}
	}

	if idx+ln < len(str) {
		b := str[idx+ln]
		if IsIdent(b) && b != '.' {
			return false
		}
	}

	return true
}

// IsIdent returns whether byte is part of an identifier
func IsIdent(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '_' || c == '.' || c >= utf8.RuneSelf
}
