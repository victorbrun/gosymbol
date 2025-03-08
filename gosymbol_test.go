package gosymbol

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestAddArgument(t *testing.T) {
	type inputArgs struct {
		args  Arguments
		v     variable
		value Expr
	}

	tests := []struct {
		input         inputArgs
		expectedError error
		expectedArgs  Arguments
	}{
		{ // Test 1: testing with duplicate name
			input: inputArgs{
				args:  Arguments{Var("X"): Int(1)},
				v:     Var("X"),
				value: Int(2),
			},
			expectedError: &DuplicateArgumentError{},
			expectedArgs:  Arguments{Var("X"): Int(1)},
		},
		{ // Test 2: testing with no duplicate names
			input: inputArgs{
				args:  Arguments{Var("X"): Int(1)},
				v:     Var("Y"),
				value: Int(2),
			},
			expectedError: nil,
			expectedArgs:  Arguments{Var("X"): Int(1), Var("Y"): Int(2)},
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			err := test.input.args.AddArgument(test.input.v, test.input.value)
			if !errors.Is(err, test.expectedError) {
				errMsg := fmt.Sprintf("Failed test: %v\nExpectedError: %v\nGot: %v", ix+1, test.expectedError, err)
				t.Error(errMsg)
			} else if !reflect.DeepEqual(test.expectedArgs, test.input.args) {
				errMsg := fmt.Sprintf("Failed test: %v\nExpectedArgs: %v\nGot: %v", ix+1, test.expectedArgs, test.input.args)
				t.Error(errMsg)
			}
		})
	}
}

func TestIsBAE(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput bool
	}{
		{
			name:           "2/4 is BAE",
			input:          Div(Int(2), Int(4)),
			expectedOutput: true,
		},
		{
			name:           "a * (x + x)",
			input:          Mul(Var("a"), Add(Var("x"), Var("x"))),
			expectedOutput: true,
		},
		{
			name:           "a + (b^3 / b)",
			input:          Add(Var("a"), Div(Pow(Var("b"), Int(3)), Var("b"))),
			expectedOutput: true,
		},
		{
			name:           "a + ( b + c ) + d",
			input:          Add(Var("a"), Add(Var("b"), Var("c")), Var("d")),
			expectedOutput: true,
		},
		{
			name:           "2 * 3 * x * x^2",
			input:          Mul(Int(2), Int(3), Var("x"), Pow(Var("x"), Int(2))),
			expectedOutput: true,
		},
		{
			name:           "0^3",
			input:          Pow(Int(0), Int(3)),
			expectedOutput: true,
		},
		{
			name:           "2 / (a - a)",
			input:          Div(Int(2), Add(Var("a"), Neg(Var("a")))),
			expectedOutput: true,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			output := IsBAE(test.input)
			if output != test.expectedOutput {
				errMsg := fmt.Sprintf(
					"Failed test: %v\nInput: %v\nExpected: %v\nGot: %v",
					test.name,
					test.input,
					test.expectedOutput,
					output,
				)
				t.Error(errMsg)
			}
		})
	}
}

