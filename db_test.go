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
