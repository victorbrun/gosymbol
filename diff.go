package gosymbol

func (e constant) D(v variable) Expr {
	return Const(0.0)
}

func (e variable) D(v variable) Expr {
	if v == e {
		return Const(1.0)
	} else {
		return Const(0.0)
	}
}

func (e add) D(v variable) Expr {
	differentiatedOps := make([]Expr, len(e.Operands))
	for ix, op := range e.Operands {
		differentiatedOps[ix] = op.D(v)
	}
	return Add(differentiatedOps...)
}

// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
func (e mul) D(v variable) Expr {
	terms := make([]Expr, len(e.Operands))
	for ix := 0; ix < len(e.Operands); ix++ {
		productOperands := make([]Expr, len(e.Operands))
		copy(productOperands, e.Operands)
		productOperands[ix] = productOperands[ix].D(v)
		terms[ix] = Mul(productOperands...)
	}
	return Add(terms...)
}

func (e exp) D(v variable) Expr {
	return Mul(e, e.Arg.D(v))
}

func (e log) D(v variable) Expr {
	return Mul(Pow(e.Arg, Const(-1)), e.Arg.D(v))
}

// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
func (e pow) D(v variable) Expr {
	if exponentTyped, ok := e.Exponent.(constant); ok {
		return Mul(e.Exponent, Pow(e.Base, Const(exponentTyped.Value-1)), e.Base.D(v))
	} else {
		exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
		return Mul(Exp(exponentLogBaseProd), exponentLogBaseProd.D(v))
	}
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(v variable) Expr {
	return Mul(Div(Const(1), Const(2)), Div(Const(1), e), e.Arg.D(v))
}
