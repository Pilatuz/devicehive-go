package rest

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// UpdateCommand updates the device command.
func (service *Service) UpdateCommand(device *dh.Device, command *dh.Command) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/command/%d", device.ID, command.ID)

	// request body (do not put all fields)
	body := &dh.Command{
		Status: command.Status,
		Result: command.Result,
	}

	// do PUT and check status is 200
	task := newTask("PUT", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/command/update", &body, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
