package gosymbol_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/victorbrun/gosymbol"
)


func correctnesCheck(t *testing.T, result, expectedOutput any, testNumber int) {
		if !reflect.DeepEqual(result, expectedOutput) {
			errMsg := fmt.Sprintf("Failed test: %v. Expected: %v, Got: %v", testNumber, expectedOutput, result)
			t.Error(errMsg)
		}
}

func TestExprPrint(t *testing.T) {
	tests := []struct{
		input gosymbol.Expr
		expectedOutput string
	} {
		{
			input: gosymbol.Var("X"),
			expectedOutput: "X",
		},
		{
			input: gosymbol.Const(-1),
			expectedOutput: "( -1 )",
		},
		{
			input: gosymbol.Const(10),
			expectedOutput: "10",
		},
		{
			input: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X + Y )",
		},
		{
			input: gosymbol.Sub(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X + ( ( -1 ) * Y ) )",
		},
		{
			input: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X * Y )",
		},
		{
			input: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: "( X * ( Y^( -1 ) ) )",
		},
		{
			input: gosymbol.Exp(gosymbol.Var("X")),
			expectedOutput: "exp( X )",
		},
		{
			input: gosymbol.Log(gosymbol.Var("X")),
			expectedOutput: "log( X )",
		},
		{
			input: gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(9)),
			expectedOutput: "( X^9 )",
		},
	}
	
	for ix, test := range tests {
		result := fmt.Sprint(test.input)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
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
				expr: gosymbol.Const(0),
				args: gosymbol.Arguments{ "X": 10 },
			},
			expectedOutput: 0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Var("X"),
				args: gosymbol.Arguments{ "X": 1.0 },
			},
			expectedOutput: 1.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Neg(gosymbol.Const(10)),
				args: gosymbol.Arguments{"X": 0},
			},
			expectedOutput: -10,
		},
		{
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 3.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Sub(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: -1.0,
	
		},
		{
			input: inputArgs{
				expr: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Y")),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 2.0,
		},
		{
			input: inputArgs{
				expr: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Y" )),
				args: gosymbol.Arguments{ "X": 1.0, "Y": 2.0 },
			},
			expectedOutput: 0.5,
		},
		{
			input: inputArgs{
				expr: gosymbol.Exp(gosymbol.Var("X")),
				args: gosymbol.Arguments{"X": 0},
			},
			expectedOutput: 1,
		},
		{
			input: inputArgs{
				expr: gosymbol.Log(gosymbol.Var("X")),
				args: gosymbol.Arguments{"X": 1},
			},
			expectedOutput: math.Log(1),
		},
		{
			input: inputArgs{
				expr: gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(-1)),
				args: gosymbol.Arguments{"X": 10},
			},
			expectedOutput: 1/10.0,
		},
	}

	for ix, test := range tests {
		result := test.input.expr.Eval()(test.input.args)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestD(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		diffVar gosymbol.VarName
	}

	tests := []struct{
		input inputArgs
		expectedOutput gosymbol.Expr
	} {
		{ // Test 1
			input: inputArgs{
				expr: gosymbol.Const(10),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Const(0),
		},
		{ // Test 2
			input: inputArgs{
				expr: gosymbol.Var("X"),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Const(1),
		},
		{ // Test 3
			input: inputArgs{
				expr: gosymbol.Var("X"),
				diffVar: "Y",
			},
			expectedOutput: gosymbol.Const(0),
		},
		/*{ // Test 4
			input: inputArgs{
				expr: gosymbol.Exp(gosymbol.Var("X")),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Exp(gosymbol.Var("X")),
		},
		{ // Test 5
			input: inputArgs{
				expr: gosymbol.Log(gosymbol.Var("X")),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Div(gosymbol.Const(1), gosymbol.Var("X")),
		},
		{ // Test 6
			input: inputArgs{
				expr: gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(2)),
				diffVar: "X",
			},
			expectedOutput: gosymbol.Mul(gosymbol.Const(2), gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(1))),
		},*/
	}

	for ix, test := range tests {
		result := test.input.expr.D(test.input.diffVar)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestContains(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		u gosymbol.Expr
	}

	tests := []struct {
		input inputArgs
		expectedOutput bool
	} { 
		{ // Test 1: testing for equal expressions
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("X"),
			},
			expectedOutput: true,
		},
		{ // Test 2: Testing for inequality
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("Y"),
			},
			expectedOutput: false,
		},
		{ // Test 3: Testing for part of n-ary operator
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Const(1), gosymbol.Const(2), gosymbol.Const(3)),
				u: gosymbol.Add(gosymbol.Const(1), gosymbol.Const(2)),
			},
			expectedOutput: false,
		},
		{ // Test 4: Testing for sub-tree equality
			input: inputArgs{
				expr: gosymbol.Add(
					gosymbol.Const(1), 
					gosymbol.Mul(gosymbol.Const(2), gosymbol.Var("X")), 
					gosymbol.Div(gosymbol.Const(1), gosymbol.Var("y")),
				),
				u: gosymbol.Mul(gosymbol.Const(2), gosymbol.Var("X")),
			},
			expectedOutput: true,
		},
	}

	for ix, test := range tests {
		result := gosymbol.Contains(test.input.expr, test.input.u)
		if result != test.expectedOutput {
			errMsg := fmt.Sprintf(
				"Failed test: %v: Expected: %v, Got: %v. expr = %v, u = %v",
				ix+1, test.expectedOutput, result, test.input.expr, test.input.u,
			)
			t.Error(errMsg)
		}
	}
}

