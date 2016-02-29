package ws

import (
	dh "github.com/pilatuz/devicehive-go"
)

// SubscribeCommand starts listening for the commands.
func (service *Service) SubscribeCommands(device *dh.Device, timestamp string) (*dh.CommandListener, error) {
	if listener := service.findCommandListener(device.ID); listener != nil {
		return listener, nil // already exists
	}

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "command/subscribe",
		"requestId": task.identifier,
	}

	// timestamp [optional]
	if len(timestamp) != 0 {
		task.dataToSend["timestamp"] = timestamp
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/command/subscribe")
	if err != nil {
		return nil, err
	}

	// create listener
	listener := dh.NewCommandListener(64) // TODO: dedicated variable for buffer size
	service.insertCommandListener(device.ID, listener)

	return listener, nil // OK
}
