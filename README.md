 # gogen (golang code generator)

 Gogen is a very simple code generator. It does not force you to write code that would not be compilable, nor wrap your code in comments. You use comments just as orientation for gogen. All you need to do to install gogen is building this repository on your os and adding the root of the repo to your path environment variable.

# annotations

Nothing can be done without annotations as parsing all your code, as possible template would be inefficient. Gogen uses annotation blocks witch restrict what gogen should care about or what to ignore. Blocks then can have further annotations. Always close your blocks or they will get ignored. Lets start with a core piece, def-block.

```go
//def(
//rules Max<int, Ident>
func Ident(a, b int) int {
	if a > b {
		return a
	}
	return b
}
//)
```

Def block has a rules annotation, rules define how your template work. In this case int is a template argument that will get replaced when generating code. Likewise Ident will be replaced with identifier. As go does not have polymorphism you have to name function your self. Next is gen-block.

```go
/*gen(
	Max<float64, MaxF64>
	Max<float32, MaxF32>
	Max<byte, MaxB>
)*/
```
after running `gogen <package import>` (gogen project/main) gogen creates new file named gogen-output.go with following content

```go
package main

func MaxF64(a, b float64) float64 {
if a > b {
return a
}
return b
}

func MaxF32(a, b float32) float32 {
if a > b {
return a
}
return b
}

func MaxB(a, b byte) byte {
if a > b {
return a
}
return b
}
```

This is how you can generate your templates, you have to tell what you need, but its already better then defining them all by hand, now you can make a change to original and just rerun generation. Cross package generation is also supported. We have imp-block for this reason:

```go
/*imp(
	templates/max
)*/
```

You can then refer to the templates from package as `(package name).(template name)<...template arguments>` (max.Max<float64, MaxF64>). in case you want to use external types in your templates you have to inform gogen about it:

```go
/*gen(
	!libs/my_types
	Max<my_types.Float64, MaxF64>
	Max<my_types.Float32, MaxF32>
)*/
```

Last type of block is ign-block that its for ignoring pieces of code. Gogen takes notes about all items in your package, so it can annotate all items with `(package name).` in case of external generation. You may be shadowing something and so if you are not willing to rename shadows you can wrap shadowed code in ign-block. In case you have some huge file in your package and you do not want gogen to bother with that you can put opened ign-block on a beginning of a file.

# todo

This is a section with listed features that should be implemented, contributors are welcomed
	* support nested definitions
	* make block syntax configurable
	* make output file name configurable





