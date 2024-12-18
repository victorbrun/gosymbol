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
Swaps operand number n1 with operand number n2 in expr.
*/
func swapOperands(expr Expr, n1, n2 int) Expr {
	op1 := Operand(expr, n1)
	op2 := Operand(expr, n2)
	expr = replaceOperand(expr, n1, op2)
	expr = replaceOperand(expr, n2, op1)
	return expr
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
	case sqrt:
		return v.Arg
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
		panic(errMsg)
	}
}

// TODO: see Computer Algebra and Symbolic Computation page 10 to understand this shit
func Map(F Expr, u ...Expr) Expr {panic("Not implemented yet")}


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

/* Automatic Simplification */

func Simplify(expr Expr) Expr {
	// Having this here makes it possible
	// to remove all rules in simplification_rules.go 
	// that basically just checks if the expression contains 
	// undefined.
	if Contains(expr, Undefined()) {
		return Undefined()
	}

	// Only sorting the top operands is sufficient
	// to sort the whole expression since in the next
	// step we recursively simplify all the operands.
	// Note that the operator must be commutative for 
	// this not to fuck shit up!
	switch expr.(type) {
	case add:
		expr = topOperandSort(expr)
	case mul:
		expr = topOperandSort(expr)
	}

	// Recusively simplify all operands.
	for ix := 1; ix <= NumberOfOperands(expr); ix++ {
		op := Operand(expr, ix)
		expr = replaceOperand(expr, ix, Simplify(op)) 
	}
	
	// Application of all Simplification rules follow this same pattern.
	// Returns the simplified expression and a boolean describing whether any
	// simplification rule has actually been applied.
	rulesApplication := func(expr Expr, ruleSlice []transformationRule) (Expr, bool) {
		atLeastOneapplied := false
		for _, rule := range ruleSlice {
			var applied bool
			expr, applied = rule.apply(expr)
			//if applied { fmt.Println("Applied rule ", ix) }
			atLeastOneapplied = atLeastOneapplied || applied
		}
		return expr, atLeastOneapplied
	}

	// Applies simplification rules depending on the operator type
	// This will extend as more rules gets added! The base cases
	// are fully simplified so we just return them.
	expressionAltered := false
	switch expr.(type) {
	case constant:
		// Fully simplified
	case variable:
		// Fully simplified 
	case constrainedVariable:
		// Fully simplified 
	case add:
		expr, expressionAltered = rulesApplication(expr, sumSimplificationRules)
	case mul:
		expr, expressionAltered = rulesApplication(expr, productSimplificationRules)	
	case pow:
		expr, expressionAltered = rulesApplication(expr, powerSimplificationRules)
	}

	// If the expression has been altered it might be possible to apply some other rule 
	// we thus recursively sort until the expression is not altered any more.
	if expressionAltered {
		return Simplify(expr)
	}
	return expr
}

