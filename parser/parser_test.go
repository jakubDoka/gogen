package parser

// most of the tests are performed on test packages but that are not included as they are full of garbage

import (
	"testing"

	"github.com/jakubDoka/gogen/dirs"
)

func TestExtractImps(t *testing.T) {
	tests := []struct {
		name    string
		content dirs.Paragraph
		last    int
		result  map[string]string
	}{
		{
			"one-line",
			dirs.NParagraph(
				`import "hello/there"`,
			),
			0,
			map[string]string{
				"there": "hello/there",
			},
		},
		{
			"multiple-line",
			dirs.NParagraph(
				`import (`,
				`"hello/there"`,
				`"mmm/ccc"`,
				``,
				`"sss/kkk"`,
				`)`,
			),
			5,
			map[string]string{
				"there": "hello/there",
				"ccc":   "mmm/ccc",
				"kkk":   "sss/kkk",
			},
		},
		{
			"combo",
			dirs.NParagraph(
				`import (`,
				`"hello/there"`,
				`)`,
				`//hello`,
				`import "memory/doubt"`,
				``,
			),
			4,
			map[string]string{
				"there": "hello/there",
				"doubt": "memory/doubt",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, last := ExtractImps(test.content)
			for k, v := range res {
				if test.result[k] != v {
					t.Errorf("%v != %v", test.result, res)
				}
			}
			if len(test.result) != len(res) || last != test.last {
				t.Errorf("%v != %v", test.result, res)
			}
		})
	}
}

func TestBuildImps(t *testing.T) {
	test := Imp{"bb": "mm/bb", "ff": "ff/cc", "fff": "fff"}
	result := "import (\n\t\"mm/bb\"\n\t\"ff/cc\"\n)\n"
	res := test.Build(SS{"fff": true})

	if res != result {
		t.Errorf("%q != %q", res, result)
	}

}

func TestCollectContent(t *testing.T) {
	test := dirs.NParagraph(
		"type Hello struct {}",
		"func Mom() {}",
		"//def(",
		"type SneakySneak struct {}",
		"type OtherSneakySneak struct {}",
		"//)",
		"var Hell = 0",
		"const Fri = 10",
		"var (",
		"All = 0",
		"Fll = 0",
		")",
	)

	result := SS{"Hello": true, "Mom": true, "Hell": true, "Fri": true, "All": true, "Fll": true}
	resultBlock := BlockSlice{
		Definition,
		dirs.Paragraph{
			{Path: nil, Idx: 3, Content: "type SneakySneak struct {}"},
			{Path: nil, Idx: 4, Content: "type OtherSneakySneak struct {}"},
		},
	}

	res, block := CollectContent(test)

	for _, k := range res {
		if !result[k] {
			t.Errorf("%v != %v", result, res)
		}
	}

	if !CompareBlocks(block[0], resultBlock) {
		t.Errorf("%v != %v", resultBlock, block[0])
	}

}

func CompareBlocks(a, b BlockSlice) (bl bool) {
	if len(a.Raw) != len(b.Raw) || a.Type != b.Type {
		return
	}

	for i, v := range a.Raw {
		if b.Raw[i] != v {
			return
		}
	}

	return true
}
