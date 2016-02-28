package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test RegisterDevice and GetDevice methods
func TestDeviceRegisterAndGet(t *testing.T) {
	service := testNewWsDevice(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-ws"
	device.Name += "-ws"

	err := service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		t.Logf("device registered: %s", device)

		err = service.Authenticate(device)
		assert.NoError(t, err, "Failed to authenticate device")

		b, err := service.GetDevice(device.ID, device.Key)
		assert.NoError(t, err, "Failed to get device")
		if assert.NotNil(t, b, "No device available") {
			device.Network.Key = ""
			b.DeviceClass.ID = 0
			b.Network.ID = 0
			//t.Logf("device/A: %s", device)
			//t.Logf("device/B: %s", b)
			assert.JSONEq(t, toJsonStr(device), toJsonStr(b), "Devices are not the same")
		}
	}
}
