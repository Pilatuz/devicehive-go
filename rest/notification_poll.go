package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"net/url"
	"time"
)

// Prepare PollNotification task
func (service *Service) preparePollNotification(device *core.Device, timestamp, names, waitTimeout string) (task Task, err error) {
	// create request
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
	var url string
	if len(query) != 0 {
		url = fmt.Sprintf("%s/device/%s/notification?%s", service.baseUrl, device.Id, query.Encode())
	} else {
		url = fmt.Sprintf("%s/device/%s/notification", service.baseUrl, device.Id)
	}

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /notification/poll request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process PollNotification task
func (service *Service) processPollNotification(task Task) (notifications []core.Notification, err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /notification/poll status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, &notifications)
	if err != nil {
		log.Warnf("REST: failed to parse /notification/poll body (error: %s)", err)
		return
	}

	return
}

// GetNotification() function poll the notifications.
func (service *Service) PollNotifications(device *core.Device, timestamp, names, waitTimeout string, timeout time.Duration) (notifications []core.Notification, err error) {
	log.Tracef("REST: polling notifications %q...", device.Id)

	task, err := service.preparePollNotification(device, timestamp, names, waitTimeout)
	if err != nil {
		log.Warnf("REST: failed to prepare /notification/poll task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /notification/poll task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		notifications, err = service.processPollNotification(task)
		if err != nil {
			log.Warnf("REST: failed to process /notification/poll task (error: %s)", err)
			return
		}
	}

	return
}
