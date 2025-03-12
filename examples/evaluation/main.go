package main

import (
	"fmt"

	"github.com/victorbrun/gosymbol"
)

func main() {
	x := gosymbol.Var("x")
	y := gosymbol.Var("y")
	f := gosymbol.FromLatex(`x + \pi + y`)
	fmt.Println("f(x, y) = ", f)
	val := gosymbol.Arguments{}
	err := val.AddArgument(x, y)
	if err != nil {
		fmt.Println(val, err)
	}
	err = val.AddArgument(y, y)
	if err != nil {
		fmt.Println(val, err)
	}
	fn := f.Eval()
	fn_eval := fn(val)
	fmt.Printf("f(%s, %s) = %s\n", val[x], val[y], fn_eval)

	a := gosymbol.FromLatex("(2+2)*6")
	if err != nil {
		println(err)
	}

	println("%s", a.String())
}
