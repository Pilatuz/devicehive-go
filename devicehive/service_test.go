package devicehive

import (
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"github.com/devicehive/devicehive-go/devicehive/rest"
	"github.com/devicehive/devicehive-go/devicehive/ws"
	"strings"
	"testing"
	"time"
)

var (
	testRestServerUrl = "http://playground.devicehive.com/api/rest"
	testWsServerUrl   = "ws://playground.devicehive.com/api/websocket"
	testAccessKey     = "<place access key here>"
	testDeviceId      = "go-test-dev-id"
	testDeviceKey     = "go-test-dev-key"
	testWaitTimeout   = 20 * time.Second
)

var (
	testRawRestService *rest.Service
	testRawWsService   *ws.Service
	testRestService    Service
	testWsService      Service
)

// initialize test environment
func init() {
	// TODO: load test parameters from command line arguments
	// or from environment variables

	// TODO: check the test parameters

	log.SetLevel(log.TRACE)
	testRawRestService, _ = rest.NewService(testRestServerUrl, testAccessKey)
	testRawWsService, _ = ws.NewService(testWsServerUrl, testAccessKey)

	// safe dereference!
	if testRawRestService != nil {
		testRestService = testRawRestService
	}
	if testRawWsService != nil {
		testWsService = testRawWsService
	}
}

