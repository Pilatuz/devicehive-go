package rest

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// GetDevice gets the device data.
func (service *Service) GetDevice(deviceID, deviceKey string) (*dh.Device, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s", deviceID)

	// result
	var res map[string]interface{}
	device := &dh.Device{
		ID:  deviceID,
		Key: deviceKey,
	}

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do200(task, "/device/get", nil, &res)
	if err != nil {
		return nil, err
	}

	// convert map to device
	if err := device.FromMap(res); err != nil {
		return nil, err
	}

	return device, nil // OK
}
