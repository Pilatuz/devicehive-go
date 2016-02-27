package rest

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// InsertCommand inserts the device command.
func (service *Service) InsertCommand(device *dh.Device, command *dh.Command) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/command", device.ID)

	// request body (do not put all fields)
	body := &dh.Command{
		Name:       command.Name,
		Parameters: command.Parameters,
		Lifetime:   command.Lifetime,
	}

	// do POST and check status is 2xx
	task := newTask("POST", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/command/insert", &body, command)
	if err != nil {
		return err
	}

	return nil // OK
}
