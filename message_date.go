package main

import (
	"strconv"
	"strings"
	"time"
)

type MessageDate time.Time

// String returns the MessageDate as a formatted string "Mon, 02 Jan 2006 15:04:05".
func (t *MessageDate) String() string {
	return time.Time(*t).Format("Mon, 02 Jan 2006 15:04:05")
}

// UnmarshalJSON parses a .NET datetime that has been serialised into JSON
// with a shape of "\/Date(1730328597000)\/", representing a UNIX timestamp
// with milliseconds
func (t *MessageDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.TrimPrefix(s, `"\/Date(`)
	s = strings.TrimSuffix(s, `)\/"`)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	parsed := time.Unix(i/1000, 0)
	*t = MessageDate(parsed)
	return nil
}

// MarshalJSON serializes the MessageDate to a JSON-formatted string using the RFC3339 time format.
// See https://en.wikipedia.org/wiki/ISO_8601
func (t *MessageDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*t).Format(time.RFC3339) + `"`), nil
}
