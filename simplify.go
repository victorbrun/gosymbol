package gosymbol

func SimplifyV2(expr Expr) Expr {
	// having this here means we do not need rules
	// to check for Undefined.
	if Contains(expr, Undefined()) {
		return Undefined()
	}

	// Sorting the top operand arguments.
	// Beacuse only addition and multiplication have multiple
	// arguments, this step only needs to be applied to 
	// those cases
	expr = MatchTransform[Expr](
		expr, 
		func(expr undefined) Expr { return expr },
		func(expr constant) Expr { return expr },
		func(expr variable) Expr { return expr },
		func(expr add) Expr { return topOperandSort(expr) },
		func(expr mul) Expr { return topOperandSort(expr) },
		func(expr pow) Expr { return expr },
		func(expr exp) Expr { return expr },
		func(expr log) Expr { return expr },
		func(expr sqrt) Expr { return expr },
	)

	// Applies simplification rules
	applicationResult := MatchTransform[transformationRuleApplicationResult](
		expr,
		func(expr undefined) transformationRuleApplicationResult {
			// Fully simplified
			return transformationRuleApplicationResult{
				expr: expr,
				exprAltered: false,
			}
		},
		func(expr constant) transformationRuleApplicationResult {
			// Fully simplified
			return transformationRuleApplicationResult{
				expr: expr,
				exprAltered: false,
			}
		},
		func(expr variable) transformationRuleApplicationResult {
			// Fully simplified
			return transformationRuleApplicationResult{
				expr: expr,
				exprAltered: false,
			}
		},
		func(expr add) transformationRuleApplicationResult {
			return applySimplificationRules(expr, sumSimplificationRules)
		},
		func(expr mul) transformationRuleApplicationResult {
			return applySimplificationRules(expr, productSimplificationRules)
		},
		func(expr pow) transformationRuleApplicationResult {
			return applySimplificationRules(expr, powerSimplificationRules)
		},
		func(expr exp) transformationRuleApplicationResult {
			return applySimplificationRules(expr, powerSimplificationRules)
		},
		func(expr log) transformationRuleApplicationResult {
			return applySimplificationRules(expr, logSimplificationRules)
		},
		func(expr sqrt) transformationRuleApplicationResult {
			return applySimplificationRules(expr, sqrtSimplifiactionRules)
		},
	)

	// When the epxression has been modified, we need to try and simplify it
	// again as there might be new possible simplification paths
	if applicationResult.exprAltered {
		return SimplifyV2(expr)
	}
	return expr

}

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
		expr = topOperandSort(expr)
	case mul:
		expr = topOperandSort(expr)
	}

	// Recursively simplify all operands.
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
	case add:
		expr, expressionAltered = rulesApplication(expr, sumSimplificationRules)
	case mul:
		expr, expressionAltered = rulesApplication(expr, productSimplificationRules)	
	case pow:
		expr, expressionAltered = rulesApplication(expr, powerSimplificationRules)
	case log:
		expr, expressionAltered = rulesApplication(expr, logSimplificationRules)
	case exp:
		expr, expressionAltered = rulesApplication(expr, expSimplificationRules)
	case sqrt:
		expr, expressionAltered = rulesApplication(expr, sqrtSimplifiactionRules)
	}

	// If the expression has been altered it might be possible to apply some other rule 
	// we thus recursively sort until the expression is not altered any more.
	if expressionAltered {
		// TODO: this will get stuck since we flatten the 
		// expr, then when simplifying it we will turn it into a 
		// binary tree and then it has been altered so we will get in here again 
		// and flatten it and then start all over. We will have inifinite loop :(
		expr = flattenTopLevel(expr)
		return Simplify(expr)
	}
	return expr
}

func applySimplificationRules(
	expr Expr, 
	rulesSlice []transformationRule, 
) transformationRuleApplicationResult {
	atLestOneRuleApplied := false 
	for _, rule := range rulesSlice {
		var applied bool 
		expr, applied = rule.apply(expr)
		atLestOneRuleApplied = atLestOneRuleApplied || applied
	}

	return transformationRuleApplicationResult{
		expr: expr,
		exprAltered: atLestOneRuleApplied,
	}
}
