package gosymbol

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestNumerator(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput integer
	}{
		{
			name:           "numerator of fraction n/m is n",
			input:          Frac(Int(1), Int(2)),
			expectedOutput: Int(1),
		},
		{
			name:           "numerator of integer n is n",
			input:          Int(1),
			expectedOutput: Int(1),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.numerator()

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestDenominator(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput integer
	}{
		{
			name:           "denominator of fraction n/m is m",
			input:          Frac(Int(1), Int(2)),
			expectedOutput: Int(2),
		},
		{
			name:           "denominator of integer n is 1",
			input:          Int(2),
			expectedOutput: Int(1),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.denominator()

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestApprox(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput float64
	}{
		{
			name:           "approx of fraction 1/2",
			input:          Frac(Int(1), Int(2)),
			expectedOutput: 0.5,
		},
		{
			name:           "approx of fraction 1/0 (infinity)",
			input:          Frac(Int(1), Int(0)),
			expectedOutput: math.Inf(1),
		},
		{
			name:           "approx of empty fraction NaN",
			input:          EmptyFraction,
			expectedOutput: math.NaN(),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.approx()

			if !math.IsNaN(result) && result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestSimplifyRational(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput rational
	}{
		{
			name:           "simplify fraction 4/2",
			input:          Frac(Int(-4), Int(2)),
			expectedOutput: Int(-2),
		},
		{
			name:           "simplify fraction 0/1",
			input:          Frac(Int(0), Int(1)),
			expectedOutput: Int(0),
		},
		{
			name:           "simplify empty fraction",
			input:          EmptyFraction,
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := test.input.simplifyRational()

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestRatMul(t *testing.T) {
	tests := []struct {
		name           string
		input1         rational
		input2         rational
		expectedOutput rational
	}{
		{
			name:           "multiply fractions 1/2 * 3/4",
			input1:         Frac(Int(1), Int(2)),
			input2:         Frac(Int(3), Int(4)),
			expectedOutput: Frac(Int(3), Int(8)),
		},
		{
			name:           "multiply integer 2 * 3",
			input1:         Int(2),
			input2:         Int(3),
			expectedOutput: Int(6),
		},
		{
			name:           "multiply by empty fraction",
			input1:         Frac(Int(1), Int(2)),
			input2:         EmptyFraction,
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratMul(test.input1, test.input2)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational1: %v, rational2: %v\nExpected: %v\nGot: %v", test.name, test.input1, test.input2, test.expectedOutput, result)
			}
		})
	}
}

func TestRatDiv(t *testing.T) {
	tests := []struct {
		name           string
		input1         rational
		input2         rational
		expectedOutput rational
	}{
		{
			name:           "divide fractions 1/2 ÷ 3/4",
			input1:         Frac(Int(1), Int(2)),
			input2:         Frac(Int(3), Int(4)),
			expectedOutput: Frac(Int(2), Int(3)),
		},
		{
			name:           "divide integer 6 ÷ 3",
			input1:         Int(6),
			input2:         Int(3),
			expectedOutput: Int(2),
		},
		{
			name:           "divide by empty fraction",
			input1:         Frac(Int(1), Int(2)),
			input2:         EmptyFraction,
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratDiv(test.input1, test.input2)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational1: %v, rational2: %v\nExpected: %v\nGot: %v", test.name, test.input1, test.input2, test.expectedOutput, result)
			}
		})
	}
}

func TestRatPow(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		exponent       integer
		expectedOutput rational
	}{
		{
			name:           "power of fraction (1/2)^2",
			input:          Frac(Int(1), Int(2)),
			exponent:       Int(2),
			expectedOutput: Frac(Int(1), Int(4)),
		},
		{
			name:           "power of integer 2^3",
			input:          Int(2),
			exponent:       Int(3),
			expectedOutput: Int(8),
		},
		{
			name:           "power of fraction (3/2)^-1",
			input:          Frac(Int(3), Int(2)),
			exponent:       Int(-1),
			expectedOutput: Frac(Int(2), Int(3)),
		},
		{
			name:           "power of empty fraction",
			input:          EmptyFraction,
			exponent:       Int(2),
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratPow(test.input, test.exponent)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v, Exponent: %v\nExpected: %v\nGot: %v", test.name, test.input, test.exponent, test.expectedOutput, result)
			}
		})
	}
}

func TestRatAdd(t *testing.T) {
	tests := []struct {
		name           string
		input1         rational
		input2         rational
		expectedOutput rational
	}{
		{
			name:           "add fractions 1/2 + 3/4",
			input1:         Frac(Int(1), Int(2)),
			input2:         Frac(Int(3), Int(4)),
			expectedOutput: Frac(Int(5), Int(4)),
		},
		{
			name:           "add integer 1 + 2",
			input1:         Int(1),
			input2:         Int(2),
			expectedOutput: Int(3),
		},
		{
			name:           "add integer 1 + 0",
			input1:         Int(1),
			input2:         Int(0),
			expectedOutput: Int(1),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratAdd(test.input1, test.input2)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational1: %v, rational2: %v\nExpected: %v\nGot: %v", test.name, test.input1, test.input2, test.expectedOutput, result)
			}
		})
	}
}

