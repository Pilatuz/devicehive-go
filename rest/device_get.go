package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetDevice gets the device data.
func (service *Service) GetDevice(deviceID, deviceKey string, timeout time.Duration) (*devicehive.Device, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s", deviceID)

	// result
	var res map[string]interface{}
	device := new(devicehive.Device)
	device.ID = deviceID
	device.Key = deviceKey

	// do GET and check status is 200
	task := newTask("GET", &URL, timeout)
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
