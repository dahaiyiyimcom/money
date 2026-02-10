package money_test

import (
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func TestAmount_StringFixed2(t *testing.T) {
	cases := []struct {
		minor int64
		want  string
	}{
		{0, "0.00"},
		{1, "0.01"},
		{10, "0.10"},
		{99, "0.99"},
		{100, "1.00"},
		{1234, "12.34"},
		{-1, "-0.01"},
		{-1234, "-12.34"},
	}

	for _, tc := range cases {
		a := money.NewMinor(tc.minor)
		if got := a.StringFixed2(); got != tc.want {
			t.Fatalf("minor=%d got=%q want=%q", tc.minor, got, tc.want)
		}
	}
}

func TestAmount_Arithmetic(t *testing.T) {
	a := money.NewMinor(1234) // 12.34
	b := money.NewMinor(66)   // 0.66

	if got := a.Add(b).Minor(); got != 1300 {
		t.Fatalf("Add got=%d want=1300", got)
	}
	if got := a.Sub(b).Minor(); got != 1168 {
		t.Fatalf("Sub got=%d want=1168", got)
	}
	if got := a.MulQty(3).Minor(); got != 3702 {
		t.Fatalf("MulQty got=%d want=3702", got)
	}
}
