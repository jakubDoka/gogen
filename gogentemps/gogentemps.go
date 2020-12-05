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
//rules Doubleset<a, b, Doubleset>
//dep Set<a, setA>
//dep Set<b, setB>

// Doubleset exists purely for testing purposes
type Doubleset struct {
	A setA
	B setB
}

// DoubleAppend ...
func (d Doubleset) DoubleAppend(valA a, valB b) {
	d.A.Add(valA)
	d.B.Add(valB)
}

//)

type c = string
type d = string
type setC = Doubleset
type setD = Doubleset

//def(
//rules Quadrupleset<a, b, c, d, Quadrupleset>
//dep Doubleset<a, b, setC>
//dep Doubleset<c, d, setD>

// Quadrupleset exists purely for testing purposes
type Quadrupleset struct {
	A setC
	B setD
}

// QuadrupleAppend ...
func (q Quadrupleset) QuadrupleAppend(valA a, valB b, valC c, valD d) {
	q.A.DoubleAppend(valA, valB)
	q.B.DoubleAppend(valC, valD)
}

//)
