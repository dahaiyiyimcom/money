package money_test

import (
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func TestDBAmount_Scan_FromString(t *testing.T) {
	var db money.DBAmount
	if err := db.Scan("12.34"); err != nil {
		t.Fatalf("scan err: %v", err)
	}
	if got := db.A.Minor(); got != 1234 {
		t.Fatalf("got=%d want=1234", got)
	}
}

func TestDBAmount_Scan_FromBytes(t *testing.T) {
	var db money.DBAmount
	if err := db.Scan([]byte("12.34")); err != nil {
		t.Fatalf("scan err: %v", err)
	}
	if got := db.A.Minor(); got != 1234 {
		t.Fatalf("got=%d want=1234", got)
	}
}

func TestDBAmount_Scan_Nil(t *testing.T) {
	var db money.DBAmount
	if err := db.Scan(nil); err != nil {
		t.Fatalf("scan err: %v", err)
	}
	if got := db.A.Minor(); got != 0 {
		t.Fatalf("got=%d want=0", got)
	}
}

func TestDBAmount_Value(t *testing.T) {
	db := money.DBAmount{A: money.NewMinor(1234)}
	v, err := db.Value()
	if err != nil {
		t.Fatalf("value err: %v", err)
	}
	if v.(string) != "12.34" {
		t.Fatalf("got=%v want=%q", v, "12.34")
	}
}

// Scan tests for invalid types and formats.

func TestDBAmount_Scan_Float64(t *testing.T) {
	var m money.DBAmount
	if err := m.Scan(419.29); err != nil {
		t.Fatalf("Scan(float64) err: %v", err)
	}
	if got := m.A.Minor(); got != 41929 {
		t.Fatalf("got=%d want=%d", got, 41929)
	}
}

func TestDBAmount_Scan_Float64_RoundingSafety(t *testing.T) {
	// Float representation risk: ensure formatting to 2 decimals is used safely.
	var m money.DBAmount
	if err := m.Scan(12.34); err != nil {
		t.Fatalf("Scan(float64) err: %v", err)
	}
	if got := m.A.Minor(); got != 1234 {
		t.Fatalf("got=%d want=%d", got, 1234)
	}
}

func TestDBAmount_Scan_Variants(t *testing.T) {
	var m money.DBAmount

	// nil
	if err := m.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) err: %v", err)
	}
	if m.A.Minor() != 0 {
		t.Fatalf("Scan(nil) minor got=%d want=0", m.A.Minor())
	}

	// string
	if err := m.Scan("419.29"); err != nil {
		t.Fatalf("Scan(string) err: %v", err)
	}
	if m.A.Minor() != 41929 {
		t.Fatalf("Scan(string) got=%d want=41929", m.A.Minor())
	}

	// []byte
	if err := m.Scan([]byte("12.34")); err != nil {
		t.Fatalf("Scan([]byte) err: %v", err)
	}
	if m.A.Minor() != 1234 {
		t.Fatalf("Scan([]byte) got=%d want=1234", m.A.Minor())
	}

	// float64 (the real-world bug)
	if err := m.Scan(12.34); err != nil {
		t.Fatalf("Scan(float64) err: %v", err)
	}
	if m.A.Minor() != 1234 {
		t.Fatalf("Scan(float64) got=%d want=1234", m.A.Minor())
	}

	// float32
	if err := m.Scan(float32(12.34)); err != nil {
		t.Fatalf("Scan(float32) err: %v", err)
	}
	if m.A.Minor() != 1234 {
		t.Fatalf("Scan(float32) got=%d want=1234", m.A.Minor())
	}
}

func TestDBAmount_Scan_UnsupportedType(t *testing.T) {
	var m money.DBAmount
	if err := m.Scan(true); err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}
