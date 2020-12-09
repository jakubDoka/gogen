package templates

//def(
//rules Set<string>

// Set is a string set
type Set map[string]bool

// Add adds value to set
func (n Set) Add(val string) {
	n[val] = true
}

// Rem removes value from set
func (n Set) Rem(val string) (ok bool) {
	ok = n[val]
	delete(n, val)
	return
}

//)

//def(
//rules Vec<interface{}>

// Vec is a standard Vector type with utility methods
type Vec []interface{}

// Clone creates new Vec copies content of v to it and returns it
func (v Vec) Clone() Vec {
	nv := make(Vec, len(v))
	copy(nv, v)

	return nv
}

// Remove removes element and returns it
func (v *Vec) Remove(idx int) (val interface{}) {
	var nil interface{}

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// PopFront removes first element and returns it
func (v *Vec) PopFront() interface{} {
	return v.Remove(0)
}

// Pop removes last element
func (v *Vec) Pop() interface{} {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *Vec) Insert(idx int, val interface{}) {
	dv := *v
	*v = append(append(append(make(Vec, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// Swap swaps two elements
func (v Vec) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// Reverse reverses content of slice
func (v Vec) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
}

//)
