package gosymbol

const ()

/* Factories */

func Undefined() undefined {
	return undefined{}
}

func Const(val float64) constant {
	return constant{Value: val}
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

/*
Constructs the expression lhs/rhs.

Note: is lhs is one, the function returns a pow type,
otherwise it returns a mul type. This is to avoid unneccessary
calls to Simplify.
*/
func Div(lhs, rhs Expr) Expr {
	if Equal(lhs, Const(1)) {
		return Pow(rhs, Const(-1))
	} else {
		return Mul(lhs, Pow(rhs, Const(-1)))
	}
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

// Returns if expr is a basic algebraic expression (BAE).
//
// Following Definition 3.19 in COHEN, Joel S. Computer algebra and
// symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003,
// the expression u is a BAE if any of the following rules are satisfied:
//
// - BAE-1: u is an integer.
//
// - BAE-2: u is a fraction.
//
// - BAE-3: u is a symbol.
//
// - BAE-4: u is a product with one or more operands that are BAEs.
//
// - BAE-5: u is a sum with one or more operands that are BAEs.
//
// - BAE-6: u is a quotient with two operands that are BAEs.
//
// - BAE-7: u is a unary or binary difference where each operand is a BAE.
//
// - BAE-8: u is a power where both operands are BAE.
//
// - BAE-9: u is a factorial where the operand is a BAE.
//
// - BAE-10: u is a function form with one or more operands that are BAEs.
//
// The following experssions are BAE:
// 1. 2/4,
//
// 2. a * (x + x),
//
// 3. a + (b^3 / b),
//
// 4. b - 3 * b,
//
// 5. a + ( b + c ) + d,
//
// 6. 2 * 3 * x * x^2,
//
// 7. f(x)^1,
//
// 8. + x^2 - x,
//
// 9. 0^3,
//
// 10. * x,
//
// 11. 2 / (a - a),
//
// 12. 3!.
func IsBAE(expr Expr) bool {
	switch expr.(type) {
	case constant:
		// BAE-1
		return true
	case undefined:
		// BAE-3
		return true
	case variable:
		// BAE-3
		return true
	case mul:
		// BAE-4
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			if IsBAE(op) {
				return true
			}
		}
		return false
	case add:
		// BAE-5
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			if IsBAE(op) {
				return true
			}
		}
		return false
	case pow:
		base := Operand(expr, 1)
		exponent := Operand(expr, 2)
		return IsBAE(base) && IsBAE(exponent)
	default:
		return false
	}
}

// Returns if expr is an automatically simplified algebraic expression (ASAE).
//
// Following Definition 3.21 in COHEN, Joel S. Computer algebra and
// symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003,
// the expression u is an ASAE if any of the following are satisfied:
//
// - ASAE-1: u is an integer.
//
// - ASAE-2: u is a fraction on standard form.
//
// - ASAE-3: u is a symbol except the Undefined symbol
func IsASAE(expr Expr) bool {
	switch exprTyped := expr.(type) {
	case constant:
		//ASAE-1
		return true
	case variable:
		// ASAE-3
		return true
	case undefined:
		// ASAE-3
		return false
	case mul:
		// ASAE-4
		return mulIsASAE(exprTyped)
	case add:
		// ASAE-5
		return addIsASAE(exprTyped)
	case pow:
		// ASAE-6
		return powIsASAE(exprTyped)
	default:
		return false
	}
}

// ASAE-4
// TODO
func mulIsASAE(expr mul) bool {
	return false
}

// ASAE-5
// TODO
func addIsASAE(expr add) bool {
	return false
}

// ASAE-6
// TODO
func powIsASAE(expr pow) bool {
	return false
}
