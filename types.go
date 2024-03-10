package gosymbol

type VarName string
type Arguments map[VarName]float64
type Func func(Arguments) float64

type Expr interface {
	// Public functions
	String() string
	Eval() Func

	// Private function
	d(VarName) Expr // Differentiation function
	compare(Expr) bool
	simplify() Expr
}

/**
Type used for storing expressions to enable advanced pattern matching.
Without constrain this is a different type from `variable` to remove problems of 
cacheing a variable as itself, and thus creating infinite recursion fuck ups.

If the Constraint field is set it ought to be checked before caching. This featrues make 
it possible to define patterns like 0^n = undefined if n < 0.
*/
type cacheVariable struct {
	Expr
	Name VarName
	Constraint func(expr Expr) bool
}

/** 
Used to define transformations from an expression
into another. The transformation can happen in two ways:
TODO: continue documentation and mention that a variable with
the same name as an constrained variable is considered to the same
variable by mathPattern.

This structure is, for example, used in the simplification
of expressions. 
*/
type transformationRule struct {
	// One of pattern and patternFunction must be defined.
	// pattern is prioritised, i.e. if pattern is matched 
	// then patternFunction will be ignored. To match on 
	// patternFunciton set pattern = nil.
	pattern Expr
	patternFunction func(Expr) bool

	// The mapping from pattern to whatever you define
	transform func(Expr) Expr 
}

type _variableCache struct {
	cache map[string]Expr
}

/**
Stuct to store the result from `applySimplificationRules()`.
**/
type transformationRuleApplicationResult struct {
	expr Expr
	exprAltered bool
}

/* Basic operators */

type undefined struct {
	Expr
}

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

type pow struct {
	Expr
	Base Expr
	Exponent Expr
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

type sqrt struct {
	Expr
	Arg Expr
}

/* Factories */

func Undefined() undefined {
	return undefined{}
}

func Const(val float64) constant {
	return constant{Value: val}
}

func CacheVar(name VarName) cacheVariable {
	return cacheVariable{Name: name}
}

func Var(name VarName) variable {
	return variable{Name: name}
}

func ConstrCacheVar(name VarName, constrFunc func(Expr) bool) cacheVariable {
	return cacheVariable{Name: name, Constraint: constrFunc}
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
