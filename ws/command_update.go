package ws

import (
	dh "github.com/pilatuz/go-devicehive"
)

// CommandUpdate updates the command result.
func (service *Service) UpdateCommand(device *dh.Device, command *dh.Command) error {
	// request data (do not put all fields)
	data := *command // deep copy
	data.ID = 0      // do not put Id inside

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "command/update",
		"requestId": task.identifier,
		"commandId": command.ID,
		"command":   data,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/command/update")
	if err != nil {
		return err
	}

	return nil // OK
}