func TestSubstitute(t *testing.T) {
	type inputArgs struct {
		expr gosymbol.Expr
		u gosymbol.Expr
		t gosymbol.Expr
	}

	tests := []struct {
		input inputArgs
		expectedOutput gosymbol.Expr
	} {
		{ // Test 1: base case 1: constant
			input: inputArgs{
				expr: gosymbol.Const(7),
				u: gosymbol.Const(7),
				t: gosymbol.Const(-7),
			},
			expectedOutput: gosymbol.Const(-7),
		},
		{ // Test 2: base case 2: variable
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("X"),
				t: gosymbol.Const(0),
			},
			expectedOutput: gosymbol.Const(0),
		},
		{ // Test 3: add operator
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Y"), gosymbol.Const(0), gosymbol.Var("Y")),
				u: gosymbol.Var("Y"),
				t: gosymbol.Var("Z"),
			},
			expectedOutput: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("Z"), gosymbol.Const(0), gosymbol.Var("Z")),
		},
		{ // Test 4: mul operator
			input: inputArgs{
				expr: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Y"), gosymbol.Const(0), gosymbol.Var("Y")),
				u: gosymbol.Var("Y"),
				t: gosymbol.Var("Z"),
			},
			expectedOutput: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("Z"), gosymbol.Const(0), gosymbol.Var("Z")),
		},
		{ // Test 5: Div
			input: inputArgs{
				expr: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Y")),
				u: gosymbol.Var("Y"),
				t: gosymbol.Var("Z"),
			},
			expectedOutput: gosymbol.Div(gosymbol.Var("X"), gosymbol.Var("Z")),
		},
		{ // Test 6: substituting whole subtree
			input: inputArgs{
				expr: gosymbol.Add(gosymbol.Div(gosymbol.Const(9), gosymbol.Var("X")), gosymbol.Var("Y"), gosymbol.Const(0), gosymbol.Var("Y")),
				u: gosymbol.Div(gosymbol.Const(9), gosymbol.Var("X")),
				t: gosymbol.Var("Z"),
			},
			expectedOutput: gosymbol.Add(gosymbol.Var("Z"), gosymbol.Var("Y"), gosymbol.Const(0), gosymbol.Var("Y")),
		},
		{ // Test 7: nested substitution, making sure that the recursion starts bottom up
			input: inputArgs{
				expr: gosymbol.Div(gosymbol.Const(1), gosymbol.Div(gosymbol.Const(1), gosymbol.Var("X"))),
				u: gosymbol.Div(gosymbol.Const(1), gosymbol.Var("X")),
				t: gosymbol.Var("X"),
			},
			expectedOutput: gosymbol.Var("X"),
		},
		{ // Test 8: no expression matching u
			input: inputArgs{
				expr: gosymbol.Var("X"),
				u: gosymbol.Var("Y"),
				t: gosymbol.Const(0),
			},
			expectedOutput: gosymbol.Var("X"),
		},
		{ // Test 9: exp operator
			input: inputArgs{
				expr: gosymbol.Exp(gosymbol.Var("X")),
				u: gosymbol.Exp(gosymbol.Var("X")),
				t: gosymbol.Var("Y"),
			},
			expectedOutput: gosymbol.Var("Y"),
		},
		{ // Test 10: log operator
			input: inputArgs{
				expr: gosymbol.Log(gosymbol.Var("X")),
				u: gosymbol.Log(gosymbol.Var("X")),
				t: gosymbol.Var("Y"),
			},
			expectedOutput: gosymbol.Var("Y"),
		},
		{ // Test 11: pow operator
			input: inputArgs{
				expr: gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(10)),
				u: gosymbol.Pow(gosymbol.Var("X"), gosymbol.Const(10)),
				t: gosymbol.Var("Y"),
			},
			expectedOutput: gosymbol.Var("Y"),
		},
		{ // Test 12: 4 level nested substitution
			input: inputArgs{
				expr: gosymbol.Div(gosymbol.Const(1), gosymbol.Div(gosymbol.Const(1),gosymbol.Div(gosymbol.Const(1), gosymbol.Div(gosymbol.Const(1), gosymbol.Var("X"))))),
				u: gosymbol.Div(gosymbol.Const(1), gosymbol.Var("X")),
				t: gosymbol.Var("X"),
			},
			expectedOutput: gosymbol.Var("X"),
			
		},
	}

	for ix, test := range tests {
		result := gosymbol.Substitute(test.input.expr, test.input.u, test.input.t)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestVariableNames(t *testing.T) {
	tests := []struct {
		input gosymbol.Expr
		expectedOutput []gosymbol.VarName
	}{
		{
			input: gosymbol.Undefined(),
			expectedOutput: []gosymbol.VarName{},
		},
		{
			input: gosymbol.Const(0),
			expectedOutput: []gosymbol.VarName{},
		},
		{
			input: gosymbol.Var("X"),
			expectedOutput: []gosymbol.VarName{"X"},
		},
		{
			input: gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: []gosymbol.VarName{"X", "Y"},
		},
		{
			input: gosymbol.Add(gosymbol.Var("X"), gosymbol.Var("X"), gosymbol.Var("Y")),
			expectedOutput: []gosymbol.VarName{"X", "Y"},
		},
		{
			input: gosymbol.Pow(gosymbol.Mul(gosymbol.Var("X"), gosymbol.Var("X"), gosymbol.Var("Y")), gosymbol.Const(10)), 
			expectedOutput: []gosymbol.VarName{"X", "Y"},
		},
	}

	for ix, test := range tests {
		result := gosymbol.VariableNames(test.input)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestSimplify(t *testing.T) {
	tests := []struct {
		input gosymbol.Expr
		expectedOutput gosymbol.Expr
	} {
		{ // undefined^y = undefined
			input: gosymbol.Pow(gosymbol.Undefined(), gosymbol.Var("y")),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // x^undefined = undefined
			input: gosymbol.Pow(gosymbol.Var("x"), gosymbol.Undefined()),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // 0^x = 0
			input: gosymbol.Pow(gosymbol.Const(0), gosymbol.Const(10)),
			expectedOutput: gosymbol.Const(0),
		},
		{ // 0^0 = undefined
			input: gosymbol.Pow(gosymbol.Const(0), gosymbol.Const(0)),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // 1^x = undefined
			input: gosymbol.Pow(gosymbol.Const(1), gosymbol.Exp(gosymbol.Const(7))),
			expectedOutput: gosymbol.Const(1),
		},
		{ // x^0 = 1
			input: gosymbol.Pow(gosymbol.Var("kuk"), gosymbol.Const(0)),
			expectedOutput: gosymbol.Const(1),
		},
		{ // (v_1 * ... * v_n)^m = v_1^m * .. * v_n^m (note that the result is also sorted)
			input: gosymbol.Pow(gosymbol.Mul(gosymbol.Var("x"), gosymbol.Const(3), gosymbol.Var("y")), gosymbol.Var("elle")),
			expectedOutput: gosymbol.Mul(gosymbol.Pow(gosymbol.Const(3), gosymbol.Var("elle")), gosymbol.Pow(gosymbol.Var("x"), gosymbol.Var("elle")), gosymbol.Pow(gosymbol.Var("y"), gosymbol.Var("elle"))),
		},
		{ // (i^j)^k = i^(j*k)
			input: gosymbol.Pow(gosymbol.Pow(gosymbol.Var("i"), gosymbol.Var("j")), gosymbol.Exp(gosymbol.Mul(gosymbol.Const(10), gosymbol.Var("k")))),
			expectedOutput: gosymbol.Pow(gosymbol.Var("i"), gosymbol.Mul(gosymbol.Var("j"), gosymbol.Exp(gosymbol.Mul(gosymbol.Const(10), gosymbol.Var("k"))))),
		},
		{ // undefined * ... = undefined
			input: gosymbol.Mul(gosymbol.Undefined(), gosymbol.Var("x"), gosymbol.Const(10)),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // 0 * ... = 0
			input: gosymbol.Mul(gosymbol.Var("x"), gosymbol.Const(-9), gosymbol.Const(0)),
			expectedOutput: gosymbol.Const(0),
		},
		{ // undefined * 0 = undefined 
			input: gosymbol.Mul(gosymbol.Undefined(), gosymbol.Const(0)),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // 0 * undefined = undefined
			input: gosymbol.Mul(gosymbol.Const(0), gosymbol.Undefined()),
			expectedOutput: gosymbol.Undefined(),
		},
		{ // Mult with only one operand simplifies to the operand
			input: gosymbol.Mul(gosymbol.Exp(gosymbol.Var("x"))),
			expectedOutput: gosymbol.Exp(gosymbol.Var("x")),
		},
		{ // Mult with no operands simplify to 1 
			input: gosymbol.Mul(),
			expectedOutput: gosymbol.Const(1),
		},
		{ // 1 * x = x 
			input: gosymbol.Mul(gosymbol.Const(1), gosymbol.Exp(gosymbol.Var("x"))),
			expectedOutput: gosymbol.Exp(gosymbol.Var("x")),
		},
		{ // x*x = x^2
			input: gosymbol.Mul(gosymbol.Const(10), gosymbol.Const(10)),
			expectedOutput: gosymbol.Pow(gosymbol.Const(10), gosymbol.Const(2)),
		},
		{ // x*x^n = x^(n+1)
			input: gosymbol.Mul(gosymbol.Const(10), gosymbol.Pow(gosymbol.Const(10), gosymbol.Const(2))),
			expectedOutput: gosymbol.Pow(gosymbol.Const(10), gosymbol.Const(3)),
		},
		{ // x * (1/x) = 1
			input: gosymbol.Mul(gosymbol.Var("x"), gosymbol.Div(gosymbol.Const(1), gosymbol.Var("x"))),
			expectedOutput: gosymbol.Const(1),
		},
		{ // x^m * x^n = x^(m+n)
			input: gosymbol.Mul(gosymbol.Pow(gosymbol.Var("x"), gosymbol.Var("n")), gosymbol.Pow(gosymbol.Var("x"), gosymbol.Var("m"))),
			expectedOutput: gosymbol.Pow(gosymbol.Var("x"), gosymbol.Add(gosymbol.Var("m"), gosymbol.Var("n"))),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			//fmt.Println("Simplifying: ", test.input)
			result := gosymbol.Simplify(test.input)
			correctnesCheck(t, result, test.expectedOutput, ix+1)
		})
	}
}

func TestDepth(t *testing.T) {
	tests := []struct {
		input gosymbol.Expr
		expectedOutput int
	} {
		{
			input: gosymbol.Undefined(),
			expectedOutput: 0,
		},
		{
			input: gosymbol.Const(1),
			expectedOutput: 0,
		},
		{
			input: gosymbol.Var("x"),
			expectedOutput: 0,
		},
		{
			input: gosymbol.ConstrVar("x", func(expr gosymbol.Expr) bool {return true}),
			expectedOutput: 0,
		},
		{
			input: gosymbol.Add(gosymbol.Const(0), gosymbol.Var("x"),gosymbol.Const(0), gosymbol.Var("x"),gosymbol.Const(0), gosymbol.Var("x")),
			expectedOutput: 1,
		},
		{
			input: gosymbol.Add(gosymbol.Mul(gosymbol.Var("x"),gosymbol.Pow(gosymbol.Const(10), gosymbol.Exp(gosymbol.Var("x")))), gosymbol.Var("x"),gosymbol.Const(0), gosymbol.Var("x"),gosymbol.Const(0), gosymbol.Var("x")),
			expectedOutput: 4,
		},
	}

	for ix, test := range tests {
		result := gosymbol.Depth(test.input)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}
}

func TestTemp(t *testing.T) {
	tests := []struct{
		input gosymbol.Expr
		expectedOutput gosymbol.Expr
	}{ 
		{	// x * (1/x) = 1
			input: gosymbol.Mul(gosymbol.Var("x"), gosymbol.Div(gosymbol.Const(1), gosymbol.Var("x"))),
			expectedOutput: gosymbol.Const(1),
		},
	}

	for ix, test := range tests {
		result := gosymbol.Simplify(test.input)
		correctnesCheck(t, result, test.expectedOutput, ix+1)
	}

	//e := gosymbol.Mul(gosymbol.Pow(gosymbol.Var("x"), gosymbol.Const(3)), gosymbol.Var("x"))
	e1 := gosymbol.Mul(gosymbol.Var("x"), gosymbol.Div(gosymbol.Const(1), gosymbol.Var("x")))
	fmt.Println("unsorted: ", e1)
	fmt.Println("sorted: ", gosymbol.TopOperandSort(gosymbol.Add(gosymbol.Var("m"), gosymbol.Const(-1), gosymbol.Const(9), gosymbol.Neg(gosymbol.Var("m")))))
}
