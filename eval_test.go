package gosymbol

import (
	"fmt"
	"strconv"
	"testing"
)

func TestExprEval(t *testing.T) {
	type inputArgs struct {
		expr Expr
		args Arguments
	}

	tests := []struct {
		input          inputArgs
		expectedOutput Expr
	}{
		{
			input: inputArgs{
				expr: Int(0),
				args: Arguments{Var("X"): Int(10)},
			},
			expectedOutput: Int(0),
		},
		{
			input: inputArgs{
				expr: Var("X"),
				args: Arguments{Var("X"): Int(1)},
			},
			expectedOutput: Int(1),
		},
		{
			input: inputArgs{
				expr: Neg(Int(10)),
				args: Arguments{Var("X"): Int(0)},
			},
			expectedOutput: Int(-10),
		},
		{
			input: inputArgs{
				expr: Add(Var("X"), Var("Y")),
				args: Arguments{Var("X"): Int(1), Var("Y"): Int(2)},
			},
			expectedOutput: Int(3),
		},
		{
			input: inputArgs{
				expr: Sub(Var("X"), Var("Y")),
				args: Arguments{Var("X"): Int(1), Var("Y"): Int(2)},
			},
			expectedOutput: Int(-1),
		},
		{
			input: inputArgs{
				expr: Mul(Var("X"), Var("Y")),
				args: Arguments{Var("X"): Int(1), Var("Y"): Int(2)},
			},
			expectedOutput: Int(2),
		},
		{
			input: inputArgs{
				expr: Div(Var("X"), Var("Y")),
				args: Arguments{Var("X"): Int(1), Var("Y"): Int(2)},
			},
			expectedOutput: Frac(Int(1), Int(2)),
		},
		{
			input: inputArgs{
				expr: Exp(Var("X")),
				args: Arguments{Var("X"): Int(0)},
			},
			expectedOutput: Int(1),
		},
		{
			input: inputArgs{
				expr: Log(Var("X")),
				args: Arguments{Var("X"): Int(1)},
			},
			expectedOutput: Int(0),
		},
		{
			input: inputArgs{
				expr: Pow(Var("X"), Int(-1)),
				args: Arguments{Var("X"): Int(10)},
			},
			expectedOutput: Frac(Int(1), Int(10)),
		},
	}

	for ix, test := range tests {
		println(ix + 1)
		result := test.input.expr.Eval()(test.input.args)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)
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
			input:          Int(-1),
			expectedOutput: "-1",
		},
		{
			input:          Int(10),
			expectedOutput: "10",
		},
		{
			input:          Add(Var("X"), Var("Y")),
			expectedOutput: "( X + Y )",
		},
		{
			input:          Sub(Var("X"), Var("Y")),
			expectedOutput: "( X + ( -1 * Y ) )",
		},
		{
			input:          Mul(Var("X"), Var("Y")),
			expectedOutput: "( X * Y )",
		},
		{
			input:          Div(Var("X"), Var("Y")),
			expectedOutput: "( X * ( Y^-1 ) )",
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
			input:          Pow(Var("X"), Int(9)),
			expectedOutput: "( X^9 )",
		},
	}

	for ix, test := range tests {
		result := fmt.Sprint(test.input)
		correctnesCheck(t, strconv.Itoa(ix+1), test.input, test.expectedOutput, result)
	}
}
