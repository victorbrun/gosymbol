package gosymbol

type VarName string
type Arguments map[variable]Expr
type Func func(Arguments) Expr

type Expr interface {
	// Public functions
	String() string
	Eval() Func
	D(variable) Expr
	Simplify() Expr
	Approx() float64
}

// The Binding type is used in patternmatching.go
// to bind pattern variables to expressions
type Binding map[VarName]Expr

// A constrainedVariable is just like the data type
// variable but it has a constraint function. When trying
// to replace the variable with another expression, this
// function is run, testing whether the new expression satisfy
// the constrain, if yes it returns true and false otherwise.
// Note that you yourself need to define this function to fit your
// needs. TODO: replace this with some sort of logic DSL :)
type constrainedVariable struct {
	Expr
	name       VarName
	Constraint func(expr Expr) bool

	// Indicating if variable is part of a pattern
	isPattern bool
}

// Used to define transformations from an expression
// into another. The transformation can happen in two ways:
// TODO: continue documentation and mention that a variable with
// the same name as an constrained variable is considered to the same
// variable by mathPattern.
//
// This structure is, for example, used in the simplifation
// of expressions.
// This structure is, for example, used in the simplifation
// of expressions.
type transformationRule struct {
	// One of pattern and patternFunction must be defined.
	// pattern is prioritised, i.e. if pattern is matched
	// then patternFunction will be ignored. To match on
	pattern         Expr
	patternFunction func(Expr) bool

	// The mapping from pattern to whatever you define
	transform func(Expr) Expr
}

/* Basic operators */

type undefined struct {
	Expr
}

type variable struct {
	Expr
	name VarName
	// Indicating if variable is part of a pattern
	isPattern bool
	// used for approximating constants
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
	Base     Expr
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

/* Const types */

type integer struct {
	value int64
}

type rational interface {
	Expr
	constant
	numerator() integer
	denominator() integer
	simplifyRational() rational
}

type fraction struct {
	num integer
	den integer
}

type real struct {
	name  VarName
	value float64
}

type constant interface {
	Value() float64
}

func (n integer) Value() float64 { return float64(n.value) }
func (u fraction) Value() float64 {
	return float64(u.numerator().Value()) / float64(u.denominator().Value())
}
func (u real) Value() float64      { return u.value }
func (u undefined) Value() float64 { return u.approx() }

type symbol interface {
	Expr
	Name() VarName
}

func (x variable) Name() VarName            { return x.name }
func (x constrainedVariable) Name() VarName { return x.name }
func (u real) Name() VarName                { return u.name }
func (u undefined) Name() VarName           { return VarName("undefined") }

type function interface {
	Expr
	functionName() string
}

func (f exp) functionName() string       { return "exp" }
func (f log) functionName() string       { return "log" }
func (f sqrt) functionName() string      { return "sqrt" }
func (u undefined) functionName() string { return "undefined" }
