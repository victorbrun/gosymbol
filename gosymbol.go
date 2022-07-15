package gosymbol

import (
	"fmt"
	"math"
	"reflect"
	"sort"

)

type VarName string
type Arguments map[VarName]float64
type Func func(Arguments) float64

type Expr interface {
	// Private functions
	contains(Expr) bool
	substitute(Expr, Expr) Expr
	variableNames(*[]string)
	numberOfOperands() int
	operand(int) Expr

	// Public functions
	String() string
	Eval() Func
	D(VarName) Expr
}

/* Basic operators */

type constant struct {
	Expr
	Value float64
}

type variable struct {
	Expr
	Name VarName
}

type add struct {
	Expr
	Operands []Expr
}

type mul struct {
	Expr
	Operands []Expr
}

/* Common Functions */

type exp struct {
	Expr
	Arg Expr
}

type log struct {
	Expr
	Arg Expr
}

type pow struct {
	Expr
	Base Expr
	Exponent constant // For expressions in exponent use exp.
}

type sqrt struct {
	Expr
	Arg Expr
}

/* Factories */

func Const(val float64) constant {
	return constant{Value: val}
}

func Var(name VarName) variable {
	return variable{Name: name}
}

func Neg(arg Expr) mul {
	return Mul(Const(-1), arg)
}

func Add(ops...Expr) add {
	return add{Operands: ops}
}

func Sub(lhs, rhs Expr) add {
	return Add(lhs, Neg(rhs))
}

func Mul(ops ...Expr) mul {
	return mul{Operands: ops}
}

func Div(lhs, rhs Expr) mul {
	return Mul(lhs, Pow(rhs, Const(-1)))
}

func Exp(arg Expr) exp {
	return exp{Arg: arg}
}

func Log(arg Expr) log {
	return log{Arg: arg}
}

func Pow(base Expr, exponent constant) pow {
	return pow{Base: base, Exponent: exponent}
}

func Sqrt(arg Expr) sqrt {
	return sqrt{Arg: arg}
}

/* Differentiation rules */

func (e constant) D(varName VarName) Expr {
	return Const(0.0)
}

func (e variable) D(varName VarName) Expr {
	if varName == e.Name {
		return Const(1.0)
	} else {
		return Const(0.0)
	}
}

func (e add) D(varName VarName) Expr {
	differentiatedOps := make([]Expr, len(e.Operands))
	for ix, op := range e.Operands {
		differentiatedOps[ix] = op.D(varName)	
	}
	return Add(differentiatedOps...)
}

// Product rule: D(fghijk...) = D(f)ghijk... + fD(g)hijk... + ....
func (e mul) D(varName VarName) Expr {
	terms := make([]Expr, len(e.Operands))
	for ix := 0; ix < len(e.Operands); ix++ {
		var productOperands []Expr
		copy(productOperands, e.Operands)
		productOperands[ix] = productOperands[ix].D(varName)
		terms[ix] = Mul(productOperands...)
	}
	return Add(terms...)
}

func (e exp) D(varName VarName) Expr {
	return Mul(e, e.Arg.D(varName))
}

func (e log) D(varName VarName) Expr {
	return Mul(Pow(e.Arg, Const(-1)), e.Arg.D(varName))
}

// Power rule: D(x^a) = ax^(a-1)
func (e pow) D(varName VarName) Expr {
	return Mul(e.Exponent, Pow(e.Base, Const(e.Exponent.Value-1)), e.Base.D(varName))
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(varName VarName) Expr {
	return Mul(Div(Const(1), Const(2)), Div(Const(1), e), e.Arg.D(varName))
}

/* Evaluation functions for operands */

func (e constant) Eval() Func {
	return func(args Arguments) float64 {return e.Value}
}

func (e variable) Eval() Func {
	return func(args Arguments) float64 {return args[e.Name]}
}

func (e add) Eval() Func {
	return func(args Arguments) float64 {
		sum := e.Operands[0].Eval()(args) // Initiate with first operand since 0 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			sum += e.Operands[ix].Eval()(args)
		}
		return sum
	}
}

