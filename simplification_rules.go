package gosymbol

import (
	"reflect"
)

/* Constrain functions */

// TODO: would be useful with a bool function negator,
// it takes a function with any input and with bool input
// and returns a function with the same input but with negated bool output
func positiveConstant(expr Expr) bool {
	exprTyped, ok := expr.(rational)
	return ok && exprTyped.approx() > Int(0).approx()
}

func negOrZeroConstant(expr Expr) bool {
	exprTyped, ok := expr.(rational)
	return ok && exprTyped.approx() <= Int(0).approx()
}

var sumSimplificationRules []transformationRule = []transformationRule{
	{ // Addition with only one operand simplify to the operand
		pattern: Add(patternVar("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Addition of no operands simplify to 0
		pattern: Add(),
		transform: func(expr Expr) Expr {
			return (Int(0))
		},
	},
	{ // x - x = 0. Due to the ordering the negative term will always be first.
		// Note that this will not work for constants since (-c) is a float
		// while -x = -1*x.
		pattern: Add(Neg(patternVar("x")), patternVar("x")),
		transform: func(expr Expr) Expr {
			return (Int(0))
		},
	},
	{ // x + x = 2x.
		pattern: Add(patternVar("x"), patternVar("x")),
		transform: func(expr Expr) Expr {
			return Mul(Int(2), Operand(expr, 1))
		},
	},
	{ // x + yx = (1+y)x.
		pattern: Add(patternVar("x"), Mul(patternVar("y"), patternVar("x"))),
		transform: func(expr Expr) Expr {
			return Mul(Int(2), Operand(expr, 1))
		},
	},
	{ // yx + zx = 2x.
		pattern: Add(Mul(patternVar("z"), patternVar("x")), Mul(patternVar("y"), patternVar("x"))),
		transform: func(expr Expr) Expr {
			return Mul(Int(2), Operand(expr, 1))
		},
	},
	{ // Sum of constants is replaced with the constant that the sum evaluates to.
		// Note that sum of some constants will replace the constants with their sum.
		patternFunction: func(expr Expr) bool {
			// Ensures expr is of type add
			_, ok := expr.(add)
			if !ok {
				return false
			}

			// Makes sure that there is at lest
			// two terms which are constants to avoid
			// getting stuck in infinite loop
			if NumberOfOperands(expr) > 1 {
				op1 := Operand(expr, 1)
				op2 := Operand(expr, 2)
				_, ok1 := op1.(rational)
				_, ok2 := op2.(rational)
				return ok1 && ok2
			} else {
				return false
			}
		},
		transform: func(expr Expr) Expr {
			// We sum all the constants in the sum
			var sum rational
			sum = Int(0)
			nons := make([]Expr, 0)
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				switch opTyped := op.(type) {
				case integer:
					sum = ratAdd(sum, opTyped)
				case fraction:
					sum = ratAdd(sum, opTyped)
				default:
					nons = append(nons, op)
				}
			}

			// Replaces the constants with their sum. Note that we
			// know that the expression is sorted, i.e. the constants
			// are at the front of the sum.
			newTerms := append([]Expr{(sum)}, nons...)
			return Add(newTerms...)
		},
	},
}

var productSimplificationRules []transformationRule = []transformationRule{
	{ // 0 * ... = 0
		patternFunction: func(expr Expr) bool {
			// Ensures expr is of type mul
			_, ok := expr.(mul)
			if !ok {
				return false
			}

			// Returns true if any operand is 0
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				if reflect.DeepEqual(op, (Int(0))) {
					return true
				}
			}

			return false
		},
		transform: func(expr Expr) Expr { return (Int(0)) },
	},
	{ // Multiplication with only one operand simplify to the operand
		pattern: Mul(patternVar("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Multiplication of no operands simplify to 1
		pattern: Mul(),
		transform: func(expr Expr) Expr {
			return (Int(1))
		},
	},
	{ // 1 * x_1 * ... * x_n = x_1 * ... * x_n
		patternFunction: func(expr Expr) bool {
			// If expr is not mul, we can directly
			// return false
			exprTyped, ok := expr.(mul)
			if !ok {
				return false
			}

			// Checks if the arguments of expr contains Const(1)
			// Note: we cannot use RecContains as it would return true
			// on any 1, e.g, even on x^1.
			one := Int(1)
			for _, elem := range exprTyped.Operands {
				if Equal(elem, one) {
					return true
				}
			}
			return false
		},
		transform: func(expr Expr) Expr {
			if NumberOfOperands(expr) == 1 {
				return expr
			}

			// To account for possibility of more than one 1
			// we append non-1s
			var newFactors []Expr
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				if !Equal(op, (Int(1))) {
					newFactors = append(newFactors, op)
				}
			}
			return Mul(newFactors...)
		},
	},
	{ // x*x = x^2
		pattern: Mul(patternVar("x"), patternVar("x")),
		transform: func(expr Expr) Expr {
			base := Operand(expr, 1)
			return Pow(base, (Int(2)))
		},
	},
	{ // x*x^n = x^(n+1) this applies to positive n due to the ordering of an expression
		pattern: Mul(patternVar("x"), Pow(patternVar("x"), patternVar("y"))),
		transform: func(expr Expr) Expr {
			newBase := Operand(expr, 1)
			oldExponent := Operand(Operand(expr, 2), 2)
			newExponent := Add(oldExponent, (Int(1)))
			return Pow(newBase, newExponent)
		},
	},
	{ // x^n * x = x^(n+1) this applies to negative n due to the ordering of an expression
		pattern: Mul(Pow(patternVar("x"), patternVar("y")), patternVar("x")),
		transform: func(expr Expr) Expr {
			newBase := Operand(expr, 2)
			oldExponent := Operand(Operand(expr, 1), 2)
			newExponent := Add(oldExponent, (Int(1)))
			return Pow(newBase, newExponent)
		},
	},
	{ // x^(n) * x^(m) = x^(n+m)
		pattern: Mul(Pow(patternVar("x"), patternVar("n")), Pow(patternVar("x"), patternVar("m"))),
		transform: func(expr Expr) Expr {
			base := Operand(Operand(expr, 1), 1)
			exponent1 := Operand(Operand(expr, 1), 2)
			exponent2 := Operand(Operand(expr, 2), 2)
			return Pow(base, Add(exponent1, exponent2))
		},
	},
	{ // Prod of constants is replaced with the constant that the product evaluates to.
		// Note that product of some constants will replace the constants with their product.
		patternFunction: func(expr Expr) bool {
			_, ok := expr.(mul)
			if !ok {
				return false
			}
			// Makes sure that there is at lest
			// two terms which are constants to avoid
			// getting stuck in infinite loop
			if NumberOfOperands(expr) > 1 {
				op1 := Operand(expr, 1)
				op2 := Operand(expr, 2)
				_, ok1 := op1.(rational)
				_, ok2 := op2.(rational)
				return ok1 && ok2
			} else {
				return false
			}
		},
		transform: func(expr Expr) Expr {
			// We multiply all the constants in the product
			var prod rational
			prod = Int(1)
			nons := make([]Expr, 0)
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				switch opTyped := op.(type) {
				case integer:
					prod = ratMul(prod, opTyped)
				case fraction:
					prod = ratMul(prod, opTyped)
				default:
					nons = append(nons, op)
				}
			}
			// the expression is sorted, i.e. the constants
			// are at the front of the sum.
			newTerms := append([]Expr{(prod)}, nons...)
			return Mul(newTerms...)
		},
	},
}

