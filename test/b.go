package test

import (
	"fmt"
	"reflect"
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
