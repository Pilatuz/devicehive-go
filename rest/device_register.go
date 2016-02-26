package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// RegisterDevice registers the device.
func (service *Service) RegisterDevice(device *devicehive.Device, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s", device.ID)

	// request body
	body := *device // copy
	body.ID = ""    // do not put ID to the request body

	// do PUT and check status is 2xx
	task := newTask("PUT", &URL, timeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/device/register", &body, device)
	if err != nil {
		return err
	}

	return nil // OK
}
