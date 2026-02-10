package money

import "sort"

// AllocateProportional distributes `discount` across `bases` proportionally.
// Returns shares with sum(shares)=discount (minor exact).
// bases: line totals etc. discount must be >=0 and <= sum(bases) typically.
func AllocateProportional(bases []Amount, discount Amount) []Amount {
	n := len(bases)
	out := make([]Amount, n)
	if n == 0 || discount == 0 {
		return out
	}

	var total int64
	for _, b := range bases {
		if b < 0 {
			// usually you shouldn't allocate against negative lines; handle at caller.
			// We'll include them to keep function simple.
		}
		total += int64(b)
	}
	if total == 0 {
		return out
	}

	type rem struct {
		i int
		r int64
	}
	rems := make([]rem, 0, n)

	var sumShares int64
	d := int64(discount)

	for i := 0; i < n; i++ {
		numer := int64(bases[i]) * d
		base := numer / total
		r := numer % total
		out[i] = Amount(base)
		sumShares += base
		rems = append(rems, rem{i: i, r: r})
	}

	left := d - sumShares
	sort.Slice(rems, func(i, j int) bool { return rems[i].r > rems[j].r })
	for k := int64(0); k < left; k++ {
		out[rems[k%int64(n)].i]++
	}

	return out
}
