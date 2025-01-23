package gosymbol

import (
	"fmt"
	"reflect"
	"testing"
)

func TestD(t *testing.T) {
	type inputArgs struct {
		expr    Expr
		diffVar variable
	}

	tests := []struct {
		name           string
		input          inputArgs
		expectedOutput Expr
	}{
		{ // Test 1
			name: "Diff of constant",
			input: inputArgs{
				expr:    Const(10),
				diffVar: Var("X"),
			},
			expectedOutput: Const(0),
		},
		{ // Test 2
			name: "Diff of variable",
			input: inputArgs{
				expr:    Var("X"),
				diffVar: Var("X"),
			},
			expectedOutput: Const(1),
		},
		{ // Test 3
			name: "Diff of variable not in expression",
			input: inputArgs{
				expr:    Var("X"),
				diffVar: Var("Y"),
			},
			expectedOutput: Const(0),
		},
		{ // Test 4
			name: "Diff of exponential function",
			input: inputArgs{
				expr:    Exp(Var("X")),
				diffVar: Var("X"),
			},
			expectedOutput: Exp(Var("X")),
		},
		{ // Test 5
			name: "Diff of log functions",
			input: inputArgs{
				expr:    Log(Var("X")),
				diffVar: Var("X"),
			},
			expectedOutput: Div(Const(1), Var("X")),
		},
		{ // Test 6
			name: "Diff of power function",
			input: inputArgs{
				expr:    Pow(Var("X"), Const(2)),
				diffVar: Var("X"),
			},
			expectedOutput: Mul(Const(2), Var("X")),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.expr.D(test.input.diffVar)

			if !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Following test failed: %s\nInput expr: %v\nInput diff var: %v\nExpected: %v\nGot: %v", test.name, test.input.expr, test.input.diffVar, test.expectedOutput, result)
			}
		})
	}
}
