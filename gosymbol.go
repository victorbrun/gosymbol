package gosymbol

import (
	"fmt"
)

type Expr interface {
	String() string
	Eval(Arguments) float64
	D(string) Expr
}

type Arguments map[string]float64

type constant struct {
	Expr
	Value float64
}

type variable struct {
	Expr
	Name string
}

type add struct {
	Expr
	LHS Expr
	RHS Expr
}

type sub struct {
	Expr
	LHS Expr
	RHS Expr
}

type mul struct {
	Expr
	LHS Expr
	RHS Expr
}

type div struct {
	Expr
	LHS Expr
	RHS Expr
}

/* Factory */

func Const(val float64) constant {
	return constant{Value: val}
}

func Var(name string) variable {
	return variable{Name: name}
}

func Add(lhs, rhs Expr) add {
	return add{LHS: lhs, RHS: rhs}
}

func Sub(lhs, rhs Expr) sub {
	return sub{LHS: lhs, RHS: rhs}
}

func Mul(lhs, rhs Expr) mul {
	return mul{LHS: lhs, RHS: rhs}
}

func Div(lhs, rhs Expr) div {
	return div{LHS: lhs, RHS: rhs}
}

/* Differentiation rules */

func (e constant) D(varName string) Expr {
	return Const(0.0)
}

func (e variable) D(varName string) Expr {
	if varName == e.Name {
		return Const(1.0)
	} else {
		return Const(0.0)
	}
}

func (e add) D(varName string) Expr {
	return Add(e.LHS.D(varName), e.RHS.D(varName))
}

func (e sub) D(varName string) Expr {
	return Sub(e.LHS.D(varName), e.RHS.D(varName))
}

// Product rule: D(fg) = D(f)g + fD(g)
func (e mul) D(varName string) Expr {
	return Add(
		Mul(e.LHS.D(varName), e.RHS),
		Mul(e.LHS, e.RHS.D(varName)),
	)
}

// Quote rule: D(f/g) = (D(f)g -fD(g))/g^2
func (e div) D(varName string) Expr {
	return Div(
		Sub(
			Mul(e.LHS.D(varName), e.RHS),
			Mul(e.LHS, e.RHS.D(varName)),
		),
		Mul(e.RHS, e.RHS),
	)
}

/* Evaluation functions for operands */

func (e constant) Eval(args Arguments) float64 {
	return e.Value
}

func (e variable) Eval(args Arguments) float64 {
	return args[e.Name]
}

func (e add) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) + e.RHS.Eval(args)
}

func (e sub) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) - e.RHS.Eval(args)
}

func (e mul) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) * e.RHS.Eval(args)
}

func (e div) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) / e.RHS.Eval(args)
}

/* Implementing String() to get nicely formated expressions upon print */

func (e constant) String() string {
	return fmt.Sprint(e.Value)
}

func (e variable) String() string {
	return e.Name
}

func (e add) String() string {
	return fmt.Sprintf("( %v ) + ( %v )", e.LHS, e.RHS)
}

func (e sub) String() string {
	return fmt.Sprintf("( %v ) - ( %v )", e.LHS, e.RHS)
}

func (e mul) String() string {
	return fmt.Sprintf("( %v ) * ( %v )", e.LHS, e.RHS)
}

func (e div) String() string {
	return fmt.Sprintf("( %v ) / ( %v )", e.LHS, e.RHS) 
}

