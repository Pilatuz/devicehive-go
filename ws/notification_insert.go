package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare InsertNotification task
func (service *Service) prepareInsertNotification(device *core.Device, notification *core.Notification) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "notification/insert",
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	ntf_data := core.Notification{Name: notification.Name,
		Parameters: notification.Parameters}
	task.dataToSend["notification"] = ntf_data

	return
}

// Process InsertNotification task
func (service *Service) processInsertNotification(task *Task, notification *core.Notification) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /notification/insert status (error: %s)", err)
		return
	}

	// parse response
	err = notification.AssignJSON(task.dataRecved["notification"])
	if err != nil {
		log.Warnf("WS: failed to parse /notification/insert response (error: %s)", err)
		return
	}

	return
}

// InsertNotification() function inserts the notification.
func (service *Service) InsertNotification(device *core.Device, notification *core.Notification, timeout time.Duration) (err error) {
	task, err := service.prepareInsertNotification(device, notification)
	if err != nil {
		log.Warnf("WS: failed to prepare /notification/insert task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /notification/insert task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processInsertNotification(task, notification)
		if err != nil {
			log.Warnf("WS: failed to process /notification/insert task (error: %s)", err)
			return
		}
	}

	return
}
