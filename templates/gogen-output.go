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

// Remove removes element and returns it
func (v *IntVec) Remove(idx int) (val int) {
var nil int

dv := *v

val = dv[idx]
*v = append(dv[:idx], dv[1+idx:]...)

dv[len(dv)-1] = nil

return val
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

// Sort is quicksort for IntVec, because this is a template comp function is necessary
func (v IntVec) Sort(comp func(a, b int) bool) {
ps := IntVec{-1, len(v)}
var p int
var l, e, s, j int
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

// Eq compares two IntVec
func (v IntVec) Eq(o IntVec) bool {
if len(v) != len(o) {
return false
}

for i := range v {
if v[i] != o[i] {
return false
}
}

return true
}

// Swap swaps two elements
func (v IntVec) Swap(a, b int) {
v[a], v[b] = v[b], v[a]
}

