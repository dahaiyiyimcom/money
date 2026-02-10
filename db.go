package money

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// DBAmount is a DB-facing wrapper for DECIMAL(10,2) columns.
// Use it in structs scanned from sql rows.
type DBAmount struct {
	A Amount
}

func (m *DBAmount) Scan(value any) error {
	if value == nil {
		m.A = 0
		return nil
	}

	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("money: unsupported scan type %T", value)
	}

	a, err := ParseString(s)
	if err != nil {
		return err
	}
	m.A = a
	return nil
}

func (m DBAmount) Value() (driver.Value, error) {
	// write back as "12.34" for DECIMAL(10,2)
	return m.A.StringFixed2(), nil
}

// ParseString parses decimal string with optional sign, 2 decimals max.
// Accepts: "12.34", "12", "-0.50", "  12.30 "
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
	var fracStr string
	if len(parts) == 2 {
		fracStr = parts[1]
	}

	// normalize fraction to 2 digits
	switch len(fracStr) {
	case 0:
		fracStr = "00"
	case 1:
		fracStr = fracStr + "0"
	case 2:
		// ok
	default:
		// If more than 2 decimals exist, you can either reject or round.
		// For strict DECIMAL(10,2), reject is safest.
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

	minor := int64(whole)*100 + int64(frac)
	return Amount(sign * minor), nil
}

func parseUint(s string) (uint64, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	var n uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("non-digit %q", c)
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}
