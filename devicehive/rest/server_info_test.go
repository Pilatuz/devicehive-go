package rest

import (
	"strings"
	"testing"
	"time"
)

// TestGetServerInfo() unit test for /info GET method
func TestGetServerInfo(t *testing.T) {
	s, err := NewService(testServerUrl, "")
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	info, err := s.GetServerInfo(testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get server info (error: %s)", err)
		return
	}
	t.Logf("server info: %s", info)

	if len(info.Version) == 0 {
		t.Error("No API version")
	}

	if len(info.Timestamp) == 0 {
		t.Error("No server timestamp")
	}

	// websocket URL might be empty
}

// TestGetServerInfoBadAddress() unit test for /info GET method (invalid server address)
func TestGetServerInfoBadAddress(t *testing.T) {
	s, err := NewService(strings.Replace(testServerUrl, ".", "_", -1), "")
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	_, err = s.GetServerInfo(100 * time.Second)
	if err == nil {
		t.Error("Expected 'unknown host' error")
	}
}

// TestGetServerInfoBadPath() unit test for /info GET method (invalid path)
func TestGetServerInfoBadPath(t *testing.T) {
	s, err := NewService(strings.Replace(testServerUrl, "rest", "reZZZt", -1), "")
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	_, err = s.GetServerInfo(100 * time.Second)
	if err == nil {
		t.Error("Expected 'invalid path' error")
	}
}
