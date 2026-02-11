package money_test

import (
	"encoding/json"
	"testing"

	"github.com/dahaiyiyimcom/money"
)

type dto struct {
	Price money.Amount `json:"price"`
}

func TestAmount_JSONMarshal(t *testing.T) {
	d := dto{Price: money.NewMinor(1234)}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal err: %v", err)
	}
	if string(b) != `{"price":"12.34"}` {
		t.Fatalf("got=%s want=%s", string(b), `{"price":"12.34"}`)
	}
}

func TestAmount_JSONUnmarshal(t *testing.T) {
	var d dto
	err := json.Unmarshal([]byte(`{"price":"12.34"}`), &d)
	if err != nil {
		t.Fatalf("unmarshal err: %v", err)
	}
	if got := d.Price.Minor(); got != 1234 {
		t.Fatalf("got=%d want=1234", got)
	}
}

func TestAmount_JSONUnmarshal_Invalid(t *testing.T) {
	var d dto
	err := json.Unmarshal([]byte(`{"price":"12.345"}`), &d)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestAmount_UnmarshalJSON_NonString(t *testing.T) {
	var a money.Amount
	// price is a number, but we only accept string
	err := json.Unmarshal([]byte(`12.34`), &a)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestAmount_UnmarshalJSON_InvalidString(t *testing.T) {
	var a money.Amount
	err := json.Unmarshal([]byte(`"12.345"`), &a)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
