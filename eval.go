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
			return Simplify(value)
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
		return Simplify(sum)
	}
}

func (e mul) Eval() Func {
	return func(args Arguments) Expr {
		prod := e.Operands[0].Eval()(args) // Initiate with first operand since 1 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			prod = Mul(prod, e.Operands[ix].Eval()(args))
		}
		return Simplify(prod)
	}
}

func (e exp) Eval() Func {
	return func(args Arguments) Expr { return Simplify(Exp(e.Arg.Eval()(args))) }
}

func (e log) Eval() Func {
	return func(args Arguments) Expr { return Simplify(Log(e.Arg.Eval()(args))) }
}

func (e pow) Eval() Func {
	return func(args Arguments) Expr {
		return Simplify(Pow(e.Base.Eval()(args), e.Exponent.Eval()(args)))
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
