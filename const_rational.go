package gosymbol

import "fmt"

var EmptyFraction = fraction{}

func (u fraction) String() string {
	return fmt.Sprintf("%s/%s", u.num.String(), u.den.String())
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

func Inv(u rational) rational {
	return Frac(u.denominator(), u.numerator()).simplifyRational()
}

func (u fraction) approx() float64 {
	return u.num.approx() / u.den.approx()
}

func (u fraction) simplifyRational() rational {
	if u == EmptyFraction {
		return EmptyFraction
	}
	r, err := intMod(u.num, u.den)
	if err != nil {
		return EmptyFraction
	}
	if r == Int(0) {
		q, err := intQuotient(u.num, u.den)
		if err != nil {
			return EmptyFraction
		}
		return q
	}
	gcd, err := gcd(u.num, u.den)
	if err != nil {
		return EmptyFraction
	}
	num, err1 := intQuotient(u.num, gcd)
	den, err2 := intQuotient(u.den, gcd)
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
	return ratMul(u, Inv(w))
}

func ratPow(u rational, n integer) rational {
	u = u.simplifyRational()
	if n.value < 0 {
		u = Inv(u)
		n = intMinus(n)
	}
	switch v := u.(type) {
	case integer:
		pow, err := intPow(v, n)
		if err != nil {
			return EmptyFraction
		}
		if n.value < 0 {
			return Inv(pow)
		}
		return pow
	case fraction:
		num, err1 := intPow(v.num, n)
		den, err2 := intPow(v.den, n)
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
		return intMinus(v)
	case fraction:
		return Frac(intMinus(u.numerator()), u.denominator())
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
