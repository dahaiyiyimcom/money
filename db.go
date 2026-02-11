package money

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type DBAmount struct {
	A Amount
}

func (m *DBAmount) Scan(value any) error {
	if value == nil {
		m.A = 0
		return nil
	}

	switch v := value.(type) {
	case []byte:
		a, err := ParseString(string(v))
		if err != nil {
			return err
		}
		m.A = a
		return nil

	case string:
		a, err := ParseString(v)
		if err != nil {
			return err
		}
		m.A = a
		return nil

	case float64:
		// Force 2 decimals, then parse strictly.
		s := strconv.FormatFloat(v, 'f', 2, 64)
		a, err := ParseString(s)
		if err != nil {
			return err
		}
		m.A = a
		return nil

	case float32:
		s := strconv.FormatFloat(float64(v), 'f', 2, 64)
		a, err := ParseString(s)
		if err != nil {
			return err
		}
		m.A = a
		return nil

	case int64:
		// If driver returns integer, assume it's already minor? (ambiguous)
		// Better to treat as major units without decimals is dangerous.
		// If you WANT to support this, define it clearly. For safety, reject:
		return fmt.Errorf("money: unsupported scan int64=%d (ambiguous units)", v)

	default:
		return fmt.Errorf("money: unsupported scan type %T", value)
	}
}

func (m DBAmount) Value() (driver.Value, error) {
	return m.A.StringFixed2(), nil
}
