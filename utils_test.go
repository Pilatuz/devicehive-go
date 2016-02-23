package devicehive

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// convert object to JSON string.
func toJsonStr(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Cannot convert %s to JSON (error: %s)", v, err)
		// return ""
	}
	return string(b)
}

// Test datetime layout
func TestTimestampFormat(t *testing.T) {
	str := "2015-10-22T14:15:16.999"
	ts, err := time.Parse(DateTimeLayout, str)
	assert.NoError(t, err)

	assert.Equal(t, ts.Year(), 2015)
	assert.Equal(t, ts.Month(), time.October)
	assert.Equal(t, ts.Day(), 22)
	assert.Equal(t, ts.Hour(), 14)
	assert.Equal(t, ts.Minute(), 15)
	assert.Equal(t, ts.Second(), 16)
	assert.Equal(t, ts.Nanosecond(), 999000000)
}
