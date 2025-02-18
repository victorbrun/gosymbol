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
			inputPattern:    patternVar("x"),
			inputExpression: Var("x"),
			expectedOutput:  true,
		},
		{
			name:            "Simple variable test 2",
			inputPattern:    patternVar("x"),
			inputExpression: Var("y"),
			expectedOutput:  true,
		},
		{
			name:            "First expression variable with same name as pattern variable",
			inputPattern:    Add(patternVar("x"), patternVar("x")),
			inputExpression: Add(Var("x"), Var("y")),
			expectedOutput:  false,
		},
		{
			name:            "Second expression variable with same name as pattern variable",
			inputPattern:    Add(patternVar("x"), patternVar("x")),
			inputExpression: Add(Var("y"), Var("x")),
			expectedOutput:  false,
		},
		/*{
			name:            "Positive variable matched against positive constant",
			inputPattern:    constraPatternVar("x", positiveConstant),
			inputExpression: (Int(1)),
			expectedOutput:  true,
		},
		{
			name:            "Positive variable matched against zero constant",
			inputPattern:    constraPatternVar("x", positiveConstant),
			inputExpression: (Int(0)),
			expectedOutput:  false,
		},
		{
			name:            "Non positive variable matched against negative constant",
			inputPattern:    constraPatternVar("x", negOrZeroConstant),
			inputExpression: (Int(-1)),
			expectedOutput:  true,
		},*/
		{
			name:            "Advanced test with same variable name in pattern and expression",
			inputPattern:    Mul(patternVar("x"), Exp(patternVar("x"))),
			inputExpression: Mul(Pow((Int(2)), Var("x")), Exp(Pow((Int(2)), Var("x")))),
			expectedOutput:  true,
		},
		{
			name:            "Advanced test with different variable name in pattern and expression",
			inputPattern:    Mul(patternVar("y"), Exp(patternVar("y"))),
			inputExpression: Mul(Pow((Int(2)), Var("x")), Exp(Pow((Int(2)), Var("x")))),
			expectedOutput:  true,
		},
		{
			name:            "Advanced test with different variable name in pattern and expression returning false",
			inputPattern:    Mul(patternVar("y"), Exp(patternVar("y"))),
			inputExpression: Mul(Pow((Int(2)), Var("x")), Exp(Pow((Int(1)), Var("x")))),
			expectedOutput:  false,
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
