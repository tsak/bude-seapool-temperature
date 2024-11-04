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

func (t *Temperature) String() string {
	return fmt.Sprintf("%.1f°C", *t)
}
