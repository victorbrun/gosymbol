package main

import (
	"fmt"

	"github.com/victorbrun/gosymbol"
)

func main() {
	x := gosymbol.Var("x")
	f := gosymbol.Mul(x, gosymbol.Exp(x))

	fmt.Println("f(x) = ", f)
	fmt.Println("f'(x) = ", f.D(x))
}
