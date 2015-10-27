package devicehive

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"testing"
	"time"
)

// Test datetime layout
func TestTimestampFormat(t *testing.T) {
	str := "2015-10-22T14:15:16.999"
	ts, err := time.Parse(core.DateTimeLayout, str)
	if err != nil {
		t.Errorf("DateTime: Cannot parse timestamp %s (error: %s)", str, err)
	}

	if ts.Year() != 2015 ||
		ts.Month() != time.October ||
		ts.Day() != 22 ||
		ts.Hour() != 14 ||
		ts.Minute() != 15 ||
		ts.Second() != 16 ||
		ts.Nanosecond() != 999000000 {
		t.Errorf("DateTime: Cannot parse timestamp %s (wrong data parsed)", str)
	}
}

// check JSON marshaling
func testCheckJson(t *testing.T, info interface{}, expectedJson string) {
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		t.Errorf("Cannot convert %+v to JSON (error: %s)", info, err)
	}
	if expectedJson != string(jsonBytes) {
		t.Errorf("Cannot convert %+v to JSON\n\t   found:%s,\n\texpected:%s",
			info, string(jsonBytes), expectedJson)
	}

	// TODO: check AssignJSON()

	//t.Logf("%+v converted to %s", info, string(jsonBytes))
}

// Test ServiceInfo JSON marshaling
func TestJsonServiceInfo(t *testing.T) {
	info := core.ServerInfo{
		Version:      "1.2.3",
		Timestamp:    "2015-10-22T14:15:16.999",
		WebsocketUrl: "ws://devicehive.com"}
	testCheckJson(t, info, `{"apiVersion":"1.2.3","serverTimestamp":"2015-10-22T14:15:16.999","webSocketServerUrl":"ws://devicehive.com"}`)

	info.WebsocketUrl = ""
	info.RestUrl = "https://devicehive.com"
	testCheckJson(t, info, `{"apiVersion":"1.2.3","serverTimestamp":"2015-10-22T14:15:16.999","restServerUrl":"https://devicehive.com"}`)
}

// Test Equipment JSON marshaling
func TestJsonEquipment(t *testing.T) {
	equipment := core.NewEquipment("eqp-name", "eqp-code", "eqp-type")
	testCheckJson(t, equipment, `{"name":"eqp-name","code":"eqp-code","type":"eqp-type"}`)

	equipment.Data = "custom data"
	equipment.Id = 100
	testCheckJson(t, equipment, `{"id":100,"name":"eqp-name","code":"eqp-code","type":"eqp-type","data":"custom data"}`)
}

// Test DeviceClass JSON marshaling
func TestJsonDeviceClass(t *testing.T) {
	deviceClass := core.NewDeviceClass("class-name", "1.2.3")
	deviceClass.OfflineTimeout = 60
	testCheckJson(t, deviceClass, `{"name":"class-name","version":"1.2.3","offlineTimeout":60}`)

	deviceClass.Data = "custom data"
	deviceClass.Id = 100
	testCheckJson(t, deviceClass, `{"id":100,"name":"class-name","version":"1.2.3","offlineTimeout":60,"data":"custom data"}`)
}

// Test Network JSON marshaling
func TestJsonNetwork(t *testing.T) {
	network := core.NewNetwork("net-name", "net-key")
	network.Description = "custom description"
	testCheckJson(t, network, `{"name":"net-name","key":"net-key","description":"custom description"}`)

	network.Description = ""
	network.Id = 100
	testCheckJson(t, network, `{"id":100,"name":"net-name","key":"net-key"}`)
}

// Test Device JSON marshaling
func TestJsonDevice(t *testing.T) {
	device := core.NewDevice("dev-id", "dev-name", nil)
	device.Key = "dev-key"
	device.Status = "Online"
	testCheckJson(t, device, `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online"}`)

	device.Data = "custom data"
	testCheckJson(t, device, `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data"}`)

	device.Network = core.NewNetwork("net-name", "net-key")
	testCheckJson(t, device, `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"}}`)

	device.DeviceClass = core.NewDeviceClass("class-name", "3.4.5")
	testCheckJson(t, device, `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"},"deviceClass":{"name":"class-name","version":"3.4.5"}}`)
}

// Test Command JSON marshaling
func TestJsonCommand(t *testing.T) {
	command := core.NewCommand("cmd-name", "hello")
	command.Timestamp = "2005-10-22"
	testCheckJson(t, command, `{"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello"}`)

	command.Id = 100
	command.Result = "custom data"
	command.Status = "done"
	testCheckJson(t, command, `{"id":100,"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello","status":"done","result":"custom data"}`)
}

// Test Notification JSON marshaling
func TestJsonNotification(t *testing.T) {
	notification := core.NewNotification("ntf-name", "hello")
	testCheckJson(t, notification, `{"notification":"ntf-name","parameters":"hello"}`)

	notification.Id = 100
	testCheckJson(t, notification, `{"id":100,"notification":"ntf-name","parameters":"hello"}`)
}
