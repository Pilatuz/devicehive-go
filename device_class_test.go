package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: DeviceClass to string
// TODO: DeviceClass from map

// Test DeviceClass JSON marshaling
func TestDeviceClassJson(t *testing.T) {
	deviceClass := NewDeviceClass("class-name", "1.2.3")
	deviceClass.OfflineTimeout = 60
	assert.JSONEq(t, toJsonStr(t, deviceClass), `{"name":"class-name","version":"1.2.3","offlineTimeout":60}`)

	deviceClass.Data = "custom data"
	deviceClass.ID = 100
	assert.JSONEq(t, toJsonStr(t, deviceClass), `{"id":100,"name":"class-name","version":"1.2.3","offlineTimeout":60,"data":"custom data"}`)
}
