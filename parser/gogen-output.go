package parser


// SS is a string set
type SS map[string]bool

// Add adds value to set
func (n SS) Add(val string) {
n[val] = true
}

// Rem removes value from set
func (n SS) Rem(val string) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// Doubleset exists purely for testing purposes
type Doubleset struct {
A setA
B setB
}

// DoubleAppend ...
func (string Doubleset) DoubleAppend(valA int, valB string) {
string.A.Add(valA)
string.B.Add(valB)
}


// setA is a int set
type setA map[int]bool

// Add adds value to set
func (n setA) Add(val int) {
n[val] = true
}

// Rem removes value from set
func (n setA) Rem(val int) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// setB is a string set
type setB map[string]bool

// Add adds value to set
func (n setB) Add(val string) {
n[val] = true
}

// Rem removes value from set
func (n setB) Rem(val string) (ok bool) {
ok = n[val]
delete(n, val)
return
}

