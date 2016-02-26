package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// InsertNotification inserts the device notification.
func (service *Service) InsertNotification(device *devicehive.Device, notification *devicehive.Notification, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/notification", device.ID)

	// request body (do not put all fields)
	body := &devicehive.Notification{
		Name:       notification.Name,
		Parameters: notification.Parameters,
	}

	// do POST and check status is 2xx
	task := newTask("POST", &URL, timeout)
	task.deviceAuth = device
	err := service.do2xx(task, "/notification/insert", &body, notification)
	if err != nil {
		return err
	}

	return nil // OK
}
