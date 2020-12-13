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
		ignored     []Block
	}{
		{
			"just sub",
			sep,
			[]string{},
			[]Block{},
		},
		{
			"no sep",
			"aasdasd",
			[]string{"aasdasd"},
			[]Block{},
		},
		{
			"bound check",
			"asep",
			[]string{"a"},
			[]Block{},
		},
		{
			"no empty string",
			"sepASsepASsep",
			[]string{"AS", "AS"},
			[]Block{},
		},
		{
			"split limit",
			"AsepAsepAsepAsep",
			[]string{"AsepA", "A", "A"},
			[]Block{},
		},
		{
			"classic",
			"ABsepACsepAD",
			[]string{"AB", "AC", "AD"},
			[]Block{},
		},
		{
			"with ignore blocks",
			"AB!sep!ACsepAD",
			[]string{"AB!sep!AC", "AD"},
			[]Block{{"!", "!"}},
		},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res := RevSplit(te.input, sep, 3, te.ignored...)
			if len(res) != len(te.result) {
				t.Errorf("t=%v, r=%v", te.result, res)
				return
			}
			for i, s := range res {
				if s != te.result[i] {
					t.Errorf("t=%v, r=%v", te.result, res)
				}
			}
		})
	}
}
