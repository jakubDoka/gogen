package str

import "testing"

func TestRemInv(t *testing.T) {
	tests := []struct {
		input, result string
	}{
		{
			"nana  as\t\n\rfra",
			"nanaasfra",
		},
		{
			"(",
			"(",
		},
	}

	for _, test := range tests {
		res := RemInv(test.input)
		if res != test.result {
			t.Errorf("%q != %q", res, test.result)
		}
	}
}

func TestGoDefNms(t *testing.T) {
	test := []string{
		"a = 0",
		"b = f{",
		"r: 0",
		"}",
		"c = fu(",
		"a",
		")",
		")",
	}

	result := []string{"a", "b", "c"}

	res := GoDefNms(test)

	for i, v := range result {
		if v != res[i] {
			t.Errorf("%q != %q", res, result)
		}
	}

}
