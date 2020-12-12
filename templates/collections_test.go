package templates

import "testing"

func TestVecSort(t *testing.T) {
	v := IntVec{1, 5, 4, 2, 6, 3, 7}
	v.Sort(func(a, b int) bool {
		return a > b
	})
	res := IntVec{7, 6, 5, 4, 3, 2, 1}
	for i, e := range res {
		if v[i] != e {
			t.Error(v, "!=", res)
		}
	}
	t.Fail()
}
