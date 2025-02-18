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
