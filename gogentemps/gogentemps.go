package gogentemps

//def(
//rules Set<string, Set>

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

type a = string
type b = string
type setA = Set
type setB = Set

//def(
//rules DoubleSet<a, b, DoubleSet>
//dep Set<a, setA>
//dep Set<b, setB>

// Doubleset exists purely for testing purposes
type Doubleset struct {
	A setA
	B setB
}

// DoubleAppend ...
func (b Doubleset) DoubleAppend(valA a, valB b) {
	b.A.Add(valA)
	b.B.Add(valB)
}

//)
