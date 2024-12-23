package gosymbol

func (e constant) D(varName VarName) Expr {
	return Const(0.0)
}

func (e variable) D(varName VarName) Expr {
	if varName == e.Name {
		return Const(1.0)
	} else {
		return Const(0.0)
	}
}

func (e add) D(varName VarName) Expr {
	differentiatedOps := make([]Expr, len(e.Operands))
	for ix, op := range e.Operands {
		differentiatedOps[ix] = op.D(varName)
	}
	return Add(differentiatedOps...)
}

// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
func (e mul) D(varName VarName) Expr {
	terms := make([]Expr, len(e.Operands))
	for ix := 0; ix < len(e.Operands); ix++ {
		var productOperands []Expr
		copy(productOperands, e.Operands)
		productOperands[ix] = productOperands[ix].D(varName)
		terms[ix] = Mul(productOperands...)
	}
	return Add(terms...)
}

func (e exp) D(varName VarName) Expr {
	return Mul(e, e.Arg.D(varName))
}

func (e log) D(varName VarName) Expr {
	return Mul(Pow(e.Arg, Const(-1)), e.Arg.D(varName))
}

// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
func (e pow) D(varName VarName) Expr {
	if exponentTyped, ok := e.Exponent.(constant); ok {
		return Mul(e.Exponent, Pow(e.Base, Const(exponentTyped.Value-1)), e.Base.D(varName))
	} else {
		exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
		return Mul(Exp(exponentLogBaseProd), exponentLogBaseProd.D(varName))
	}
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(varName VarName) Expr {
	return Mul(Div(Const(1), Const(2)), Div(Const(1), e), e.Arg.D(varName))
}
