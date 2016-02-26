package rest

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
)

// Test GetServerInfo method
func TestGetServerInfoOK(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	info, err := service.GetServerInfo(testWaitTimeout)
	assert.NoError(t, err, "Failed to get server info")
	assert.NotNil(t, info, "No service info avaialble")
	// t.Logf("server info: %s", info)

	assert.False(t, len(info.Version) == 0, "No API version")
	assert.False(t, len(info.Timestamp) == 0, "No server timestamp")
	// websocket URL might be empty
}

// Test GetServerInfo method (invalid server address)
func TestGetServerInfoBadAddress(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, ".", "_", -1), "")
	assert.NoError(t, err, "Failed to create service")
	assert.NotNil(t, service, "No service created")

	info, err := service.GetServerInfo(testWaitTimeout)
	assert.Error(t, err, `No "unknown host" expected error`)
	assert.Nil(t, info, "No service info expected")
}

// Test GetServerInfo method (invalid path)
func TestGetServerInfoBadPath(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, "rest", "reZZZt", -1), "")
	assert.NoError(t, err, "Failed to create service")
	assert.NotNil(t, service, "No service created")

	info, err := service.GetServerInfo(testWaitTimeout)
	assert.Error(t, err, `No "invalid path" expected error`)
	assert.Nil(t, info, "No service info expected")
}
