package str

import (
	"github.com/jakubDoka/gogen/templates"
)

// Vec is a standard Vector type with utility methods
type Vec []string

// Clone creates new Vec copies content of v to it and returns it
func (v Vec) Clone() Vec {
	nv := make(Vec, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *Vec) Clear() {
	v.Truncate(0)
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// Vec will have length you specify
func (v *Vec) Truncate(l int) {
	var nil string
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Remove removes element and returns it
func (v *Vec) Remove(idx int) (val string) {
	var nil string

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *Vec) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *Vec) PopFront() string {
	return v.Remove(0)
}

// Pop removes last element
func (v *Vec) Pop() string {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *Vec) Insert(idx int, val string) {
	dv := *v
	*v = append(append(append(make(Vec, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *Vec) InsertSlice(idx int, val []string) {
	dv := *v
	*v = append(append(append(make(Vec, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v Vec) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v Vec) Last() string {
	return v[len(v)-1]
}

// Sort is quicksort for Vec, because this is part of a template comp function is necessary
func (v Vec) Sort(comp func(a, b string) bool) {
	if len(v) < 2 {
		return
	}
	ps := make(templates.IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p string

		l, e, s, j int
	)

	for {
		l = len(ps)

		e = ps[l-1] - 1
		if e <= 0 {
			return
		}

		s = ps[l-2] + 1
		p = v[e]

		if s < e {
			for j = s; j < e; j++ {
				if comp(v[j], p) {
					v.Swap(s, j)
					s++
				}
			}

			v.Swap(s, e)
			ps.Insert(l-1, s)
		} else {
			ps = ps[:l-1]
		}
	}
}

// Swap swaps two elements
func (v Vec) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v Vec) ForEach(con func(i int, e string) string) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *Vec) Filter(filter func(e string) bool) {
	dv := *v

	var i int
	for _, e := range dv {
		if filter(e) {
			dv[i] = e
			i++
		}
	}

	v.Truncate(i)
}

// Find returns first element for which find returns true along with index,
// if there is none, index equals -1
func (v Vec) Find(find func(e string) bool) (idx int, res string) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}

// String is a standard Vector type with utility methods
type String []rune

// Clone creates new String copies content of v to it and returns it
func (v String) Clone() String {
	nv := make(String, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *String) Clear() {
	v.Truncate(0)
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// String will have length you specify
func (v *String) Truncate(l int) {
	var nil rune
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Remove removes element and returns it
func (v *String) Remove(idx int) (val rune) {
	var nil rune

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *String) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *String) PopFront() rune {
	return v.Remove(0)
}

// Pop removes last element
func (v *String) Pop() rune {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *String) Insert(idx int, val rune) {
	dv := *v
	*v = append(append(append(make(String, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *String) InsertSlice(idx int, val []rune) {
	dv := *v
	*v = append(append(append(make(String, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v String) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v String) Last() rune {
	return v[len(v)-1]
}

// Sort is quicksort for String, because this is part of a template comp function is necessary
func (v String) Sort(comp func(a, b rune) bool) {
	if len(v) < 2 {
		return
	}
	ps := make(templates.IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p rune

		l, e, s, j int
	)

	for {
		l = len(ps)

		e = ps[l-1] - 1
		if e <= 0 {
			return
		}

		s = ps[l-2] + 1
		p = v[e]

		if s < e {
			for j = s; j < e; j++ {
				if comp(v[j], p) {
					v.Swap(s, j)
					s++
				}
			}

			v.Swap(s, e)
			ps.Insert(l-1, s)
		} else {
			ps = ps[:l-1]
		}
	}
}

// Swap swaps two elements
func (v String) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v String) ForEach(con func(i int, e rune) rune) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *String) Filter(filter func(e rune) bool) {
	dv := *v

	var i int
	for _, e := range dv {
		if filter(e) {
			dv[i] = e
			i++
		}
	}

	v.Truncate(i)
}

// Find returns first element for which find returns true along with index,
// if there is none, index equals -1
func (v String) Find(find func(e rune) bool) (idx int, res rune) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}
