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
	test2.Clamp<float64, Clamp64>
	test2.Doubleset<int, float64, BB>
	test2.Doubleset<int, string, HH>
	test2.Quadrupleset<int, float64, bool, int, KK>
	test2.Quadrupleset<int, string, int8, float32, NN>
)*/

//def(
//rules Hello<a>

// Hello ...
func Hello(b Bar, value a) a {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

//)
