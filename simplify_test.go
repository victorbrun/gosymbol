package gosymbol

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSimplify(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput Expr
	}{
		{
			name:           "undefined^y = undefined",
			input:          Pow(Undefined(), Var("y")),
			expectedOutput: Undefined(),
		},
		{
			name:           "x^undefined = undefined",
			input:          Pow(Var("x"), Undefined()),
			expectedOutput: Undefined(),
		},
		{
			name:           "0^x = 0",
			input:          Pow(Const(0), Const(10)),
			expectedOutput: Const(0),
		},
		{
			name:           "0^0 = undefined",
			input:          Pow(Const(0), Const(0)),
			expectedOutput: Undefined(),
		},
		{
			name:           "1^x = undefined",
			input:          Pow(Const(1), Exp(Const(7))),
			expectedOutput: Const(1),
		},
		{
			name:           "x^0 = 1",
			input:          Pow(Var("kuk"), Const(0)),
			expectedOutput: Const(1),
		},
		{
			name:           "(v_1 * ... * v_n)^m = v_1^m * .. * v_n^m (note that the result is also sorted)",
			input:          Pow(Mul(Var("x"), Const(3), Var("y")), Var("elle")),
			expectedOutput: Mul(Pow(Const(3), Var("elle")), Pow(Var("x"), Var("elle")), Pow(Var("y"), Var("elle"))),
		},
		{
			name:           "(i^j)^k = i^(j*k)",
			input:          Pow(Pow(Var("i"), Var("j")), Exp(Mul(Const(10), Var("k")))),
			expectedOutput: Pow(Var("i"), Mul(Var("j"), Exp(Mul(Const(10), Var("k"))))),
		},
		{
			name:           "undefined * ... = undefined",
			input:          Mul(Undefined(), Var("x"), Const(10)),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * ... = 0",
			input:          Mul(Var("x"), Const(-9), Const(0)),
			expectedOutput: Const(0),
		},
		{
			name:           "undefined * 0 = undefined",
			input:          Mul(Undefined(), Const(0)),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * undefined = undefined",
			input:          Mul(Const(0), Undefined()),
			expectedOutput: Undefined(),
		},
		{
			name:           "Mult with only one operand simplifies to the operand",
			input:          Mul(Exp(Var("x"))),
			expectedOutput: Exp(Var("x")),
		},
		{
			name:           "Mult with no operands simplify to 1",
			input:          Mul(),
			expectedOutput: Const(1),
		},
		{
			name:           "1 * x = x",
			input:          Mul(Const(1), Exp(Var("x"))),
			expectedOutput: Exp(Var("x")),
		},
		{
			name:           "x * x = x^2",
			input:          Mul(Const(10), Const(10)),
			expectedOutput: Pow(Const(10), Const(2)),
		},
		{
			name:           "x * x^n = x^(n+1)",
			input:          Mul(Const(10), Pow(Const(10), Const(2))),
			expectedOutput: Pow(Const(10), Const(3)),
		},
		{
			name:           "x * (1/x) = 1",
			input:          Mul(Var("x"), Div(Const(1), Var("x"))),
			expectedOutput: Const(1),
		},
		{
			name:           "x^m * x^n = x^(m+n)",
			input:          Mul(Pow(Var("x"), Var("n")), Pow(Var("x"), Var("m"))),
			expectedOutput: Pow(Var("x"), Add(Var("m"), Var("n"))),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			//fmt.Println("Simplifying: ", test.input)
			result := Simplify(test.input)

			if !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Following test failed: %s\nInput: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
			correctnesCheck(t, test.name, test.input, test.expectedOutput, result)

		})
	}
}
