package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetServerInfo method
func TestServerInfo(t *testing.T) {
	service := testNewWS(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	if assert.NotNil(t, info, "No service info available") {
		assert.NotEmpty(t, info.Version, "No API version")
		assert.NotEmpty(t, info.Timestamp, "No server timestamp")
		// REST URL might be empty
		t.Logf("server info: %s", info)
	}
}
