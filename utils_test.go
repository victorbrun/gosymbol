package gosymbol

import (
	"fmt"
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	type inputArgs struct {
		expr Expr
		u    Expr
	}

	tests := []struct {
		input          inputArgs
		expectedOutput bool
	}{
		{ // Test 1: testing for equal expressions
			input: inputArgs{
				expr: Var("X"),
				u:    Var("X"),
			},
			expectedOutput: true,
		},
		{ // Test 2: Testing for inequality
			input: inputArgs{
				expr: Var("X"),
				u:    Var("Y"),
			},
			expectedOutput: false,
		},
		{ // Test 3: Testing for part of n-ary operator
			input: inputArgs{
				expr: Add((Int(1)), (Int(2)), (Int(3))),
				u:    Add((Int(1)), (Int(2))),
			},
			expectedOutput: false,
		},
		{ // Test 4: Testing for sub-tree equality
			input: inputArgs{
				expr: Add(
					(Int(1)),
					Mul((Int(2)), Var("X")),
					Div((Int(1)), Var("y")),
				),
				u: Mul((Int(2)), Var("X")),
			},
			expectedOutput: true,
		},
	}

	for ix, test := range tests {
		result := Contains(test.input.expr, test.input.u)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf(
				"Failed test: %v: Expected: %v, Got: %v. expr = %v, u = %v",
				ix+1, test.expectedOutput, result, test.input.expr, test.input.u,
			)
			t.Error(errMsg)
		}
	}
}

func TestSubstitute(t *testing.T) {
	type inputArgs struct {
		expr Expr
		u    Expr
		t    Expr
	}

	tests := []struct {
		input          inputArgs
		expectedOutput Expr
	}{
		{ // Test 1: base case 1: constant
			input: inputArgs{
				expr: (Int(7)),
				u:    (Int(7)),
				t:    (Int(-7)),
			},
			expectedOutput: (Int(-7)),
		},
		{ // Test 2: base case 2: variable
			input: inputArgs{
				expr: Var("X"),
				u:    Var("X"),
				t:    (Int(0)),
			},
			expectedOutput: (Int(0)),
		},
		{ // Test 3: add operator
			input: inputArgs{
				expr: Add(Var("X"), Var("Y"), (Int(0)), Var("Y")),
				u:    Var("Y"),
				t:    Var("Z"),
			},
			expectedOutput: Add(Var("X"), Var("Z"), (Int(0)), Var("Z")),
		},
		{ // Test 4: mul operator
			input: inputArgs{
				expr: Mul(Var("X"), Var("Y"), (Int(0)), Var("Y")),
				u:    Var("Y"),
				t:    Var("Z"),
			},
			expectedOutput: Mul(Var("X"), Var("Z"), (Int(0)), Var("Z")),
		},
		{ // Test 5: Div
			input: inputArgs{
				expr: Div(Var("X"), Var("Y")),
				u:    Var("Y"),
				t:    Var("Z"),
			},
			expectedOutput: Div(Var("X"), Var("Z")),
		},
		{ // Test 6: substituting whole subtree
			input: inputArgs{
				expr: Add(Div((Int(9)), Var("X")), Var("Y"), (Int(0)), Var("Y")),
				u:    Div((Int(9)), Var("X")),
				t:    Var("Z"),
			},
			expectedOutput: Add(Var("Z"), Var("Y"), (Int(0)), Var("Y")),
		},
		{ // Test 7: nested substitution, making sure that the recursion starts bottom up
			input: inputArgs{
				expr: Div((Int(1)), Div((Int(1)), Var("X"))),
				u:    Div((Int(1)), Var("X")),
				t:    Var("X"),
			},
			expectedOutput: Var("X"),
		},
		{ // Test 8: no expression matching u
			input: inputArgs{
				expr: Var("X"),
				u:    Var("Y"),
				t:    (Int(0)),
			},
			expectedOutput: Var("X"),
		},
		{ // Test 9: exp operator
			input: inputArgs{
				expr: Exp(Var("X")),
				u:    Exp(Var("X")),
				t:    Var("Y"),
			},
			expectedOutput: Var("Y"),
		},
		{ // Test 10: log operator
			input: inputArgs{
				expr: Log(Var("X")),
				u:    Log(Var("X")),
				t:    Var("Y"),
			},
			expectedOutput: Var("Y"),
		},
		{ // Test 11: pow operator
			input: inputArgs{
				expr: Pow(Var("X"), (Int(10))),
				u:    Pow(Var("X"), (Int(10))),
				t:    Var("Y"),
			},
			expectedOutput: Var("Y"),
		},
		{ // Test 12: 4 level nested substitution
			input: inputArgs{
				expr: Div((Int(1)), Div((Int(1)), Div((Int(1)), Div((Int(1)), Var("X"))))),
				u:    Div((Int(1)), Var("X")),
				t:    Var("X"),
			},
			expectedOutput: Var("X"),
		},
	}

	for ix, test := range tests {
		result := Substitute(test.input.expr, test.input.u, test.input.t)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)
	}
}

func TestVariableNames(t *testing.T) {
	tests := []struct {
		input          Expr
		expectedOutput []VarName
	}{
		{
			input:          Undefined(),
			expectedOutput: []VarName{},
		},
		{
			input:          (Int(0)),
			expectedOutput: []VarName{},
		},
		{
			input:          Var("X"),
			expectedOutput: []VarName{"X"},
		},
		{
			input:          Mul(Var("X"), Var("X"), Var("Y")),
			expectedOutput: []VarName{"X", "Y"},
		},
		{
			input:          Add(Var("X"), Var("X"), Var("Y")),
			expectedOutput: []VarName{"X", "Y"},
		},
		{
			input:          Pow(Mul(Var("X"), Var("X"), Var("Y")), (Int(10))),
			expectedOutput: []VarName{"X", "Y"},
		},
	}

	for ix, test := range tests {
		result := VariableNames(test.input)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)
	}
}

func TestDepth(t *testing.T) {
	tests := []struct {
		input          Expr
		expectedOutput int
	}{
		{
			input:          Undefined(),
			expectedOutput: 0,
		},
		{
			input:          (Int(1)),
			expectedOutput: 0,
		},
		{
			input:          Var("x"),
			expectedOutput: 0,
		},
		{
			input:          ConstrVar("x", func(expr Expr) bool { return true }),
			expectedOutput: 0,
		},
		{
			input:          Add((Int(0)), Var("x"), (Int(0)), Var("x"), (Int(0)), Var("x")),
			expectedOutput: 1,
		},
		{
			input:          Add(Mul(Var("x"), Pow((Int(10)), Exp(Var("x")))), Var("x"), (Int(0)), Var("x"), (Int(0)), Var("x")),
			expectedOutput: 4,
		},
	}

	for ix, test := range tests {
		result := Depth(test.input)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)
	}
}
