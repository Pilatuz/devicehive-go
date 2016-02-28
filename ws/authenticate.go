package ws

import (
	dh "github.com/pilatuz/go-devicehive"
)

// Authenticate authenticates the device.
func (service *Service) Authenticate(device *dh.Device) error {
	const OP = "/authenticate"

	// Prepare Authenticate task
	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "authenticate",
		"accessKey": service.accessKey,
		"requestId": task.identifier,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, OP)
	if err != nil {
		return err
	}

	return nil // OK
}
