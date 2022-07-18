package gosymbol

var additionSimplificationRules []simplificationRule = []simplificationRule{}
var multiplicationSimplificationRules []simplificationRule = []simplificationRule{}

var powerSimplificationRules []simplificationRule = []simplificationRule{
	{ // (Undefined)^x = Undefined
		lhs: Pow(Undefined(), Var("x")),
		rhs: Undefined(),
	},
	{ // 0^x = 0 for x in R_+
		lhs: Pow(Const(0), constrainedVariable{
			Name: "x",
			Constraint: func(expr Expr) bool {
				exprTyped, ok := expr.(constant)
				return ok && exprTyped.Value > 0
			},
		}),
		rhs: Const(0),
	},
	{ // 0^x = Undefined for x <= 0
		lhs: Pow(Const(0), constrainedVariable{
			Name: "x",
			Constraint: func(expr Expr) bool {
				exprTyped, ok := expr.(constant)
				return ok && exprTyped.Value <= 0
			},
		}),
		rhs: Undefined(),
	},
	{ // 1^x = 1
		lhs: Pow(Const(1), Var("x")),
		rhs: Const(1),
	},
	{ // x^0 = 1
		lhs: Pow(Var("x"), Const(0)),
		rhs: Const(1),
	},
	{ // (x^y)^z = x^(y*z)
		lhs: Pow(Pow(Var("x"), Var("y")), Var("z")),
		rhs: Pow(Var("x"), Mul(Var("y"), Var("z"))),
	},
}


var exponentialSimplificationRules []simplificationRule = []simplificationRule{}
var logarithmicSimplificationRules []simplificationRule = []simplificationRule{}