// Applies rule to expr and returns the transformed expression.
// If expression does not match rule the ingoing expression 
// will just be returned.
func (rule transformationRule) apply(expr Expr) (Expr, bool) {
	if rule.match(expr) {
		return rule.transform(expr), true
	}
	return expr, false
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

/*
Recursively checks if expr matches pattern. varCache is an empty
map internally used to keep track of what the variables in pattern 
corresponds to in expr. The function expects that no variable has the 
same name as a constrained variable.
*/
func patternMatch(pattern, expr Expr, varCache map[VarName]Expr) bool {
	switch v := pattern.(type) {
	case undefined:
		_, ok := expr.(undefined)
		return ok
	case constant:
		exprTyped, ok := expr.(constant)
		return ok && v.Value == exprTyped.Value 
	case variable:
		e, cacheOk := varCache[v.Name]
		eTyped, varOk := e.(variable)
		if cacheOk && Equal(e, expr) {
			return true
		} else if cacheOk && varOk && Equal(v, eTyped) {
			return false
		} else if cacheOk && varOk {
			return patternMatch(e, expr, varCache)	
		} else if cacheOk {
			return patternMatch(e, expr, varCache)
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

/*
Checks whether the ordering e1 < e2 is true.
The function returns true if e1 "comes before" e2 and false otherwis false.
"comes before" is defined using the order relation defined in [1] (with some 
extensions to include functions like exp, sin, etc).
E.g.:
O-1: if e1 and e2 are constants then compare(e1, e2) -> e1 < e2
O-2: if e1 and e2 are variables compare(e1, e2) is defined by the
lexographical order of the symbols.
O-3: etc.

NOTE: the function assumes that e1 and e2 are automatically simplified algebraic expressions (ASAEs)

NOTE: When TypeOf(e1) != TypeOf(e2) the recursive evaluation pattern creates a new expression 
of the same type as either e1 or e2. For TypeOf(e1) = add, TypeOf(e2) = mul this looks like
	return compare(Mul(e1), e2)
What this mean in practice is that the type of e2 is prioritised higher than the type of e1.
When extending this function you utilise this to specifiy, e.g. that a < x^3 and not the other 
way around.

[1] COHEN, Joel S. Computer algebra and symbolic computation: Mathematical methods. AK Peters/CRC Press, 2003. Figure 3.9.
*/
func compare(e1, e2 Expr) bool {
	switch e1Typed := e1.(type) {
	case constant:
		switch e2Typed := e2.(type) {
		case constant:
			return orderRule1(e1Typed, e2Typed) 
		default:
			return true
		}
	case variable:
		switch e2Typed := e2.(type) {
		case constant:
			return false
		case variable:
			return orderRule2(e1Typed, e2Typed)
		case constrainedVariable:
			return e1Typed.Name < e2Typed.Name // This is very ugly :(
		case add:
			return compare(Add(e1), e2)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			return compare(Exp(e1), e2)
		case log:
			return compare(Log(e1), e2)
		case sqrt:
			return compare(Sqrt(e1), e2)
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case constrainedVariable:
		switch e2Typed := e2.(type) {
		case constant:
			return false
		case variable:
			return e1Typed.Name < e2Typed.Name // This is very ugly :(
		case constrainedVariable:
			return orderRule2_1(e1Typed, e2Typed)
		case add:
			return compare(Add(e1), e2)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			return compare(Exp(e1), e2)
		case log:
			return compare(Log(e1), e2)
		case sqrt:
			return compare(Sqrt(e1), e2)
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case add:
		switch e2Typed := e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Add(e2))
		case constrainedVariable:
			return compare(e1, Add(e2))
		case add:
			return orderRule3(e1Typed, e2Typed)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			return compare(e1, Add(e2))
		case log:
			return compare(e1, Add(e2))
		case sqrt:
			return compare(e1, Add(e2))
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case mul:
		switch e2Typed := e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Mul(e2))
		case constrainedVariable:
			return compare(e1, Mul(e2))
		case add: 
			return compare(e1, Mul(e2))
		case mul:
			return orderRule3_1(e1Typed, e2Typed)
		case pow:
			return compare(e1, Mul(e2))
		case exp:
			return compare(e1, Mul(e2))
		case log:
			return compare(e1, Mul(e2))
		case sqrt:
			return compare(e1, Mul(e2))
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case pow:
		switch e2Typed := e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Pow(e2, Const(1)))
		case constrainedVariable:
			return compare(e1, Pow(e2, Const(1)))
		case add: 
			return compare(e1, Pow(e2, Const(1)))
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return orderRule4(e1Typed, e2Typed)
		case exp:
			return compare(e1, Pow(e2, Const(1)))
		case log:
			return compare(e1, Pow(e2, Const(1)))
		case sqrt:
			return compare(e1, Pow(e2, Const(1)))
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case exp:
		switch e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Exp(e2))
		case constrainedVariable:
			return compare(e1, Exp(e2))
		case add:
			return compare(Add(e1), e2)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case log:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case sqrt:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case log:
		switch e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Exp(e2))
		case constrainedVariable:
			return compare(e1, Exp(e2))
		case add:
			return compare(Add(e1), e2)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case log:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case sqrt:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	case sqrt:
		switch e2.(type) {
		case constant:
			return false
		case variable:
			return compare(e1, Exp(e2))
		case constrainedVariable:
			return compare(e1, Exp(e2))
		case add:
			return compare(Add(e1), e2)
		case mul:
			return compare(Mul(e1), e2)
		case pow:
			return compare(Pow(e1, Const(1)), e2)
		case exp:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case log:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		case sqrt:
			e1Arg := Operand(e1, 1)
			e2Arg := Operand(e2, 1)
			return compare(e1Arg, e2Arg)
		default:
			errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
			panic(errMsg)
		}
	default:
		errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(e1Typed))
		panic(errMsg)
	}
}

