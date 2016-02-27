package rest

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// InsertNotification inserts the device notification.
func (service *Service) InsertNotification(device *dh.Device, notification *dh.Notification) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/notification", device.ID)

	// request body (do not put all fields)
	body := &dh.Notification{
		Name:       notification.Name,
		Parameters: notification.Parameters,
	}

	// do POST and check status is 2xx
	task := newTask("POST", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/notification/insert", &body, notification)
	if err != nil {
		return err
	}

	return nil // OK
}
