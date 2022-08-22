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

func ConstrVar(name VarName, constrFunc func(Expr) bool) constrainedVariable {
	return constrainedVariable{Name: name, Constraint: constrFunc}
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

func TransformationRule(pattern Expr, transform func(Expr) Expr) transformationRule {
	return transformationRule{pattern: pattern, transform: transform}
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

/* Helper Functionality */

// Substitutes u for t in expr.
func Substitute(expr, u, t Expr) Expr {
	if Equal(u, t) {
		return u
	} else if Equal(expr, u) {
		return t
	} else if Contains(expr, u) {	
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			processedOp := Substitute(Operand(expr, ix), u, t)
			expr = replaceOperand(expr, ix, processedOp)
		}
		// Subsituting here again in order to continue to Substitute
		// if the result from above substitution is something that can be
		// Substituted again. Se test 7 and test 12 in gosymbol_test.go.
		return Substitute(expr, u, t)
	} else {
		return expr
	}
}

/*
Replaces operand number n in t with u and returns the resulting
expression. The function panics if n is larger than 
the NumberOfOperands(t).
*/
func replaceOperand(t Expr, n int, u Expr) Expr {
	nop := NumberOfOperands(t)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	} else if n <= 0 {
		errMsg := fmt.Sprintf("ERROR: there exists no non-positive indexed operands, you are trying to replace operand: %v", nop)
		panic(errMsg)
	}

	// Since the first cases has no operands
	// we consider replacing one as just returning
	// the original expression
	switch v := t.(type) {
	case undefined:
		return v
	case constant:
		return v
	case variable:
		return v
	case constrainedVariable:
		return v
	case add:
		v.Operands[n-1] = u
		return v
	case mul:
		v.Operands[n-1] = u
		return v
	case pow:
		if n == 1 {
			v.Base = u
		} else {
			v.Exponent = u
		}
		return v
	case exp:
		v.Arg = u
		return v
	case log:
		v.Arg = u
		return v
	case sqrt:
		v.Arg = u
		return v
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

/* 
Recursively checks exact syntactical equality between t and u,
i.e. it does not simplify any expression nor does it 
take any properties, e.g. commutativity, into account.
*/
func Equal(t, u Expr) bool {
	switch v := t.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case constant:
		uTyped, ok := u.(constant)
		return ok && v.Value == uTyped.Value
	case variable:
		uTyped, ok := u.(variable)
		return ok && v.Name == uTyped.Name
	case constrainedVariable:
		// TODO: how do we check equality of constrain??
		return false
	case add:
		_, ok := u.(add)
		if !ok {
			return false
		} else if NumberOfOperands(v) != NumberOfOperands(u) {
			return false
		}

		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			exprOp := Operand(v, ix)
			uOp := Operand(u, ix)
			if !Equal(exprOp, uOp) {
				return false
			}
		}
		return true
	case mul:
		_, ok := u.(mul)
		if !ok {
			return false
		} else if NumberOfOperands(v) != NumberOfOperands(u) {
			return false
		}

		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			exprOp := Operand(v, ix)
			uOp := Operand(u, ix)
			if !Equal(exprOp, uOp) {
				return false
			}
		}
		return true
	case pow:
		_, ok := u.(pow)
		return ok && Equal(Operand(v, 1), Operand(u, 1)) && Equal(Operand(v, 2), Operand(u, 2))
	case exp:
		_, ok := u.(exp)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	case log:
		_, ok := u.(log)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	case sqrt:
		_, ok := u.(sqrt)
		return ok && Equal(Operand(v, 1), Operand(u, 1))
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
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

/* 
Recursively checks exact equality between expr and u,
i.e. it does not simplify any expression nor does it 
take any properties, e.g. commutativity, into account.
*/
func Contains(expr, u Expr) bool {
	switch v := expr.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case constant:
		uTyped, ok := u.(constant)
		return ok && v.Value == uTyped.Value
	case variable:
		uTyped, ok := u.(variable)
		return ok && v.Name == uTyped.Name
	case constrainedVariable:
		// TODO: how do we check equality of constrain??
		return false
	default:
		if Equal(v,u) {return true}
		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			vOp := Operand(v, ix)
			if Contains(vOp, u) {
				return true
			}
		}
		return false
	}
}

// Returns the different variable names 
// present in the given expression.
func VariableNames(expr Expr) []VarName {
	var stringSlice []string
	variableNames(expr, &stringSlice)
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

/* 
Recursively travesrses the whole AST and appends
the variable names to targetSlice. 
*/
func variableNames(expr Expr, targetSlice *[]string) {
	switch v := expr.(type) {
	case undefined:
		return 
	case constant:
		return
	case variable:
		*targetSlice = append(*targetSlice, string(v.Name))
	case constrainedVariable:
		*targetSlice = append(*targetSlice, string(v.Name))
	default:
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			variableNames(op, targetSlice)
		}
	}
}

