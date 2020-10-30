 # gogen (golang code generator)

 Gogen is a very simple code generator, It does not force you to write code that would
 not be compilable, you can use comments to annotate your templates. Main reason why i 
 have created gogen is to substitute missing generics golang needs so match. All you need
 to do to install gogen is building this repository on your os and adding the root to your
 path.
 
 Sor those that are not sure:
 ```
 git clone https://github.com/jakubDoka/gogen
 cd gogen
 go build gogen.go
 ``` 


 ## global state

 Simplest wey to use gogen is to add directory with executable to your path environment
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

 And thats about it, just write this wherever you like, within the package you need it.
 Lets also check out what will happen if you use folloving commands.
```
gogen add
gogen gen
```
 If you done this in the same directory where you have templates and requests, new file
 named "gogen-output.go" should appear. Now lets see whats inside.

```go
package example

import (
  	"fmt"
)

// IPrint prints an value with addition of annoying message
// its also absolutely useless
func IPrint(value int) {
	fmt.Println("look at me i can print int like this:", value)
}

// VecF on the other hand is lot more useful
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
argument was substituted for a strange the part of a name, and even 
comments and strings. Its shows beautifully how stupid yet useful 
gogen is. 

Finally if you change your template just use `gogen gen` and 
files will be regenerated.

## local state

Other way around is making local configuration. That is not a problem with gogen, 
all you need to do us add `-l` flag to your commands. Gogen will create local config 
file you can then edit. 
