package main

import (
	"fmt"

	"github.com/victorbrun/gosymbol"
)

func main() {
	x := gosymbol.Var("x")
	f := gosymbol.Exp(x)

	val := gosymbol.Arguments{}
	err := val.AddArgument(gosymbol.Var("x"), 8)
	if err != nil {
		fmt.Println(val, err)
	}
	err = val.AddArgument(gosymbol.Var("x"), 9)
	if err != nil {
		fmt.Println(val, err)
	}
	fn := f.Eval()
	fn_eval := fn(val)
	fmt.Println("f(x) = ", fn_eval)
}
