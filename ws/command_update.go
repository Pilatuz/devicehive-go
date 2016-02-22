package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare UpdateCommand task
func (service *Service) prepareUpdateCommand(device *core.Device, command *core.Command) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "command/update",
		"commandId": command.Id,
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	cmd_data := *command // deep copy
	cmd_data.Id = 0      // do not put Id inside
	task.dataToSend["command"] = cmd_data

	return
}

// Process UpdateCommand task
func (service *Service) processUpdateCommand(task *Task) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /command/update status (error: %s)", err)
		return
	}

	return
}

// CommandUpdate() function updates the command.
func (service *Service) UpdateCommand(device *core.Device, command *core.Command, timeout time.Duration) (err error) {
	task, err := service.prepareUpdateCommand(device, command)
	if err != nil {
		log.Warnf("WS: failed to prepare /command/update task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /command/update task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processUpdateCommand(task)
		if err != nil {
			log.Warnf("WS: failed to process /command/update task (error: %s)", err)
			return
		}
	}

	return
}
