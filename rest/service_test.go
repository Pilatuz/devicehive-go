package rest

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testServerURL   = "http://playground.devicehive.com/api/rest"
	testAccessKey   = ""
	testWaitTimeout = 60 * time.Second
	testLogLevel    = "debug"
)

// initialize test environment
func init() {
	flag.StringVar(&testServerURL, "url", testServerURL, "REST service URL")
	flag.StringVar(&testAccessKey, "access-key", testAccessKey, "key to access playground")
	flag.StringVar(&testLogLevel, "log-level", testLogLevel, "Logging level")
	flag.Parse()

	SetLogLevel(testLogLevel)
}

// creates new REST service
func testNewRest(t *testing.T) *Service {
	if len(testServerURL) == 0 {
		return nil
	}

	service, err := NewService(testServerURL, testAccessKey)
	assert.NoError(t, err, "Failed to create REST service")
	assert.NotNil(t, service, "No service created")
	return service
}

// convert object to JSON string.
func toJsonStr(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		// t.Errorf("Cannot convert %s to JSON (error: %s)", v, err)
		return "-" // bad JSON
	}
	return string(b)
}
