package gosymbol

import (
	"errors"
	"fmt"
)

func (u integer) String() string {
	return fmt.Sprintf("%d", u.value)
}
func (u integer) numerator() integer {
	return u
}
func (u integer) denominator() integer {
	return Int(1)
}

func (u integer) simplifyRational() rational {
	return u
}

func (n integer) approx() float64 {
	return float64(n.value)
}

func intQuotient(a integer, b integer) (integer, error) {
	if b == Int(0) {
		return integer{}, errors.New("Division by 0")
	}
	return Int(a.value / b.value), nil
}

func intMul(a integer, b integer) integer {
	return Int(a.value * b.value)
}

func intAdd(a integer, b integer) integer {
	return Int(a.value + b.value)
}

func intMod(a integer, b integer) (integer, error) {
	if b == Int(0) {
		return integer{}, errors.New("Division by 0")
	}
	return Int(a.value % b.value), nil
}

func intMinus(a integer) integer {
	return Int(-a.value)
}

func intAbs(a integer) integer {
	if a.value < 0 {
		return intMinus(a)
	}
	return a
}

func intPow(a integer, b integer) (integer, error) {
	if a == Int(0) && b == Int(0) {
		return integer{}, errors.New("Undefined 0^0")
	}
	if b.value < 0 {
		return integer{}, errors.New("exponent in intPow must be non negative")
	}
	if a == Int(0) {
		return Int(0), nil
	}
	if b == Int(0) {
		return Int(1), nil
	}
	partialPow, err := intPow(a, intSubtract(b, Int(1)))
	if err != nil {
		return integer{}, err
	}
	return intMul(partialPow, a), nil
}

func intSubtract(a integer, b integer) integer {
	return intAdd(a, intMinus(b))
}

func gcd(a integer, b integer) (integer, error) {
	if a.value <= 0 || b.value <= 0 {
		return integer{}, errors.New("GCD only accepts positive values")
	}
	for b != Int(0) {
		r, err := intMod(a, b)
		if err != nil {
			return integer{}, nil
		}
		a, b = b, r
	}
	return intAbs(a), nil
}