func (e mul) Eval() Func {
	return func(args Arguments) float64 {
		prod := e.Operands[0].Eval()(args) // Initiate with first operand since 1 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			prod *= e.Operands[ix].Eval()(args)
		}
		return prod
	}
}

func (e exp) Eval() Func {
	return func(args Arguments) float64 {return math.Exp(e.Arg.Eval()(args))}
}

func (e log) Eval() Func {
	return func(args Arguments) float64 {return math.Log(e.Arg.Eval()(args))}
}

func (e pow) Eval() Func {
	return func(args Arguments) float64 {return math.Pow(e.Base.Eval()(args), e.Exponent.Eval()(args))}
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
	return string(e.Name)
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

func (e exp) String() string {
	return fmt.Sprintf("exp( %v )", e.Arg)
}

func (e log) String() string {
	return fmt.Sprintf("log( %v )", e.Arg)
}

func (e pow) String() string {
	return fmt.Sprintf("( %v^%v )", e.Base, e.Exponent)
}

// Substitutes u for t in expr.
func Substitute(expr, u, t Expr) Expr {
	// While expr still contains u we 
	// continue to substitute. This is catch
	// the case when a substitution leads to 
	// another substitutable sub-expression 
	// (see test 7 in gosymbol_test.go).
	for Contains(expr, u) {
		expr = expr.substitute(u, t)
	}
	return expr.substitute(u, t)
}

func (e constant) substitute(u, t Expr) Expr {
	if uTyped, ok := u.(constant); ok && e.Value == uTyped.Value {
		return t
	} else {
		return e
	}
}

func (e variable) substitute(u, t Expr) Expr {
	if uTyped, ok := u.(variable); ok && e.Name == uTyped.Name {
		return t
	} else {
		return e
	}
}

func (e add) substitute(u, t Expr) Expr {
	// If e equals u we return t, otherwise
	// we run substitute to possibly alter every
	// operand of e and then returns e.
	if reflect.DeepEqual(e, u) {
		return t
	} else {
		for ix, op := range e.Operands {
			e.Operands[ix] = op.substitute(u, t)
		}
		return e
	}
}

func (e mul) substitute(u, t Expr) Expr { 
	// If e equals u we return t, otherwise
	// we run substitute to possibly alter every
	// operand of e and then returns e.
	if reflect.DeepEqual(e, u) {
		return t
	} else {
		for ix, op := range e.Operands {
			e.Operands[ix] = op.substitute(u, t)
		}
		return e
	}
}

func (e exp) substitute(u, t Expr) Expr {
	if reflect.DeepEqual(e, u) {
		return t
	}
	e.Arg = e.Arg.substitute(u, t)
	return e
}

func (e log) substitute(u, t Expr) Expr {
	if reflect.DeepEqual(e, u) {
		return t
	}
	e.Arg = e.Arg.substitute(u, t)
	return e
}

func (e pow) substitute(u, t Expr) Expr {
	if reflect.DeepEqual(e, u) {
		return t
	} else if tTyped, ok := t.(constant); reflect.DeepEqual(e.Exponent, u) && ok {
		e.Exponent = tTyped
		return e
	} else {
		e.Base = e.Base.substitute(u, t)
		return e
	}
}

// Checks if expr contains u by formating expr and
// u to strings and running a sub-string check.
func Contains(expr, u Expr) bool {
	return expr.contains(u)
}

func (e constant) contains(u Expr) bool {	
	if uTyped, ok := u.(constant); ok && e.Value == uTyped.Value {
		return true
	} else {
		return false
	}
}

func (e variable) contains(u Expr) bool {	
	if uTyped, ok := u.(variable); ok && e.Name == uTyped.Name {
		return true
	} else {
		return false
	}
}

func (e add) contains(u Expr) bool {
	if reflect.DeepEqual(e, u) {
		return true
	}
	
	cumBool := false
	for _, op := range e.Operands {
		cumBool = cumBool || op.contains(u)	
	}
	return cumBool
}

func (e mul) contains(u Expr) bool {
	if reflect.DeepEqual(e, u) {
		return true
	}
	
	cumBool := false
	for _, op := range e.Operands {
		cumBool = cumBool || op.contains(u)	
	}
	return cumBool
}

func (e exp) contains(u Expr) bool {
	if reflect.DeepEqual(e, u) {
		return true
	}
	return e.Arg.contains(u)
} 

func (e log) contains(u Expr) bool {
	if reflect.DeepEqual(e, u) {
		return true
	}
	return e.Arg.contains(u)
} 

func (e pow) contains(u Expr) bool {
	if reflect.DeepEqual(e, u) {
		return true
	}
	return e.Base.contains(u)
}

// Returns the different variable names 
// present in the given expression.
func VariableNames(expr Expr) []VarName {
	var stringSlice []string
	expr.variableNames(&stringSlice)
	if len(stringSlice) == 0 {
		return []VarName{}
	}

	// Checks for and deletes duplicates of variable names
	// TODO: speed up by checking for duplicates during sorting
	sort.Strings(stringSlice)
	var variableNamesSlice []VarName
	variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[0]))
	jx := 0
	for ix := 1; ix < len(stringSlice); ix++ {
		if variableNamesSlice[jx] != VarName(stringSlice[ix]) {
			variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[ix]))
			jx++
		}
	}

	return variableNamesSlice
}

