package gosymbol

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string
	Contains(Expr) bool
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
	Operands []Expr
}

type mul struct {
	Expr
	Operands []Expr
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

func Neg(arg Expr) mul {
	return Mul(Const(-1), arg)
}

func Add(ops...Expr) add {
	return add{Operands: ops}
}

func Sub(lhs, rhs Expr) add {
	return Add(lhs, Mul(Const(-1), rhs))
}

func Mul(ops ...Expr) mul {
	return mul{Operands: ops}
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
	differentiatedOps := make([]Expr, len(e.Operands))
	for ix, op := range e.Operands {
		differentiatedOps[ix] = op.D(varName)	
	}
	return Add(differentiatedOps...)
}

// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
func (e mul) D(varName string) Expr {
	terms := make([]Expr, len(e.Operands))
	for ix := 0; ix < len(e.Operands); ix++ {
		var productOperands []Expr
		copy(productOperands, e.Operands)
		productOperands[ix] = productOperands[ix].D(varName)
		terms[ix] = Mul(productOperands...)
	}
	return Add(terms...)
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
	sum := e.Operands[0].Eval(args) // Initiate with first operand since 0 may not always be identity
	for ix := 1; ix < len(e.Operands); ix++ {
		sum += e.Operands[ix].Eval(args)
	}
	return sum
}

func (e mul) Eval(args Arguments) float64 {
	prod := e.Operands[0].Eval(args) // Initiate with first operand since 1 may not always be identity
	for ix := 1; ix < len(e.Operands); ix++ {
		prod *= e.Operands[ix].Eval(args)
	}
	return prod
}

func (e div) Eval(args Arguments) float64 {
	return e.LHS.Eval(args) / e.RHS.Eval(args)
}

/* Implementing String() to get nicely formated expressions upon print */

func (e constant) String() string {
	if e.Value < 0 {
		return fmt.Sprintf("( %v )", e.Value)	
	} else {
		return fmt.Sprint(e.Value)
	}
}

func (e variable) String() string {
	return e.Name
}

func (e add) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" + %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e mul) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" * %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e div) String() string {
	return fmt.Sprintf("( %v / %v )", e.LHS, e.RHS) 
}

// Substitutes u for t in expr.
func Substitute(expr, u, t Expr) Expr {
	return nil
}

// Checks if expr contains u by formating expr and
// u to strings and running a sub-string check.
func Contains(expr, u Expr) bool {
	return strings.Contains(fmt.Sprint(expr), fmt.Sprint(u))
}


