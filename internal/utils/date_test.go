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

func TestDate_MarshalJSON(t *testing.T) {
	parsedTime, _ := time.Parse("2006-01-02", "2023-01-01")
	d := Date(parsedTime)

	bytes, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, "\"2023-01-01\"", string(bytes))
}

func TestDate_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"date": "2023-12-31"}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", time.Time(ts.Date).Format("2006-01-02"))
}

func TestDate_UnmarshalJSON_Null(t *testing.T) {
	jsonStr := `{"date": null}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)

	assert.NoError(t, err)
	assert.True(t, time.Time(ts.Date).IsZero())
}

func TestDate_IsZero(t *testing.T) {
	var d Date
	assert.True(t, d.IsZero())

	d = Date(time.Now())
	assert.False(t, d.IsZero())
}
