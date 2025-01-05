package gosymbol

func (expr undefined) Simplify() Expr {
	return simplify(expr)
}

func (expr constant) Simplify() Expr {
	return simplify(expr)
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
	case constant:
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
		expr, appliedRuleIdx = rulesApplicator(expr, powerSimplificationRules)
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
