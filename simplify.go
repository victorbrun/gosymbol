package gosymbol

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

var mulSimplificationRules []simplificationRule = []simplificationRule{}
var addSimplificationRules []simplificationRule = []simplificationRule{}

func (r simplificationRule) match(expr Expr) bool {
	// Rule concerns the top most operator in the 
	// tree so these need to have matching types.
	if !isSameType(r.lhs, expr) {
		return false
	} 

	// Iterate though each operand of both expr and r.lhs,
	// checking if they match "good enough".
	for ix := 1; ix <= NumberOfOperands(r.lhs); ix++ {
		ruleOperand := Operand(r.lhs, ix)
		exprOperand := Operand(expr, ix)
		if _, ok := ruleOperand.(variable); ok {
			continue
		} else if opTyped, ok := ruleOperand.(constrainedVariable); ok && opTyped.Constraint(exprOperand) {
			continue
		} else if Equal(ruleOperand, exprOperand) {
			continue
		} else {
			return false
		}
	}
	return true
}





