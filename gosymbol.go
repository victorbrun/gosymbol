package gosymbol

import "fmt"

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
	case integer:
		// BAE-1
		return true
	case fraction:
		// BAE-2
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
	case integer:
		//ASAE-1
		return true
	case fraction:
		//ASAE-2
		return Equal(exprTyped, exprTyped.simplifyRational())
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
// u is a product satisfying all of the following properties:
//
// 0. u has two or more operands u1 * u2 * ...
//
// 1. u has all admissible factors
//
// 2. At most one operand ui is a constant (integer or fraction)
//
// 3. If i != j, then AsaeBase(ui) != AsaeBase(uj)
//
// 4. If i < j, then compare(ui, uj) = true
func mulIsASAE(u mul) bool {
	// 0.
	if len(u.Operands) < 2 {
		return false
	}

	// 1.
	if !hasAllAdmissibleFactors(u) {
		return false
	}

	// 2.
	constCount := 0
	for _, factor := range u.Operands {
		switch factor.(type) {
		case integer:
			constCount++
		case fraction:
			constCount++
		}
	}
	if constCount > 0 {
		return false
	}

	// 3. and 4.
	for ix := 1; ix <= len(u.Operands); ix++ {
		ui := Operand(u, ix)
		for jx := ix + 1; jx <= len(u.Operands); jx++ {
			uj := Operand(u, jx)
			if Equal(AsaeBase(ui), AsaeBase(uj)) {
				return false
			} else if !compare(ui, uj) {
				return false
			}
		}
	}

	// If we have not returned before arriving here,
	// every property is satisfied
	return true
}

// Returns if expr has all admissible factors.
//
// A factor is said to be admissible if it
// is an ASAE which can be either an integer
// (!= 0, 1), fraction, symbol (except undefined),
// sum, power, function.
//
// Note: the factor of a product cannot be a product
// for it to be admissible.
func hasAllAdmissibleFactors(expr mul) bool {
	// Iterating over each factor and checking if
	// it is a admissible factor. Returning false
	// at first non-admissible factor
	for _, factor := range expr.Operands {
		switch factor.(type) {
		case undefined:
			return false
		case mul:
			return false
		case integer:
			if Equal(factor, Int(0)) || Equal(factor, Int(1)) {
				return false
			}
		default:
			if !IsASAE(factor) {
				return false
			}
		}
	}
	return true
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

// Returns the base of an ASAE expression
//
// Examples:
//
// 1. AsaeBase(x^2) = x
//
// 2. AsaeBase(x) = x
func AsaeBase(expr Expr) Expr {
	switch expr.(type) {
	case pow:
		return Operand(expr, 1)
	case integer:
		return Undefined()
	case fraction:
		return Undefined()
	default:
		return expr
	}
}

// Returns the base of an ASAE expression
//
// Examples:
//
// 1. AsaeExponent(x^2) = 2
//
// 2. AsaeExponent(x) = 1
func AsaeExponent(expr Expr) Expr {
	switch expr.(type) {
	case pow:
		return Operand(expr, 2)
	case integer:
		return Undefined()
	case fraction:
		return Undefined()
	default:
		return Int(1)
	}
}

// Returns the term from a ASAE multiplication.
//
// Examples:
//
// 1. Term(x) = *x
//
// 2. Term(2*y) = *y
//
// 3. Term(x*y) = x*y
func AsaeTerm(expr Expr) Expr {
	// The steps taken in this function only works
	// if the input expression is already an ASAE
	if !IsASAE(expr) {
		panic(fmt.Sprintf("Expression is not an ASAE: %v", expr))
	}

	switch exprTyped := expr.(type) {
	case integer:
		return Undefined()
	case fraction:
		return Undefined()
	case mul:
		op1 := Operand(expr, 1)
		switch op1.(type) {
		case integer:
			return Mul(exprTyped.Operands[1:]...)
		case fraction:
			return Mul(exprTyped.Operands[1:]...)
		default:
			return expr
		}
	default:
		return Mul(expr)
	}

}

// Returns the constant from a ASAE multiplication.
//
// Examples:
//
// 1. Term(x) = 1
//
// 2. Term(2*y) = 2
//
// 3. Term(x*y) = 1
func AsaeConst(expr Expr) Expr {
	if !IsASAE(expr) {
		panic(fmt.Sprintf("Expression is not an ASAE: %v", expr))
	}

	switch expr.(type) {
	case integer:
		return Undefined()
	case fraction:
		return Undefined()
	case mul:
		op1 := Operand(expr, 1)
		switch op1.(type) {
		case integer:
			return op1
		case fraction:
			return op1
		default:
			return Int(1)
		}
	default:
		return Int(1)
	}
}
