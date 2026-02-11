package money_test

import (
	"strings"
	"testing"

	"github.com/dahaiyiyimcom/money"
)

// FuzzParseString tries to find crashes/panics and parsing inconsistencies.
func FuzzParseString(f *testing.F) {
	// Seed corpus (important real-world cases)
	seeds := []string{
		"0", "0.00", "1", "1.0", "1.00", "12.34", "-12.34",
		"99999999.99", "00012.34", " 12.34 ", "+12.34",
		".50", "12.", "-0.01",
		"12.345", "abc", "12..3", "--1.00", "", "   ",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// Safety: ParseString must never panic.
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("ParseString panicked for %q: %v", s, r)
			}
		}()

		a, err := money.ParseString(s)
		if err != nil {
			return
		}

		// Property: formatting is stable and round-trips.
		out := a.StringFixed2()

		a2, err := money.ParseString(out)
		if err != nil {
			t.Fatalf("round-trip parse failed: in=%q out=%q err=%v", s, out, err)
		}
		if a2.Minor() != a.Minor() {
			t.Fatalf("round-trip mismatch: in=%q out=%q a=%d a2=%d", s, out, a.Minor(), a2.Minor())
		}

		// Optional: ensure output contains only digits, optional '-', and one '.'
		if strings.Count(out, ".") != 1 {
			t.Fatalf("unexpected format: %q", out)
		}
	})
}
