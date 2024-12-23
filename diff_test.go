package gosymbol

import (
	"strconv"
	"testing"
)

func TestD(t *testing.T) {
	type inputArgs struct {
		expr    Expr
		diffVar VarName
	}

	tests := []struct {
		input          inputArgs
		expectedOutput Expr
	}{
		{ // Test 1
			input: inputArgs{
				expr:    Const(10),
				diffVar: "X",
			},
			expectedOutput: Const(0),
		},
		{ // Test 2
			input: inputArgs{
				expr:    Var("X"),
				diffVar: "X",
			},
			expectedOutput: Const(1),
		},
		{ // Test 3
			input: inputArgs{
				expr:    Var("X"),
				diffVar: "Y",
			},
			expectedOutput: Const(0),
		},
		/*{ // Test 4
			input: inputArgs{
				expr: Exp(Var("X")),
				diffVar: "X",
			},
			expectedOutput: Exp(Var("X")),
		},
		{ // Test 5
			input: inputArgs{
				expr: Log(Var("X")),
				diffVar: "X",
			},
			expectedOutput: Div(Const(1), Var("X")),
		},
		{ // Test 6
			input: inputArgs{
				expr: Pow(Var("X"), Const(2)),
				diffVar: "X",
			},
			expectedOutput: Mul(Const(2), Pow(Var("X"), Const(1))),
		},*/
	}

	for ix, test := range tests {
		result := test.input.expr.D(test.input.diffVar)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)

	}
}
