package money

type RoundingMode int

const (
	RoundHalfUp RoundingMode = iota // 0.5 yukarı
	RoundFloor
	RoundCeil
)

// MulRatio computes round(a * num / den).
func (a Amount) MulRatio(num, den int64, mode RoundingMode) Amount {
	if den == 0 {
		panic("denominator cannot be zero")
	}

	// NOTE: If you expect extremely large values, add an overflow-safe path with math/big.
	x := int64(a) * num
	q := x / den
	r := x % den
	if r == 0 {
		return Amount(q)
	}

	switch mode {
	case RoundFloor:
		if x < 0 {
			return Amount(q - 1)
		}
		return Amount(q)
	case RoundCeil:
		if x > 0 {
			return Amount(q + 1)
		}
		return Amount(q)
	case RoundHalfUp:
		absR := r
		if absR < 0 {
			absR = -absR
		}
		if absR*2 >= den {
			if x > 0 {
				return Amount(q + 1)
			}
			return Amount(q - 1)
		}
		return Amount(q)
	default:
		return Amount(q)
	}
}

// Percent returns round(a * percent / 100).
func (a Amount) Percent(percent int64, mode RoundingMode) Amount {
	return a.MulRatio(percent, 100, mode)
}
