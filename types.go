package gosymbol

type VarName string
type Arguments map[VarName]float64
type Func func(Arguments) float64

type Expr interface {
	// Private functions
	equal(Expr) bool
	contains(Expr) bool
	substitute(Expr, Expr) Expr
	variableNames(*[]string)
	numberOfOperands() int
	operand(int) Expr
	simplify() Expr

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

// Simplification always occurs left to right.
// This means that the left expression aout to be the
// more complicated one, e.g. x^a * x^b = x^(a+b).
type simplificationRule struct {
	lhs Expr
	rhs Expr
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
	Exponent Expr
}

type sqrt struct {
	Expr
	Arg Expr
}
