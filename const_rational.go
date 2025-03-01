package gosymbol

import (
	"fmt"
	"math"
)

func (u fraction) String() string {
	return fmt.Sprintf("%s/%s", u.numerator().String(), u.denominator().String())
}

func (u fraction) numerator() integer {
	return u.num
}

func (u fraction) denominator() integer {
	return u.den
}

func ratInv(u rational) rational {
	return Div(u.denominator(), u.numerator()).(rational).simplifyRational()
}

func (u fraction) simplifyRational() rational {
	r, err := intMod(u.numerator(), u.denominator())
	if err != nil {
		return Undefined()
	}
	if r == Int(0) {
		q, err := intQuotient(u.numerator(), u.den)
		if err != nil {
			return Undefined()
		}
		return q
	}
	gcd, err := gcd(u.numerator(), u.denominator())
	if err != nil {
		return Undefined()
	}
	num, err1 := intQuotient(u.numerator(), gcd)
	den, err2 := intQuotient(u.denominator(), gcd)
	if err1 != nil || err2 != nil {
		return Undefined()
	}
	return Div(num, den).(rational)
}

func ratMul(u rational, w rational) rational {
	if w.numerator() == Int(0) {
		return Undefined()
	}
	return Div(
		intMul(u.numerator(), w.numerator()),
		intMul(u.denominator(), w.denominator()),
	).(rational).simplifyRational()
}

func ratDiv(u rational, w rational) rational {
	return ratMul(u, ratInv(w))
}

func ratPow(u rational, n integer) rational {
	u = u.simplifyRational()
	if n.value < 0 {
		u = ratInv(u)
		n = intNeg(n)
	}
	switch v := u.(type) {
	case undefined:
		return Undefined()
	case integer:
		pow, err := intPow(v, n)
		if err != nil {
			return Undefined()
		}
		if n.value < 0 {
			return ratInv(pow)
		}
		return pow
	case fraction:
		num, err1 := intPow(v.numerator(), n)
		den, err2 := intPow(v.denominator(), n)
		if err1 != nil || err2 != nil {
			return Undefined()
		}
		if n.value < 0 {
			return Div(den, num).(rational)
		}
		return Div(num, den).(rational)
	}
	return nil
}

func ratAdd(u rational, w rational) rational {
	return Div(
		intAdd(
			intMul(u.numerator(), w.denominator()),
			intMul(u.denominator(), w.numerator()),
		),
		intMul(u.denominator(), w.denominator()),
	).(rational).simplifyRational()
}

func ratMinus(u rational) rational {
	switch v := u.(type) {
	case undefined:
		return Undefined()
	case integer:
		return intNeg(v)
	case fraction:
		return Div(intNeg(u.numerator()), u.denominator()).(rational)
	}
	return nil
}

func ratSubtract(u rational, w rational) rational {
	return ratAdd(u, ratMinus(w))
}

func ratAbs(u rational) rational {
	switch v := u.(type) {
	case undefined:
		return Undefined()
	case integer:
		return intAbs(v)
	case fraction:
		return Div(intAbs(u.numerator()), u.denominator()).(rational)
	}
	return nil
}

func (u undefined) approx() float64 {
	return math.NaN()
}

func (u undefined) numerator() integer {
	return Int(0)
}

func (u undefined) denominator() integer {
	return Int(0)
}

func (u undefined) simplifyRational() rational {
	return u
}
