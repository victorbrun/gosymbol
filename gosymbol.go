package gosymbol

import (
	"fmt"
	"math"
	"reflect"
	"sort"

)

/* Factories */

func Undefined() undefined {
	return undefined{}
}

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

func Pow(base Expr, exponent Expr) pow {
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

// IF EXPONENT IS CONSTANT: Power rule: D(x^a) = ax^(a-1)
// IF EXPONENT IS NOT CONSTANT: Exponential deriv: D(f^g) = D(exp(g*log(f))) = exp(g*log(f))*D(g*log(f))
func (e pow) D(varName VarName) Expr {
	if exponentTyped, ok := e.Exponent.(constant); ok {
		return Mul(e.Exponent, Pow(e.Base, Const(exponentTyped.Value-1)), e.Base.D(varName))
	} else {
		exponentLogBaseProd := Mul(e.Exponent, Log(e.Base))
		return Mul(Exp(exponentLogBaseProd), exponentLogBaseProd.D(varName))
	}
}

// D(sqrt(f)) = (1/2)*(1/sqrt(f))*D(f)
func (e sqrt) D(varName VarName) Expr {
	return Mul(Div(Const(1), Const(2)), Div(Const(1), e), e.Arg.D(varName))
}

/* Evaluation */

func (e undefined) Eval() Func {
	return func(args Arguments) float64 {return math.NaN()}
}

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

/* Evaluation to string */

func (e undefined) String() string {
	return "Undefined"
}

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

func (e constrainedVariable) String() string {
	return fmt.Sprintf("%v_CONSTRAINED", e.Name)
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

func (e undefined) substitutes(u, t Expr) Expr {
	if _, ok := u.(undefined); ok {
		return t
	} else {
		return e
	}
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

func Equal(t, u Expr) bool {
	return reflect.DeepEqual(t, u)
}

// Returns true if t and u are equal up to type 
// for every element in resp. syntax tree, otherwise 
// false. This means that two constants of different 
// value, or variables with different names, would
// return true.
func TypeEqual(t, u Expr) bool {
	if !isSameType(t, u) {
		return false
	} else if NumberOfOperands(t) != NumberOfOperands(u) {
		return false
	}

	// Base cases are the leaf node types
	switch t.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case constant:
		_, ok := u.(constant)
		return ok
	case variable: 
	_, ok := u.(variable)
		return ok
	}

	// Recusively checks if every operand is type equal 
	// as well. Breaks and returns false if any of the
	// operands are not equal.
	// TODO: this does not take associativty of e.g. add and mul into account.
	for ix := 1; ix <= NumberOfOperands(t); ix++ {
		tOperand := Operand(t, ix)
		uOperand := Operand(u, ix)
		if !TypeEqual(tOperand, uOperand) {
			return false
		}
	}

	return true
}

// Checks if expr contains u by formating expr and
// u to strings and running a sub-string check.
func Contains(expr, u Expr) bool {
	return expr.contains(u)
}

func (e undefined) contains(u Expr) bool {
	_, ok := u.(undefined)
	return ok
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

func (e undefined) variableNames(targetSlice *[]string) {}
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

func (e undefined) numberOfOperands() int {return 0}
func (e constant) numberOfOperands() int {return 0}
func (e variable) numberOfOperands() int {return 0}
func (e add) numberOfOperands() int {return len(e.Operands)}
func (e mul) numberOfOperands() int {return len(e.Operands)}
func (e pow) numberOfOperands() int {return 2} 

// Returns the n:th (starting at 1) operand (left to right) of expr.
// If expr has no operands it returns nil.
// If n is larger than NumberOfOperands(expr)-1 it will panic.
func Operand(expr Expr, n int) Expr {
	return expr.operand(n)
}

func (e undefined) operand(n int) Expr {return nil}
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

func isSameType(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

/* Automatic Simplification */

func Simplify(expr Expr) Expr {
	return expr.simplify()
}

func (e undefined) simplify() Expr {return e}
func (e constant) simplify() Expr {return e}
func (e variable) simplify() Expr {return e} 

func (e add) simplify() Expr {
	// Iterating though every simplification rule
	// and if one matches it is applied and we exit
	// this function.
	for _, rule := range powerSimplificationRules {
		if rule.match(e) {
			return rule.rhs
		} 
	}

	// If no match is found we return input expression.
	return e
}

func (e mul) simplify() Expr {
	// Iterating though every simplification rule
	// and if one matches it is applied and we exit
	// this function.
	for _, rule := range powerSimplificationRules {
		if rule.match(e) {
			return rule.rhs
		} 
	}

	// If no match is found we return input expression.
	return e
}

func (e pow) simplify() Expr {
	// Iterating though every simplification rule
	// and if one matches it is applied and we exit
	// this function.
	for ix, rule := range powerSimplificationRules {	
		fmt.Printf("RULE %v:\n", ix)
		if rule.match(e) {
			// TODO: USE SUBSTITUTE FUNCTION HERE!!!
			return rule.rhs
		} 
	}

	// If no match is found we return input expression.
	return e
}

func (e exp) simplify() Expr {return nil}
func (e log) simplify() Expr {return nil}

// Returns true if expr matches the pattern defined
// by r, else it returns false.
func (rule simplificationRule) match(expr Expr) bool {
	// Rule concerns the top most operator in the 
	// tree so these need to have matching types.
	if !isSameType(rule.lhs, expr) {
		return false
	} 

	// Iterate though each operand of both expr and r.lhs,
	// checking if they match "good enough".
	for ix := 1; ix <= NumberOfOperands(rule.lhs); ix++ {
		ruleOperand := Operand(rule.lhs, ix)
		exprOperand := Operand(expr, ix)
		
		fmt.Printf("	ruleOperand: %v\n", ruleOperand)
		fmt.Printf("	exprOperand: %v\n", exprOperand) 

		if _, ok := ruleOperand.(variable); ok {
			continue
		} else if opTyped, ok := ruleOperand.(constrainedVariable); ok && opTyped.Constraint(exprOperand) {
			continue
		} else if TypeEqual(ruleOperand, exprOperand) {
			continue
		} else if Equal(ruleOperand, exprOperand) {
			continue
		} else {
			return false
		}
	}
	return true
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	return nil
}


