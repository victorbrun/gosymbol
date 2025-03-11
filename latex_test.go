package gosymbol

import (
	"reflect"
	"testing"
)

func TestParseExpression(t *testing.T) {
	inputs := []string{
		`\var`,
		`5`,
		`10`,
		`-5`,
		`5 + 10`,
		`5 * 10`,
		`5 / 10`,
		`5 ^ 10`,
		`5 + 5 / 10`,
		`(5 + 5) / 10`,
		`\log(10)`,
		`\exp(10)`,
		`5 + 5 a + \alpha`,
		`ab`,
		`(a)(b)`,
		`a \alpha`,
		`\alpha a`,
		`\pi`,
		`e`,
		`\frac{1}{2}`,
	}

	tests := []Expr{
		Var(VarName("\\var")),
		Int(5),
		Int(10),
		Add(Int(0), Mul(Int(-1), Int(5))),
		Add(Int(5), Int(10)),
		Mul(Int(5), Int(10)),
		Div(Int(5), Int(10)),
		Pow(Int(5), Int(10)),
		Add(Int(5), Div(Int(5), Int(10))),
		Div(Add(Int(5), Int(5)), Int(10)),
		Log(Int(10)),
		Exp(Int(10)),
		Add(Add(Int(5), Mul(Int(5), Var("a"))), Var("\\alpha")),
		Mul(Var(VarName("a")), Var(VarName("b"))),
		Mul(Var(VarName("a")), Var(VarName("b"))),
		Mul(Var(VarName("a")), Var(VarName("\\alpha"))),
		Mul(Var(VarName("\\alpha")), Var(VarName("a"))),
		PI,
		E,
		Div(Int(1), Int(2)),
	}

	for i, input := range inputs {
		expr := FromLatex(input)
		if !reflect.DeepEqual(tests[i], expr) {
			t.Errorf("Expr not '%s'. got=%s", tests[i].String(), expr.String())
			return
		}
	}
}
