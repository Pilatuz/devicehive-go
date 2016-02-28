package ws

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
	testServerURL          = "ws://playground.devicehive.com/api/websocket/"
	testAccessKey          = ""
	testDeviceId           = "go-test-dev-id"
	testDeviceName         = "go-test-dev-name"
	testDeviceKey          = "go-test-dev-key"
	testDeviceClassName    = "go-device-class"
	testDeviceClassVersion = "1.2.3"
	testNetworkName        = ""
	testNetworkKey         = ""
	testNetworkDesc        = ""
	testWaitTimeout        = 60 * time.Second
	testLogLevel           = "debug"
)

// initialize test environment
func init() {
	flag.StringVar(&testServerURL, "url", testServerURL, "WS service URL")
	flag.StringVar(&testAccessKey, "access-key", testAccessKey, "key to access playground")
	flag.StringVar(&testLogLevel, "log-level", testLogLevel, "logging level")

	flag.StringVar(&testDeviceId, "device-id", testDeviceId, "test Device identifier")
	flag.StringVar(&testDeviceName, "device-name", testDeviceName, "test Device name")
	flag.StringVar(&testDeviceKey, "device-key", testDeviceKey, "test Device key")

	flag.StringVar(&testDeviceClassName, "device-class-name", testDeviceClassName, "test Device class name")
	flag.StringVar(&testDeviceClassVersion, "device-class-version", testDeviceClassVersion, "test Device class version")

	flag.StringVar(&testNetworkName, "network-name", testNetworkName, "test Network name")
	flag.StringVar(&testNetworkKey, "network-key", testNetworkKey, "test Network key")
	flag.StringVar(&testNetworkDesc, "network-desc", testNetworkDesc, "test Network description")

	flag.Parse()

	SetLogLevel(testLogLevel)
}

// creates new test Device with device class initialized
func testNewDevice() *dh.Device {
	dc := dh.NewDeviceClass(testDeviceClassName, testDeviceClassVersion)
	device := dh.NewDevice(testDeviceId, testDeviceName, dc)
	device.Key = testDeviceKey
	return device
}

// creates new test Network
func testNewNetwork() *dh.Network {
	if len(testNetworkName) != 0 {
		network := dh.NewNetwork(testNetworkName, testNetworkKey)
		network.Description = testNetworkDesc
		return network
	}
	return nil
}

// creates new WS /device service
func testNewWsDevice(t *testing.T) *Service {
	if len(testServerURL) == 0 {
		return nil
	}

	service, err := NewDeviceService(testServerURL, testAccessKey)
	assert.NoError(t, err, "Failed to create WS device service")
	if assert.NotNil(t, service, "No service created") {
		service.SetTimeout(testWaitTimeout)

		// check DeviceService is implemented
		_ = dh.DeviceService(service)
	}

	return service
}

// creates new WS /client service
func testNewWsClient(t *testing.T) *Service {
	if len(testServerURL) == 0 {
		return nil
	}

	service, err := NewClientService(testServerURL, testAccessKey)
	assert.NoError(t, err, "Failed to create WS client service")
	if assert.NotNil(t, service, "No service created") {
		service.SetTimeout(testWaitTimeout)

		// check ClientService is implemented
		_ = dh.ClientService(service)
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

	service, err := NewDeviceService(strings.Replace(testServerURL, ".", "_", -1), "")
	assert.Error(t, err, `No "unknown host" expected error`)
	assert.Nil(t, service, "No service expected")
}

// Test GetServerInfo method (invalid path)
func TestServiceBadPath(t *testing.T) {
	if len(testServerURL) == 0 {
		return // nothing to test
	}

	service, err := NewDeviceService(strings.Replace(testServerURL, "websocket", "webZZZocket", -1), "")
	assert.Error(t, err, `No "invalid path" expected error`)
	assert.Nil(t, service, "No service expected")
}

// Test service.Stop method
func TestServiceStop(t *testing.T) {
	service := testNewWsDevice(t)
	if service == nil {
		return // nothing to test
	}

	N := 5 // requests at the same time
	ch := make(chan int, N)

	go func() {
		time.Sleep(100 * time.Millisecond)
		service.Stop()
		ch <- -1
	}()

	for i := 0; i < N; i++ {
		go func(i int) {
			info, err := service.GetServerInfo()
			assert.Error(t, err, `No "stopped" expected error`)
			assert.Nil(t, info, "No service info expected")
			ch <- i
		}(i)
	}

	// wait all
	for i := 0; i < N+1; i++ {
		<-ch
	}
}
