package core

import (
	"encoding/json"
	"testing"
)

// Test Device JSON marshaling
func TestJsonDevice(t *testing.T) {
	check := func(device Device, expectedJson string) {
		jsonBytes, err := json.Marshal(device)
		if err != nil {
			t.Errorf("Device: Cannot convert %+v to JSON (error: %s)", device, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("Device: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				device, string(jsonBytes), expectedJson)
		}

		t.Logf("Device: %+v converted to %s", device, string(jsonBytes))
	}

	device := Device{
		Name:   "eqp-name",
		Key:    "1.2.3",
		Status: "Online"}

	check(device, `{"name":"eqp-name","key":"1.2.3","status":"Online"}`)

	device.Data = "custom data"
	device.Id = "100"
	check(device, `{"id":"100","name":"eqp-name","key":"1.2.3","status":"Online","data":"custom data"}`)

	device.Network = &Network{Name: "net-name", Key: "net-key"}
	check(device, `{"id":"100","name":"eqp-name","key":"1.2.3","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"}}`)

	device.DeviceClass = &DeviceClass{Name: "class-name", Version: "3.4.5"}
	check(device, `{"id":"100","name":"eqp-name","key":"1.2.3","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"},"deviceClass":{"name":"class-name","version":"3.4.5"}}`)
}
