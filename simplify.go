package gosymbol

func (expr undefined) Simplify() Expr {
	return simplify(expr)
}

func (expr integer) Simplify() Expr {
	return expr
}

func (expr fraction) Simplify() Expr {
	return (expr.simplifyRational())
}

func (expr variable) Simplify() Expr {
	return simplify(expr)
}

func (expr constrainedVariable) Simplify() Expr {
	return simplify(expr)
}

func (expr add) Simplify() Expr {
	return simplify(expr)
}

func (expr mul) Simplify() Expr {
	return simplify(expr)
}

func (expr pow) Simplify() Expr {
	return simplify(expr)
}

func (expr exp) Simplify() Expr {
	return simplify(expr)
}

func (expr log) Simplify() Expr {
	return simplify(expr)
}

func (expr sqrt) Simplify() Expr {
	return simplify(expr)
}

func simplify(expr Expr) Expr {
	// Having this here makes it possible
	// to remove all rules in simplification_rules.go
	// that basically just checks if the expression contains
	// undefined.
	if RecContains(expr, Undefined()) {
		return Undefined()
	}

	// Only sorting the top operands is sufficient
	// to sort the whole expression since in the next
	// step we recursively simplify all the operands.
	// Note that the operator must be commutative for
	// this not to fuck shit up!
	switch expr.(type) {
	case add:
		expr = TopOperandSort(expr)
	case mul:
		expr = TopOperandSort(expr)
	}

	// Recusively simplify all operands.
	for ix := 1; ix <= NumberOfOperands(expr); ix++ {
		op := Operand(expr, ix)
		expr = replaceOperand(expr, ix, op.Simplify())
	}

	// Applies simplification rules depending on the operator type
	// This will extend as more rules gets added! The base cases
	// are fully simplified so we just return them.
	appliedRuleIdx := -1
	switch expr.(type) {
	case rational:
		// Fully simplified
	case variable:
		// Fully simplified
	case constrainedVariable:
		// Fully simplified
	case add:
		expr, appliedRuleIdx = rulesApplicator(expr, sumSimplificationRules)
	case mul:
		expr, appliedRuleIdx = rulesApplicator(expr, productSimplificationRules)
	case pow:
		expr, appliedRuleIdx = rulesApplicator(expr, powerSimplificationRules2)
	case exp:
		expr, appliedRuleIdx = rulesApplicator(expr, expSimplificationRules)
	case log:
		expr, appliedRuleIdx = rulesApplicator(expr, logSimplificationRules)
	}

	// If the expression has been altered it might be possible to apply some other rule
	// we thus recursively simplify until the expression is not altered any more.
	if appliedRuleIdx > -1 {
		return simplify(expr)
	}
	return expr
}

/*
Tries to apply the transformation rules in ruleSlice to expr. If expr matches
the pattern of the transformation rule, the transformed expression is returned
together with the index of the rule that was applied.

Note: the function returns after application of the first matching rule,
or after all rules in ruleSlice have been tried. In the latter case,
-1 is returned instead of a rule index.
*/
func rulesApplicator(expr Expr, ruleSlice []transformationRule) (Expr, int) {
	for ix, rule := range ruleSlice {
		// Tries to apply rule
		transformedExpr, applied := rule.apply(expr)

		if applied {
			return transformedExpr, ix
		}

	}

	// If function did not return above no rule was applied
	return expr, -1
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	panic("Not implemented yet")
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
	if constCount > 1 {
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
// u is a sum satisfying all of the following properties:
//
// 0. u has two or more operands u1 * u2 * ...
//
// 1. u has all admissible terms
//
// 2. At most one operand ui is a constant (integer or fraction)
//
// 3. If i != j, then AsaeTerm(ui) != AsaeTerm(uj)
//
// 4. If i < j, then compare(ui, uj) = true
func addIsASAE(u add) bool {
	// 0.
	if len(u.Operands) < 2 {
		return false
	}

	// 1.
	if !hasAllAdmissibleTerms(u) {
		return false
	}

	// 2.
	constCount := 0
	for _, term := range u.Operands {
		switch term.(type) {
		case integer:
			constCount++
		case fraction:
			constCount++
		}
	}
	if constCount > 1 {
		return false
	}

	// 3. and 4.
	for ix := 1; ix <= len(u.Operands); ix++ {
		ui := Operand(u, ix)
		for jx := ix + 1; jx <= len(u.Operands); jx++ {
			uj := Operand(u, jx)
			if Equal(AsaeTerm(ui), AsaeTerm(uj)) {
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

// Returns if expr has all admissible terms.
//
// A term is said to be admissible if it
// is an ASAE which can be either an integer
// (!= 0), fraction, symbol (except undefined),
// product, power, function.
//
// Note: the term of a sum cannot be a sum
// for it to be admissible.
func hasAllAdmissibleTerms(expr add) bool {
	// Iterating over each factor and checking if
	// it is a admissible factor. Returning false
	// at first non-admissible factor
	for _, factor := range expr.Operands {
		switch factor.(type) {
		case undefined:
			return false
		case add:
			return false
		case integer:
			if Equal(factor, Int(0)) {
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

// ASAE-6
// u is a power v^w satisfying all of the following properties:
//
// 1. the expressions v and w are ASAE
//
// 2. The exponent w is not 0 or 1
//
// 3. If w is an integer, then the base v is an ASAE
// which is a symbol (except undefined), sum, or function
//
// 4. If w is not an integer, then the base v is any ASAE
// except 0 or 1
func powIsASAE(u pow) bool {
	v := Operand(u, 1)
	w := Operand(u, 2)

	// 1.
	if !IsASAE(v) || !IsASAE(w) {
		return false
	}

	// 2.
	if Equal(w, Int(0)) || Equal(w, Int(1)) {
		return false
	}

	switch w.(type) {
	// 3.
	case integer:
		switch v.(type) {
		case undefined:
			return false
		case integer:
			return false
		case fraction:
			return false
		case mul:
			return false
		case pow:
			return false
		}
	// 4.
	default:
		if Equal(v, Int(0)) || Equal(v, Int(1)) {
			return false
		}
	}

	// If we have not returned before arriving here,
	// every property is satisfied
	return true
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
