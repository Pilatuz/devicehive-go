package ws

import (
	"testing"

	dh "github.com/pilatuz/go-devicehive"
	"github.com/stretchr/testify/assert"
)

// Test RegisterDevice and GetDevice methods
func TestDeviceRegisterAndGet(t *testing.T) {
	service := testNewWS(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	device := dh.NewDevice("go-unit-test-device-ws", "go test device ws",
		dh.NewDeviceClass("go-test-deviceclass-ws", "0.0.1"))
	err := service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		t.Logf("device registered: %s", device)

		b, err := service.GetDevice(device.ID, device.Key)
		assert.NoError(t, err, "Failed to get device")
		if assert.NotNil(t, b, "No device available") {
			//t.Logf("device-%d/A: %s", i, device)
			//t.Logf("device-%d/B: %s", i, b)
			assert.JSONEq(t, toJsonStr(device), toJsonStr(b), "Devices are not the same")
		}
	}
}
