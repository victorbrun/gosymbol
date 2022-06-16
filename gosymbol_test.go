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
			input: gosymbol.Var{ Name: "X" },
			expectedOutput: "X",
		},
		{
			input: gosymbol.Var{ Name: "X" },
			expectedOutput: "X",
		},
		{
			input: gosymbol.Var{ Name: "X" },
			expectedOutput: "X",
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
