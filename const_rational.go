package gosymbol

import (
	"fmt"
	"math"
)

var EmptyFraction = fraction{}

func (u fraction) String() string {
	return fmt.Sprintf("%s/%s", u.numerator().String(), u.denominator().String())
}

func (u fraction) numerator() integer {
	if u == EmptyFraction {
		return Int(0)
	}
	return u.num
}

func (u fraction) denominator() integer {
	if u == EmptyFraction {
		return Int(0)
	}
	return u.den
}

func ratInv(u rational) rational {
	return Frac(u.denominator(), u.numerator()).simplifyRational()
}

func (u fraction) approx() float64 {
	if u.denominator() == Int(0) {
		if u.numerator().value > 0 {
			return math.Inf(1)
		}
		if u.numerator().value < 0 {
			return math.Inf(0)
		}
		return math.NaN()
	}
	return u.numerator().approx() / u.denominator().approx()
}

func (u fraction) simplifyRational() rational {
	if u == EmptyFraction {
		return EmptyFraction
	}
	r, err := intMod(u.numerator(), u.denominator())
	if err != nil {
		return EmptyFraction
	}
	if r == Int(0) {
		q, err := intQuotient(u.numerator(), u.den)
		if err != nil {
			return EmptyFraction
		}
		return q
	}
	gcd, err := gcd(u.numerator(), u.denominator())
	if err != nil {
		return EmptyFraction
	}
	num, err1 := intQuotient(u.numerator(), gcd)
	den, err2 := intQuotient(u.denominator(), gcd)
	if err1 != nil || err2 != nil {
		return EmptyFraction
	}
	return Frac(num, den)
}

func ratMul(u rational, w rational) rational {
	if w.numerator() == Int(0) {
		return EmptyFraction
	}
	return Frac(
		intMul(u.numerator(), w.numerator()),
		intMul(u.denominator(), w.denominator()),
	).simplifyRational()
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
	case integer:
		pow, err := intPow(v, n)
		if err != nil {
			return EmptyFraction
		}
		if n.value < 0 {
			return ratInv(pow)
		}
		return pow
	case fraction:
		num, err1 := intPow(v.numerator(), n)
		den, err2 := intPow(v.denominator(), n)
		if err1 != nil || err2 != nil {
			return EmptyFraction
		}
		if n.value < 0 {
			return Frac(den, num)
		}
		return Frac(num, den)
	}
	return nil
}

func ratAdd(u rational, w rational) rational {
	return Frac(
		intAdd(
			intMul(u.numerator(), w.denominator()),
			intMul(u.denominator(), w.numerator()),
		),
		intMul(u.denominator(), w.denominator()),
	).simplifyRational()
}

func ratMinus(u rational) rational {
	switch v := u.(type) {
	case integer:
		return intNeg(v)
	case fraction:
		return Frac(intNeg(u.numerator()), u.denominator())
	}
	return nil
}

func ratSubtract(u rational, w rational) rational {
	return ratAdd(u, ratMinus(w))
}

func ratAbs(u rational) rational {
	switch v := u.(type) {
	case integer:
		return intAbs(v)
	case fraction:
		return Frac(intAbs(u.numerator()), u.denominator())
	}
	return nil
}
