package gosymbol

import (
	"fmt"
	"reflect"
	"testing"
)

func correctnesCheck(t *testing.T, result, expectedOutput any, testNumber int) {
		if !reflect.DeepEqual(result, expectedOutput) {
			errMsg := fmt.Sprintf("Failed test: %v. Expected: %v, Got: %v", testNumber, expectedOutput, result)
			t.Error(errMsg)
		}
}

func TestComapre(t *testing.T) {
	type inputArgs struct {
		expr1 Expr
		expr2 Expr
	}

	tests := []struct{
		input inputArgs
		expectedOutput bool
	} { // Below I denote the generalised order relation with <
		{ // Test 1: a + b < a + c
			input: inputArgs{
				expr1: Add(Var("a"), Var("b")),
				expr2: Add(Var("a"), Var("c")),
			},
			expectedOutput: true,
		},
		{ // Test 2: a + b |> a + c
			input: inputArgs{
				expr2: Add(Var("a"), Var("b")),
				expr1: Add(Var("a"), Var("c")),
			},
			expectedOutput: false,
		},
		{ // Test 3: a + c + d < b + c + d
			input: inputArgs{
				expr1: Add(Var("a"), Var("c"), Var("d")),
				expr2: Add(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: true,
		},
		{ // Test 4: a + c + d |> b + c + d
			input: inputArgs{
				expr2: Add(Var("a"), Var("c"), Var("d")),
				expr1: Add(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: false,
		},
		{ // Test 5: c + d < b + c + d
			input: inputArgs{
				expr1: Add(Var("c"), Var("d")),
				expr2: Add(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: true,
		},
		{ // Test 6: c + d |> b + c + d
			input: inputArgs{
				expr2: Add(Var("c"), Var("d")),
				expr1: Add(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: false,
		},
		{ // Test 7: a * b < a * c
			input: inputArgs{
				expr1: Mul(Var("a"), Var("b")),
				expr2: Mul(Var("a"), Var("c")),
			},
			expectedOutput: true,
		},
		{ // Test 8: a * b |> a * c
			input: inputArgs{
				expr2: Mul(Var("a"), Var("b")),
				expr1: Mul(Var("a"), Var("c")),
			},
			expectedOutput: false,
		},
		{ // Test 9: a * c * d < b * c * d
			input: inputArgs{
				expr1: Mul(Var("a"), Var("c"), Var("d")),
				expr2: Mul(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: true,
		},
		{ // Test 10: a * c * d |> b * c * d
			input: inputArgs{
				expr2: Mul(Var("a"), Var("c"), Var("d")),
				expr1: Mul(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: false,
		},
		{ // Test 11: c * d < b * c * d
			input: inputArgs{
				expr1: Mul(Var("c"), Var("d")),
				expr2: Mul(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: true,
		},
		{ // Test 12: c * d |> b * c * d
			input: inputArgs{
				expr2: Mul(Var("c"), Var("d")),
				expr1: Mul(Var("b"), Var("c"), Var("d")),
			},
			expectedOutput: false,
		},
		{ // Test 13: (1 + x)^2 < (1 + x)^3
			input: inputArgs{
				expr1: Pow(Add(Const(1), Var("x")), Const(2)),
				expr2: Pow(Add(Const(1), Var("x")), Const(3)),
			},
			expectedOutput: true,
		},
		{ // Test 14: (1 + x)^2 |> (1 + x)^3
			input: inputArgs{
				expr2: Pow(Add(Const(1), Var("x")), Const(2)),
				expr1: Pow(Add(Const(1), Var("x")), Const(3)),
			},
			expectedOutput: false,
		},
		{ // Test 15: (1 + x)^2 < (1 + y)^2
			input: inputArgs{
				expr1: Pow(Add(Const(1), Var("x")), Const(2)),
				expr2: Pow(Add(Const(1), Var("y")), Const(2)),
			},
			expectedOutput: true,
		},
		{ // Test 16: (1 + x)^2 |> (1 + y)^2
			input: inputArgs{
				expr2: Pow(Add(Const(1), Var("x")), Const(2)),
				expr1: Pow(Add(Const(1), Var("y")), Const(2)),
			},
			expectedOutput: false,
		},
		{ // Test 17: (1 + x)^3 < (1 + y)^2
			input: inputArgs{
				expr1: Pow(Add(Const(1), Var("x")), Const(3)),
				expr2: Pow(Add(Const(1), Var("y")), Const(2)),
			},
			expectedOutput: true,
		},
		{ // Test 18: (1 + x)^3 |> (1 + y)^2
			input: inputArgs{
				expr2: Pow(Add(Const(1), Var("x")), Const(3)),
				expr1: Pow(Add(Const(1), Var("y")), Const(2)),
			},
			expectedOutput: false,
		},
		{ // Test 19: a * x^2 < x^3
			input: inputArgs{
				expr1: Mul(Var("a"), Pow(Var("x"), Const(2))),
				expr2: Pow(Var("x"), Const(3)),
			},
			expectedOutput: true,
		},
		{ // Test 20: a * x^2 |> x^3
			input: inputArgs{
				expr2: Mul(Var("a"), Pow(Var("x"), Const(2))),
				expr1: Pow(Var("x"), Const(3)),
			},
			expectedOutput: false,
		},
		{ // Test 21: x < x^2
			input: inputArgs{
				expr1: Var("x"),
				expr2: Pow(Var("x"), Const(2)),
			},
			expectedOutput: true,
		},
		{ // Test 22: x |> x^2
			input: inputArgs{
				expr2: Var("x"),
				expr1: Pow(Var("x"), Const(2)),
			},
			expectedOutput: false,
		},
		{ // Test 23: x < Exp(y)
			input: inputArgs{
				expr1: Var("x"),
				expr2: Exp(Var("y")),
			},
			expectedOutput: true,
		},
		{ // Test 24: x |> Exp(y)
			input: inputArgs{
				expr2: Var("x"),
				expr1: Exp(Var("y")),
			},
			expectedOutput: false,
		},
		{ // Test 25: Exp(x) < Exp(x^2)
			input: inputArgs{
				expr1: Exp(Var("x")),
				expr2: Exp(Pow(Var("x"), Const(2))),
			},
			expectedOutput: true,
		},
		{ // Test 26: Exp(x) |> Exp(x^2)
			input: inputArgs{
				expr2: Exp(Var("x")),
				expr1: Exp(Pow(Var("x"), Const(2))),
			},
			expectedOutput: false,
		},
		{ // Test 27: x < Log(y)
			input: inputArgs{
				expr1: Var("x"),
				expr2: Log(Var("y")),
			},
			expectedOutput: true,
		},
		{ // Test 28: x |> Log(y)
			input: inputArgs{
				expr2: Var("x"),
				expr1: Log(Var("y")),
			},
			expectedOutput: false,
		},
		{ // Test 29: Log(x) < Log(x^2)
			input: inputArgs{
				expr1: Log(Var("x")),
				expr2: Log(Pow(Var("x"), Const(2))),
			},
			expectedOutput: true,
		},
		{ // Test 30: Log(x) |> Log(x^2)
			input: inputArgs{
				expr2: Log(Var("x")),
				expr1: Log(Pow(Var("x"), Const(2))),
			},
			expectedOutput: false,
		},
		{ // Test 31: x < Sqrt(y)
			input: inputArgs{
				expr1: Var("x"),
				expr2: Sqrt(Var("y")),
			},
			expectedOutput: true,
		},
		{ // Test 32: x |> Sqrt(y)
			input: inputArgs{
				expr2: Var("x"),
				expr1: Sqrt(Var("y")),
			},
			expectedOutput: false,
		},
		{ // Test 33: Sqrt(x) < Sqrt(x^2)
			input: inputArgs{
				expr1: Sqrt(Var("x")),
				expr2: Sqrt(Pow(Var("x"), Const(2))),
			},
			expectedOutput: true,
		},
		{ // Test 34: Sqrt(x) |> Sqrt(x^2)
			input: inputArgs{
				expr2: Sqrt(Var("x")),
				expr1: Sqrt(Pow(Var("x"), Const(2))),
			},
			expectedOutput: false,
		},
	}

	for ix, test := range tests {
		result := compare(test.input.expr1, test.input.expr2)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}
