package parser

import (
	"fmt"
	"os"
)

/*gen(

)*/

//template (

func hello() {
	huf()
	fmt.Println("hello")
	print(os.SEEK_CUR)
}

//)

//gen self.hello<int, float64>
//gen self.hello<float64, hello>

func huf() {}
