package templates


// IntVec is a standard Vector type with utility methods
type IntVec []int

// Clone creates new IntVec copies content of v to it and returns it
func (v IntVec) Clone() IntVec {
	nv := make(IntVec, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *IntVec) Clear() {
	v.Truncate(0)
}

// Rewrite revrites elements from index to o values
func (v IntVec) Rewrite(o IntVec, idx int) {
	copy(v[idx:], o)
}

// Len implements VertexData interface
func (v IntVec) Len() int {
	return len(v)
}

// UClear does not care about memory leaks, it just sets length to 0
func (v *IntVec) UClear() {
	*v = (*v)[:0]
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// IntVec will have length you specify
func (v *IntVec) Truncate(l int) {
	var nil int
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Extend extends vec size by amount so then len(IntVec) = len(IntVec) + amount
func (v *IntVec) Extend(amount int) {
	vv := *v
	l := len(vv) + amount
	if cap(vv) >= l { // no need to allocate
		*v = vv[:l]
		return
	}

	nv := make(IntVec, l)
	copy(nv, vv)
	*v = nv
}

// Remove removes element and returns it
func (v *IntVec) Remove(idx int) (val int) {
	var nil int

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *IntVec) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *IntVec) PopFront() int {
	return v.Remove(0)
}

// Pop removes last element
func (v *IntVec) Pop() int {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *IntVec) Insert(idx int, val int) {
	dv := *v
	*v = append(append(append(make(IntVec, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *IntVec) InsertSlice(idx int, val []int) {
	dv := *v
	*v = append(append(append(make(IntVec, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v IntVec) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v IntVec) Last() int {
	return v[len(v)-1]
}

// Sort is quicksort for IntVec, because this is part of a template comp function is necessary
func (v IntVec) Sort(comp func(a, b int) bool) {
	if len(v) < 2 {
		return
	}
	// Template is part  of its self, how amazing
	ps := make(IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p int

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
func (v IntVec) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v IntVec) ForEach(con func(i int, e int) int) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *IntVec) Filter(filter func(e int) bool) {
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
func (v IntVec) Find(find func(e int) bool) (idx int, res int) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}

//BiSearch performs a binary search on Ves assuming it is sorted. cmp consumer should
// return 0 if a == b equal, 1 if a > b and 2 if b > a, even if value wos not found it returns
// it returns closest index and false
func (v IntVec) BiSearch(value int, cmp func(a, b int) uint8) (int, bool) {
	start, end := 0, len(v)
	for {
		mid := start + (end-start)/2
		switch cmp(v[mid], value) {
		case 0:
			return mid, true
		case 1:
			end = mid + 0
		case 2:
			start = mid + 1
		}

		if start == end {
			return start, start < len(v) && cmp(v[start], value) == 0
		}
	}
}

// BiInsert inserts inserts value in a way that keeps vec sorted, binary search is used to determinate
// where to insert
func (v *IntVec) BiInsert(value int, cmp func(a, b int) uint8) {
	i, _ := v.BiSearch(value, cmp)
	v.Insert(i, value)
}