func (e constant) variableNames(targetSlice *[]string) {}

func (e variable) variableNames(targetSlice *[]string) {
	*targetSlice = append(*targetSlice, string(e.Name))
}

func (e add) variableNames(targetSlice *[]string) {
	for _, op := range e.Operands {
		op.variableNames(targetSlice)
	}
}

func (e mul) variableNames(targetSlice *[]string) {
	for _, op := range e.Operands {
		op.variableNames(targetSlice)
	}
}

func (e pow) variableNames(targetSlice *[]string) {
	e.Base.variableNames(targetSlice)
}

func (e exp) variableNames(targetSlice *[]string) {
	e.Arg.variableNames(targetSlice)
}

func (e log) variableNames(targetSlice *[]string) {
	e.Arg.variableNames(targetSlice)
}

// Returns the number of operands for top level operation.
func NumberOfOperands(expr Expr) int {
	return expr.numberOfOperands()
}

func (e constant) numberOfOperands() int {return 0}
func (e variable) numberOfOperands() int {return 0}
func (e add) numberOfOperands() int {return len(e.Operands)}
func (e mul) numberOfOperands() int {return len(e.Operands)}
func (e pow) numberOfOperands() int {return 2} // TODO: Should this be 1 or 2?

// Returns the n:th (starting at 1) operand (left to right) of expr.
// If expr has no operands it returns nil.
// If n is larger than NumberOfOperands(expr)-1 it will panic.
func Operand(expr Expr, n int) Expr {
	return expr.operand(n)
}

func (e constant) operand(n int) Expr {return nil}
func (e variable) operand(n int) Expr {return nil}
func (e add) operand(n int) Expr {
	if n > NumberOfOperands(e) {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, len(e.Operands))
		panic(errMsg)
	}
	return e.Operands[n-1]
}
func (e mul) operand(n int) Expr {
	if n > NumberOfOperands(e) {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, len(e.Operands))
		panic(errMsg)
	}
	return e.Operands[n-1]
}
func (e pow) operand(n int) Expr {	
	if n > NumberOfOperands(e) {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, 2)
		panic(errMsg)
	} else if n == 1 {
		return e.Base
	} else {
		return e.Exponent
	}
}


// TODO: figure this out
func Simplify(expr Expr) Expr {
	// To implement:
	// 
	return nil
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	return nil
}

