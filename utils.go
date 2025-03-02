package gosymbol

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func correctnesCheck(t *testing.T, testName string, testInput, testExpectedOutput, result any) {
	if !reflect.DeepEqual(result, testExpectedOutput) {
		t.Errorf("Test %s failed.\nInput: %v\nExpected: %v\nGot: %v", testName, testInput, testExpectedOutput, result)
	}
}

// Substitutes u for t in expr.
func Substitute(expr, u, t Expr) Expr {
	if Equal(u, t) {
		return u
	} else if Equal(expr, u) {
		return t
	} else if RecContains(expr, u) {
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			processedOp := Substitute(Operand(expr, ix), u, t)
			expr = replaceOperand(expr, ix, processedOp)
		}
		// Subsituting here again in order to continue to Substitute
		// if the result from above substitution is something that can be
		// Substituted again. Se test 7 and test 12 in gosymbol_test.go.
		return Substitute(expr, u, t)
	} else {
		return expr
	}
}

/*
Replaces operand number n in t with u and returns the resulting
expression. The function panics if n is larger than
the NumberOfOperands(t).
*/
func replaceOperand(t Expr, n int, u Expr) Expr {
	nop := NumberOfOperands(t)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	} else if n <= 0 {
		errMsg := fmt.Sprintf("ERROR: there exists no non-positive indexed operands, you are trying to replace operand: %v", nop)
		panic(errMsg)
	}

	// Since the first cases has no operands
	// we consider replacing one as just returning
	// the original expression
	switch v := t.(type) {
	case undefined:
		return v
	case rational:
		return v
	case variable:
		return v
	case constrainedVariable:
		return v
	case add:
		v.Operands[n-1] = u
		return v
	case mul:
		v.Operands[n-1] = u
		return v
	case pow:
		if n == 1 {
			v.Base = u
		} else {
			v.Exponent = u
		}
		return v
	case exp:
		v.Arg = u
		return v
	case log:
		v.Arg = u
		return v
	case sqrt:
		v.Arg = u
		return v
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

/*
Swaps operand number n1 with operand number n2 in expr.
*/
func swapOperands(expr Expr, n1, n2 int) Expr {
	op1 := Operand(expr, n1)
	op2 := Operand(expr, n2)
	expr = replaceOperand(expr, n1, op2)
	expr = replaceOperand(expr, n2, op1)
	return expr
}

/*
Recursively checks exact syntactical equality between t and u,
i.e. it does not simplify any expression nor does it
take any properties, e.g. commutativity, into account.
*/
func Equal(t, u Expr) bool {
	switch v := t.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case rational:
		uTyped, ok := u.(rational)
		return ok && v == uTyped
	case variable:
		uTyped, ok := u.(variable)
		return ok && v.Name == uTyped.Name
	case constrainedVariable:
		// TODO: how do we check equality of constrain??
		return false
	case add:
		_, ok := u.(add)
		if !ok {
			return false
		} else if NumberOfOperands(v) != NumberOfOperands(u) {
			return false
		}

		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			exprOp := Operand(v, ix)
			uOp := Operand(u, ix)
			if !Equal(exprOp, uOp) {
				return false
			}
		}
		return true
	case mul:
		_, ok := u.(mul)
		if !ok {
			return false
		} else if NumberOfOperands(v) != NumberOfOperands(u) {
			return false
		}

		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			exprOp := Operand(v, ix)
			uOp := Operand(u, ix)
			if !Equal(exprOp, uOp) {
				return false
			}
		}
		return true
	case pow:
		_, ok := u.(pow)
		return ok && Equal(Operand(v, 1), Operand(u, 1)) && Equal(Operand(v, 2), Operand(u, 2))
	case exp:
		_, ok := u.(exp)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	case log:
		_, ok := u.(log)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	case sqrt:
		_, ok := u.(sqrt)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// Returns true if t and u are equal up to type
// for every element in resp. syntax tree, otherwise
// false. This means that two rationals of different
// value, or variables with different names, would
// return true.
func TypeEqual(t, u Expr) bool {
	if !isSameType(t, u) {
		return false
	} else if NumberOfOperands(t) != NumberOfOperands(u) {
		return false
	}

	// Base cases are the leaf node types
	switch t.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case rational:
		_, ok := u.(rational)
		return ok
	case variable:
		_, ok := u.(variable)
		return ok
	}

	// Recusively checks if every operand is type equal
	// as well. Breaks and returns false if any of the
	// operands are not equal.
	// TODO: this does not take associativty of e.g. add and mul into account.
	for ix := 1; ix <= NumberOfOperands(t); ix++ {
		tOperand := Operand(t, ix)
		uOperand := Operand(u, ix)
		if !TypeEqual(tOperand, uOperand) {
			return false
		}
	}

	return true
}

