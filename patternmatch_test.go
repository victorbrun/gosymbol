package gosymbol

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	tests := []struct {
		name            string
		inputPattern    Expr
		inputExpression Expr
		expectedOutput  bool
	}{
		{
			name:            "Simple variable test 1",
			inputPattern:    Var("x"),
			inputExpression: Var("x"),
			expectedOutput:  true,
		},
		{
			name:            "Simple variable test 2",
			inputPattern:    Var("x"),
			inputExpression: Var("y"),
			expectedOutput:  true,
		},
		{
			name:            "First expression variable with same name as pattern variable",
			inputPattern:    Add(Var("x"), Var("x")),
			inputExpression: Add(Var("x"), Var("y")),
			expectedOutput:  false,
		},
		{
			name:            "Second expression variable with same name as pattern variable",
			inputPattern:    Add(Var("x"), Var("x")),
			inputExpression: Add(Var("y"), Var("x")),
			expectedOutput:  false,
		},
		{
			name:            "Positive variable matched against positive constant",
			inputPattern:    ConstrVar("x", positiveConstant),
			inputExpression: Const(1),
			expectedOutput:  true,
		},
		{
			name:            "Positive variable matched against zero constant",
			inputPattern:    ConstrVar("x", positiveConstant),
			inputExpression: Const(0),
			expectedOutput:  false,
		},
		{
			name:            "Negative variable matched against zero constant",
			inputPattern:    ConstrVar("x", negOrZeroConstant),
			inputExpression: Const(-1),
			expectedOutput:  true,
		},
		{
			name:            "Advanced test 1",
			inputPattern:    Mul(Var("x"), Exp(Var("x"))),
			inputExpression: Mul(Pow(Const(2), Var("x")), Exp(Pow(Const(2), Var("x")))),
			expectedOutput:  true,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := patternMatch(test.inputExpression, test.inputPattern, make(Binding))

			if !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Following test failed: %s\nInput expression: %v\nInput pattern: %v\nExpected: %v\nGot: %v", test.name, test.inputExpression, test.inputPattern, test.expectedOutput, result)
			}
		})
	}
}
