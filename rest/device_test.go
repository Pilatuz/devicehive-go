package rest

import (
	"testing"

	//dh "github.com/pilatuz/devicehive-go"
	"github.com/stretchr/testify/assert"
)

// Test GetDeviceList and GetDevice methods
func TestDeviceListAndGet(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	devices, err := service.GetDeviceList(0, 0)
	assert.NoError(t, err, "Failed to get list of devices")
	// assert.NotEmpty(t, devices, "No any device available")
	//	for i, d := range devices {
	//		t.Logf("device-%d: %s", i, d)
	//	}

	for i, a := range devices {
		b, err := service.GetDevice(a.ID, a.Key)
		assert.NoError(t, err, "Failed to get device")
		assert.NotNil(t, b, "No device available")
		t.Logf("device-%d/A: %s", i, a)
		t.Logf("device-%d/B: %s", i, b)
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "Devices are not the same")
	}
}

// Test RegisterDevice and DeleteDevice methods
func TestDeviceRegisterAndDelete(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-rest"
	device.Name += "-rest"

	err := service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		t.Logf("device registered: %s", device)

		devices, err := service.GetDeviceList(0, 0)
		assert.NoError(t, err, "Failed to get list of devices")
		for _, d := range devices {
			if d.ID == device.ID {
				err = service.DeleteDevice(device)
				assert.NoError(t, err, "Failed to delete device")
				return // OK
			}
		}

		assert.Fail(t, "No new device found in the device list")
	}
}
