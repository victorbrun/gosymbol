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
			name:           "Sum between rationals is simplified to a single rational",
			input:          Add(Int(6), Frac(Int(2), Int(3))),
			expectedOutput: Frac(Int(20), Int(3)),
		},
		{
			name:           "Sum between integers is simpliied to a single integer",
			input:          Add(Int(2), Int(3)),
			expectedOutput: Int(5),
		},
		{
			name:           "Sum with two integers and a real in the middle is simplified to a single integers and a real",
			input:          Add(Int(2), Real("e"), Int(3)),
			expectedOutput: Add(Int(5), Real("e")),
		},
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
			input:          Pow((Int(0)), (Int(1))),
			expectedOutput: (Int(0)),
		},
		{
			name:           "0^0 = undefined",
			input:          Pow((Int(0)), (Int(0))),
			expectedOutput: Undefined(),
		},
		{
			name:           "1^x = 1",
			input:          Pow((Int(1)), Exp((Int(7)))),
			expectedOutput: (Int(1)),
		},
		{
			name:           "x^0 = 1",
			input:          Pow(Var("kuk"), (Int(0))),
			expectedOutput: (Int(1)),
		},
		{
			name:           "(v_1 * ... * v_n)^m = v_1^m * .. * v_n^m (note that the result is also sorted)",
			input:          Pow(Mul(Var("x"), (Int(3)), Var("y")), Var("elle")),
			expectedOutput: Mul(Pow((Int(3)), Var("elle")), Pow(Var("x"), Var("elle")), Pow(Var("y"), Var("elle"))),
		},
		{
			name:           "(i^j)^k = i^(j*k)",
			input:          Pow(Pow(Var("i"), Var("j")), Exp(Mul((Int(1)), Var("k")))),
			expectedOutput: Pow(Var("i"), Mul(Var("j"), Exp(Var("k")))),
		},
		{
			name:           "undefined * ... = undefined",
			input:          Mul(Undefined(), Var("x"), (Int(1))),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * ... = 0",
			input:          Mul(Var("x"), (Int(-9)), (Int(0))),
			expectedOutput: (Int(0)),
		},
		{
			name:           "undefined * 0 = undefined",
			input:          Mul(Undefined(), (Int(0))),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * undefined = undefined",
			input:          Mul((Int(0)), Undefined()),
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
			expectedOutput: Int(1),
		},
		{
			name:           "Mult between rationals is simplified to a single rational",
			input:          Mul(Int(6), Frac(Int(2), Int(3))),
			expectedOutput: Int(4),
		},
		{
			name:           "Mult between integers is simpliied to a single integer",
			input:          Mul(Int(2), Int(3)),
			expectedOutput: Int(6),
		},
		{
			name:           "Mult with two integers and a real in the middle is simplified to a single integers and a real",
			input:          Mul(Int(2), Real("e"), Int(3)),
			expectedOutput: Mul(Int(6), Real("e")),
		},
		{
			name:           "1 * x = x",
			input:          Mul((Int(1)), Exp(Var("x"))),
			expectedOutput: Exp(Var("x")),
		},
		{
			name:           "x * x = x^2",
			input:          Mul(Var(VarName("x")), Var(VarName("x"))),
			expectedOutput: Pow(Var(VarName("x")), (Int(2))),
		},
		{
			name:           "x * x^n = x^(n+1)",
			input:          Mul(Var(VarName("x")), Pow(Var(VarName("x")), Int(2))),
			expectedOutput: Pow(Var(VarName("x")), Int(3)),
		},
		{
			name:           "x * (1/x) = 1",
			input:          Mul(Var("x"), Div((Int(1)), Var("x"))),
			expectedOutput: (Int(1)),
		},
		{
			name:           "x^m * x^n = x^(m+n)",
			input:          Mul(Pow(Var("x"), Var("n")), Pow(Var("x"), Var("m"))),
			expectedOutput: Pow(Var("x"), Add(Var("m"), Var("n"))),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			//fmt.Println("Simplifying: ", test.input)
			result := Simplify(test.input)

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
