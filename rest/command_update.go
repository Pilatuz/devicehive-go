package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// UpdateCommand updates the device command.
func (service *Service) UpdateCommand(device *devicehive.Device, command *devicehive.Command, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/command/%d", device.ID, command.ID)

	// request body (do not put all fields)
	body := &devicehive.Command{
		Status: command.Status,
		Result: command.Result,
	}

	// do PUT and check status is 200
	task := newTask("PUT", &URL, timeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/command/update", &body, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
