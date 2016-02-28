package ws

import (
	dh "github.com/pilatuz/go-devicehive"
)

// SubscribeNotifications starts listening for the notifications.
// /client service only.
func (service *Service) SubscribeNotifications(device *dh.Device, timestamp string) (*dh.NotificationListener, error) {
	if listener := service.findNotificationListener(device.ID); listener != nil {
		return listener, nil // already exists
	}

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "notification/subscribe",
		"requestId": task.identifier,
	}

	// timestamp [optional]
	if len(timestamp) != 0 {
		task.dataToSend["timestamp"] = timestamp
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/notification/subscribe")
	if err != nil {
		return nil, err
	}

	// create listener
	listener := dh.NewNotificationListener(64) // TODO: dedicated variable for buffer size
	service.insertNotificationListener(device.ID, listener)

	return listener, nil // OK
}
