package templates

import (
	"strconv"
	"testing"
)

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
}

func TestVecRemoveSlice(t *testing.T) {
	v := IntVec{1, 5, 4, 2, 6, 3, 7}

	v.RemoveSlice(2, 5)
	res := IntVec{1, 5, 3, 7}

	if len(v) != len(res) {
		t.Error(v, "!=", res)
	}

	for i, e := range res {
		if v[i] != e {
			t.Error(v, "!=", res)
			return
		}
	}
}

func TestBiSearch(t *testing.T) {
	vec := make(IntVec, 10000)
	for i := range vec {
		vec[i] = i
	}

	for _, tC := range vec {
		t.Run(strconv.Itoa(tC), func(t *testing.T) {
			res, ok := vec.BiSearch(tC, func(a, b int) uint8 {
				if a == b {
					return 0
				} else if a > b {
					return 1
				}
				return 2
			})

			if !ok || res != tC {
				t.Error(res, tC)
			}
		})
	}
}

func TestBiSearchFail(t *testing.T) {
	vec := make(IntVec, 10000)
	for i := range vec {
		vec[i] = i
	}

	for i, tC := range vec {
		t.Run(strconv.Itoa(tC), func(t *testing.T) {
			v := make(IntVec, 10000)
			copy(v, vec)
			v[i] = -1
			res, ok := v.BiSearch(tC, func(a, b int) uint8 {
				if a == b {
					return 0
				} else if a > b {
					return 1
				}
				return 2
			})

			if ok {
				t.Error(res, tC)
			}
		})
	}
}

func TestBiInsert(t *testing.T) {
	testCases := []struct {
		ins, start, res IntVec
	}{
		{IntVec{4}, IntVec{1, 2, 3, 4, 5}, IntVec{1, 2, 3, 4, 4, 5}},
		{IntVec{4}, IntVec{1, 2, 3, 5}, IntVec{1, 2, 3, 4, 5}},
		{IntVec{4, 8, 0, 6}, IntVec{1, 2, 3, 5}, IntVec{0, 1, 2, 3, 4, 5, 6, 8}},
	}

	for i, tC := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for _, i := range tC.ins {
				tC.start.BiInsert(i, func(a, b int) uint8 {
					if a == b {
						return 0
					} else if a > b {
						return 1
					}
					return 2
				})
			}

			for i, v := range tC.start {
				if v != tC.res[i] {
					t.Error(tC)
					break
				}
			}
		})
	}
}
