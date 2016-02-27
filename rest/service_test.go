package rest

import (
	"encoding/json"
	"flag"
	"strings"
	"testing"
	"time"

	dh "github.com/pilatuz/go-devicehive"
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
	if assert.NotNil(t, service, "No service created") {
		service.SetTimeout(testWaitTimeout)

		// check DeviceService is implemented
		_ = dh.DeviceService(service)
	}
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

// Test GetServerInfo method (invalid server address)
func TestServiceBadAddress(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, ".", "_", -1), "")
	assert.NoError(t, err, "Failed to create service")
	if assert.NotNil(t, service, "No service created") {
		info, err := service.GetServerInfo()
		assert.Error(t, err, `No "unknown host" expected error`)
		assert.Nil(t, info, "No service info expected")
	}
}

// Test GetServerInfo method (invalid path)
func TestServiceBadPath(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewService(strings.Replace(testServerURL, "rest", "reZZZt", -1), "")
	assert.NoError(t, err, "Failed to create service")
	if assert.NotNil(t, service, "No service created") {
		info, err := service.GetServerInfo()
		assert.Error(t, err, `No "invalid path" expected error`)
		assert.Nil(t, info, "No service info expected")
	}
}

// Test service.Stop method
func TestServiceStop(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		service.Stop()
	}()

	N := 5
	ch := make(chan int, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			info, err := service.GetServerInfo()
			assert.Error(t, err, `No "stopped" expected error`)
			assert.Nil(t, info, "No service info expected")
			ch <- i
		}(i)
	}

	// wait all
	for i := 0; i < N; i++ {
		<-ch
	}
}
