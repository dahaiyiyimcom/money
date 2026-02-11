package money_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/dahaiyiyimcom/money"
)

func TestAllocateProportional_Property_SumExact(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 5000; iter++ {
		n := r.Intn(15) // 0..14
		bases := make([]money.Amount, n)

		var total int64
		for i := 0; i < n; i++ {
			// Keep bases mostly positive; include some zeros
			v := int64(r.Intn(50000)) // up to 500.00
			bases[i] = money.NewMinor(v)
			total += v
		}

		// Choose discount in [0, total] (if total==0 => 0)
		var d int64
		if total > 0 {
			d = int64(r.Intn(int(total + 1)))
		}
		discount := money.NewMinor(d)

		out := money.AllocateProportional(bases, discount)

		// Invariant 1: output length matches input length
		if len(out) != len(bases) {
			t.Fatalf("len(out)=%d len(bases)=%d", len(out), len(bases))
		}

		// Invariant 2: sum(out) == discount
		var sum int64
		for _, x := range out {
			sum += x.Minor()
		}
		if sum != discount.Minor() {
			t.Fatalf("sum(out)=%d discount=%d bases=%v out=%v", sum, discount.Minor(), bases, out)
		}

		// Optional invariant 3: if discount==0 -> all zeros
		if discount.Minor() == 0 {
			for i, x := range out {
				if x.Minor() != 0 {
					t.Fatalf("expected all zeros when discount=0; i=%d got=%d", i, x.Minor())
				}
			}
		}
	}
}
