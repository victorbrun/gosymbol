package gosymbol

import "fmt"

type Expr interface {
	Eval(Arguments) float64
	D(string) Expr
	String() string
}

type Arguments map[string]float64

type Const struct {
	Expr
	Value float64
}

type Var struct {
	Expr
	Name string
}

type Add struct {
	Expr
	LHS Expr
	RHS Expr
}

type Sub struct {
	Expr
	LHS Expr
	RHS Expr
}

type Mul struct {
	Expr
	LHS Expr
	RHS Expr
}

type Div struct {
	Expr
	LHS Expr
	RHS Expr
}

/* Derivateive rules */

func (e Const) D(varName string) Expr {
	return Const{ Value: 0 }
}

func (e Var) D(varName string) Expr {
	if varName == e.Name {
		return Const{ Value: 1.0 }
	} else {
		return Const{ Value: 0.0 }
	}
}

func (e Add) D(varName string) Expr {
	return Add{ RHS: e.RHS.D(varName), LHS: e.LHS.D(varName) }
}

func (e Sub) D(varName string) Expr {
	return Sub{ RHS: e.RHS.D(varName), LHS: e.LHS.D(varName) }
}

// Product rule
func (e Mul) D(varName string) Expr {
	return Add{ 
		LHS: Mul{ 
			LHS: e.LHS.D(varName), 
			RHS: e.RHS, 
		}, 
		RHS: Mul{ 
			LHS: e.LHS, 
			RHS: e.RHS.D(varName), 
		}, 
	}
}

func (e Div) D(varName string) Expr {
	panic("Not implemented")
}

/* Evaluation functions for operands */

func (e Const) Eval(args Arguments) float64 {
	return e.Value
}

func (e Var) Eval(args Arguments) float64 {
	return args[e.Name]
}

func (e Add) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) + e.RHS.Eval(args)
}

func (e Sub) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) - e.RHS.Eval(args)
}

func (e Mul) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) * e.RHS.Eval(args)
}

func (e Div) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) / e.RHS.Eval(args)
}

/* Implementing String() to get nicely formated expressions upon print */

func (e Const) String() string {
	return fmt.Sprint(e.Value)
}

func (e Var) String() string {
	return e.Name
}

func (e Add) String() string {
	return fmt.Sprintf("( %v ) + ( %v )", e.LHS, e.RHS)
}

func (e Sub) String() string {
	return fmt.Sprintf("( %v ) - ( %v )", e.LHS, e.RHS)
}

func (e Mul) String() string {
	return fmt.Sprintf("( %v ) * ( %v )", e.LHS, e.RHS)
}

func (e Div) String() string {
	return fmt.Sprintf("( %v ) / ( %v )", e.LHS, e.RHS) 
}

