package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare UnsubscribeCommand task
func (service *Service) prepareUnsubscribeCommand(device *core.Device) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "command/unsubscribe",
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	return
}

// Process UnsubscribeCommand task
func (service *Service) processUnsubscribeCommand(task *Task) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /command/unsubscribe status (error: %s)", err)
		return
	}

	return
}

// UnsubscribeCommand() function updates the command.
func (service *Service) UnsubscribeCommands(device *core.Device, timeout time.Duration) (err error) {
	task, err := service.prepareUnsubscribeCommand(device)
	if err != nil {
		log.Warnf("WS: failed to prepare /command/unsubscribe task (error: %s)", err)
		return
	}

	service.removeCommandListener(device.Id)

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /command/unsubscribe task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processUnsubscribeCommand(task)
		if err != nil {
			log.Warnf("WS: failed to process /command/unsubscribe task (error: %s)", err)
			return
		}
	}

	return
}
