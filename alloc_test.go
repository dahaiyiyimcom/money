package money_test

import (
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func sum(xs []money.Amount) int64 {
	var s int64
	for _, x := range xs {
		s += x.Minor()
	}
	return s
}

func TestAllocateProportional_SumExact(t *testing.T) {
	// Line totals: 10.00, 20.00, 30.00 => total 60.00
	bases := []money.Amount{
		money.NewMinor(1000),
		money.NewMinor(2000),
		money.NewMinor(3000),
	}
	discount := money.NewMinor(100) // 1.00

	shares := money.AllocateProportional(bases, discount)

	if got := sum(shares); got != discount.Minor() {
		t.Fatalf("sum(shares) got=%d want=%d", got, discount.Minor())
	}
	// Oransal beklenen ~ 0.17, 0.33, 0.50 => 17,33,50 (tam)
	if shares[0].Minor() != 17 || shares[1].Minor() != 33 || shares[2].Minor() != 50 {
		t.Fatalf("unexpected shares: %d,%d,%d", shares[0].Minor(), shares[1].Minor(), shares[2].Minor())
	}
}

func TestAllocateProportional_LeftoverDistribution(t *testing.T) {
	// 3 eşit line: 1.00, 1.00, 1.00; discount 0.01
	bases := []money.Amount{
		money.NewMinor(100),
		money.NewMinor(100),
		money.NewMinor(100),
	}
	discount := money.NewMinor(1)

	shares := money.AllocateProportional(bases, discount)
	if got := sum(shares); got != 1 {
		t.Fatalf("sum(shares) got=%d want=1", got)
	}

	// 1 kuruş sadece birine gider, diğerleri 0 olur
	nonZero := 0
	for _, s := range shares {
		if s.Minor() != 0 {
			nonZero++
		}
	}
	if nonZero != 1 {
		t.Fatalf("expected exactly 1 non-zero share, got %d shares=%v", nonZero, shares)
	}
}

func TestAllocateProportional_EmptyOrZero(t *testing.T) {
	if got := money.AllocateProportional(nil, money.NewMinor(100)); len(got) != 0 {
		t.Fatalf("expected empty")
	}

	bases := []money.Amount{money.NewMinor(100)}
	got := money.AllocateProportional(bases, 0)
	if len(got) != 1 || got[0].Minor() != 0 {
		t.Fatalf("expected [0], got=%v", got)
	}
}

func TestAllocateProportional_WithNegativeBase(t *testing.T) {
	bases := []money.Amount{
		money.NewMinor(-1000),
		money.NewMinor(2000),
	}
	discount := money.NewMinor(100)

	out := money.AllocateProportional(bases, discount)

	var sum int64
	for _, v := range out {
		sum += v.Minor()
	}
	if sum != 100 {
		t.Fatalf("sum got=%d want=100", sum)
	}
}
