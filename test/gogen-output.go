package test

import (
	"fmt"
	"reflect"
)

func e(b Bar, value int) int {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

func r(b Bar, value float64) float64 {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

func h(b Bar, value Bar) Bar {
	fmt.Print(b)
	reflect.TypeOf("brb")
	return value
}

// MinF64 ...
func MinF64(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}