func TestIsASAE(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput bool
	}{
		// Test ASAE-4
		{
			name:           "2 * x * y * z^2",
			input:          Mul(Int(2), Var("x"), Var("y"), Pow(Var("z"), Int(2))),
			expectedOutput: true,
		},
		{
			name:           "2 * (x * y) * z^2 (ASAE-4-1 is not satisfied)",
			input:          Mul(Int(2), Mul(Var("x"), Var("y")), Pow(Var("z"), Int(2))),
			expectedOutput: false,
		},
		{
			name:           "1 * x * y * z^2 (ASAE-4-1 is not satisfied)",
			input:          Mul(Int(1), Var("x"), Var("y"), Pow(Var("z"), Int(2))),
			expectedOutput: false,
		},
		{
			name:           "2 * x * y * z * z^2 (ASAE-4-3 is not satisfied)",
			input:          Mul(Int(2), Var("x"), Var("y"), Var("z"), Pow(Var("z"), Int(2))),
			expectedOutput: false,
		},
		{
			name:           "2 * x + 3 * y + 4 * z",
			input:          Add(Mul(Int(2), Var("x")), Mul(Int(3), Var("y")), Mul(Int(4), Var("z"))),
			expectedOutput: true,
		},

		// Test ASAE-5
		{
			name:           "1 + ( x + y ) + z (ASAE-5-1 is not satisfied)",
			input:          Add(Int(1), Add(Var("x"), Var("y")), Var("z")),
			expectedOutput: false,
		},
		{
			name:           "1 + 2 + x (ASAE-5-2 is not satisfied)",
			input:          Add(Int(1), Int(2), Var("x")),
			expectedOutput: false,
		},
		{
			name:           "1 + x + 2 * x (ASAE-5-3 is not satisfied)",
			input:          Add(Int(1), Var("x"), Mul(Int(2), Var("x"))),
			expectedOutput: false,
		},
		{
			name:           "z + y + x (ASAE-5-4 is not satisfied)",
			input:          Add(Var("z"), Var("y"), Var("x")),
			expectedOutput: false,
		},

		// Test ASAE-6
		{
			name:           "x^2",
			input:          Pow(Var("x"), Int(2)),
			expectedOutput: true,
		},
		{
			name:           "( 1 + x )^3",
			input:          Pow(Add(Int(1), Var("x")), Int(3)),
			expectedOutput: true,
		},
		{
			name:           "2^m",
			input:          Pow(Int(2), Var("m")),
			expectedOutput: true,
		},
		{
			name:           "(x * y)^(1/2)",
			input:          Pow(Mul(Var("x"), Var("y")), Frac(Int(1), Int(2))),
			expectedOutput: true,
		},
		{
			name:           "(x^(1/2))^(1/2)",
			input:          Pow(Pow(Var("x"), Frac(Int(1), Int(2))), Frac(Int(1), Int(2))),
			expectedOutput: true,
		},
		{
			name:           "2^3",
			input:          Pow(Int(2), Int(3)),
			expectedOutput: false,
		},
		{
			name:           "(x^2)^3",
			input:          Pow(Pow(Var("x"), Int(2)), Int(3)),
			expectedOutput: false,
		},
		{
			name:           "( x * y )^2",
			input:          Pow(Mul(Var("x"), Var("y")), Int(2)),
			expectedOutput: false,
		},
		{
			name:           "( 1 + x )^1",
			input:          Pow(Add(Int(1), Var("x")), Int(1)),
			expectedOutput: false,
		},
		{
			name:           "1^m",
			input:          Pow(Int(1), Var("m")),
			expectedOutput: false,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			output := IsASAE(test.input)
			if output != test.expectedOutput {
				errMsg := fmt.Sprintf(
					"Failed test: %v\nInput: %v\nExpected: %v\nGot: %v",
					test.name,
					test.input,
					test.expectedOutput,
					output,
				)
				t.Error(errMsg)
			}
		})
	}
}

func TestAsaeTerm(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput Expr
	}{
		{
			name:           "term(x) = *x",
			input:          Var("x"),
			expectedOutput: Mul(Var("x")),
		},
		{
			name:           "term(2*y) = *y",
			input:          Mul(Int(2), Var("y")),
			expectedOutput: Mul(Var("y")),
		},
		{
			name:           "term(x*y) = x*y",
			input:          Mul(Var("x"), Var("y")),
			expectedOutput: Mul(Var("x"), Var("y")),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			output := AsaeTerm(test.input)
			if !Equal(output, test.expectedOutput) {
				errMsg := fmt.Sprintf(
					"Failed test: %v\nInput: %v\nExpected: %v\nGot: %v",
					test.name,
					test.input,
					test.expectedOutput,
					output,
				)
				t.Error(errMsg)
			}
		})
	}
}

func TestAsaeConst(t *testing.T) {
	tests := []struct {
		name           string
		input          Expr
		expectedOutput Expr
	}{
		{
			name:           "const(x) = 1",
			input:          Var("x"),
			expectedOutput: Int(1),
		},
		{
			name:           "const(2*y) = 2",
			input:          Mul(Int(2), Var("y")),
			expectedOutput: Int(2),
		},
		{
			name:           "const(x*y) = x*y",
			input:          Mul(Var("x"), Var("y")),
			expectedOutput: Int(1),
		},
		{
			name:           "const(69) = undefined",
			input:          Int(69),
			expectedOutput: Undefined(),
		},
		{
			name:           "const(6/9) = undefined",
			input:          Frac(Int(6), Int(9)),
			expectedOutput: Undefined(),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			output := AsaeConst(test.input)
			if !Equal(output, test.expectedOutput) {
				errMsg := fmt.Sprintf(
					"Failed test: %v\nInput: %v\nExpected: %v\nGot: %v",
					test.name,
					test.input,
					test.expectedOutput,
					output,
				)
				t.Error(errMsg)
			}
		})
	}
}
