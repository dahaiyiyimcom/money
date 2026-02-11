package money_test

import (
	"encoding/json"
	"testing"

	"github.com/dahaiyiyimcom/money"
)

func BenchmarkParseString(b *testing.B) {
	in := "419.29"
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = money.ParseString(in)
	}
}

func BenchmarkStringFixed2(b *testing.B) {
	a := money.NewMinor(41929)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = a.StringFixed2()
	}
}

func BenchmarkJSONMarshalAmount(b *testing.B) {
	type DTO struct {
		Price money.Amount `json:"price"`
	}
	dto := DTO{Price: money.NewMinor(41929)}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(dto)
	}
}

func BenchmarkAllocateProportional(b *testing.B) {
	bases := make([]money.Amount, 50)
	for i := 0; i < len(bases); i++ {
		bases[i] = money.NewMinor(int64(1000 + i*3))
	}
	discount := money.NewMinor(1234)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = money.AllocateProportional(bases, discount)
	}
}

func BenchmarkMulRatio(b *testing.B) {
	a := money.NewMinor(1234567) // 12,345.67
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = a.MulRatio(18, 100, money.RoundHalfUp)
	}
}
