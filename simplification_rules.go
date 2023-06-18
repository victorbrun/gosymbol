package gosymbol

import (
	"reflect"
)

/* Constrain functions */

// TODO: would be useful with a bool function negator,
// it takes a function with any input and with bool input
// and returns a function with the same input but with negated bool output
func positiveConstant(expr Expr) bool {
	exprTyped, ok := expr.(constant)
	return ok && exprTyped.Value > 0
} 

func negOrZeroConstant(expr Expr) bool {
	exprTyped, ok := expr.(constant)
	return ok && exprTyped.Value <= 0
}


var sumSimplificationRules []transformationRule = []transformationRule{
	{ // Addition with only one operand simplify to the operand
		pattern: Add(CacheVar("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Addition of no operands simplify to 0
		pattern: Add(),
		transform: func(expr Expr) Expr {
			return Const(0)		
		},
	},
	{ // x - x = 0. Due to the ordering the negative term will always be first.
	  // Note that this will not work for constants since Const(-c) is a float
	  // while -x = -1*x. 
		pattern: Add(Neg(CacheVar("x")), CacheVar("x")),
		transform: func(expr Expr) Expr {
			return Const(0)
		},
	},
	{ // Sum of constants is replaced with the constant that the sum evaluates to.
	  // Note that sum of some constants will replace the constants with their sum. 
		patternFunction: func(expr Expr) bool {
			// Ensures expr is of type add 
			_, ok := expr.(add)
			if !ok {return false}

			// Makes sure that there is at lest
			// two terms which are constants to avoid
			// getting stuck in infinite loop
			if NumberOfOperands(expr) > 1 {
				op1 := Operand(expr, 1)
				op2 := Operand(expr, 2)
				_, ok1 := op1.(constant)
				_, ok2 := op2.(constant)
				return ok1 && ok2
			} else {
				return false
			}
		},
		transform: func(expr Expr) Expr {
			// We sum all the constants in the sum
			sum := 0.0
			nonConsts := make([]Expr, 0)
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				opTyped, ok := op.(constant)
				if ok {
					sum += opTyped.Value
				} else {
					nonConsts = append(nonConsts, op)
				}
			}
			
			// Replaces the constants with their sum. Note that we 
			// know that the expression is sorted, i.e. the constants 
			// are at the front of the sum.
			newTerms := append([]Expr{Const(sum)}, nonConsts...)
			return Add(newTerms...)
		},
	},
	{ // x + x = 2*x
		pattern: Add(CacheVar("x"), CacheVar("x")),
		transform: func(expr Expr) Expr { 
			x := Operand(expr, 1)
			return Mul(Const(2), x)
		},	
	},
	{ // n*x + x = (1+n)*x this applies to positive n due to the ordering of an expression
		pattern: Add(Mul(CacheVar("n"), CacheVar("x")), CacheVar("x")),
		transform: func(expr Expr) Expr {
			multiplication := Operand(expr, 1)
			commonVar := Operand(expr, 2)
			oldFactor := Operand(multiplication, 1)
			newFactor := Add(oldFactor, Const(1))
			return Mul(newFactor, commonVar)
		},
	},
	{ // x + n*x = (1+n)*x this applies to negative n due to the ordering of an expression
		pattern: Add(CacheVar("x"), Mul(CacheVar("n"), CacheVar("x"))),
		transform: func(expr Expr) Expr {
			multiplication := Operand(expr, 2)
			commonVar := Operand(expr, 1)
			oldFactor := Operand(multiplication, 1)
			newFactor := Add(oldFactor, Const(1))
			return Mul(newFactor, commonVar)
		},
	},
	{ // n*x + m*x = (n+m)*x
		pattern: Add(Mul(CacheVar("n"), CacheVar("x")), Mul(CacheVar("m"), CacheVar("x"))),
		transform: func(expr Expr) Expr {
			term1 := Operand(expr, 1)
			factor1 := Operand(term1, 1)
			term2 := Operand(expr, 2)
			factor2 := Operand(term2, 1)
			commonVar := Operand(term1, 2)
			return Mul(Add(factor1, factor2), commonVar)
		},
	},
	/*{ //  When applied recursively, this turns a flat sum into a binary tree sum. 
	  // This is needed for other rules to do the simplification: \sum_{i=1}^{n} x = n*x.
		patternFunction: func(expr Expr) bool {
			if _, ok := expr.(add); !ok {
				return false 
			} else if NumberOfOperands(expr) <= 2 {
				return false
			}
			return true
		},
		transform: func(expr Expr) Expr {
			newOperands := make([]Expr, 0)
			// Groups the operands into pairs and creates add of them
			for ix := 1; ix < NumberOfOperands(expr); ix+=2 {
				op1 := Operand(expr, ix)
				op2 := Operand(expr, ix+1)
				newOperands = append(newOperands, Add(op1, op2))
			}

			// If number of operands is odd there will be one
			// trailing operand so we add it here
			if NumberOfOperands(expr) % 2 != 0 {
				op := Operand(expr, NumberOfOperands(expr))
				newOperands = append(newOperands, op)
			}
			return Add(newOperands...)
		},
	},*/
}
var productSimplificationRules []transformationRule = []transformationRule{
	{ // 0 * ... = 0
		patternFunction: func(expr Expr) bool {
			// Ensures expr is of type mul
			_, ok := expr.(mul)
			if !ok {return false}

			// Returns true if any operand is 0
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				if reflect.DeepEqual(op, Const(0)) {
					return true
				}
			}

			return false
		},
		transform: func(expr Expr) Expr {return Const(0)},
	},
	{ // Multiplication with only one operand simplify to the operand
		pattern: Mul(CacheVar("x")),
		transform: func(expr Expr) Expr {
			return Operand(expr, 1)
		},
	},
	{ // Multiplication of no operands simplify to 1
		pattern: Mul(),
		transform: func(expr Expr) Expr {
			return Const(1)		
		},
	},
	{ // 1 * x_1 * ... * x_n = x_1 * ... * x_n 
		patternFunction: func(expr Expr) bool {
			// Ensures that the operator is multiplication 
			if _, ok := expr.(mul); !ok { return false }
			
			// Checks if any of the top level operands are 
			// Const(1), if so we have a match
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				op := Operand(expr, ix)
				if Equal(op, Const(1)) { return true }
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
				if !Equal(op, Const(1)) {
					newFactors = append(newFactors, op)
				}
			}
			return Mul(newFactors...)
		},
	},
	{ // x*x = x^2
		pattern: Mul(CacheVar("x"), CacheVar("x")),
		transform: func(expr Expr) Expr {
			base := Operand(expr, 1)
			return Pow(base, Const(2))
		},
	},
	{ // x*x^n = x^(n+1) this applies to positive n due to the ordering of an expression
		pattern: Mul(CacheVar("x"), Pow(CacheVar("x"), CacheVar("y"))),
		transform: func(expr Expr) Expr {
			newBase := Operand(expr, 1)
			oldExponent := Operand(Operand(expr, 2), 2)
			newExponent := Add(oldExponent, Const(1))
			return Pow(newBase, newExponent)
		},
	},
	{ // x^n * x = x^(n+1) this applies to negative n due to the ordering of an expression
		pattern: Mul(Pow(CacheVar("x"), CacheVar("y")), CacheVar("x")),
		transform: func(expr Expr) Expr {
			newBase := Operand(expr, 2)
			oldExponent := Operand(Operand(expr, 1), 2)
			newExponent := Add(oldExponent, Const(1))
			return Pow(newBase, newExponent)
		},
	},
	{ // x^(n) * x^(m) = x^(n+m)
		pattern: Mul(Pow(CacheVar("x"), CacheVar("n")), Pow(CacheVar("x"), CacheVar("m"))),
		transform: func(expr Expr) Expr {
			base := Operand(Operand(expr, 1), 1)
			exponent1 := Operand(Operand(expr, 1), 2)
			exponent2 := Operand(Operand(expr, 2), 2)
			return Pow(base, Add(exponent1, exponent2))
		},
	},
	/*{ //  When applied recursively, this turna a flat product into a binary tree product. 
	  // This is needed for other rules to do the simplification: \prod_{i=1}^{n} x = x^n.
		patternFunction: func(expr Expr) bool {
			if _, ok := expr.(mul); !ok {
				return false 
			} else if NumberOfOperands(expr) <= 2 {
				return false
			}
			return true
		},
		transform: func(expr Expr) Expr {
			newOperands := make([]Expr, 0)
			// Groups the operands into pairs and creates mul of them
			for ix := 1; ix < NumberOfOperands(expr); ix+=2 {
				op1 := Operand(expr, ix)
				op2 := Operand(expr, ix+1)
				newOperands = append(newOperands, Mul(op1, op2))
			}

			// If number of operands is odd there will be one
			// trailing operand so we add it here
			if NumberOfOperands(expr) % 2 != 0 {
				op := Operand(expr, NumberOfOperands(expr))
				newOperands = append(newOperands, op)
			}
			return Mul(newOperands...)
		},
	},*/
}

var powerSimplificationRules []transformationRule = []transformationRule{
	{ // 0^x = 0 for x in R_+
		pattern: Pow(Const(0), ConstrCacheVar("x", positiveConstant)),
		transform: func(expr Expr) Expr {
			return Const(0) 
		},
	},
	{ // 0^x = Undefined for x <= 0
		pattern: Pow(Const(0), ConstrCacheVar("x", negOrZeroConstant)),
		transform: func(expr Expr) Expr {
			return Undefined()
		},
	},
	{ // 1^x = 1
		pattern: Pow(Const(1), CacheVar("x")),
		transform: func(expr Expr) Expr {
			return Const(1)
		},
	},
	{ // x^0 = 1
		pattern: Pow(CacheVar("x"), Const(0)),
		transform: func(expr Expr) Expr {
			return Const(1)
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
		pattern: Pow(Pow(CacheVar("x"), CacheVar("y")), CacheVar("z")),
		transform: func(expr Expr) Expr {	
			x := Operand( Operand(expr, 1), 1)
			y := Operand( Operand(expr, 1), 2)
			z := Operand(expr, 2)
			return Pow(x, Mul(y, z))
		},
	},
} 

