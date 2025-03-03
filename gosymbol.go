package gosymbol

import "math"

/* Factories */

func Undefined() undefined {
	return undefined{}
}

func Var(name VarName) variable {
	return variable{Name: name, isPattern: false}
}

func patternVar(name VarName) variable {
	return variable{Name: name, isPattern: true}
}

func ConstrVar(name VarName, constrFunc func(Expr) bool) constrainedVariable {
	return constrainedVariable{Name: name, Constraint: constrFunc, isPattern: false}
}

func constraPatternVar(name VarName, constrFunc func(Expr) bool) constrainedVariable {
	return constrainedVariable{Name: name, Constraint: constrFunc, isPattern: true}
}

func Neg(arg Expr) mul {
	return Mul(Int(-1), arg)
}

func Add(ops ...Expr) add {
	var newOps []Expr
	for _, op := range ops {
		switch opTyped := op.(type) {
		case add:
			newOps = append(newOps, opTyped.Operands...)
		default:
			newOps = append(newOps, op)
		}
	}
	return add{Operands: newOps}
}

func Sub(lhs, rhs Expr) add {
	return Add(lhs, Neg(rhs))
}

func Mul(ops ...Expr) mul {
	var newOps []Expr
	for _, op := range ops {
		switch opTyped := op.(type) {
		case mul:
			newOps = append(newOps, opTyped.Operands...)
		default:
			newOps = append(newOps, op)
		}
	}
	return mul{Operands: newOps}
}

/*
Constructs the expression lhs/rhs.

Note: is lhs is one, the function returns a pow type,
otherwise it returns a mul type. This is to avoid unneccessary
calls to Simplify.
*/
func Div(lhs, rhs Expr) Expr {
	num, okNum := lhs.(integer)
	den, okDen := rhs.(integer)
	if okNum && okDen {
		if intMul(num, den).value >= 0 {
			return fraction{num: intAbs(num), den: intAbs(den)}
		}
		return fraction{num: intNeg(intAbs(num)), den: intAbs(den)}
	}
	if Equal(rhs, Int(0)) {
		return undefined{}
	}
	if Equal(rhs, Int(1)) {
		return lhs
	}
	if Equal(lhs, Int(1)) {
		return Pow(rhs, Int(-1))
	}
	return Mul(lhs, Pow(rhs, Int(-1)))
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

func (args Arguments) AddArgument(v variable, value Expr) error {
	for arg := range args {
		if arg.Name == v.Name {
			return &DuplicateArgumentError{}
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
		bindings := make(Binding)
		return patternMatch(expr, rule.pattern, bindings)
	} else if rule.patternFunction != nil {
		return rule.patternFunction(expr)
	} else {
		return false
	}
}

func Int(value int64) integer {
	return integer{value: value}
}

func Real(symbol string, approxValue float64) variable {
	return variable{Name: VarName(symbol), isPattern: false, isConstant: true, constValue: approxValue}
}

var PI = Real("π", math.Pi)
var E = Exp(Int(1))
