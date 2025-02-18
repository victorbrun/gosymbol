package gosymbol

func (e integer) D(v variable) Expr {
	return Int(0)
}

func (e fraction) D(v variable) Expr {
	return Int(0)
}

func (e variable) D(v variable) Expr {
	if v == e {
		return Int(1)
	} else {
		return Int(0)
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
	return Mul(Pow(e.Arg, Int(-1)), e.Arg.D(v))
}

// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
func (e pow) D(v variable) Expr {
	switch exponentTyped := e.Exponent.(type) {
	case integer:
		return Mul(e.Exponent, Pow(e.Base, Add(exponentTyped, Int(-1))), e.Base.D(v))
	case fraction:
		return Mul(e.Exponent, Pow(e.Base, Add(exponentTyped, Int(-1))), e.Base.D(v))
	default:
		exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
		return Mul(Exp(exponentLogBaseProd), exponentLogBaseProd.D(v))
	}
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(v variable) Expr {
	return Mul(Div(Int(1), Int(2)), Div(Int(1), e), e.Arg.D(v))
}
