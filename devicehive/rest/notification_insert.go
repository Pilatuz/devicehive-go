package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare InsertNotification task
func (service *Service) prepareInsertNotification(device *core.Device, notification *core.Notification) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s/notification", service.baseUrl, device.Id)

	// do not put some fields to the request body
	notification = &core.Notification{Name: notification.Name,
		Parameters: notification.Parameters}

	body, err := json.Marshal(notification)
	if err != nil {
		log.Warnf("REST: failed to format /notification/insert request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /notification/insert request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process InsertNotification task
func (service *Service) processInsertNotification(task Task, notification *core.Notification) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /notification/insert status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, notification)
	if err != nil {
		log.Warnf("REST: failed to parse /notification/insert body (error: %s)", err)
		return
	}

	return
}

// InsertNotification() function inserts the device notification.
func (service *Service) InsertNotification(device *core.Device, notification *core.Notification, timeout time.Duration) (err error) {
	log.Tracef("REST: inserting notification %q to %q...", notification.Name, device.Id)

	task, err := service.prepareInsertNotification(device, notification)
	if err != nil {
		log.Warnf("REST: failed to prepare /notification/insert task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /notification/insert task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processInsertNotification(task, notification)
		if err != nil {
			log.Warnf("REST: failed to process /notification/insert task (error: %s)", err)
			return
		}
	}

	return
}
