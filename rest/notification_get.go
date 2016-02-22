package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare GetNotification task
func (service *Service) prepareGetNotification(device *core.Device, notificationId uint64) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s/notification/%d", service.baseUrl, device.Id, notificationId)

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /notification/get request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process GetNotification task
func (service *Service) processGetNotification(task Task, notification *core.Notification) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /notification/get status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, notification)
	if err != nil {
		log.Warnf("REST: failed to parse /notification/get body (error: %s)", err)
		return
	}

	return
}

// GetNotification() function get the notification data.
func (service *Service) GetNotification(device *core.Device, notificationId uint64, timeout time.Duration) (notification *core.Notification, err error) {
	log.Tracef("REST: getting notification %q/%d...", device.Id, notificationId)

	task, err := service.prepareGetNotification(device, notificationId)
	if err != nil {
		log.Warnf("REST: failed to prepare /notification/get task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /notification/get task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		notification = &core.Notification{Id: notificationId}
		err = service.processGetNotification(task, notification)
		if err != nil {
			log.Warnf("REST: failed to process /notification/get task (error: %s)", err)
			return
		}
	}

	return
}
