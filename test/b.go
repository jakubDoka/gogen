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

//def(
//rules Hello<a>

func f(b Bar, value a) a {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

//)
