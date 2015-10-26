package core

import (
	"testing"
	"time"
)

// Test datetime layout
func TestTimestamp(t *testing.T) {
	str := "2015-10-22T14:15:16.999"
	ts, err := time.Parse(DateTimeLayout, str)
	if err != nil {
		t.Errorf("DateTime: Cannot parse timestamp %s (error: %s)", str, err)
	}

	if ts.Year() != 2015 ||
		ts.Month() != time.October ||
		ts.Day() != 22 ||
		ts.Hour() != 14 ||
		ts.Minute() != 15 ||
		ts.Second() != 16 ||
		ts.Nanosecond() != 999000000 {
		t.Errorf("DateTime: Cannot parse timestamp %s (wrong data parsed)", str)
	}
}
