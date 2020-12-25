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

// Join joins o with n
func (n Set) Join(o Set) {
	for k := range o {
		n[k] = true
	}
}

//)

//def(
//rules OrderedMap<string, interface{}>
//dep Vec<interface{}, Vec>

// OrderedMap stores its items in underlying slice and map just keeps indexes
type OrderedMap struct {
	m map[string]int
	s Vec
}

// NOrderedMap initializes inner map
func NOrderedMap() OrderedMap {
	return OrderedMap{
		m: map[string]int{},
	}
}

// Value ...
func (o *OrderedMap) Value(key string) (val interface{}, ok bool) {
	idx, k := o.m[key]
	if !k {
		return
	}
	return o.s[idx], true
}

// Put ...
func (o *OrderedMap) Put(key string, value interface{}) {
	if i, ok := o.m[key]; ok {
		o.s[i] = value
	} else {
		o.m[key] = len(o.s)
		o.s = append(o.s, value)
	}
}

// Remove can be very slow if map is huge
func (o *OrderedMap) Remove(key string) (val interface{}, ok bool) {
	val, ok = o.Value(key)
	if ok {
		idx := o.m[key]
		delete(o.m, key)
		o.s.Remove(idx)
		for k, v := range o.m {
			if v > idx {
				o.m[k] = v - 1
			}
		}
	}
	return
}

// Slice returns underlying slice
func (o *OrderedMap) Slice() Vec {
	return o.s
}

// Index returns index of a key's value
func (o *OrderedMap) Index(name string) (int, bool) {
	val, ok := o.m[name]
	return val, ok
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

// Clear is equivalent to Truncate(0)
func (v *Vec) Clear() {
	v.Truncate(0)
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// Vec will have length you specify
func (v *Vec) Truncate(l int) {
	var nil interface{}
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
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

// RemoveSlice removes sequence of slice
func (v *Vec) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
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
	*v = append(append(append(make(Vec, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *Vec) InsertSlice(idx int, val []interface{}) {
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
func (v Vec) Last() interface{} {
	return v[len(v)-1]
}

// Sort is quicksort for Vec, because this is part of a template comp function is necessary
func (v Vec) Sort(comp func(a, b interface{}) bool) {
	if len(v) < 2 {
		return
	}
	ps := make(IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p interface{}

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
func (v Vec) ForEach(con func(i int, e interface{}) interface{}) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *Vec) Filter(filter func(e interface{}) bool) {
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
func (v Vec) Find(find func(e interface{}) bool) (idx int, res interface{}) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}

//)

//def(
//rules Storage<interface{}>

// Storage generates ids witch makes no need to use hashing
type Storage struct {
	vec      []interface{}
	freeIDs  []int
	occupied []int
	count    int
	outdated bool
}

// Insert inserts an value
func (s *Storage) Insert(value interface{}) int {
	s.count++
	s.outdated = true

	l := len(s.freeIDs)
	if l != 0 {
		id := s.freeIDs[l-1]
		s.freeIDs = s.freeIDs[:l-1]
		s.vec[id] = value
		return id
	}
	s.vec = append(s.vec, value)
	return len(s.vec) - 1
}

// Remove removes a value and frees memory for something else
func (s *Storage) Remove(id int) (val interface{}) {
	var nil interface{}
	val = s.vec[id]
	if val == nil {
		panic("removeing already removed value")
	}

	s.count--
	s.outdated = true

	s.freeIDs = append(s.freeIDs, id)
	s.vec[id] = nil

	return val
}

// Item returns value under the "id"
func (s *Storage) Item(id int) interface{} {
	return s.vec[id]
}

// Used returns whether id is used
func (s *Storage) Used(id int) bool {
	var nil interface{}
	return s.vec[id] != nil
}

// Len returns size of storage
func (s *Storage) Len() int {
	return len(s.vec)
}

// Update updates state of occupied slice, every time you remove or add
// element, storage gets outdated, this makes it up to date
func (s *Storage) Update() {
	var nil interface{}

	s.outdated = false

	s.occupied = s.occupied[:0]
	for i, v := range s.vec {
		if v != nil {
			s.occupied = append(s.occupied, i)
		}
	}
}

// Occupied return all occupied ids in storage, this method panics if Storage is outdated
// See Update method.
func (s *Storage) Occupied() []int {
	if s.outdated {
		panic("accessing occupied when storage is outdated")
	}

	return s.occupied
}

// Clear ...
func (s *Storage) Clear() {
	var nil interface{}

	for i := range s.vec {
		s.vec[i] = nil
	}

	s.vec = s.vec[:0]
	s.freeIDs = s.freeIDs[:0]
	s.count = 0
}

//)

/*gen(
	Vec<int, IntVec>
)*/
