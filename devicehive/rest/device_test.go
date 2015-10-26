package rest

import (
	"github.com/devicehive/devicehive-go/devicehive/core"
	"testing"
)

// TestRegisterDevice() unit test for /device/register PUT method
func TestRegisterDevice(t *testing.T) {
	s, err := NewService(testServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	device := &core.Device{Id: testDeviceId, Key: testDeviceKey, Name: "test-name", Status: "Online"}
	device.DeviceClass = &core.DeviceClass{Name: "go-device-class", Version: "1.2.3"}

	err = s.RegisterDevice(device, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to register device (error: %s)", err)
		return
	}
}

// TestGetDevice() unit test for /device/get GET method
func TestGetDevice(t *testing.T) {
	s, err := NewService(testServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	deviceKey := ""
	device, err := s.GetDevice(testDeviceId, deviceKey, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get device (error: %s)", err)
		return
	}
	t.Logf("device: %s", device)
}

// TestDeleteDevice() unit test for /device/delete DELETE method
func TestDeleteDevice(t *testing.T) {
	s, err := NewService(testServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	device := &core.Device{Id: testDeviceId, Key: testDeviceKey}
	err = s.DeleteDevice(device, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get device (error: %s)", err)
		return
	}
}
