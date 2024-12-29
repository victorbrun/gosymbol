package gosymbol

import (
	"fmt"
	"math"
	"testing"
)

func TestExprEval(t *testing.T) {
	type inputArgs struct {
		expr Expr
		args Arguments
	}

	tests := []struct {
		input          inputArgs
		expectedOutput float64
	}{
		{
			input: inputArgs{
				expr: Const(0),
				args: Arguments{Var("X"): 10},
			},
			expectedOutput: 0,
		},
		{
			input: inputArgs{
				expr: Var("X"),
				args: Arguments{Var("X"): 1.0},
			},
			expectedOutput: 1.0,
		},
		{
			input: inputArgs{
				expr: Neg(Const(10)),
				args: Arguments{Var("X"): 0},
			},
			expectedOutput: -10,
		},
		{
			input: inputArgs{
				expr: Add(Var("X"), Var("Y")),
				args: Arguments{Var("X"): 1.0, Var("Y"): 2.0},
			},
			expectedOutput: 3.0,
		},
		{
			input: inputArgs{
				expr: Sub(Var("X"), Var("Y")),
				args: Arguments{Var("X"): 1.0, Var("Y"): 2.0},
			},
			expectedOutput: -1.0,
		},
		{
			input: inputArgs{
				expr: Mul(Var("X"), Var("Y")),
				args: Arguments{Var("X"): 1.0, Var("Y"): 2.0},
			},
			expectedOutput: 2.0,
		},
		{
			input: inputArgs{
				expr: Div(Var("X"), Var("Y")),
				args: Arguments{Var("X"): 1.0, Var("Y"): 2.0},
			},
			expectedOutput: 0.5,
		},
		{
			input: inputArgs{
				expr: Exp(Var("X")),
				args: Arguments{Var("X"): 0},
			},
			expectedOutput: 1,
		},
		{
			input: inputArgs{
				expr: Log(Var("X")),
				args: Arguments{Var("X"): 1},
			},
			expectedOutput: math.Log(1),
		},
		{
			input: inputArgs{
				expr: Pow(Var("X"), Const(-1)),
				args: Arguments{Var("X"): 10},
			},
			expectedOutput: 1 / 10.0,
		},
	}

	for ix, test := range tests {
		result := test.input.expr.Eval()(test.input.args)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestExprString(t *testing.T) {
	tests := []struct {
		input          Expr
		expectedOutput string
	}{
		{
			input:          Var("X"),
			expectedOutput: "X",
		},
		{
			input:          Const(-1),
			expectedOutput: "( -1 )",
		},
		{
			input:          Const(10),
			expectedOutput: "10",
		},
		{
			input:          Add(Var("X"), Var("Y")),
			expectedOutput: "( X + Y )",
		},
		{
			input:          Sub(Var("X"), Var("Y")),
			expectedOutput: "( X + ( ( -1 ) * Y ) )",
		},
		{
			input:          Mul(Var("X"), Var("Y")),
			expectedOutput: "( X * Y )",
		},
		{
			input:          Div(Var("X"), Var("Y")),
			expectedOutput: "( X * ( Y^( -1 ) ) )",
		},
		{
			input:          Exp(Var("X")),
			expectedOutput: "exp( X )",
		},
		{
			input:          Log(Var("X")),
			expectedOutput: "log( X )",
		},
		{
			input:          Pow(Var("X"), Const(9)),
			expectedOutput: "( X^9 )",
		},
	}

	for ix, test := range tests {
		result := fmt.Sprint(test.input)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}
