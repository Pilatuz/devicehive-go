package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Device to string
// TODO: Device from map

// Test Device JSON marshaling
func TestDeviceJson(t *testing.T) {
	device := NewDevice("dev-id", "dev-name", nil)
	device.Key = "dev-key"
	device.Status = "Online"
	assert.JSONEq(t, toJsonStr(t, device), `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online"}`)

	device.Data = "custom data"
	assert.JSONEq(t, toJsonStr(t, device), `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data"}`)

	device.Network = NewNetwork("net-name", "net-key")
	assert.JSONEq(t, toJsonStr(t, device), `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"}}`)

	device.DeviceClass = NewDeviceClass("class-name", "3.4.5")
	assert.JSONEq(t, toJsonStr(t, device), `{"id":"dev-id","name":"dev-name","key":"dev-key","status":"Online","data":"custom data","network":{"name":"net-name","key":"net-key"},"deviceClass":{"name":"class-name","version":"3.4.5"}}`)
}
