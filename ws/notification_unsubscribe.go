package ws

import (
	dh "github.com/pilatuz/devicehive-go"
)

// UnsubscribeNotifications removes the notification listener.
// /client service only.
func (service *Service) UnsubscribeNotifications(device *dh.Device) error {
	if listener := service.findNotificationListener(device.ID); listener == nil {
		return nil // nothing to unsubscribe
	}

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "notification/unsubscribe",
		"requestId": task.identifier,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, "/notification/unsubscribe")
	if err != nil {
		return err
	}

	service.removeNotificationListener(device.ID)
	return nil // OK
}
