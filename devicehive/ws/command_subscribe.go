package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare SubscribeCommand task
func (service *Service) prepareSubscribeCommand(device *core.Device, timestamp string) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "command/subscribe",
		"requestId": task.id}

	// timestamp [optional]
	if len(timestamp) != 0 {
		task.dataToSend["timestamp"] = timestamp
	}

	// prepare authorization
	task.prepareAuthorization(device)

	return
}

// Process SubscribeCommand task
func (service *Service) processSubscribeCommand(task *Task) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /command/subscribe status (error: %s)", err)
		return
	}

	return
}

// SubscribeCommand() function updates the command.
func (service *Service) SubscribeCommands(device *core.Device, timestamp string, timeout time.Duration) (listener *core.CommandListener, err error) {
	task, err := service.prepareSubscribeCommand(device, timestamp)
	if err != nil {
		log.Warnf("WS: failed to prepare /command/subscribe task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /command/subscribe task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processSubscribeCommand(task)
		if err != nil {
			log.Warnf("WS: failed to process /command/subscribe task (error: %s)", err)
			return
		}

		// done, create listener
		listener = core.NewCommandListener()
		service.insertCommandListener(device.Id, listener)
	}

	return
}
