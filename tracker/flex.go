package tracker

import (
	"encoding/json"
	"fmt"
)

// FlexString represents an ID that the Yandex Tracker API may return as either
// a JSON string or a JSON number, depending on the endpoint. It stores the
// value as a string and always marshals back to a JSON string.
type FlexString string

// UnmarshalJSON implements the json.Unmarshaler interface.
// It accepts both JSON strings ("42") and JSON numbers (42).
func (f *FlexString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexString(n.String())
		return nil
	}

	return fmt.Errorf("FlexString: cannot unmarshal %s", string(data))
}

// MarshalJSON implements the json.Marshaler interface.
// It always outputs a JSON string.
func (f FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}
