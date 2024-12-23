package main

import (
	"fmt"

	"github.com/victorbrun/gosymbol"
)

func main() {
	x := gosymbol.Var("x")
	f := gosymbol.Pow(x, gosymbol.Const(2))

	fmt.Println("f(x) = ", f)
	fmt.Println("f'(x) = ", f.D("x"))
}
