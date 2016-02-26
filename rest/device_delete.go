package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// DeleteDevice deletes the device.
func (service *Service) DeleteDevice(device *devicehive.Device, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s", device.ID)

	// do DELETE and check status is 2xx
	task := newTask("DELETE", &URL, timeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/device/delete", nil, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
