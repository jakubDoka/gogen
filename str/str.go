package str

/*imp(
	gogen/templates
)*/

/*gen(
	templates.Vec<string, Vec>
)*/

// IsUpper returns whether byte is upper case
func IsUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// ImpNm returns name of import
func ImpNm(imp string) string {
	idx := LastByte(imp, '/')
	if idx == -1 {
		return imp
	}

	return imp[idx+1:]
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

// StartsWith ...
func StartsWith(str, sub string) bool {
	return len(str) >= len(sub) && str[:len(sub)] == sub
}

// EndsWith ...
func EndsWith(str, sub string) bool {
	return len(str) >= len(sub) && str[len(str)-len(sub):] == sub
}

// SplitToTwo splits string to two on first occurrence of byte
func SplitToTwo(str string, sep byte) (a, b string) {
	for i := 0; i < len(str); i++ {
		if str[i] == sep {
			return str[:i], str[i+1:]
		}
	}
	return str, ""
}

// RevSplit is similar to strings.Split but it starts from back of a string and
// will split string to at most count substrings
func RevSplit(str, sep string, count int) (res Vec) {
	l, sl := len(str), len(sep)
	for i := l; i >= sl && count != 1; i-- {
		if sep == str[i-sl:i] {
			if l != i {
				res = append(res, str[i:l])
				count--
			}
			i -= sl
			l = i
		}
	}

	if 0 != l {
		res = append(res, str[:l])
	}

	res.Reverse()

	return
}
