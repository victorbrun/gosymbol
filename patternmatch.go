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
func patternMatch(pattern, expr Expr, varCache map[VarName]Expr) bool {
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
			return true
		} else if cacheOk && varOk && Equal(v, eTyped) {
			return false
		} else if cacheOk && varOk {
			return patternMatch(e, expr, varCache)
		} else if cacheOk {
			return patternMatch(e, expr, varCache)
		} else {
			varCache[v.Name] = expr
			return true
		}
	case constrainedVariable:
		// Does just as above but before assigning an expression
		// to a variable the constraint function is checked as well
		if e, ok := varCache[v.Name]; ok {
			return patternMatch(e, expr, varCache)
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
		if !patternMatch(patternOp, exprOp, varCache) {
			return false
		}
	}
	return true
}
