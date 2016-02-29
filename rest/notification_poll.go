package rest

import (
	"fmt"
	"net/url"

	dh "github.com/pilatuz/devicehive-go"
)

// PollNotifications polls the notifications.
func (service *Service) PollNotifications(device *dh.Device, timestamp, names, waitTimeout string) ([]*dh.Notification, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/notification/poll", device.ID)
	query := url.Values{}
	if len(timestamp) != 0 {
		query.Set("timestamp", timestamp)
	}
	if len(names) != 0 {
		query.Set("names", names)
	}
	if len(waitTimeout) != 0 {
		query.Set("waitTimeout", waitTimeout)
	}
	URL.RawQuery = query.Encode()

	// result
	var notifications []*dh.Notification

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do200(task, "/notification/poll", nil, &notifications)
	if err != nil {
		return nil, err
	}

	// convert map to notifications
	//	notifications := make([]*dh.Notification, 0, len(res))
	//	for _, data := range res {
	//		n := new(dh.Notification)
	//		if err := n.FromMap(data); err != nil {
	//			return nil, err
	//		}
	//		notifications = append(notifications, n)
	//	}

	return notifications, nil // OK
}
