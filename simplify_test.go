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
			name:           "0^(-1) = undefined",
			input:          Pow(Const(0), Const(-1)),
			expectedOutput: Undefined(),
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
			name:           "1^x = 1",
			input:          Pow(Const(1), Var("x")),
			expectedOutput: Const(1),
		},
		{
			name:           "x^0 = 1",
			input:          Pow(Var("kuk"), Const(0)),
			expectedOutput: Const(1),
		},
		{
			name:           "x^1 = x",
			input:          Pow(Var("x"), Const(1)),
			expectedOutput: Var("x"),
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
		{
			name:           "x + 0 = x",
			input:          Add(Var("x"), Const(0)),
			expectedOutput: Var("x"),
		},
		{
			name:           "0 + x = x",
			input:          Add(Const(0), Var("x")),
			expectedOutput: Var("x"),
		},

		// Below tests are taken from pages 64 - 68 in
		// COHEN, Joel S. Computer algebra and symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003.
		{
			name:           "( a + b ) * c = a * c + b * c",
			input:          Mul(Add(Var("a"), Var("b")), Var("c")),
			expectedOutput: Add(Mul(Var("a"), Var("c")), Mul(Var("b"), Var("c"))),
		},
		{
			name:           "2x + y + (3/2)x = (7/2)x + y",
			input:          Add(Mul(Const(2), Var("x")), Var("y"), Mul(Div(Const(3), Const(2)), Var("x"))),
			expectedOutput: Add(Mul(Div(Const(7), Const(2)), Var("x")), Var("y")),
		},
		{
			name:           "x + 1 + (-1)(x + 1) = 0",
			input:          Add(Var("x"), Const(1), Mul(Const(-1), Add(Var("x"), Const(1)))),
			expectedOutput: Const(0),
		},
		{
			name:           "( ( ( u + v ) * w ) * x + y ) + z = ( u + v ) * w * x + y + z ",
			input:          Add(Add(Mul(Mul(Add(Var("u"), Var("v")), Var("w")), Var("x")), Var("y")), Var("z")),
			expectedOutput: Add(Mul(Add(Var("u"), Var("v")), Var("w"), Var("x")), Var("y"), Var("z")),
		},
		{
			name:           "2(xyz) + 3x(yz) + 4(xy)z = 9xyz",
			input:          Add(Mul(Const(2), Mul(Var("x"), Var("y"), Var("z"))), Mul(Const(3), Var("x"), Mul(Var("y"), Var("z"))), Mul(Const(4), Mul(Var("x"), Var("y")), Var("z"))),
			expectedOutput: Mul(Const(9), Var("x"), Var("y"), Var("z")),
		},

		// Below tests are taken from pages 74 - 76 in
		// COHEN, Joel S. Computer algebra and symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003.
		{
			name:           "( x * y ) * ( x * y )^2 = x^3 * y^3",
			input:          Mul(Mul(Var("x"), Var("y")), Pow(Mul(Var("x"), Var("y")), Const(2))),
			expectedOutput: Mul(Pow(Var("x"), Const(3)), Pow(Var("y"), Const(3))),
		},
		{
			name:           "( x * y ) * ( x * y )^(1/2) = x * y * ( x  * y )^(1/2)",
			input:          Mul(Mul(Var("x"), Var("y")), Pow(Mul(Var("x"), Var("y")), Div(Const(1), Const(2)))),
			expectedOutput: Mul(Pow(Var("x"), Const(3)), Pow(Var("y"), Const(3))),
		},
		{
			name:           "( a * x^3 ) / x = a * x^2",
			input:          Div(Mul(Var("a"), Pow(Var("x"), Const(3))), Var("x")),
			expectedOutput: Mul(Var("a"), Pow(Var("x"), Const(2))),
		},

		// Below tests are taken from pages 77 in
		// COHEN, Joel S. Computer algebra and symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003.
		{
			name:           "u / 0 = undefined",
			input:          Div(Var("u"), Const(0)),
			expectedOutput: Undefined(),
		},
		{
			name:           "0 / u = 0 (u ~= 0)",
			input:          Div(Const(0), Const(1)),
			expectedOutput: Const(0),
		},
		{
			name:           "0 / u = 0 (u ~= 0)",
			input:          Div(Const(0), Const(-1)),
			expectedOutput: Const(0),
		},
		{
			name:           "u / 1",
			input:          Div(Var("u"), Const(1)),
			expectedOutput: Var("u"),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
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
