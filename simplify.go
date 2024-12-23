package gosymbol

func Simplify(expr Expr) Expr {
	// Having this here makes it possible
	// to remove all rules in simplification_rules.go
	// that basically just checks if the expression contains
	// undefined.
	if Contains(expr, Undefined()) {
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
		expr = replaceOperand(expr, ix, Simplify(op))
	}

	// Application of all Simplification rules follow this same pattern.
	// Returns the simplified expression and a boolean describing whether any
	// simplification rule has actually been applied.
	rulesApplication := func(expr Expr, ruleSlice []transformationRule) (Expr, bool) {
		atLeastOneapplied := false
		for _, rule := range ruleSlice {
			var applied bool
			expr, applied = rule.apply(expr)
			//if applied { fmt.Println("Applied rule ", ix) }
			atLeastOneapplied = atLeastOneapplied || applied
		}
		return expr, atLeastOneapplied
	}

	// Applies simplification rules depending on the operator type
	// This will extend as more rules gets added! The base cases
	// are fully simplified so we just return them.
	expressionAltered := false
	switch expr.(type) {
	case constant:
		// Fully simplified
	case variable:
		// Fully simplified
	case constrainedVariable:
		// Fully simplified
	case add:
		expr, expressionAltered = rulesApplication(expr, sumSimplificationRules)
	case mul:
		expr, expressionAltered = rulesApplication(expr, productSimplificationRules)
	case pow:
		expr, expressionAltered = rulesApplication(expr, powerSimplificationRules)
	}

	// If the expression has been altered it might be possible to apply some other rule
	// we thus recursively sort until the expression is not altered any more.
	if expressionAltered {
		return Simplify(expr)
	}
	return expr
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	panic("Not implemented yet")
}
