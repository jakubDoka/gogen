package test

import (
	"fmt"
	"reflect"
)

const (
	hell int = iota
	mel
	cel
	gel
)

// Bar ...
type Bar struct {
	reflect.Type
}

/*

 */

//ign(
type a = interface{}

//)

/*gen(
	Hello<int, e>
	Hello<float64, r>
	Hello<Bar, h>
	test2.min<float64, MinF64>
)*/

//def(
//rules Hello<a, f>

func f(b Bar, value a) a {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

//)
