package money

import "fmt"

// Amount represents money in minor units (e.g., kuruş).
// 12.34 TL => 1234
type Amount int64

func NewMinor(minor int64) Amount { return Amount(minor) }
func (a Amount) Minor() int64     { return int64(a) }

func (a Amount) Add(b Amount) Amount { return a + b }
func (a Amount) Sub(b Amount) Amount { return a - b }

// MulQty multiplies by quantity (e.g., unit price * qty).
func (a Amount) MulQty(qty int64) Amount { return Amount(int64(a) * qty) }

func (a Amount) IsNegative() bool { return a < 0 }

// StringFixed2 formats as "12.34" (always 2 decimals).
func (a Amount) StringFixed2() string {
	min := int64(a)
	sign := ""
	if min < 0 {
		sign = "-"
		min = -min
	}
	whole := min / 100
	frac := min % 100
	return fmt.Sprintf("%s%d.%02d", sign, whole, frac)
}