var powerSimplificationRules []transformationRule = []transformationRule{
	{ // 0^x = 0 for x in R_+
		pattern: Pow((Int(0)), constraPatternVar("x", positiveConstant)),
		transform: func(expr Expr) Expr {
			return (Int(0))
		},
	},
	{ // 0^x = Undefined for x <= 0
		pattern: Pow((Int(0)), constraPatternVar("x", negOrZeroConstant)),
		transform: func(expr Expr) Expr {
			return Undefined()
		},
	},
	{ // 1^x = 1
		pattern: Pow((Int(1)), patternVar("x")),
		transform: func(expr Expr) Expr {
			return (Int(1))
		},
	},
	{ // x^1 = x
		pattern: Pow(patternVar("x"), Int(1)),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // x^0 = 1
		pattern: Pow(patternVar("x"), Int(0)),
		transform: func(expr Expr) Expr {
			return Int(1)
		},
	},
	{ // (v_1 * ... * v_m)^n = v_1^n * ... * v_m^n
		patternFunction: func(expr Expr) bool {
			if _, ok := expr.(pow); ok {
				base := Operand(expr, 1)
				_, ok := base.(mul)
				return ok
			}
			return false
		},
		transform: func(expr Expr) Expr {
			base := Operand(expr, 1)
			exponent := Operand(expr, 2)
			for ix := 1; ix <= NumberOfOperands(base); ix++ {
				factor := Operand(base, ix)
				base = replaceOperand(base, ix, Pow(factor, exponent))
			}
			return base
		},
	},
	{ // (x^y)^z = x^(y*z)
		pattern: Pow(Pow(patternVar("x"), patternVar("y")), patternVar("z")),
		transform: func(expr Expr) Expr {
			x := Operand(Operand(expr, 1), 1)
			y := Operand(Operand(expr, 1), 2)
			z := Operand(expr, 2)
			return Pow(x, Mul(y, z))
		},
	},
	{ // Prod of constants is replaced with the constant that the product evaluates to.
		// Note that product of some constants will replace the constants with their product.
		patternFunction: func(expr Expr) bool {
			power, ok := expr.(pow)
			if ok {
				_, ok1 := power.Base.(rational)
				_, ok2 := power.Exponent.(integer)
				return ok1 && ok2
			}
			return false
		},
		transform: func(expr Expr) Expr {
			power := expr.(pow)
			return ratPow(power.Base.(rational), power.Exponent.(integer))
		},
	},
}

var expSimplificationRules []transformationRule = []transformationRule{
	{ // e^0 = 1
		pattern:   Exp(Int(0)),
		transform: func(expr Expr) Expr { return Int(1) },
	},
}

var logSimplificationRules []transformationRule = []transformationRule{
	{ // log(1) =0
		pattern:   Log(Int(1)),
		transform: func(expr Expr) Expr { return Int(0) },
	},
}
