package gosymbol

import "fmt"

func (e integer) D(v variable) Expr {
	return differentiate(e, v)
}

func (e fraction) D(v variable) Expr {
	return differentiate(e, v)
}

func (e variable) D(v variable) Expr {
	return differentiate(e, v)
}

func (e add) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

func (e mul) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

func (e exp) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

func (e log) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

func (e pow) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(v variable) Expr {
	return differentiate(e, v).Simplify()
}

/*
Differentiates expr w.r.t. v.
*/
func differentiate(expr Expr, v variable) Expr {
	switch e := expr.(type) {
	case integer:
		return Int(0)
	case fraction:
		return Int(0)
	case variable:
		if v == e {
			return Int(1)
		} else {
			return Int(0)
		}

	case add:
		differentiatedOps := make([]Expr, len(e.Operands))
		for ix, op := range e.Operands {
			differentiatedOps[ix] = differentiate(op, v)
		}
		return Add(differentiatedOps...)

	case mul:
		// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
		terms := make([]Expr, len(e.Operands))
		for ix := 0; ix < len(e.Operands); ix++ {
			productOperands := make([]Expr, len(e.Operands))
			copy(productOperands, e.Operands)
			productOperands[ix] = differentiate(productOperands[ix], v)
			terms[ix] = Mul(productOperands...)
		}
		return Add(terms...)

	case pow:
		// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
		// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
		switch exponentTyped := e.Exponent.(type) {
		case integer:
			return Mul(e.Exponent, Pow(e.Base, intAdd(exponentTyped, Int(-1))), differentiate(e.Base, v))
		case fraction:
			return Mul(e.Exponent, Pow(e.Base, ratAdd(exponentTyped, Int(-1))), differentiate(e.Base, v))
		default:
			exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
			return Mul(Exp(exponentLogBaseProd), differentiate(exponentLogBaseProd, v))
		}
	case exp:
		return Mul(e, differentiate(e.Arg, v))

	case log:
		return Mul(Pow(e.Arg, Int(-1)), differentiate(e.Arg, v))

	default:
		errMsg := fmt.Errorf("ERROR: expression %#v have no differentiation pattern case implemented", e)
		panic(errMsg)
	}
}
