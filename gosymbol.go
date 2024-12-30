package gosymbol

import (
	"errors"
)

/* Factories */

func Undefined() undefined {
	return undefined{}
}

func Const(val float64) constant {
	return constant{Value: val}
}

func Var(name VarName) variable {
	return variable{Name: name}
}

func ConstrVar(name VarName, constrFunc func(Expr) bool) constrainedVariable {
	return constrainedVariable{Name: name, Constraint: constrFunc}
}

func Neg(arg Expr) mul {
	return Mul(Const(-1), arg)
}

func Add(ops ...Expr) add {
	return add{Operands: ops}
}

func Sub(lhs, rhs Expr) add {
	return Add(lhs, Neg(rhs))
}

func Mul(ops ...Expr) mul {
	return mul{Operands: ops}
}

func Div(lhs, rhs Expr) mul {
	return Mul(lhs, Pow(rhs, Const(-1)))
}

func Exp(arg Expr) exp {
	return exp{Arg: arg}
}

func Log(arg Expr) log {
	return log{Arg: arg}
}

func Pow(base Expr, exponent Expr) pow {
	return pow{Base: base, Exponent: exponent}
}

func Sqrt(arg Expr) sqrt {
	return sqrt{Arg: arg}
}

func TransformationRule(pattern Expr, transform func(Expr) Expr) transformationRule {
	return transformationRule{pattern: pattern, transform: transform}
}

func (args Arguments) AddArgument(v variable, value float64) error {
	for arg := range args {
		if arg.Name == v.Name {
			return errors.New("multiple variables have the same name")
		}
	}
	args[v] = value
	return nil
}

// Applies rule to expr and returns the transformed expression.
// If expression does not match rule the ingoing expression
// will just be returned.
func (rule transformationRule) apply(expr Expr) (Expr, bool) {
	if rule.match(expr) {
		return rule.transform(expr), true
	}
	return expr, false
}

func (rule transformationRule) match(expr Expr) bool {
	// Fisrt check if pattern is defined. If not
	// we execute patternFunction if it exists.
	// If no pattern or patternFunction exists we return false
	if rule.pattern != nil {
		varCache := make(map[VarName]Expr)
		return patternMatch(rule.pattern, expr, varCache)
	} else if rule.patternFunction != nil {
		return rule.patternFunction(expr)
	} else {
		return false
	}
}
