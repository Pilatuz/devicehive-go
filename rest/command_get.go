package rest

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// GetCommand gets the command data.
func (service *Service) GetCommand(device *dh.Device, commandID uint64) (*dh.Command, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/command/%d", device.ID, commandID)

	// result
	command := &dh.Command{
		ID: commandID,
	}

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do200(task, "/command/get", nil, command)
	if err != nil {
		return nil, err
	}

	return command, nil // OK
}
