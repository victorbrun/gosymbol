package gosymbol

var PI = Real("π")
var E = Exp(Int(1))

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

func Int(value int64) integer {
	return integer{value: value}
}

func Real(symbol string) variable {
	return variable{Name: VarName(symbol), isPattern: false}
}

func Neg(arg Expr) mul {
	return Mul(Int(-1), arg)
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

// Constructs fraction
func Frac(nom integer, denom integer) fraction {
	return Div(nom, denom).(fraction)
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
