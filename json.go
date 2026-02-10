package money

import (
	"encoding/json"
)

func (a Amount) MarshalJSON() ([]byte, error) {
	// JSON output: "12.34"
	return json.Marshal(a.StringFixed2())
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	// Accept: "12.34" (string)
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := ParseString(s)
	if err != nil {
		return err
	}
	*a = parsed
	return nil
}
