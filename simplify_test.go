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
			input:          FromLatex(`6 + \frac{2}{3}`),
			expectedOutput: FromLatex(`\frac{20}{3}`),
		},
		{
			name:           "Sum between integers is simpliied to a single integer",
			input:          FromLatex(`2 + 3`),
			expectedOutput: Int(5),
		},
		{
			name:           "Sum with two integers and a real in the middle is simplified to a single integers and a real",
			input:          FromLatex(`2 + e + 3`),
			expectedOutput: FromLatex(`5 + e`),
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
			name:           "0^1 = 0",
			input:          FromLatex(`0^1`),
			expectedOutput: Int(0),
		},
		{
			name:           "0^0 = undefined",
			input:          Pow(Int(0), Int(0)),
			expectedOutput: Undefined(),
		},
		{
			name:           "1^x = 1",
			input:          FromLatex(`1^{e^7}`),
			expectedOutput: Int(1),
		},
		{
			name:           "x^0 = 1",
			input:          Pow(Var("kuk"), Int(0)),
			expectedOutput: Int(1),
		},
		{
			name:           "(v_1 * ... * v_n)^m = v_1^m * .. * v_n^m (note that the result is also sorted)",
			input:          FromLatex(`(3xy)^l`),
			expectedOutput: FromLatex(`3^l x^l y^l`),
		},
		{
			name:           "(i^j)^k = i^(j*k)",
			input:          FromLatex(`(i^j)^{\exp(1*k)}`),
			expectedOutput: FromLatex(`i^{\exp(k)j}`),
		},
		{
			name:           "undefined * ... = undefined",
			input:          Mul(Undefined(), Var("x"), Int(1)),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * ... = 0",
			input:          FromLatex("x * (-9) * 0"),
			expectedOutput: Int(0),
		},
		{
			name:           "undefined * 0 = undefined",
			input:          Mul(Undefined(), Int(0)),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 * undefined = undefined",
			input:          Mul(Int(0), Undefined()),
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
			input:          FromLatex(`6\frac{2}{3}`),
			expectedOutput: Int(4),
		},
		{
			name:           "Mult between integers is simpliied to a single integer",
			input:          FromLatex("2*3"),
			expectedOutput: Int(6),
		},
		{
			name:           "Mult with two integers and a real in the middle is simplified to a single integers and a real",
			input:          FromLatex("2e3"),
			expectedOutput: FromLatex("6e"),
		},
		{
			name:           "1 * x = x",
			input:          FromLatex("1e^x"),
			expectedOutput: FromLatex("e^x"),
		},
		{
			name:           "x * x = x^2",
			input:          FromLatex("x*x"),
			expectedOutput: FromLatex("x^2"),
		},
		{
			name:           "x * x^n = x^(n+1)",
			input:          FromLatex("x * x^2"),
			expectedOutput: FromLatex("x^3"),
		},
		{
			name:           "x * (1/x) = 1",
			input:          FromLatex(`x \frac{1}{x}`),
			expectedOutput: Int(1),
		},
		{
			name:           "x^m * x^n = x^(m+n)",
			input:          FromLatex("x^m x^n"),
			expectedOutput: FromLatex("x^{m+n}"),
		},
		{
			name:           "2 * 1",
			input:          FromLatex("2*1"),
			expectedOutput: FromLatex("2"),
		},
		{
			name:           "2 * x^1 * 1",
			input:          FromLatex("2 * x^1 * 1"),
			expectedOutput: FromLatex("2x"),
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
