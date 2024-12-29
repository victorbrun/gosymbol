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
		value float64
	}

	tests := []struct {
		input         inputArgs
		expectedError error
		expectedArgs  Arguments
	}{
		{ // Test 1: testing with duplicate name
			input: inputArgs{
				args:  Arguments{Var("X"): 1},
				v:     Var("X"),
				value: 2,
			},
			expectedError: errors.New("multiple variables have the same name"),
			expectedArgs:  Arguments{Var("X"): 1},
		},
		{ // Test 1: testing with no duplicate names
			input: inputArgs{
				args:  Arguments{Var("X"): 1},
				v:     Var("Y"),
				value: 2,
			},
			expectedError: nil,
			expectedArgs:  Arguments{Var("X"): 1, Var("Y"): 2},
		},
	}

	for ix, test := range tests {
		err := test.input.args.AddArgument(test.input.v, test.input.value)
		if !errors.Is(err, test.expectedError) || !reflect.DeepEqual(test.expectedArgs, test.input.args) {
			errMsg := fmt.Sprintf("Failed test: %v. ExpectedError: %v, Got: %v; ExpectedArgs: %v, Got: %v", ix+1, test.expectedError, err, test.expectedArgs, test.input.args)
			t.Error(errMsg)
		}
	}
}
