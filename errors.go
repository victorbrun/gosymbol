package gosymbol

type DuplicateArgumentError struct{}

func (e *DuplicateArgumentError) Error() string { return "multiple variables have the same name" }
