package money_test

import (
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func TestParseString_OK(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"0", 0},
		{"0.00", 0},
		{"12", 1200},
		{"12.3", 1230},
		{"12.30", 1230},
		{"12.34", 1234},
		{"  12.34  ", 1234},
		{"+12.34", 1234},
		{"-12.34", -1234},
		{"-0.50", -50},
		{".00", 0}, // wholeStr -> "0" (bu davranışı istiyorsan)
		{"0.", 0},
	}

	for _, tc := range cases {
		a, err := money.ParseString(tc.in)
		if err != nil {
			t.Fatalf("in=%q unexpected err: %v", tc.in, err)
		}
		if got := a.Minor(); got != tc.want {
			t.Fatalf("in=%q got=%d want=%d", tc.in, got, tc.want)
		}
	}
}

func TestParseString_Errors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"abc",
		"12.a",
		"12.345",  // too many decimals (strict)
		"12..34",  // invalid
		"--12.34", // invalid
		"12-34",   // invalid
	}

	for _, in := range cases {
		_, err := money.ParseString(in)
		if err == nil {
			t.Fatalf("in=%q expected error, got nil", in)
		}
	}
}

// Parse edge tests
func TestParseString_LeadingDot(t *testing.T) {
	a, err := money.ParseString(".50")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if a.Minor() != 50 {
		t.Fatalf("got=%d want=50", a.Minor())
	}
}

func TestParseString_TrailingDot(t *testing.T) {
	a, err := money.ParseString("12.")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if a.Minor() != 1200 {
		t.Fatalf("got=%d want=1200", a.Minor())
	}
}
