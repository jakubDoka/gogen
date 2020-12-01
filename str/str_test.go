package str

import "testing"

func TestRemInv(t *testing.T) {
	test := struct {
		input, result string
	}{
		"nana  as\t\n\rfra",
		"nanaasfra",
	}

	res := RemInv(test.input)
	if res != test.result {
		t.Errorf("%q != %q", res, test.result)
	}

}
