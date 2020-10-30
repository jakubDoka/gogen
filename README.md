 Welcome to the gogen tutorial, here i will try to explain how
 to use this simple tool to use golang like it has generics

 i assume that if you opened this file you already installed gogen
 and used some commands, so lets talk about how gogen can be used

 1. global state
 Simplest wey to use gogen is to add executable to your path environment
 variable and use it in every your project. You can keep your templates
 global this way and use them in any project. Because gogen cannot scan
 all storage in your computer (or cna but it would be slow), you have to
 specify directories where it should search for them. Another thing you have
 to do is annotating all your template files, so lets show how template file
 should look like

//gogen_template // notice this comment annotation, to speed things up gogen looks for this first and ignores file if its missing
package example // comment does not apply here, i just commented it so ide will not show me an error


```go
//gogen_template
package example
// T is template argument
type T = int // the type of T is not important as long as you can compile signatures that use it
//<<< Print<T, prefix>

// prefixPrint prints an value with addition of annoying message
// its also absolutely useless
func prefixPrint(value T) {
	fmt.Println("look at me i can print T like this:", value)
}

//<<< Vec<T, __>

// Vec__ on the other hand is lot more usefull
type Vec__ struct {
	slice []T
}

// MakeVec__ creates new Vec__ with given cap and len
func MakeVec__(cap, len int) Vec__ {
	return Vec__{make([]T, cap, len)}
}

// Push appends value
func (v *Vec__) Push(value T) {
	v.slice = append(v.slice, value)
}

// Clear clears all values but should preserve cap
func (v *Vec__) Clear() {
	v.slice = v.slice[:0]
}

//>>> // this specifies end of template
```

 You can also customize how you annotate by changing config file of gogen, 
 that is always created after first run. We can for example change prefix 
 to just "__". Mind that gogen, for now, does not differentiate between 
 part template specifier of fraction fo a name of your variable. For 
 instance if you have template T and you also use variable hashTable, 
 gogen will make it hashintable, but that should not be a problem because 
 you can name your template arguments how ewer you like and need.

 Now thats nice and all but how do we actually use our templates?

```go
//--snip-- <- snips are optional
//gen Print<int, I>
//gen Vec<float64, F>
//--snip--
```

 And thats about it, just write this wherever you like, within the package you need it
 lets also check out what will happen if you use folloving commands.
```
$ gogen add
$ gogen gen
```
 If you done this all in the same directory where you have all your files, new file
 named "gogen-output.go" should appear. Now lets see whats inside.

```go
package example

import (
  	"fmt"
)

// IPrint prints an value with adition of annoing message
// its also absolutely useles
func IPrint(value int) {
	fmt.Println("look at me i can print int like this:", value)
}

// VecF on the other hand is lot more usefull
type VecF struct {
	slice []float64
}

// MakeVecF creates new VecF with given cap and len
func MakeVecF(cap, len int) VecF {
	return VecF{make([]float64, cap, len)}
}

// Push appends value
func (v *VecF) Push(value float64) {
	v.slice = append(v.slice, value)
}

// Clear clears all values but should preserve cap
func (v *VecF) Clear() {
	v.slice = v.slice[:0]
}
```

Lovely, our templates are used as they should and notice that last
argument wos substituted for a strange the part of a name, and even 
comments and strings. Its shows beautifully how stupid yet useful 
gogen is. 

Finally if you change your template just use `gogen gen` and 
files will be regenerated.