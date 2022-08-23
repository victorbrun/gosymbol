package gosymbol

import (
	"reflect"
)

/* Constrain functions */

// TODO: would be useful with a bool function negator,
// it takes a function with any input and with bool input
// and returns a function with the same input but with negated bool output
func positiveConstant(expr Expr) bool {
	exprTyped, ok := expr.(constant)
	return ok && exprTyped.Value > 0
} 

func negOrZeroConstant(expr Expr) bool {
	exprTyped, ok := expr.(constant)
	return ok && exprTyped.Value <= 0
}


var sumSimplificationRules []transformationRule = []transformationRule{
	{ // (Undefined) + ... = Undefined 
		patternFunction: func(expr Expr) bool {
			_, ok := expr.(add)
			return  ok && Contains(expr, Undefined())
		},
		transform: func(expr Expr) Expr {return Undefined()},
	},
	{ // Addition with only one operand simplify to the operand
		pattern: Add(Var("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Addition of no operands simplify to 0
		pattern: Add(),
		transform: func(expr Expr) Expr {
			return Const(0)		
		},
	},
}

var productSimplificationRules []transformationRule = []transformationRule{
	{ // (Undefined) * ... = Undefined 
		patternFunction: func(expr Expr) bool {
			_, ok := expr.(mul)
			return  ok && Contains(expr, Undefined())
		},
		transform: func(expr Expr) Expr {return Undefined()},
	},
	{ // 0 * ... = 0
		patternFunction: func(expr Expr) bool {
			// Ensures expr is of type mul
			_, ok := expr.(mul)
			if !ok {return false}

			// Returns true if any operand is 0
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				if reflect.DeepEqual(op, Const(0)) {
					return true
				}
			}

			return false
		},
		transform: func(expr Expr) Expr {return Const(0)},
	},
	{ // Multiplication with only one operand simplify to the operand
		pattern: Mul(Var("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Multiplication of no operands simplify to 1
		pattern: Mul(),
		transform: func(expr Expr) Expr {
			return Const(1)		
		},
	},
	{ // 1 * x = x
		pattern: Mul(Const(1), Var("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 2)
		},
	},
	{ // x*x = x^2
		pattern: Mul(Var("x"), Var("x")),
		transform: func(expr Expr) Expr {
			base := Operand(expr, 1)
			return Pow(base, Const(2))
		},
	},
	{ // x*x^n = x^(n+1)
		pattern: Mul(Var("x"), Pow(Var("x"), Var("y"))),
		transform: func(expr Expr) Expr {
			newBase := Operand(expr, 1)
			oldExponent := Operand(Operand(expr, 2), 2)
			newExponent := Add(oldExponent, Const(1))
			return Pow(newBase, newExponent)
		},
	},
}

var powerSimplificationRules []transformationRule = []transformationRule{
	{ // (Undefined)^x = Undefined
		pattern: Pow(Undefined(), Var("x")),
		transform: func(expr Expr) Expr {
			return Undefined()
		},
	},
	{ // 0^x = 0 for x in R_+
		pattern: Pow(Const(0), ConstrVar("x", positiveConstant)),
		transform: func(expr Expr) Expr {
			return Const(0) 
		},
	},
	{ // 0^x = Undefined for x <= 0
		pattern: Pow(Const(0), ConstrVar("x", negOrZeroConstant)),
		transform: func(expr Expr) Expr {
			return Undefined()
		},
	},
	{ // 1^x = 1
		pattern: Pow(Const(1), Var("x")),
		transform: func(expr Expr) Expr {
			return Const(1)
		},
	},
	{ // x^0 = 1
		pattern: Pow(Var("x"), Const(0)),
		transform: func(expr Expr) Expr {
			return Const(1)
		},
	},
	{ // (x^y)^z = x^(y*z)
		pattern: Pow(Pow(Var("x"), Var("y")), Var("z")),
		transform: func(expr Expr) Expr {	
			x := Operand( Operand(expr, 1), 1)
			y := Operand( Operand(expr, 1), 2)
			z := Operand(expr, 2)
			return Pow(x, Mul(y, z))
		},
	},
} 

