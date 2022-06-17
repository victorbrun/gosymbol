package gosymbol_test

import (
	"fmt"
	"testing"
	"github.com/victorbrun/gosymbol"
)

func TestExprPrint(t *testing.T) {
	tests := []struct{
		input gosymbol.Expr
		expectedOutput string
	} {
		{
			input: gosymbol.Var{ Name: "X" },
			expectedOutput: "X",
		},
		{
			input: gosymbol.Add{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
			expectedOutput: "( X ) + ( Y )",
		},
		{
			input: gosymbol.Sub{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
			expectedOutput: "( X ) - ( Y )",
		},
		{
			input: gosymbol.Mul{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
			expectedOutput: "( X ) * ( Y )",
		},
		{
			input: gosymbol.Div{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
			expectedOutput: "( X ) / ( Y )",
		},
	}
	
	for ix, test := range tests {
		result := fmt.Sprint(test.input)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf("Failed test: %v. Expected: %v \nGot: %v", ix+1, test.expectedOutput, result)
			t.Error(errMsg)
		}
	}
}

func TestExprEval(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		args gosymbol.Arguments
	}

	tests := []struct{
		input inputArgs
		expectedOutput float64
	} {
		{
			input: inputArgs{
				expr: gosymbol.Var{ Name: "X" },
				args: gosymbol.Arguments{ "X": 1.0 },
			},
			expectedOutput: 1.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Add{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 3.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Sub{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: -1.0,
	
		},
		{
			input: inputArgs{
				expr: gosymbol.Mul{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 2.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Div{ LHS: gosymbol.Var{ Name: "X" }, RHS: gosymbol.Var{ Name: "Y" } },
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 0.5,
		},
	}

	for ix, test := range tests {
		result := test.input.expr.Eval(test.input.args)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf("Failed test: %v\nExpected: %v\nGot: %v", ix+1, test.expectedOutput, result)
			t.Error(errMsg)
		}
	}
}




