# money

Safe, deterministic money type for Go services.

Production-grade money model for microservices and modular systems:

* Minor unit (int64) calculations — **no float errors**
* MySQL DECIMAL(10,2) scan/write support
* JSON transport as string ("12.34") — HTTP & Kafka safe
* Deterministic rounding helpers
* Proportional allocation (discount distribution) with exact sum guarantee
* Single-currency systems supported (no currency field required)
* Zero external dependencies
* Immutable value type (goroutine-safe)

Designed for: marketplace, basket, order, payment, pricing, discount, tax services.

---

![CI](https://github.com/dahaiyiyimcom/money/actions/workflows/ci.yml/badge.svg)
![Go Reference](https://pkg.go.dev/badge/github.com/dahaiyiyimcom/money.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/dahaiyiyimcom/money)

---

## Requirements

* Go 1.21+
* Tested with: go-sql-driver/mysql

---

## Install
```
go get github.com/dahaiyiyimcom/money@latest
```
Import:
```
import money "github.com/dahaiyiyimcom/money"
```
---

## Core Concept

All money values are stored and computed in minor units:

Display → Stored
12.34 → 1234

This eliminates floating-point rounding errors completely.

---

## Why Not float64?

Floating point numbers cannot exactly represent decimal fractions.

Example:

0.1 + 0.2 != 0.3

This causes rounding drift in:

* basket totals
* tax calculations
* discount distribution
* settlement calculations

This package guarantees deterministic integer math instead.

---

## Quick Start

Create amount:
```
a := money.NewMinor(1234) // 12.34
```
Format:
```
s := a.StringFixed2()
// "12.34"
```
---

## Arithmetic

Add / Subtract:
```
a := money.NewMinor(1000) // 10.00
b := money.NewMinor(250)  // 2.50

c := a.Add(b) // 12.50
d := a.Sub(b) // 7.50
```
Multiply by quantity:
```
unit := money.NewMinor(1234)
total := unit.MulQty(3)
// 37.02
```
---

## Percent & Ratio Calculations

VAT example:
```
price := money.NewMinor(1234)

vat := price.Percent(18, money.RoundHalfUp)
// 2.22

grand := price.Add(vat)
// 14.56
```
Generic ratio:
```
result := price.MulRatio(1, 3, money.RoundHalfUp)
```
---

## Rounding Modes
```
money.RoundHalfUp
money.RoundFloor
money.RoundCeil
```
Recommendation: RoundHalfUp for commerce pricing and tax logic.

---

## MySQL Integration (DECIMAL(10,2))

Repository struct:
```
type ProductRow struct {
Price money.DBAmount
}
```
Scan:
```
var r ProductRow
err := row.Scan(&r.Price)

price := r.Price.A // money.Amount

Write:

_, err := db.Exec(
"INSERT INTO products(price) VALUES(?)",
money.DBAmount{A: price},
)
```
No float conversion anywhere.

---

## JSON / HTTP / Kafka Transport

money.Amount marshals as string.
```
type DTO struct {
Price money.Amount `json:"price"`
}

dto := DTO{Price: money.NewMinor(1234)}

b, _ := json.Marshal(dto)
// {"price":"12.34"}
```
Unmarshal:
```
var dto DTO
_ = json.Unmarshal([]byte(`{"price":"12.34"}`), &dto)

dto.Price.Minor() == 1234
```
Invalid precision (more than 2 decimals) returns error.

---

## Parsing from String
```
a, err := money.ParseString("12.34")
```
Accepted:
```
"12"
"12.3"
"12.30"
"-0.50"
" 12.34 "
```
Rejected:
```
"12.345"
"abc"
"12..3"
```
Strict by design (max 2 fractional digits).

---

## Proportional Allocation (Discount Distribution)
```
lines := []money.Amount{
money.NewMinor(1000),
money.NewMinor(2000),
money.NewMinor(3000),
}

discount := money.NewMinor(100)

shares := money.AllocateProportional(lines, discount)
```
Guarantees:

* Sum(shares) == discount
* No minor-unit loss
* Deterministic remainder distribution
* Stable ordering

Ideal for basket-level discount distribution.

---

## Concurrency

money.Amount is immutable and safe to use across goroutines.

---

## Guarantees

* No float usage
* Deterministic arithmetic
* Minor-unit integer math
* MySQL DECIMAL safe conversion
* JSON safe string transport
* Exact allocation totals
* Explicit rounding policy
* Strict decimal parsing
* Stable public API
* Zero hidden rounding

---

## Recommended Data Flow

DB (DECIMAL) → DBAmount → Amount(minor int64) → calculations → JSON string

Never use float for money in any layer.

---

## Testing

Run tests:
```
go test ./... -count=1
```
CI runs:

* unit tests
* vet
* build check

---

## Versioning

Semantic versioning is used.

Breaking changes will be released under:
```
github.com/dahaiyiyimcom/money/v2
```
---

## Non-Goals

This package intentionally does NOT provide:

* Multi-currency support
* FX conversion
* Currency symbols
* Localization / formatting
* Accounting reports

Scope is limited to safe money arithmetic and transport.

---

## CI Example

.github/workflows/ci.yml
```yml
name: CI
on:
push: { branches: ["main"] }
pull_request: {}

jobs:
test:
runs-on: ubuntu-latest
steps:
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
  with:
  go-version: "1.22"
- run: go test ./... -count=1
- run: go vet ./...
```
---

## Contributing

* Keep API backward compatible
* Add tests for every behavior change
* Do not introduce float usage
* Do not change rounding defaults without major version bump

---

## License

MIT
