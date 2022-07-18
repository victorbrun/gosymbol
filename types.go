package gosymbol

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

type simplificationRule struct {
	// Simplification always occurs left to right.
	// This means that the left expression aout to be the
	// more complicated one, e.g. x^a * x^b = x^(a+b).
	lhs Expr
	rhs Expr
}
