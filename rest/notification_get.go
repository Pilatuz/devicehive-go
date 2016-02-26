package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetNotification get the notification data.
func (service *Service) GetNotification(device *devicehive.Device, notificationID uint64, timeout time.Duration) (*devicehive.Notification, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/notification/%d", device.ID, notificationID)

	// result
	notification := &devicehive.Notification{
		ID: notificationID,
	}

	// do GET and check status is 200
	task := newTask("GET", &URL, timeout)
	task.deviceAuth = device
	err := service.do200(task, "/notification/get", nil, notification)
	if err != nil {
		return nil, err
	}

	return notification, nil // OK
}
