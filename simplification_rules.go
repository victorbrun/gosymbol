package gosymbol

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


var additionSimplificationRules []transformationRule = []transformationRule{}
var multiplicationSimplificationRules []transformationRule = []transformationRule{}

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

