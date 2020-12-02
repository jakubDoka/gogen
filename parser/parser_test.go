package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExtractImps(t *testing.T) {
	tests := []struct {
		name    string
		content []string
		last    int
		result  map[string]string
	}{
		{
			"one-line",
			[]string{
				`import "hello/there"`,
			},
			0,
			map[string]string{
				"there": "hello/there",
			},
		},
		{
			"multiple-line",
			[]string{
				`import (`,
				`"hello/there"`,
				`"mmm/ccc"`,
				``,
				`"sss/kkk"`,
				`)`,
			},
			5,
			map[string]string{
				"there": "hello/there",
				"ccc":   "mmm/ccc",
				"kkk":   "sss/kkk",
			},
		},
		{
			"combo",
			[]string{
				`import (`,
				`"hello/there"`,
				`)`,
				`//hello`,
				`import "memory/doubt"`,
				``,
			},
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

	if BuildImps(test, "fff") != result {
		t.Errorf("%q != %q", BuildImps(test, "fff"), result)
	}

}

func TestCollectContent(t *testing.T) {
	test := []string{
		"type Hello struct {}",
		"func Mom() {}",
		"//def(",
		"type SneakySneak struct {}",
		"type OtherSneakySneak struct {}",
		"//)",
		"var Hell = 0",
		"const Fri = 10",
	}

	result := []string{"Hello", "Mom", "Hell", "Fri"}
	resultBlock := BlockSlice{
		Definition,
		3,
		5,
		[]string{
			"type SneakySneak struct {}",
			"type OtherSneakySneak struct {}",
		},
	}

	res, block := CollectContent(test)

	for i := range res {
		if res[i] != result[i] {
			t.Errorf("%v != %v", result, res)
		}
	}

	if !CompareBlocks(block[0], resultBlock) {
		t.Errorf("%q != %q", resultBlock, block[0])
	}

}

func CompareBlocks(a, b BlockSlice) (bl bool) {
	if len(a.Raw) != len(b.Raw) || a.Start != b.Start || a.End != b.End || a.Type != b.Type {
		return
	}

	for i, v := range a.Raw {
		if b.Raw[i] != v {
			return
		}
	}

	return true
}

func TestNPack(t *testing.T) {
	pck, err := NPack("gogen/test")
	if err != nil {
		panic(err)
	}
	bts, _ := json.MarshalIndent(pck, "", "  ")
	fmt.Println(pck.Defs[0].Produce(true, "int"))
	fmt.Println(pck.Defs[0].Produce(true, "float64"))
	t.Error(string(bts))
}
