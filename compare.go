package gosymbol


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

func (e1 undefined) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return true	
	}
	constAction := func(e2 constant) bool {
		return true
	}
	varAction := func(e2 variable) bool {
		return true
	}
	addAction := func(e2 add) bool {
		return true
	}
	mulAction := func(e2 mul) bool {
		return true
	}
	powAction := func(e2 pow) bool {
		return true
	}
	expAction := func(e2 exp) bool {
		return true
	}
	logAction := func(e2 log) bool {
		return true
	}
	sqrtAction := func(e2 sqrt) bool {
		return true
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

func (e1 constant) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return e1.Value < e2.Value
	}
	varAction := func(e2 variable) bool {
		return false
	}
	addAction := func(e2 add) bool {
		return false
	}
	mulAction := func(e2 mul) bool {
		return false
	}
	powAction := func(e2 pow) bool {
		return false
	}
	expAction := func(e2 exp) bool {
		return false
	}
	logAction := func(e2 log) bool {
		return false
	}
	sqrtAction := func(e2 sqrt) bool {
		return false
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}


func (e1 variable) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return orderRule2(e1, e2)
	}
	addAction := func(e2 add) bool {
		return Add(e1).compare(e2)
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return Pow(e1, Const(1)).compare(e2)
	}
	expAction := func(e2 exp) bool {
		return Exp(e1).compare(e2)
	}
	logAction := func(e2 log) bool {
		return Log(e1).compare(e2)
	}
	sqrtAction := func(e2 sqrt) bool {
		return Sqrt(e1).compare(e2)
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

func (e1 add) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Add(e2))
	}
	addAction := func(e2 add) bool {
		return orderRule3(e1, e2)
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return Pow(e1, Const(1)).compare(e2)
	}
	expAction := func(e2 exp) bool {
		return e1.compare(Add(e2))
	}
	logAction := func(e2 log) bool {
		return e1.compare(Add(e2))
	}
	sqrtAction := func(e2 sqrt) bool {
		return e1.compare(Add(e2))
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

func (e1 mul) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Mul(e2))
	}
	addAction := func(e2 add) bool {
		return e1.compare(Mul(e2))
	}
	mulAction := func(e2 mul) bool {
		return orderRule3_1(e1, e2)
	}
	powAction := func(e2 pow) bool {
		return e1.compare(Mul(e2))
	}
	expAction := func(e2 exp) bool {
		return e1.compare(Mul(e2))
	}
	logAction := func(e2 log) bool {
		return e1.compare(Mul(e2))
	}
	sqrtAction := func(e2 sqrt) bool {
		return e1.compare(Mul(e2))
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}


func (e1 pow) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Pow(e2, Const(1)))
	}
	addAction := func(e2 add) bool {
		return e1.compare(Pow(e2, Const(1)))
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return orderRule4(e1, e2)
	}
	expAction := func(e2 exp) bool {
		return e1.compare(Pow(e2, Const(1)))
	}
	logAction := func(e2 log) bool {
		return e1.compare(Pow(e2, Const(1)))
	}
	sqrtAction := func(e2 sqrt) bool {
		return e1.compare(Pow(e2, Const(1)))
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

func (e1 exp) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Exp(e2))
	}
	addAction := func(e2 add) bool {
		return Add(e1).compare(e2)
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return Pow(e1, Const(1)).compare(e2)
	}
	expAction := func(e2 exp) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}
	logAction := func(e2 log) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}
	sqrtAction := func(e2 sqrt) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

func (e1 log) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Log(e2))
	}
	addAction := func(e2 add) bool {
		return Add(e1).compare(e2)
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return Pow(e1, Const(1)).compare(e2)
	}
	expAction := func(e2 exp) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}
	logAction := func(e2 log) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}
	sqrtAction := func(e2 sqrt) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}


func (e1 sqrt) compare(e2 Expr) bool {
	undefAction := func(e2 undefined) bool {
		return false	
	}
	constAction := func(e2 constant) bool {
		return false
	}
	varAction := func(e2 variable) bool {
		return e1.compare(Sqrt(e2))
	}
	addAction := func(e2 add) bool {
		return Add(e1).compare(e2)
	}
	mulAction := func(e2 mul) bool {
		return Mul(e1).compare(e2)
	}
	powAction := func(e2 pow) bool {
		return Pow(e1, Const(1)).compare(e2)
	}
	expAction := func(e2 exp) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 2)
		return e1Arg.compare(e2Arg)
	}
	logAction := func(e2 log) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}
	sqrtAction := func(e2 sqrt) bool {
		e1Arg := Operand(e1, 1)
		e2Arg := Operand(e2, 1)
		return e1Arg.compare(e2Arg)
	}

	return MatchTransform[bool](
		e2, 
		undefAction,
		constAction,
		varAction,
		addAction,
		mulAction,
		powAction,
		expAction,
		logAction,
		sqrtAction,
	)
}

/* Helper functions */

func orderRule1(e1, e2 constant) bool {return e1.Value < e2.Value}

func orderRule2(e1, e2 variable) bool {return e1.Name < e2.Name}

func orderRule3(e1, e2 add) bool {
	e1NumOp := NumberOfOperands(e1)
	e2NumOp := NumberOfOperands(e2)
	e1LastOp := Operand(e1, e1NumOp)
	e2LastOp := Operand(e2, e2NumOp)
	
	if !Equal(e1LastOp, e2LastOp) {
		return  e1LastOp.compare(e2LastOp)
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
			return e1Op.compare(e2Op)
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
		return e1LastOp.compare(e2LastOp)
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
			return e1Op.compare(e2Op)
		}
	}
	return e1NumOp < e2NumOp
}

func orderRule4(e1, e2 pow) bool {
	e1Base := Operand(e1, 1)
	e2Base := Operand(e2, 1)
	if !Equal(e1Base, e2Base) {
		return e1Base.compare(e2Base)
	} else {
		e1Exponent := Operand(e1, 2)
		e2Exponent := Operand(e2, 2)
		return e1Exponent.compare(e2Exponent)
	}
}

func orderRule5(e1, e2 Expr) bool {
	panic("rule dedicated to factorial which is not implemented")
}
