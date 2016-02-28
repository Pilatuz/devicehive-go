package ws

import (
	dh "github.com/pilatuz/go-devicehive"
)

// Authenticate authenticates the device.
func (service *Service) Authenticate(device *dh.Device) error {
	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "authenticate",
		"requestId": task.identifier,
		"accessKey": service.accessKey,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/authenticate")
	if err != nil {
		return err
	}

	return nil // OK
}
