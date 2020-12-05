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
		"A = 0",
		"B = f{",
		"r: 0,",
		"}",
		"C = fu(",
		"a",
		")",
		"D, E int",
		")",
	}

	result := []string{"A", "B", "C", "E", "D"}

	res := ParseMultilineGoDef(test)

	for i, v := range result {
		if v != res[i] {
			t.Errorf("%q != %q", res, result)
		}
	}

}

func TestIsIdent(t *testing.T) {
	t.Error(IsTheIdent("[]tp", "tp", 2))
}
