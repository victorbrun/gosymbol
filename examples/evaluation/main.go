package main

import (
	"fmt"

	"github.com/victorbrun/gosymbol"
)

func main() {
	x := gosymbol.Var("x")
	f := gosymbol.Exp(x)

	val := gosymbol.Arguments{"x": 8}
	fn := f.Eval()
	fn_eval := fn(val)
	fmt.Println("f(x) = ", fn_eval)
}
