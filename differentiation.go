package gosymbol

/* Differentiation rules */

/*
Differentiates `e` w.r.t. `var` and returns the simplified resulting expression.
*/
func D(e Expr, varName VarName) Expr {
	derivative := e.d(varName)
	return Simplify(derivative)
}

func (e constant) d(varName VarName) Expr {
	return Const(0.0)
}

func (e variable) d(varName VarName) Expr {
	if varName == e.Name {
		return Const(1.0)
	} else {
		return Const(0.0)
	}
}

func (e add) d(varName VarName) Expr {
	differentiatedOps := make([]Expr, len(e.Operands))
	for ix, op := range e.Operands {
		differentiatedOps[ix] = op.d(varName)	
	}
	return Add(differentiatedOps...)
}

// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
func (e mul) d(varName VarName) Expr {
	terms := make([]Expr, len(e.Operands))
	for ix := 0; ix < len(e.Operands); ix++ {
		productOperands := make([]Expr, len(e.Operands))
		copy(productOperands, e.Operands)
		productOperands[ix] = productOperands[ix].d(varName)
		terms[ix] = Mul(productOperands...)
	}
	return Add(terms...)
}

func (e exp) d(varName VarName) Expr {
	return Mul(e, e.Arg.d(varName))
}

func (e log) d(varName VarName) Expr {
	return Mul(Pow(e.Arg, Const(-1)), e.Arg.d(varName))
}

// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
func (e pow) d(varName VarName) Expr {
	if exponentTyped, ok := e.Exponent.(constant); ok {
		return Mul(e.Exponent, Pow(e.Base, Const(exponentTyped.Value-1)), e.Base.d(varName))
	} else {
		exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
		return Mul(Exp(exponentLogBaseProd), exponentLogBaseProd.d(varName))
	}
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) d(varName VarName) Expr {
	return Mul(Div(Const(1), Const(2)), Div(Const(1), e), e.Arg.d(varName))
}
