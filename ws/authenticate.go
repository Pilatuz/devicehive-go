package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare Authenticate task
func (service *Service) prepareAuthenticate(device *core.Device) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "authenticate",
		"accessKey": service.accessKey,
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	return
}

// Process Authenticate task
func (service *Service) processAuthenticate(task *Task) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /authenticate status (error: %s)", err)
		return
	}

	return
}

// Authenticate() function authenticates the device.
func (service *Service) Authenticate(device *core.Device, timeout time.Duration) (err error) {
	task, err := service.prepareAuthenticate(device)
	if err != nil {
		log.Warnf("WS: failed to prepare /authenticate task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /authenticate task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processAuthenticate(task)
		if err != nil {
			log.Warnf("WS: failed to process /authenticate task (error: %s)", err)
			return
		}
	}

	return
}