func orderRule1(e1, e2 constant) bool {return e1.Value < e2.Value}
func orderRule2(e1, e2 variable) bool {return e1.Name < e2.Name}
func orderRule2_1(e1, e2 constrainedVariable) bool {return e1.Name < e2.Name}
func orderRule3(e1, e2 add) bool {
	e1NumOp := NumberOfOperands(e1)
	e2NumOp := NumberOfOperands(e2)
	e1LastOp := Operand(e1, e1NumOp)
	e2LastOp := Operand(e2, e2NumOp)
	
	if !Equal(e1LastOp, e2LastOp) {
		return compare(e1LastOp, e2LastOp)
	}

	bnd := 0
	if e1NumOp < e2NumOp {
		bnd = e1NumOp
	} else {
		bnd = e2NumOp
	}

	for ix := 1; ix < bnd; ix++ {
		e1Op := Operand(e1, e1NumOp - ix)
		e2Op := Operand(e2, e2NumOp - ix)
		if !Equal(e1Op, e2Op) {
			return compare(e1Op, e2Op)
		}
	}
	return e1NumOp < e2NumOp
}
func orderRule3_1(e1, e2 mul) bool {
	e1NumOp := NumberOfOperands(e1)
	e2NumOp := NumberOfOperands(e2)
	e1LastOp := Operand(e1, e1NumOp)
	e2LastOp := Operand(e2, e2NumOp)
	
	if !Equal(e1LastOp, e2LastOp) {
		return compare(e1LastOp, e2LastOp)
	}

	bnd := 0
	if e1NumOp < e2NumOp {
		bnd = e1NumOp
	} else {
		bnd = e2NumOp
	}

	for ix := 1; ix < bnd; ix++ {
		e1Op := Operand(e1, e1NumOp - ix)
		e2Op := Operand(e2, e2NumOp - ix)
		if !Equal(e1Op, e2Op) {
			return compare(e1Op, e2Op)
		}
	}
	return e1NumOp < e2NumOp
}
func orderRule4(e1, e2 pow) bool {
	e1Base := Operand(e1, 1)
	e2Base := Operand(e2, 1)
	if !Equal(e1Base, e2Base) {
		return compare(e1Base, e2Base)
	} else {
		e1Exponent := Operand(e1, 2)
		e2Exponent := Operand(e2, 2)
		return compare(e1Exponent, e2Exponent)
	}
}
func orderRule5(e1, e2 Expr) bool {
	panic("rule dedicated to factorial which is not implemented")
}

/*
Returns a new expression where the
terms in in s1 has been prepended
to the terms in s2, i.e. the order 
of terms is not changed.
*/
func mergeSums(s1, s2 add) add {
	panic("not implemented yet")
}

/*
Returns a new expression where the
factors in in p1 has been prepended
to the factors in p2, i.e. the order 
of factors is not changed.
*/
func mergeProducts(p1, p2 mul) mul {
	panic("not implemented yet")
}

/*
Sorts the operands of expr in increasing order in accordance
with the order relation defined by compare(e1,e2 Expr). It does
not recursively sort the operands operands etc. This should only 
be applied when the operator is commutative!

To make a complete sort using this it needs to be recursive
called on the operands.

NOTE: Using insertion sort so worst case time complexity is O(n^2).
*/
func topOperandSort(expr Expr) Expr {
	for ix := 1; ix <= NumberOfOperands(expr)-1; ix++ {
		op1 := Operand(expr, ix)
		op2 := Operand(expr, ix+1)
		
		n := ix
		for !compare(op1, op2) {
			expr = swapOperands(expr, n, n+1)
			
			// As long as n > 1 we are not at the
			// first operand and there is no risk of 
			// index out of bounds error. If we are at 
			// the last operand we need to break since 
			// compare will continue to return false if 
			// op1 == op2.
			if n > 1 {
				op1 = Operand(expr, n-1)
				n--
			} else {
				break
			}
		}
	}
	return expr
}

func TopOperandSort(expr Expr) Expr { return topOperandSort(expr) }
