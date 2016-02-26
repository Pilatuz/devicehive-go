package rest

import (
	"strings"
	"testing"

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
	assert.NotNil(t, info, "No service info available")
	assert.NotEmpty(t, info.Version, "No API version")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp")
	// websocket URL might be empty

	// t.Logf("server info: %s", info)
}

// Test GetServerInfo method (invalid server address)
func TestGetServerInfoBadAddress(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, ".", "_", -1), "")
	ok := assert.NoError(t, err, "Failed to create service") &&
		assert.NotNil(t, service, "No service created")

	if ok {
		info, err := service.GetServerInfo(testWaitTimeout)
		assert.Error(t, err, `No "unknown host" expected error`)
		assert.Nil(t, info, "No service info expected")
	}
}

// Test GetServerInfo method (invalid path)
func TestGetServerInfoBadPath(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, "rest", "reZZZt", -1), "")
	ok := assert.NoError(t, err, "Failed to create service") &&
		assert.NotNil(t, service, "No service created")

	if ok {
		info, err := service.GetServerInfo(testWaitTimeout)
		assert.Error(t, err, `No "invalid path" expected error`)
		assert.Nil(t, info, "No service info expected")
	}
}