func TestRatMinus(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput rational
	}{
		{
			name:           "minus of integer 2",
			input:          Int(2),
			expectedOutput: Int(-2),
		},
		{
			name:           "minus of fraction 1/2",
			input:          Frac(Int(1), Int(2)),
			expectedOutput: Frac(Int(-1), Int(2)),
		},
		{
			name:           "minus of empty fraction",
			input:          EmptyFraction,
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratMinus(test.input)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestRatSubtract(t *testing.T) {
	tests := []struct {
		name           string
		input1         rational
		input2         rational
		expectedOutput rational
	}{
		{
			name:           "subtract fractions 3/4 - 1/2",
			input1:         Frac(Int(3), Int(4)),
			input2:         Frac(Int(1), Int(2)),
			expectedOutput: Frac(Int(1), Int(4)),
		},
		{
			name:           "subtract integer 5 - 2",
			input1:         Int(5),
			input2:         Int(2),
			expectedOutput: Int(3),
		},
		{
			name:           "subtract from empty fraction",
			input1:         EmptyFraction,
			input2:         Frac(Int(1), Int(2)),
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratSubtract(test.input1, test.input2)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational1: %v, rational2: %v\nExpected: %v\nGot: %v", test.name, test.input1, test.input2, test.expectedOutput, result)
			}
		})
	}
}

func TestRatAbs(t *testing.T) {
	tests := []struct {
		name           string
		input          rational
		expectedOutput rational
	}{
		{
			name:           "absolute value of integer -2",
			input:          Int(-2),
			expectedOutput: Int(2),
		},
		{
			name:           "absolute value of fraction -1/2",
			input:          Frac(Int(-1), Int(2)),
			expectedOutput: Frac(Int(1), Int(2)),
		},
		{
			name:           "absolute value of empty fraction",
			input:          EmptyFraction,
			expectedOutput: EmptyFraction,
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := ratAbs(test.input)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v", test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestIntMinus(t *testing.T) {
	tests := []struct {
		name           string
		input          integer
		expectedOutput integer
	}{
		{
			name:           "minus of 5",
			input:          Int(5),
			expectedOutput: Int(-5),
		},
		{
			name:           "minus of 0",
			input:          Int(0),
			expectedOutput: Int(0),
		},
		{
			name:           "minus of -3",
			input:          Int(-3),
			expectedOutput: Int(3),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := intNeg(test.input)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v",
					test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestIntAbs(t *testing.T) {
	tests := []struct {
		name           string
		input          integer
		expectedOutput integer
	}{
		{
			name:           "absolute value of 5",
			input:          Int(5),
			expectedOutput: Int(5),
		},
		{
			name:           "absolute value of -3",
			input:          Int(-3),
			expectedOutput: Int(3),
		},
		{
			name:           "absolute value of 0",
			input:          Int(0),
			expectedOutput: Int(0),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result := intAbs(test.input)

			if result != test.expectedOutput {
				t.Errorf("Following test failed: %s\nInput rational: %v\nExpected: %v\nGot: %v",
					test.name, test.input, test.expectedOutput, result)
			}
		})
	}
}

func TestIntPow(t *testing.T) {
	tests := []struct {
		name           string
		base           integer
		exponent       integer
		expectedOutput integer
		expectedErr    error
	}{
		{
			name:           "power of 2^3",
			base:           Int(2),
			exponent:       Int(3),
			expectedOutput: Int(8),
			expectedErr:    nil,
		},
		{
			name:           "power of 0^0 (undefined)",
			base:           Int(0),
			exponent:       Int(0),
			expectedOutput: integer{},
			expectedErr:    errors.New("Undefined 0^0"),
		},
		{
			name:           "power of negative base -2^3",
			base:           Int(-2),
			exponent:       Int(3),
			expectedOutput: Int(-8),
			expectedErr:    nil,
		},
		{
			name:           "power of 2^-1",
			base:           Int(2),
			exponent:       Int(-1),
			expectedOutput: integer{}, // Since negative exponents are not handled here, assuming an error is returned
			expectedErr:    errors.New("exponent in intPow must be non negative"),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result, err := intPow(test.base, test.exponent)

			if result != test.expectedOutput || (err != nil && err.Error() != test.expectedErr.Error()) {
				t.Errorf("Following test failed: %s\nInput base: %v, Exponent: %v\nExpected: %v, Error: %v\nGot: %v, Error: %v",
					test.name, test.base, test.exponent, test.expectedOutput, test.expectedErr, result, err)
			}
		})
	}
}

func TestGcd(t *testing.T) {
	tests := []struct {
		name           string
		input1         integer
		input2         integer
		expectedOutput integer
		expectedErr    error
	}{
		{
			name:           "gcd of 12 and 15",
			input1:         Int(12),
			input2:         Int(15),
			expectedOutput: Int(3),
			expectedErr:    nil,
		},
		{
			name:           "gcd of 17 and 19",
			input1:         Int(17),
			input2:         Int(19),
			expectedOutput: Int(1),
			expectedErr:    nil,
		},
		{
			name:           "gcd of 0 and 5",
			input1:         Int(0),
			input2:         Int(5),
			expectedOutput: integer{},
			expectedErr:    errors.New("GCD only accepts positive values"),
		},
		{
			name:           "gcd of 5 and 0",
			input1:         Int(5),
			input2:         Int(0),
			expectedOutput: integer{},
			expectedErr:    errors.New("GCD only accepts positive values"),
		},
	}

	for ix, test := range tests {
		t.Run(fmt.Sprint(ix+1), func(t *testing.T) {
			result, err := gcd(test.input1, test.input2)

			if result != test.expectedOutput || (err != nil && err.Error() != test.expectedErr.Error()) {
				t.Errorf("Following test failed: %s\nInput a: %v, b: %v\nExpected: %v, Error: %v\nGot: %v, Error: %v",
					test.name, test.input1, test.input2, test.expectedOutput, test.expectedErr, result, err)
			}
		})
	}
}