/*
Recursively checks exact equality between expr and u,
i.e. it does not simplify any expression nor does it
take any properties, e.g. commutativity, into account.
*/
func RecContains(expr, u Expr) bool {
	switch v := expr.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case rational:
		uTyped, ok := u.(rational)
		return ok && v == uTyped
	case variable:
		uTyped, ok := u.(variable)
		return ok && v.Name == uTyped.Name
	case constrainedVariable:
		// TODO: how do we check equality of constrain??
		return false
	default:
		if Equal(v, u) {
			return true
		}
		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			vOp := Operand(v, ix)
			if RecContains(vOp, u) {
				return true
			}
		}
		return false
	}
}

// Returns the different variable names
// present in the given expression.
func VariableNames(expr Expr) []VarName {
	exprSlice := Variables(expr)
	if len(exprSlice) == 0 {
		return []VarName{}
	}

	// Extracting variable names from the list of expressions
	stringSlice := make([]string, len(exprSlice))
	for ix, e := range exprSlice {
		// Typing expression
		switch v := e.(type) {
		case variable:
			stringSlice[ix] = string(v.Name)
		case constrainedVariable:
			stringSlice[ix] = string(v.Name)
		default:
			err := fmt.Errorf("somthing went wrong: %#v is expected to be of type variable or constrainedVariable", v)
			panic(err)
		}
	}

	// Checks for and deletes duplicates of variable names
	// TODO: speed up by checking for duplicates during sorting
	sort.Strings(stringSlice)
	var variableNamesSlice []VarName
	variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[0]))
	jx := 0
	for ix := 1; ix < len(stringSlice); ix++ {
		if variableNamesSlice[jx] != VarName(stringSlice[ix]) {
			variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[ix]))
			jx++
		}
	}

	return variableNamesSlice
}

/*
Returns the list of variables (and contrained variables) in expr
as a list without duplicates
*/
func Variables(expr Expr) []Expr {
	var exprSlice []Expr
	variables(expr, &exprSlice)

	return exprSlice
}

/*
Recursively travesrses the whole AST and appends
the variables to targetSlice.
*/
func variables(expr Expr, targetSlice *[]Expr) {
	switch v := expr.(type) {
	case undefined:
		return
	case rational:
		return
	case variable:
		*targetSlice = append(*targetSlice, v)
	case constrainedVariable:
		*targetSlice = append(*targetSlice, v)
	default:
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			variables(op, targetSlice)
		}
	}
}

// Returns the number of operands for top level operation.
func NumberOfOperands(expr Expr) int {
	switch v := expr.(type) {
	case undefined:
		return 0
	case rational:
		return 0
	case variable:
		return 0
	case constrainedVariable:
		return 0
	case add:
		return len(v.Operands)
	case mul:
		return len(v.Operands)
	case pow:
		return 2
	case exp:
		return 1
	case log:
		return 1
	case sqrt:
		return 1
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// Returns the n:th (starting at 1) operand (left to right) of expr.
// If expr has no operands it returns nil.
// If n is larger than NumberOfOperands(expr)-1 it will panic.
func Operand(expr Expr, n int) Expr {
	nop := NumberOfOperands(expr)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	}

	switch v := expr.(type) {
	case undefined:
		return nil
	case rational:
		return nil
	case variable:
		return nil
	case constrainedVariable:
		return nil
	case add:
		return v.Operands[n-1]
	case mul:
		return v.Operands[n-1]
	case pow:
		if n == 1 {
			return v.Base
		} else {
			return v.Exponent
		}
	case exp:
		return v.Arg
	case log:
		return v.Arg
	case sqrt:
		return v.Arg
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// TODO: see Computer Algebra and Symbolic Computation page 10 to understand this shit
func Map(F Expr, u ...Expr) Expr { panic("Not implemented yet") }

func isSameType(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// Returns the deepest depth of expr
func Depth(expr Expr) int {
	switch expr.(type) {
	case undefined:
		return 0
	case rational:
		return 0
	case variable:
		return 0
	case constrainedVariable:
		return 0
	default:
		maxDepth := 0
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			opDepth := Depth(op)
			if opDepth > maxDepth {
				maxDepth = opDepth
			}
		}
		return maxDepth + 1
	}
}

/*
Returns a new expression where the
terms in in s1 has been prepended
to the terms in s2, i.e. the order
of terms is not changed.
*/
func mergeSums(s1, s2 add) add {
	panic("not implemented yet")
}

/*
Returns a new expression where the
factors in in p1 has been prepended
to the factors in p2, i.e. the order
of factors is not changed.
*/
func mergeProducts(p1, p2 mul) mul {
	panic("not implemented yet")
}

/*
Returns all elements at indexes `vals` which has true
in `selectors`.
*/
func filterByBoolSlice[T any](vals []T, selectors []bool) []T {
	if len(vals) != len(selectors) {
		panic("vals and selectors must be slices of the same length")
	}

	var result []T
	for ix, selected := range selectors {
		if selected {
			result = append(result, vals[ix])
		}
	}
	return result
}

/*
Negates slice of bool
*/
func negateSlice(boolSlice []bool) []bool {
	for ix, val := range boolSlice {
		boolSlice[ix] = !val
	}
	return boolSlice
}
