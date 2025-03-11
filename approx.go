package gosymbol

import "math"

func (n integer) Approx() float64 {
	return float64(n.value)
}

func (f fraction) Approx() float64 {
	if f.denominator() == Int(0) {
		if f.numerator().value > 0 {
			return math.Inf(1)
		}
		if f.numerator().value < 0 {
			return math.Inf(0)
		}
		return math.NaN()
	}
	return f.numerator().Approx() / f.denominator().Approx()
}

func (u undefined) Approx() float64 {
	return math.NaN()
}

func (x variable) Approx() float64 {
	if x.isConstant {
		return x.constValue
	}
	return math.NaN()
}

func (s add) Approx() float64 {
	var sum float64 = 0
	for _, op := range s.Operands {
		sum += op.Approx()
	}
	return sum
}

func (p mul) Approx() float64 {
	var prod float64 = 1
	for _, op := range p.Operands {
		prod *= op.Approx()
	}
	return prod
}

func (p pow) Approx() float64 {
	return math.Pow(p.Base.Approx(), p.Exponent.Approx())
}

func (e exp) Approx() float64 {
	return math.Exp(e.Arg.Approx())
}

func (l log) Approx() float64 {
	return math.Log(l.Arg.Approx())
}

func (s sqrt) Approx() float64 {
	return math.Sqrt(s.Arg.Approx())
}
