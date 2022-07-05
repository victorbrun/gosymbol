package gosymbol_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/victorbrun/gosymbol"
)

func TestExprPrint(t *testing.T) {
	tests := []struct{
		input gosymbol.Expr
		expectedOutput string
	} {
		{
			input: gosymbol.Var("X"),
			expectedOutput: "X",
		},
		{
			input: gosymbol.Const(-1),
			expectedOutput: "( -1 )",
		},
		{
			input: gosymbol.Const(10),
			expectedOutput: "10",
		},
		{
			input: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X + Y )",
		},
		{
			input: gosymbol.Sub(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X + ( ( -1 ) * Y ) )",
		},
		{
			input: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X * Y )",
		},
		{
			input: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X / Y )",
		},
	}
	
	for ix, test := range tests {
		result := fmt.Sprint(test.input)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf("Failed test: %v. Expected: %v \nGot: %v", ix+1, test.expectedOutput, result)
			t.Error(errMsg)
		}
	}
}

func TestExprEval(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		args gosymbol.Arguments
	}

	tests := []struct{
		input inputArgs
		expectedOutput float64
	} {
		{
			input: inputArgs{
				expr: gosymbol.Var("X"),
				args: gosymbol.Arguments{ "X": 1.0 },
			},
			expectedOutput: 1.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 3.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Sub(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: -1.0,
	
		},
		{
			input: inputArgs{
				expr: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 2.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Y" )),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 0.5,
		},
	}

	for ix, test := range tests {
		result := test.input.expr.Eval(test.input.args)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf("Failed test: %v\nExpected: %v\nGot: %v", ix+1, test.expectedOutput, result)
			t.Error(errMsg)
		}
	}
}

func TestD(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		diffVar string
	}

	tests := []struct{
		input inputArgs
		expectedOutput gosymbol.Expr
	} {
		{ // Test 1
			input: inputArgs{
				expr: gosymbol.Const(10),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Const(0),
		},
		{ // Test 2
			input: inputArgs{
				expr: gosymbol.Var("X"),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Const(1),
		},
		{ // Test 3
			input: inputArgs{
				expr: gosymbol.Var("X"),
				diffVar: "Y",
			},
			expectedOutput: gosymbol.Const(0),
		},	
	}

	for ix, test := range tests {
		deriv := test.input.expr.D(test.input.diffVar)
		if !reflect.DeepEqual(deriv, test.expectedOutput) {
			errMsg := fmt.Sprintf("Failed test: %v\nExpected: %v\nGot: %v", ix+1, test.expectedOutput, deriv)
			t.Error(errMsg)
		}
	}
}

func TestContains(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		u gosymbol.Expr
	}

	tests := []struct {
		input inputArgs
		expectedOutput bool
	} { 
		{ // Test 1: testing for equal expressions
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("X"),
			},
			expectedOutput: true,
		},
		{ // Test 2: Testing for inequality
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("Y"),
			},
			expectedOutput: false,
		},
		{ // Test 3: Testing for part of n-ary operator
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Const(1), gosymbol.Const(2), gosymbol.Const(3)),
				u: gosymbol.Add(gosymbol.Const(1), gosymbol.Const(2)),
			},
			expectedOutput: false,
		},
		{ // Test 4: Testing for sub-tree equality
			input: inputArgs{
				expr: gosymbol.Add(
					gosymbol.Const(1), 
					gosymbol.Mul(gosymbol.Const(2), gosymbol.Var("X")), 
					gosymbol.Div(gosymbol.Const(1), gosymbol.Var("y")),
				),
				u: gosymbol.Mul(gosymbol.Const(2), gosymbol.Var("X")),
			},
			expectedOutput: true,
		},
	}

	for ix, test := range tests {
		result := gosymbol.Contains(test.input.expr, test.input.u)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf(
				"Failed test: %v: Expected: %v, Got: %v. expr = %v, u = %v",
				ix+1, test.expectedOutput, result, test.input.expr, test.input.u,
			)
			t.Error(errMsg)
		}
	}
}
