package gosymbol

import (
	"fmt"
)

/*
Recursively checks if expr matches pattern. bindings is an empty
map internally used to keep track of what the variables in pattern
corresponds to in expr.
*/
func patternMatch(expr, pattern Expr, bindings Binding) bool {
	switch p := pattern.(type) {
	case undefined:
		if _, ok := expr.(undefined); ok {
			return true
		}
		return false

	case constant:
		if c, ok := expr.(constant); ok {
			return c.Value == p.Value
		}
		return false

	case variable:
		// If the varaible in the pattern is not
		// a pattern variable, it means that we are
		// matching expr against a binding. This is a
		// base case and we check for equality
		if !p.isPattern {
			return Equal(p, expr)
		}

		// Extracts binding to pattern variable if it exists
		// and gets if expression is a variable
		boundExpr, bindingExists := bindings[p.Name]

		// If no expression is bound to this variable
		// we bound the current expression to it and return
		// true
		if !bindingExists {
			bindings[p.Name] = expr
			return true
		}

		// If bound expression exist, we recursively
		// match expr against it
		return patternMatch(expr, boundExpr, bindings)

	case constrainedVariable:
		// If the varaible in the pattern is not
		// a pattern variable, it means that we are
		// matching expr against a binding. This is a
		// base case and we check for equality
		if !p.isPattern {
			return Equal(p, expr)
		}

		// Extracts binding to pattern variable if it exists
		// and gets if expression is a variable
		boundExpr, bindingExists := bindings[p.Name]

		// If no expression is bound to this variable
		// we bound the current expression to it and return
		// true
		if !bindingExists && p.Constraint(expr) {
			bindings[p.Name] = expr
			return true
		} else if !p.Constraint(expr) {
			return false
		}

		// If bound expression exists, we recursively
		// match expr against it
		return patternMatch(expr, boundExpr, bindings)

	case add:
		if _, ok := expr.(add); !ok {
			return false
		}

		// Checks that expr and pattern contains the
		// same number of operands
		nOpPattern := NumberOfOperands(pattern)
		nOpExpr := NumberOfOperands(expr)
		if nOpPattern != nOpExpr {
			return false
		}

		// If any of the operands between expr nd pattern
		// does not match, we return false directly
		for ix := 1; ix <= nOpExpr; ix++ {
			if !patternMatch(Operand(expr, ix), Operand(pattern, ix), bindings) {
				return false
			}
		}
		return true

	case mul:
		if _, ok := expr.(mul); !ok {
			return false
		}

		// Checks that expr and pattern contains the
		// same number of operands
		nOpPattern := NumberOfOperands(pattern)
		nOpExpr := NumberOfOperands(expr)
		if nOpPattern != nOpExpr {
			return false
		}

		// If any of the operands between expr nd pattern
		// does not match, we return false directly
		for ix := 1; ix <= nOpExpr; ix++ {
			if !patternMatch(Operand(expr, ix), Operand(pattern, ix), bindings) {
				return false
			}
		}
		return true

	case pow:
		if pw, ok := expr.(pow); ok {
			return patternMatch(pw.Base, p.Base, bindings) && patternMatch(pw.Exponent, p.Exponent, bindings)
		}
		return false

	case exp:
		if e, ok := expr.(exp); ok {
			return patternMatch(e.Arg, p.Arg, bindings)
		}
		return false

	case log:
		if l, ok := expr.(log); ok {
			return patternMatch(l.Arg, p.Arg, bindings)
		}
		return false

	default:
		errMsg := fmt.Errorf("ERROR: expression %#v have no match pattern case implemented", p)
		panic(errMsg)
	}
}