// Returns the number of operands for top level operation.
func NumberOfOperands(expr Expr) int {
	switch v := expr.(type) {
	case undefined:
		return 0 
	case constant:
		return 0
	case variable:
		return 0
	case constrainedVariable:
		return 0
	case add:
		return len(v.Operands)
	case mul:
		return len(v.Operands)
	case pow:
		return 2
	case exp:
		return 1
	case log:
		return 1
	case sqrt:
		return 1
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// Returns the n:th (starting at 1) operand (left to right) of expr.
// If expr has no operands it returns nil.
// If n is larger than NumberOfOperands(expr)-1 it will panic.
func Operand(expr Expr, n int) Expr {
	nop := NumberOfOperands(expr)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	}

	switch v := expr.(type) {
	case undefined:
		return nil
	case constant:
		return nil
	case variable:
		return nil
	case constrainedVariable:
		return nil
	case add:
		return v.Operands[n-1]
	case mul:
		return v.Operands[n-1]
	case pow:
		if n == 1 {
			return v.Base
		} else {
			return v.Exponent
		}
	case exp:
		return v.Arg
	case log:
		return v.Arg
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// TODO: see Computer Algebra and Symbolic Computation page 10 to understand this shit
func Map(F Expr, u ...Expr) Expr {return nil}


func isSameType(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// Returns the deepest depth of expr
func Depth(expr Expr) int {
	switch expr.(type) {
	case undefined:
		return 0
	case constant:
		return 0
	case variable:
		return 0
	case constrainedVariable:
		return 0
	default:
		maxDepth := 0
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			opDepth := Depth(op)
			if opDepth > maxDepth {
				maxDepth = opDepth
			}
		}
		return maxDepth + 1
	}
}

// A concurrent implementation of Depth 
func DepthConcurrent(expr Expr) int {
	return -1
}

/* Automatic Simplification */

func Simplify(expr Expr) Expr {
	// Recusively simplify all operands
	for ix := 1; ix <= NumberOfOperands(expr); ix++ {
		op := Operand(expr, ix)
		expr = replaceOperand(expr, ix, Simplify(op)) 
	}

	// Application of all Simplification rules follow this same pattern
	rulesApplication := func(expr Expr, ruleSlice []transformationRule) Expr {
		for _, rule := range ruleSlice {
			expr = rule.apply(expr)
		}
		return expr
	}

	// Applies simplification rules depending on the operator type
	switch expr.(type) {
	case undefined:
		return expr
	case constant:
		return expr
	case variable:
		return expr
	case add:
		return rulesApplication(expr, sumSimplificationRules)
	case mul:
		return rulesApplication(expr, productSimplificationRules)	
	case pow:
		return rulesApplication(expr, powerSimplificationRules)
	default:
		return expr
	}
}

// Applies rule to expr and returns the transformed expression.
// If expression does not match rule the ingoing expression 
// will just be returned.
func (rule transformationRule) apply(expr Expr) Expr {
	if rule.match(expr) {
		return rule.transform(expr)
	}
	return expr
}

func (rule transformationRule) match(expr Expr) bool {
	// Fisrt check if pattern is defined. If not
	// we execute patternFunction if it exists. 
	// If no pattern or patternFunction exists we return false 
	if rule.pattern != nil {
		varCache := make(map[VarName]Expr)
		return patternMatch(rule.pattern, expr, varCache)
	} else if rule.patternFunction != nil {
		return rule.patternFunction(expr)
	} else {
		return false
	}
}

// Recursively checks if expr matches pattern. varCache is an empty
// map internally used to keep track of what the variables in pattern 
// corresponds to in expr. The function expects that no variable has the 
// same name as a constrained variable.
func patternMatch(pattern, expr Expr, varCache map[VarName]Expr) bool {
	switch v := pattern.(type) {
	case undefined:
		_, ok := expr.(undefined)
		return ok
	case constant:
		exprTyped, ok := expr.(constant)
		return ok && v.Value == exprTyped.Value 
	case variable:
		// If we have come accross the variable name
		// before it have an assigned expression, so 
		// we check if current expression matches. If
		// we have not come accross the variable name 
		// we cache the current expression.
		if e, ok := varCache[v.Name]; ok {
			return patternMatch(e, expr, varCache) // Should this be patternMatch or Equal??
		} else {
			varCache[v.Name] = expr
			return true
		}
	case constrainedVariable:
		// Does just as above but before assigning an expression
		// to a variable the constraint function is checked as well
		if e, ok := varCache[v.Name]; ok {
			return patternMatch(e, expr, varCache)
		} else if v.Constraint(expr) {
			varCache[v.Name] = expr
			return true
		} else {
			return false
		}
	case add:
		_, ok := expr.(add)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache) 
	case mul:
		_, ok := expr.(mul)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case pow:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case exp:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	case log:
		_, ok := expr.(pow)
		if !ok {
			return false
		}
		return patternMatchOperands(v, expr, varCache)
	default:
		errMsg := fmt.Errorf("ERROR: expression of type: %v have no matchPattern case implemented", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// Checks if the operands of pattern and expr match.
// This function does not check if the main operator
// of pattern and expr match.
func patternMatchOperands(pattern, expr Expr, varCache map[VarName]Expr) bool {
	if NumberOfOperands(pattern) != NumberOfOperands(expr) {
		return false
	}

	// Recursively checks if each operand matches
	for ix := 1; ix <= NumberOfOperands(pattern); ix++ {
		patternOp := Operand(pattern, ix)
		exprOp := Operand(expr, ix)
		if !patternMatch(patternOp, exprOp, varCache) {
			return false
		}
	}
	return true
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	return nil
}


