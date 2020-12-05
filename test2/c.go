package test2

import "gogen/gogentemps"

//def(
//rules Max<int>

// Max ...
func Max(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//)

//def(
//rules Min<int>

// Min ...
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//)

//def(
//rules Clamp<int>
//dep Max<int, Max>
//dep Min<int, Min>

// Clamp for testing
func Clamp(val, min, max int) int {
	return Max(min, Min(max, val))
}

//)

type a = string
type b = string
type setA = gogentemps.Set
type setB = gogentemps.Set

//def(
//rules Doubleset<a, b>
//dep gogentemps.Set<a, setA>
//dep gogentemps.Set<b, setB>

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
//rules Quadrupleset<a, b, c, d>
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
