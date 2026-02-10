package money_test

import (
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func TestMulRatio_Rounding(t *testing.T) {
	a := money.NewMinor(100) // 1.00

	// 1.00 * 1/3 = 0.333... => 33 (floor), 34 (ceil), 33 (half-up? 0.33) remainder not >=0.5
	if got := a.MulRatio(1, 3, money.RoundFloor).Minor(); got != 33 {
		t.Fatalf("floor got=%d want=33", got)
	}
	if got := a.MulRatio(1, 3, money.RoundCeil).Minor(); got != 34 {
		t.Fatalf("ceil got=%d want=34", got)
	}
	if got := a.MulRatio(1, 3, money.RoundHalfUp).Minor(); got != 33 {
		t.Fatalf("half-up got=%d want=33", got)
	}

	// 1.00 * 1/2 = 0.50 => 50 half-up => 50 (tam)
	if got := a.MulRatio(1, 2, money.RoundHalfUp).Minor(); got != 50 {
		t.Fatalf("half-up got=%d want=50", got)
	}

	// 1.00 * 1/6 = 0.1666 => 17 half-up? 16.66 -> 17? hayır, 0.16... => 17 only if >=0.5
	// 100/6 = 16 r4 (16.66) => half-up => 17 (çünkü r*2=8 >=6)
	if got := a.MulRatio(1, 6, money.RoundHalfUp).Minor(); got != 17 {
		t.Fatalf("half-up got=%d want=17", got)
	}
}

func TestPercent(t *testing.T) {
	a := money.NewMinor(1234) // 12.34
	vat := a.Percent(18, money.RoundHalfUp)
	// 1234*18/100=222.12 => 222
	if got := vat.Minor(); got != 222 {
		t.Fatalf("got=%d want=222", got)
	}
}
