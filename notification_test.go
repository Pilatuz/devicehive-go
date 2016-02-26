package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Notification to string
// TODO: Notification from map

// Test Notification JSON marshaling
func TestNotificationJson(t *testing.T) {
	notification := NewNotification("ntf-name", "hello")
	assert.JSONEq(t, toJsonStr(t, notification), `{"notification":"ntf-name","parameters":"hello"}`)

	notification.ID = 100
	assert.JSONEq(t, toJsonStr(t, notification), `{"id":100,"notification":"ntf-name","parameters":"hello"}`)
}
