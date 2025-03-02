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
			input:          Div(Const(2), Const(4)),
			expectedOutput: true,
		},
		{
			name:           "a * (x + x)",
			input:          Mul(Var("a"), Add(Var("x"), Var("x"))),
			expectedOutput: true,
		},
		{
			name:           "a + (b^3 / b)",
			input:          Add(Var("a"), Div(Pow(Var("b"), Const(3)), Var("b"))),
			expectedOutput: true,
		},
		{
			name:           "a + ( b + c ) + d",
			input:          Add(Var("a"), Add(Var("b"), Var("c")), Var("d")),
			expectedOutput: true,
		},
		{
			name:           "2 * 3 * x * x^2",
			input:          Mul(Const(2), Const(3), Var("x"), Pow(Var("x"), Const(2))),
			expectedOutput: true,
		},
		{
			name:           "0^3",
			input:          Pow(Const(0), Const(3)),
			expectedOutput: true,
		},
		{
			name:           "2 / (a - a)",
			input:          Div(Const(2), Add(Var("a"), Neg(Var("a")))),
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
	}{}

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
