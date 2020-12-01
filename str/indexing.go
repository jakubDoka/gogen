package str

// Index returns index of sub in str or -1, start specifies point from
// witch search starts
func Index(str, f string, start int) int {
	l := len(f)
	for i := start; i < len(str)-l; i++ {
		if str[i:i+l] == f {
			return i
		}
	}
	return -1
}

// LastByte returns last occurrence of b in str or -1
func LastByte(str string, b byte) int {
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == b {
			return i
		}
	}
	return -1
}
