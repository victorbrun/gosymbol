package gosymbol

import "fmt"

type Expr interface {
	Eval(Arguments) float64
	String() string
}

type Arguments map[string]float64

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

func (e Var) String() string {
	return e.Name
}
func (e Add) String(args Arguments) string {
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
