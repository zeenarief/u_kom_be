package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Date represents a date without time component (YYYY-MM-DD)
type Date time.Time

const DateLayout = "2006-01-02"

// UnmarshalJSON parses JSON string "YYYY-MM-DD" into Date
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		*d = Date(time.Time{})
		return nil
	}
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON formats Date as JSON string "YYYY-MM-DD"
func (d Date) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(time.Time(d).Format(DateLayout))
}

// Scanc implements the Scanner interface for database values
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date(v)
	case []byte:
		return d.parseString(string(v))
	case string:
		return d.parseString(v)
	default:
		return errors.New("failed to scan Date")
	}
	return nil
}

func (d *Date) parseString(s string) error {
	// Try parsing standard date layout
	t, err := time.Parse(DateLayout, s)
	if err == nil {
		*d = Date(t)
		return nil
	}

	// Try parsing full timestamp just in case DB returns it
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		*d = Date(t)
		return nil
	}

	return fmt.Errorf("could not parse date: %s", s)
}

// Value implements the driver Valuer interface for database storage
func (d Date) Value() (driver.Value, error) {
	if time.Time(d).IsZero() {
		return nil, nil
	}
	return time.Time(d).Format(DateLayout), nil
}

// String returns the date as a string
func (d Date) String() string {
	if time.Time(d).IsZero() {
		return ""
	}
	return time.Time(d).Format(DateLayout)
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (d Date) IsZero() bool {
	return time.Time(d).IsZero()
}

// Format returns a textual representation of the time value formatted according
// to the layout defined by the argument.
func (d Date) Format(layout string) string {
	return time.Time(d).Format(layout)
}

// ToTime converts Date back to time.Time
func (d Date) ToTime() time.Time {
	return time.Time(d)
}
