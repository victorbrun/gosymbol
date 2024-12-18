package gosymbol

type VarName string
type Arguments map[VarName]float64
type Func func(Arguments) float64

type Expr interface {
	// Public functions
	String() string
	Eval() Func
	D(VarName) Expr
}

// A constrainedVariable is just like the data type
// variable but it has a constraint function. When trying
// to replace the variable with another expression, this
// function is run, testing whether the new expression satisfy
// the constrain, if yes it returns true and false otherwise.
// Note that you yourself need to define this function to fit your
// needs. TODO: replace this with some sort of logic DSL :)
type constrainedVariable struct {
	Expr
	Name VarName
	Constraint func(expr Expr) bool
}

// Used to define transformations from an expression
// into another. The transformation can happen in two ways:
// TODO: continue documentation and mention that a variable with
// the same name as an constrained variable is considered to the same
// variable by mathPattern.
//
//This structure is, for example, used in the simplifation
// of expressions. 
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
