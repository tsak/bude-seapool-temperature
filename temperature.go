package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Temperature float64

// UnmarshalJSON implements the json.Unmarshaler interface for the Temperature type.
// It converts a JSON-encoded string to a Temperature (float64) value.
func (t *Temperature) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.TrimPrefix(s, `"`)
	s = strings.TrimSuffix(s, `"`)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*t = Temperature(f)
	return nil
}

// MarshalJSON formats a temperature value as a stringified float with one decimal.
func (t *Temperature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.1f", *t)), nil
}

// String formats the Temperature value as a string with one decimal place, followed by the Celsius symbol (°C).
func (t *Temperature) String() string {
	return fmt.Sprintf("%.1f°C", *t)
}
