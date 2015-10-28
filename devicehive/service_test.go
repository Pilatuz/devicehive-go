package devicehive

import (
	"flag"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"github.com/devicehive/devicehive-go/devicehive/rest"
	"github.com/devicehive/devicehive-go/devicehive/ws"
	"strings"
	"testing"
	"time"
)

var (
	testRestServerUrl      = "http://playground.devicehive.com/api/rest"
	testWsServerUrl        = "ws://playground.devicehive.com/api/websocket"
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
	testLogLevel           = "NOLOG"
)

// initialize test environment
func init() {
	flag.StringVar(&testRestServerUrl, "rest-url", testRestServerUrl, "REST service URL")
	flag.StringVar(&testWsServerUrl, "ws-url", testWsServerUrl, "Websocket service URL")
	flag.StringVar(&testAccessKey, "access-key", testAccessKey, "key to access playground")

	flag.StringVar(&testDeviceId, "device-id", testDeviceId, "test Device identifier")
	flag.StringVar(&testDeviceName, "device-name", testDeviceName, "test Device name")
	flag.StringVar(&testDeviceKey, "device-key", testDeviceKey, "test Device key")

	flag.StringVar(&testDeviceClassName, "device-class-name", testDeviceClassName, "test Device class name")
	flag.StringVar(&testDeviceClassVersion, "device-class-version", testDeviceClassVersion, "test Device class version")

	flag.StringVar(&testNetworkName, "network-name", testNetworkName, "test Network name")
	flag.StringVar(&testNetworkKey, "network-key", testNetworkKey, "test Network key")
	flag.StringVar(&testNetworkDesc, "network-desc", testNetworkDesc, "test Network description")

	flag.StringVar(&testLogLevel, "log-level", testLogLevel, "Logging level: WARN INFO DEBUG TRACE or NOLOG")
	flag.Parse()

	log.SetLevelByName(testLogLevel)
}

// creates new REST service
func testNewRest(t *testing.T) (service *rest.Service) {
	service, err := rest.NewService(testRestServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create REST service (error: %s)", err)
	}
	return
}

// creates new Websocket service
func testNewWs(t *testing.T) (service *ws.Service) {
	service, err := ws.NewService(testWsServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create WS service (error: %s)", err)
	}
	return
}

// creates new REST service (abstract interface)
func testNewRestService(t *testing.T) (service Service) {
	if rs := testNewRest(t); rs != nil {
		service = rs
	}
	return
}

// creates new Websocket service (abstract interface)
func testNewWsService(t *testing.T) (service Service) {
	if wss := testNewWs(t); wss != nil {
		service = wss
	}
	return
}

// creates new test Device with device class initialized
func testNewDevice() (device *core.Device) {
	dc := core.NewDeviceClass(testDeviceClassName, testDeviceClassVersion)
	device = core.NewDevice(testDeviceId, testDeviceName, dc)
	device.Key = testDeviceKey
	return
}

// creates new test Network
func testNewNetwork() (network *core.Network) {
	if len(testNetworkName) != 0 {
		network = core.NewNetwork(testNetworkName, testNetworkKey)
		network.Description = testNetworkDesc
	}
	return
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
	//t.Logf("server info: %s", info)

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
	testCheckGetServerInfo(t, testNewRestService(t))
	testCheckGetServerInfo(t, testNewWsService(t))
}

// Test GetServerInfo method (invalid server address)
func TestGetServerInfoBadAddress(t *testing.T) {
	// REST
	if len(testRestServerUrl) != 0 {
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
	if len(testWsServerUrl) != 0 {
		_, err := ws.NewService(strings.Replace(testWsServerUrl, ".", "_", -1), "")
		if err == nil {
			t.Errorf("Expected 'unknown host' error")
			return
		}
	}
}

// Test GetServerInfo method (invalid path)
func TestGetServerInfoBadPath(t *testing.T) {
	// REST
	if len(testRestServerUrl) != 0 {
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

	if deletePrevious {
		rs := testNewRest(t)
		if !t.Failed() {
			_ = rs.DeleteDevice(&device, testWaitTimeout)
			// ignore possible errors
		}
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

	_ = device2 // TODO: compare device and device2
}

// check the RegisterDevice method
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
	device.Name += "-new"

	// update device
	testCheckRegisterDevice1(t, service, device, false)
	testCheckGetDevice(t, service, device)
}

//// Test RegisterDevice method (no device class, no network)
//func TestRegisterDeviceNoClassNoNet(t *testing.T) {
//	device := testNewDevice()
//	device.DeviceClass = nil
//
//	testCheckRegisterDevice(t, testNewRestService(t), *device, "-1a")
//	testCheckRegisterDevice(t, testNewWsService(t), *device, "-1b")
//}

// Test RegisterDevice method (no network)
func TestRegisterDeviceNoNet(t *testing.T) {
	device := testNewDevice()

	testCheckRegisterDevice(t, testNewRestService(t), *device, "-2a")
	testCheckRegisterDevice(t, testNewWsService(t), *device, "-2b")
}

// Test RegisterDevice method
func TestRegisterDevice1(t *testing.T) {
	device := testNewDevice()
	device.Network = testNewNetwork()

	//	if rs := testNewRest(t); rs != nil {
	//		_ = rs.InsertNetwork(device.Network, testWaitTimeout)
	//		// ignore possible errors
	//	}

	testCheckRegisterDevice(t, testNewRestService(t), *device, "-3a")
	testCheckRegisterDevice(t, testNewWsService(t), *device, "-3b")

	//	if rs := testNewRest(t); rs != nil {
	//		_ = rs.DeleteNetwork(device.Network, testWaitTimeout)
	//		// ignore possible errors
	//	}
}

//// Test InsertNetwork method
//func TestInsertNetwork(t *testing.T) {
//	network := testNewNetwork()
//
//	if rs := testNewRest(t); rs != nil {
//		err := rs.InsertNetwork(network, testWaitTimeout)
//		if err != nil {
//			t.Errorf("Failed to create network (error: %s)", err)
//		}
//	}
//}

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
	device := testNewDevice()
	device.Network = testNewNetwork()
	testCheckRegisterDevice1(t, testNewRestService(t), *device, true)
	if t.Failed() {
		return // nothing to test without device
	}

	notification := core.NewNotification("ntf-test", 12345)
	testCheckInsertNotification(t, testNewRestService(t), device, *notification)
	testCheckInsertNotification(t, testNewWsService(t), device, *notification)
}
