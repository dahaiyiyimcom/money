package money

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// ParseString parses a decimal money string with max 2 fractional digits into minor units.
// Examples:
// "12" -> 1200
// "12.3" -> 1230
// "12.34" -> 1234
// "-0.50" -> -50
func ParseString(s string) (Amount, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("money: empty string")
	}

	sign := int64(1)
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("money: invalid format: %q", s)
	}

	wholeStr := parts[0]
	if wholeStr == "" {
		wholeStr = "0"
	}

	fracStr := ""
	if len(parts) == 2 {
		fracStr = parts[1]
	}

	switch len(fracStr) {
	case 0:
		fracStr = "00"
	case 1:
		fracStr = fracStr + "0"
	case 2:
		// ok
	default:
		return 0, fmt.Errorf("money: too many decimal places: %q", s)
	}

	whole, err := parseUint(wholeStr)
	if err != nil {
		return 0, fmt.Errorf("money: invalid whole part: %w", err)
	}
	frac, err := parseUint(fracStr)
	if err != nil {
		return 0, fmt.Errorf("money: invalid fractional part: %w", err)
	}
	if frac > 99 {
		return 0, fmt.Errorf("money: invalid fractional part: %q", fracStr)
	}

	// ---- overflow guard (critical) ----
	// We need: whole*100 + frac <= MaxInt64 for positive, and <= MaxInt64 for negative magnitude too.
	// Since sign is applied at the end, just check absolute magnitude fits into int64.
	if whole > uint64(math.MaxInt64/100) {
		return 0, fmt.Errorf("money: overflow: %q", s)
	}
	minorBase := int64(whole) * 100
	// minorBase is safe now, but minorBase+frac might still overflow if whole == MaxInt64/100 and frac pushes it
	if minorBase > math.MaxInt64-int64(frac) {
		return 0, fmt.Errorf("money: overflow: %q", s)
	}
	minor := minorBase + int64(frac)
	// -----------------------------------

	return Amount(sign * minor), nil
}

func parseUint(s string) (uint64, error) {
	if s == "" {
		return 0, errors.New("money: empty number")
	}
	var n uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("money: non-digit %q", c)
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}
