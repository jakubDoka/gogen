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
	if !IsTheIdent("[]tp", "tp", 2) {
		t.Fail()
	}
}

func TestRevSplit(t *testing.T) {
	sep := "sep"
	tests := []struct {
		name, input string
		result      []string
	}{
		{
			"just sub",
			sep,
			[]string{},
		},
		{
			"bound check",
			"asep",
			[]string{"a"},
		},
		{
			"no empty string",
			"sepASsepASsep",
			[]string{"AS", "AS"},
		},
		{
			"split limit",
			"AsepAsepAsepAsep",
			[]string{"AsepA", "A", "A"},
		},
		{
			"classic",
			"ABsepACsepAD",
			[]string{"AB", "AC", "AD"},
		},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res := RevSplit(te.input, sep, 3)
			for i, s := range res {
				if s != te.result[i] {
					t.Errorf("t=%v, r=%v", te.result, res)
				}
			}
		})
	}
}
