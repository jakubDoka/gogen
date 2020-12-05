package str

// RemInvStart removes all invisible character on the beginning of a string
func RemInvStart(str string) string {
	for i := 0; i < len(str); i++ {
		if IsVisible(str[i]) {
			return str[i:]
		}
	}
	return ""
}

// RemInvStart removes all invisible character at the end of a string
func RemInvEnd(str string) string {
	for i := len(str) - 1; i >= 0; i-- {
		if IsVisible(str[i]) {
			return str[:i+1]
		}
	}
	return ""
}

// RemInv removes invisible parts of a string
func RemInv(str string) (res string) {
	var j int
	for i := 0; i < len(str); i++ {
		if !IsVisible(str[i]) {
			if j != i {
				res += str[j:i]
			}
			j = i + 1
		}
	}

	res += str[j:]

	return
}

// IsVisible returns whether byte is visible
func IsVisible(r byte) bool {
	return r != ' ' && r != '\n' && r != '\r' && r != '\t'
}
