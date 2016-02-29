package ws

import (
	dh "github.com/pilatuz/devicehive-go"
)

// UnsubscribeCommand removes the command listener
func (service *Service) UnsubscribeCommands(device *dh.Device) error {
	if listener := service.findCommandListener(device.ID); listener == nil {
		return nil // nothing to unsubscribe
	}

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "command/unsubscribe",
		"requestId": task.identifier,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/command/unsubscribe")
	if err != nil {
		return err
	}

	service.removeCommandListener(device.ID)
	return nil // OK
}
