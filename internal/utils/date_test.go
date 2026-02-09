package utils

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Date Date `json:"date"`
}

type TestStructPtr struct {
	Date *Date `json:"date"`
}

// Test MarshalJSON outputs ISO 8601 with UTC timezone
func TestDate_MarshalJSON_ISO8601(t *testing.T) {
	parsedTime, _ := time.Parse("2006-01-02", "2000-01-15")
	d := Date(parsedTime)

	bytes, err := json.Marshal(d)
	assert.NoError(t, err)
	// Should output ISO 8601 with UTC timezone and milliseconds
	assert.Equal(t, "\"2000-01-15T00:00:00.000Z\"", string(bytes))
}

// Test MarshalJSON outputs null for zero date
func TestDate_MarshalJSON_Null(t *testing.T) {
	var d Date // zero value

	bytes, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(bytes))
}

// Test UnmarshalJSON accepts YYYY-MM-DD format (from frontend)
func TestDate_UnmarshalJSON_YYYYMMDD(t *testing.T) {
	jsonStr := `{"date": "2023-12-31"}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", time.Time(ts.Date).Format("2006-01-02"))
}

// Test UnmarshalJSON accepts ISO 8601 format
func TestDate_UnmarshalJSON_ISO8601(t *testing.T) {
	jsonStr := `{"date": "2023-12-31T00:00:00Z"}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", time.Time(ts.Date).Format("2006-01-02"))
}

// Test UnmarshalJSON accepts ISO 8601 with milliseconds
func TestDate_UnmarshalJSON_ISO8601WithMs(t *testing.T) {
	jsonStr := `{"date": "2023-12-31T00:00:00.000Z"}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", time.Time(ts.Date).Format("2006-01-02"))
}

// Test UnmarshalJSON handles null
func TestDate_UnmarshalJSON_Null(t *testing.T) {
	jsonStr := `{"date": null}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.True(t, time.Time(ts.Date).IsZero())
}

// Test UnmarshalJSON handles empty string
func TestDate_UnmarshalJSON_EmptyString(t *testing.T) {
	jsonStr := `{"date": ""}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.True(t, time.Time(ts.Date).IsZero())
}

// Test UnmarshalJSON rejects invalid format
func TestDate_UnmarshalJSON_InvalidFormat(t *testing.T) {
	jsonStr := `{"date": "invalid-date"}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

// Test IsZero
func TestDate_IsZero(t *testing.T) {
	var d Date
	assert.True(t, d.IsZero())

	d = Date(time.Now())
	assert.False(t, d.IsZero())
}

// Test DatePtr returns nil for zero date
func TestDatePtr_ZeroDate(t *testing.T) {
	var d Date
	ptr := DatePtr(d)
	assert.Nil(t, ptr)
}

// Test DatePtr returns pointer for non-zero date
func TestDatePtr_NonZeroDate(t *testing.T) {
	parsedTime, _ := time.Parse("2006-01-02", "2000-01-15")
	d := Date(parsedTime)
	ptr := DatePtr(d)

	assert.NotNil(t, ptr)
	assert.Equal(t, d, *ptr)
}

// Test DateValue returns zero for nil pointer
func TestDateValue_NilPointer(t *testing.T) {
	var ptr *Date
	d := DateValue(ptr)
	assert.True(t, d.IsZero())
}

// Test DateValue dereferences non-nil pointer
func TestDateValue_NonNilPointer(t *testing.T) {
	parsedTime, _ := time.Parse("2006-01-02", "2000-01-15")
	original := Date(parsedTime)
	ptr := &original
	d := DateValue(ptr)

	assert.False(t, d.IsZero())
	assert.Equal(t, original, d)
}

// Test roundtrip: marshal then unmarshal
func TestDate_Roundtrip(t *testing.T) {
	parsedTime, _ := time.Parse("2006-01-02", "2000-01-15")
	original := TestStruct{Date: Date(parsedTime)}

	// Marshal to JSON
	bytes, err := json.Marshal(original)
	assert.NoError(t, err)

	// Unmarshal back
	var result TestStruct
	err = json.Unmarshal(bytes, &result)
	assert.NoError(t, err)

	// Check date is preserved (comparing as strings since time zones may differ)
	assert.Equal(t,
		time.Time(original.Date).Format("2006-01-02"),
		time.Time(result.Date).Format("2006-01-02"),
	)
}
