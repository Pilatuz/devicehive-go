package core

import (
	"encoding/json"
	"testing"
)

// Test Notification JSON marshaling
func TestJsonNotification(t *testing.T) {
	check := func(notification Notification, expectedJson string) {
		jsonBytes, err := json.Marshal(notification)
		if err != nil {
			t.Errorf("Notification: Cannot convert %+v to JSON (error: %s)", notification, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("Notification: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				notification, string(jsonBytes), expectedJson)
		}

		t.Logf("Notification: %+v converted to %s", notification, string(jsonBytes))
	}

	notification := Notification{
		Name:       "ntf-name",
		Parameters: "hello"}

	check(notification, `{"notification":"ntf-name","parameters":"hello"}`)

	notification.Id = 100
	check(notification, `{"id":100,"notification":"ntf-name","parameters":"hello"}`)
}
