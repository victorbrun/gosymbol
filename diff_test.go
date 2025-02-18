package gosymbol

import (
	"strconv"
	"testing"
)

func TestD(t *testing.T) {
	type inputArgs struct {
		expr    Expr
		diffVar variable
	}

	tests := []struct {
		input          inputArgs
		expectedOutput Expr
	}{
		{ // Test 1
			input: inputArgs{
				expr:    (Int(10)),
				diffVar: Var("X"),
			},
			expectedOutput: (Int(0)),
		},
		{ // Test 2
			input: inputArgs{
				expr:    Var("X"),
				diffVar: Var("X"),
			},
			expectedOutput: (Int(1)),
		},
		{ // Test 3
			input: inputArgs{
				expr:    Var("X"),
				diffVar: Var("Y"),
			},
			expectedOutput: (Int(0)),
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
			expectedOutput: Div((1), Var("X")),
		},
		{ // Test 6
			input: inputArgs{
				expr: Pow(Var("X"), (2)),
				diffVar: "X",
			},
			expectedOutput: Mul((2), Pow(Var("X"), (1))),
		},*/
	}

	for ix, test := range tests {
		result := test.input.expr.D(test.input.diffVar)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)

	}
}
