package test

import (
	"fmt"
	"reflect"
)

// e ...
func e(b Bar, value int) int {
fmt.Print(b)
reflect.TypeOf("brb")
return value
}


// r ...
func r(b Bar, value float64) float64 {
fmt.Print(b)
reflect.TypeOf("brb")
return value
}


// h ...
func h(b Bar, value Bar) Bar {
fmt.Print(b)
reflect.TypeOf("brb")
return value
}


// Clamp64 for testing
func Clamp64(val, min, max float64) float64 {
return Max(min, Min(max, val))
}


// Max ...
func Max(a, b float64) float64 {
if a > b {
return b
}
return a
}


// Min ...
func Min(a, b float64) float64 {
if a < b {
return a
}
return b
}


// BB exists purely for testing purposes
type BB struct {
A setA
B setB
}

// DoubleAppend ...
func (d BB) DoubleAppend(valA int, valB float64) {
d.A.Add(valA)
d.B.Add(valB)
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


// setB is a float64 set
type setB map[float64]bool

// Add adds value to set
func (n setB) Add(val float64) {
n[val] = true
}

// Rem removes value from set
func (n setB) Rem(val float64) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// HH exists purely for testing purposes
type HH struct {
A setA
B setB1
}

// DoubleAppend ...
func (d HH) DoubleAppend(valA int, valB string) {
d.A.Add(valA)
d.B.Add(valB)
}


// setB1 is a string set
type setB1 map[string]bool

// Add adds value to set
func (n setB1) Add(val string) {
n[val] = true
}

// Rem removes value from set
func (n setB1) Rem(val string) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// KK exists purely for testing purposes
type KK struct {
A BB
B setD
}

// QuadrupleAppend ...
func (q KK) QuadrupleAppend(valA int, valB float64, valC bool, valD int) {
q.A.DoubleAppend(valA, valB)
q.B.DoubleAppend(valC, valD)
}


// setD exists purely for testing purposes
type setD struct {
A setA1
B setA
}

// DoubleAppend ...
func (d setD) DoubleAppend(valA bool, valB int) {
d.A.Add(valA)
d.B.Add(valB)
}


// setA1 is a bool set
type setA1 map[bool]bool

// Add adds value to set
func (n setA1) Add(val bool) {
n[val] = true
}

// Rem removes value from set
func (n setA1) Rem(val bool) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// NN exists purely for testing purposes
type NN struct {
A HH
B setD1
}

// QuadrupleAppend ...
func (q NN) QuadrupleAppend(valA int, valB string, valC int8, valD float32) {
q.A.DoubleAppend(valA, valB)
q.B.DoubleAppend(valC, valD)
}


// setD1 exists purely for testing purposes
type setD1 struct {
A setA2
B setB2
}

// DoubleAppend ...
func (d setD1) DoubleAppend(valA int8, valB float32) {
d.A.Add(valA)
d.B.Add(valB)
}


// setA2 is a int8 set
type setA2 map[int8]bool

// Add adds value to set
func (n setA2) Add(val int8) {
n[val] = true
}

// Rem removes value from set
func (n setA2) Rem(val int8) (ok bool) {
ok = n[val]
delete(n, val)
return
}


// setB2 is a float32 set
type setB2 map[float32]bool

// Add adds value to set
func (n setB2) Add(val float32) {
n[val] = true
}

// Rem removes value from set
func (n setB2) Rem(val float32) (ok bool) {
ok = n[val]
delete(n, val)
return
}

