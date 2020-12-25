package str

import "unicode/utf8"

/*imp(
	gogen/templates
)*/

/*gen(
	templates.Vec<rune, String>
)*/

// NString creates String from string
func NString(s string) String {
	bs := []byte(s)
	n := make(String, 0, len(bs)) // yes we may allocate more then we need but our ram big enough to handle that
	for utf8.FullRune(bs) {
		r, size := utf8.DecodeRune(bs)
		bs = bs[size:]
		n = append(n, r)
	}

	return n
}
