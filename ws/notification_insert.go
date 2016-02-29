package ws

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// InsertNotification() function inserts the notification.
func (service *Service) InsertNotification(device *dh.Device, notification *dh.Notification) error {
	const OP = "/notification/insert"

	// request data (do not put all fields)
	data := &dh.Notification{
		Name:       notification.Name,
		Parameters: notification.Parameters,
	}

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":       "notification/insert",
		"requestId":    task.identifier,
		"notification": data,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, OP)
	if err != nil {
		return err
	}

	// parse response
	err = notification.FromMap(task.dataReceived["notification"])
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to parse %s response", TAG, OP)
		return fmt.Errorf("failed to parse %s response: %s", OP, err)
	}

	return nil // OK
}
