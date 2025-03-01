package gosymbol

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestApprox(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput float64
	}{
		{
			name:           "approx of fraction 1/2",
			input:          Div(Int(1), Int(2)).(rational),
			expectedOutput: 0.5,
		},
		{
			name:           "approx of fraction 1/0 (infinity)",
			input:          Div(Int(1), Int(0)).(rational),
			expectedOutput: math.Inf(1),
		},
		{
			name:           "approx of empty fraction NaN",
			input:          Undefined(),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx value of pi",
			input:          PI,
			expectedOutput: math.Pi,
		},
		{
			name:           "approx of var(x) is NaN",
			input:          Var("x"),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx of 1 + 2 is 3",
			input:          Add(Int(1), Int(2)),
			expectedOutput: 3.0,
		},
		{
			name:           "approx of 1 + x is NaN",
			input:          Add(Int(1), Var("x")),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx of 2 * 3 is 6",
			input:          Mul(Int(2), Int(3)),
			expectedOutput: 6.0,
		},
		{
			name:           "approx of 2 * x is NaN",
			input:          Add(Int(2), Var("x")),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx of 2 ^ 3 is 8",
			input:          Pow(Int(2), Int(3)),
			expectedOutput: 8.0,
		},
		{
			name:           "approx of 2 ^ x is NaN",
			input:          Pow(Int(2), Var("x")),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx E is e",
			input:          E,
			expectedOutput: math.E,
		},
		{
			name:           "approx Exp(2) is e^2 from math package",
			input:          Exp(Int(2)),
			expectedOutput: math.Exp(2),
		},
		{
			name:           "approx of Exp(x) is NaN",
			input:          Exp(Var("x")),
			expectedOutput: math.NaN(),
		},
		{
			name:           "approx Log(2) is log(2) for math package",
			input:          Log(Int(2)),
			expectedOutput: math.Log(2),
		},
		{
			name:           "approx of Log(x) is NaN",
			input:          Log(Var("x")),
			expectedOutput: math.NaN(),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.Approx()
			bothNans := math.IsNaN(result) && math.IsNaN(test.expectedOutput)
			if !bothNans && !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Following test failed: %s\nInput expr: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}