// check the GetServerInfo method
func testCheckGetServerInfo(t *testing.T, service Service) {
	if service == nil {
		return // do nothing
	}

	info, err := service.GetServerInfo(testWaitTimeout)
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

// Test GetServerInfo method
func TestGetServerInfo(t *testing.T) {
	testCheckGetServerInfo(t, testRestService)
	testCheckGetServerInfo(t, testWsService)
}

// Test GetServerInfo method (invalid server address)
func TestGetServerInfoBadAddress(t *testing.T) {
	// REST
	if testRestService != nil {
		rs, err := rest.NewService(strings.Replace(testRestServerUrl, ".", "_", -1), "")
		if err != nil {
			t.Errorf("Failed to create service (error: %s)", err)
			return
		}

		_, err = rs.GetServerInfo(testWaitTimeout)
		if err == nil {
			t.Error("Expected 'unknown host' error")
		}
	}

	// Websocket
	if testWsService != nil {
		wss, err := ws.NewService(strings.Replace(testWsServerUrl, ".", "_", -1), "")
		if err != nil {
			t.Errorf("Failed to create service (error: %s)", err)
			return
		}

		_, err = wss.GetServerInfo(testWaitTimeout)
		if err == nil {
			t.Error("Expected 'unknown host' error")
		}
	}
}

// Test GetServerInfo method (invalid path)
func TestGetServerInfoBadPath(t *testing.T) {
	// REST
	if testRestService != nil {
		rs, err := NewService(strings.Replace(testRestServerUrl, "rest", "reZZZt", -1), "")
		if err != nil {
			t.Errorf("Failed to create service (error: %s)", err)
			return
		}

		_, err = rs.GetServerInfo(testWaitTimeout)
		if err == nil {
			t.Error("Expected 'invalid path' error")
		}
	}
}

// check the RegisterDevice method (helper)
func testCheckRegisterDevice1(t *testing.T, service Service, device core.Device, deletePrevious bool) {
	if service == nil {
		return // do nothing
	}

	if deletePrevious && testRawRestService != nil {
		_ = testRawRestService.DeleteDevice(&device, testWaitTimeout)
		// ignore possible errors
	}

	err := service.RegisterDevice(&device, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to register device %v (error: %s)", device, err)
		return
	}
}

// check the GetDevice method
func testCheckGetDevice(t *testing.T, service Service, device core.Device) {
	if service == nil {
		return // do nothing
	}

	device2, err := service.GetDevice(device.Id, device.Key, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get device (error: %s)", err)
		return
	}

	// TODO: compare device and device2
	_ = device2
}

// Test RegisterDevice method (no device class, no network)
func testCheckRegisterDevice(t *testing.T, service Service, device core.Device, suffix string) {
	if service == nil {
		return // do nothing
	}

	device.Id += suffix
	if len(device.Key) != 0 {
		device.Key += suffix
	}
	device.Name += suffix

	// create new device
	testCheckRegisterDevice1(t, service, device, true)
	testCheckGetDevice(t, service, device)

	// change some data
	device.Data = "new data"
	device.Status = "Bad"
	device.Name = "new test-name"

	// update device
	testCheckRegisterDevice1(t, service, device, false)
	testCheckGetDevice(t, service, device)
}

// Test RegisterDevice method (no device class, no network)
func TestRegisterDeviceNoClassNoNet(t *testing.T) {
	device := core.NewDevice(testDeviceId, "test-name", nil)

	testCheckRegisterDevice(t, testRestService, *device, "-1a")
	testCheckRegisterDevice(t, testWsService, *device, "-1b")
}

// Test RegisterDevice method (no device class)
func TestRegisterDeviceNoClass(t *testing.T) {
	device := core.NewDevice(testDeviceId, "test-name",
		core.NewDeviceClass("go-device-class", "1.2.3"))

	testCheckRegisterDevice(t, testRestService, *device, "-2a")
	testCheckRegisterDevice(t, testWsService, *device, "-2b")
}

// Test RegisterDevice method
func TestRegisterDevice(t *testing.T) {
	device := core.NewDevice(testDeviceId, "test-name",
		core.NewDeviceClass("go-device-class", "1.2.3"))
	device.Network = core.NewNetwork("go-net-name", "net-key")

	testCheckRegisterDevice(t, testRestService, *device, "-3a")
	testCheckRegisterDevice(t, testWsService, *device, "-3b")
}

//// TestInsertCommand() unit test for /command/insert POST method,
//// /command/update PUT method, /command/get GET method
//// test device should be already registered!
//func TestInsertCommand(t *testing.T) {
//	TestRegisterDevice(t)
//	if t.Failed() {
//		return // nothing to test without device
//	}

//	s, err := NewService(testServerUrl, testAccessKey)
//	if err != nil {
//		t.Errorf("Failed to create service (error: %s)", err)
//		return
//	}

//	device := &core.Device{Id: testDeviceId, Key: testDeviceKey}
//	command := &core.Command{Name: "cmd-test", Parameters: 123, Lifetime: 600}
//	err = s.InsertCommand(device, command, testWaitTimeout)
//	if err != nil {
//		t.Errorf("Failed to insert command (error: %s)", err)
//		return
//	}
//	t.Logf("command: %s", command)

//	command.Status = "Done"
//	command.Result = 12345
//	err = s.UpdateCommand(device, command, testWaitTimeout)
//	if err != nil {
//		t.Errorf("Failed to update command (error: %s)", err)
//		return
//	}

//	*command, err = s.GetCommand(device, command.Id, testWaitTimeout)
//	if err != nil {
//		t.Errorf("Failed to get command (error: %s)", err)
//		return
//	}
//	t.Logf("command: %s", command)
//}

// TODO: TestPollCommand

// check InsertNotification method
// device should be already registered
func testCheckInsertNotification(t *testing.T, service Service, device *core.Device, notification core.Notification) {
	if service == nil {
		return // do nothing
	}

	err := service.InsertNotification(device, &notification, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to insert notification (error: %s)", err)
		return
	}
	//t.Logf("notification: %s", notification)

	notification2, err := service.GetNotification(device, notification.Id, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get notification (error: %s)", err)
		return
	}
	//t.Logf("notification: %s", notification)

	// TODO: compare notification & notification2
	_ = notification2
}

// Test InsertNotification method
func TestInsertNotification(t *testing.T) {
	// create device (REST)
	device := core.NewDevice(testDeviceId, "test-name",
		core.NewDeviceClass("go-device-class", "1.2.3"))
	device.Network = core.NewNetwork("go-net-name", "net-key")
	testCheckRegisterDevice1(t, testRestService, *device, false)
	if t.Failed() {
		return // nothing to test without device
	}

	notification := core.NewNotification("ntf-test", 12345)
	testCheckInsertNotification(t, testRestService, device, *notification)
	testCheckInsertNotification(t, testWsService, device, *notification)
}
