package core

import (
	"encoding/json"
	"testing"
)

// Test DeviceClass JSON marshaling
func TestJsonDeviceClass(t *testing.T) {
	check := func(deviceClass DeviceClass, expectedJson string) {
		jsonBytes, err := json.Marshal(deviceClass)
		if err != nil {
			t.Errorf("DeviceClass: Cannot convert %+v to JSON (error: %s)", deviceClass, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("DeviceClass: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				deviceClass, string(jsonBytes), expectedJson)
		}

		t.Logf("DeviceClass: %+v converted to %s", deviceClass, string(jsonBytes))
	}

	deviceClass := DeviceClass{
		Name:           "eqp-name",
		Version:        "1.2.3",
		OfflineTimeout: 60}

	check(deviceClass, `{"name":"eqp-name","version":"1.2.3","offlineTimeout":60}`)

	deviceClass.Data = "custom data"
	deviceClass.Id = 100
	check(deviceClass, `{"id":100,"name":"eqp-name","version":"1.2.3","offlineTimeout":60,"data":"custom data"}`)
}
