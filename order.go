package gosymbol

import (
	"fmt"
	"reflect"
)

/*
Sorts the operands of expr in increasing order in accordance
with the order relation defined by compare(e1,e2 Expr). It does
not recursively sort the operands operands etc. This should only
be applied when the operator is commutative!

To make a complete sort using this it needs to be recursive
called on the operands.

NOTE: Using insertion sort so worst case time complexity is O(n^2).
*/
func TopOperandSort(expr Expr) Expr {
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

func orderRule1(e1, e2 constant) bool              { return e1.Value < e2.Value }
func orderRule2(e1, e2 variable) bool              { return e1.Name < e2.Name }
func orderRule2_1(e1, e2 constrainedVariable) bool { return e1.Name < e2.Name }
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
		e1Op := Operand(e1, e1NumOp-ix)
		e2Op := Operand(e2, e2NumOp-ix)
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
		e1Op := Operand(e1, e1NumOp-ix)
		e2Op := Operand(e2, e2NumOp-ix)
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
