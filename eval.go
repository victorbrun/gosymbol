package gosymbol

import (
	"fmt"
)

/* Evaluation */

func (e undefined) Eval() Func {
	return func(args Arguments) Expr { return undefined{} }
}

func (e integer) Eval() Func {
	return func(args Arguments) Expr { return e }
}

func (e fraction) Eval() Func {
	return func(args Arguments) Expr { return e.simplifyRational() }
}

func (e variable) Eval() Func {
	return func(args Arguments) Expr {
		value, ok := args[e]
		if ok {
			return value.Simplify()
		}
		return e
	}
}

func (e add) Eval() Func {
	return func(args Arguments) Expr {
		sum := e.Operands[0].Eval()(args) // Initiate with first operand since 0 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			sum = Add(sum, e.Operands[ix].Eval()(args))
		}
		return sum.Simplify()
	}
}

func (e mul) Eval() Func {
	return func(args Arguments) Expr {
		prod := e.Operands[0].Eval()(args) // Initiate with first operand since 1 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			prod = Mul(prod, e.Operands[ix].Eval()(args))
		}
		return prod.Simplify()
	}
}

func (e exp) Eval() Func {
	return func(args Arguments) Expr { return Exp(e.Arg.Eval()(args)).Simplify() }
}

func (e log) Eval() Func {
	return func(args Arguments) Expr { return Log(e.Arg.Eval()(args)).Simplify().Simplify() }
}

func (e pow) Eval() Func {
	return func(args Arguments) Expr {
		return Pow(e.Base.Eval()(args), e.Exponent.Eval()(args)).Simplify()
	}
}

/* Evaluation to string */

func (e undefined) String() string {
	return "Undefined"
}

func (e variable) String() string {
	return string(e.Name)
}

func (e constrainedVariable) String() string {
	return fmt.Sprintf("%v_CONSTRAINED", e.Name)
}

func (e add) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" + %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e mul) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" * %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e exp) String() string {
	return fmt.Sprintf("exp( %v )", e.Arg)
}

func (e log) String() string {
	return fmt.Sprintf("log( %v )", e.Arg)
}

func (e pow) String() string {
	return fmt.Sprintf("( %v^%v )", e.Base, e.Exponent)
}
