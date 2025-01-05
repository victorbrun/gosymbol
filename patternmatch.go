package gosymbol

import (
	"fmt"
	"reflect"
)

/*
Recursively checks if expr matches pattern. varCache is an empty
map internally used to keep track of what the variables in pattern
corresponds to in expr. The function expects that no variable has the
same name as a constrained variable.
*/
func patternMatchOld(pattern, expr Expr, varCache map[VarName]Expr) bool {
	switch v := pattern.(type) {
	case undefined:
		_, ok := expr.(undefined)
		return ok
	case constant:
		exprTyped, ok := expr.(constant)
		return ok && v.Value == exprTyped.Value
	case variable:
		e, cacheOk := varCache[v.Name]
		eTyped, varOk := e.(variable)
		if cacheOk && Equal(e, expr) {
			// If the expression e, cached as the
			// pattern's variable v, equals expr
			// (the expression we are trying to match)
			// we return true.
			//
			// Example
			// -------
			// pattern := x
			// varCache[x] := y^2
			// x == y^2 --> varCache[x] == y^2 --> y^2 == y^2 --> true
			return true
		} else if cacheOk && varOk && Equal(v, eTyped) {
			// If the expression e, cached as the
			// pattern's variable v, is itself a variable
			// equalling the pattern's variable we return false
			return false
		} else if cacheOk && varOk {
			return patternMatchOld(e, expr, varCache)
		} else if cacheOk {
			return patternMatchOld(e, expr, varCache)
		} else {
			varCache[v.Name] = expr
			return true
		}
	case constrainedVariable:
		// Does just as above but before assigning an expression
		// to a variable the constraint function is checked as well
		if e, ok := varCache[v.Name]; ok {
			return patternMatchOld(e, expr, varCache)
		} else if v.Constraint(expr) {
			varCache[v.Name] = expr
			return true
		} else {
			return false
		}
	case add:
		_, ok := expr.(add)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case mul:
		_, ok := expr.(mul)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case pow:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case exp:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case log:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	default:
		errMsg := fmt.Errorf("ERROR: expression of type: %v have no matchPattern case implemented", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// Checks if the operands of pattern and expr match.
// This function does not check if the main operator
// of pattern and expr match.
func patternMatchOperands(pattern, expr Expr, varCache map[VarName]Expr) bool {
	if NumberOfOperands(pattern) != NumberOfOperands(expr) {
		return false
	}

	// Recursively checks if each operand matches
	for ix := 1; ix <= NumberOfOperands(pattern); ix++ {
		patternOp := Operand(pattern, ix)
		exprOp := Operand(expr, ix)
		if !patternMatchOld(patternOp, exprOp, varCache) {
			return false
		}
	}
	return true
}

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

		// If bound expression exist, we recursively
		// match expr against it
		return patternMatch(expr, boundExpr, bindings)

	case add:
		if a, ok := expr.(add); ok {
			// Checks that expr and pattern contains the
			// same number of operands
			nOpPattern := NumberOfOperands(p)
			nOpExpr := NumberOfOperands(a)
			if nOpPattern != nOpExpr {
				return false
			}

			allOperandsMatch := true
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				patternOp := Operand(p, ix)
				exprOp := Operand(a, ix)
				allOperandsMatch = allOperandsMatch && patternMatch(exprOp, patternOp, bindings)
			}
			return allOperandsMatch
		}
		return false

	case mul:
		if m, ok := expr.(mul); ok {
			// Checks that expr and pattern contains the
			// same number of operands
			nOpPattern := NumberOfOperands(p)
			nOpExpr := NumberOfOperands(m)
			if nOpPattern != nOpExpr {
				return false
			}

			allOperandsMatch := true
			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				patternOp := Operand(p, ix)
				exprOp := Operand(m, ix)
				allOperandsMatch = allOperandsMatch && patternMatch(exprOp, patternOp, bindings)
			}
			return allOperandsMatch
		}
		return false

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
