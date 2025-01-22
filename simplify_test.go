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
			name:           "1^x = 1",
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
		{
			name:           "2 * 1",
			input:          Mul(Const(2), Const(1)),
			expectedOutput: Const(2),
		},
		{
			name:           "2 * x^1 * 1",
			input:          Mul(Const(2), Pow(Var("x"), Const(1)), Const(1)),
			expectedOutput: Mul(Const(2), Var("x")),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			//fmt.Println("Simplifying: ", test.input)
			result := test.input.Simplify()

			if !reflect.DeepEqual(result, test.expectedOutput) {
				t.Errorf("Following test failed: %s\nInput: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestSimplificationRulesForPatternVariables(t *testing.T) {
	executeTests := func(t *testing.T, ruleSlice []transformationRule, testName string) {
		for ix, rule := range ruleSlice {
			// This test only applies for rules
			// without patternfunc
			if rule.pattern == nil {
				continue
			}

			t.Run(fmt.Sprintf("%s-%d", testName, ix+1), func(t *testing.T) {
				notOkVarsInPattern := nonPatternVariablesIn(rule.pattern)

				// Fails on all the rules with more than zero not ok variables
				if len(notOkVarsInPattern) != 0 {
					t.Errorf("Pattern %s contains the following non-pattern varaibles: %v", rule.pattern, notOkVarsInPattern)
				}

			})
		}
	}

	executeTests(t, sumSimplificationRules, "Sum")
	executeTests(t, productSimplificationRules, "Product")
	executeTests(t, powerSimplificationRules, "Power")
}

/* HELPER FUNCTIONS */

func nonPatternVariablesIn(expr Expr) []Expr {
	// Extracts all variables and if these are
	// patterns or not
	varsInExpr := Variables(expr)
	okVars := make([]bool, len(varsInExpr))
	for jx, v := range varsInExpr {
		switch v := v.(type) {
		case variable:
			okVars[jx] = v.isPattern
		case constrainedVariable:
			okVars[jx] = v.isPattern
		default:
			panic("something went wrong this code should not be accessed")
		}

	}

	// Negates bool slice to get the indexes which
	// are not ok
	notOkVars := negateSlice(okVars)
	notOkVarsInExpr := filterByBoolSlice(varsInExpr, notOkVars)

	return notOkVarsInExpr
}
