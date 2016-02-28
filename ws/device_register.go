package ws

import (
	dh "github.com/pilatuz/go-devicehive"
)

// RegisterDevice registers the device.
func (service *Service) RegisterDevice(device *dh.Device) error {
	const OP = "/device/register"

	data := *device // deep copy
	data.ID = ""    // do not put Id inside

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "device/save",
		"requestId": task.identifier,
		"device":    data,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, OP)
	if err != nil {
		return err
	}

	return nil // OK
}
