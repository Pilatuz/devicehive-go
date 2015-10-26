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

// Prepare UpdateCommand task
func (service *Service) prepareUpdateCommand(device *core.Device, command *core.Command) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s/command/%d", service.baseUrl, device.Id, command.Id)

	// do not put some fields to the request body
	command = &core.Command{
		Status: command.Status,
		Result: command.Result}

	body, err := json.Marshal(command)
	if err != nil {
		log.Warnf("REST: failed to format /command/update request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /command/update request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process UpdateCommand task
func (service *Service) processUpdateCommand(task Task, command *core.Command) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /command/update status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	_ = command // unused
	// unmarshal
	//	err = json.Unmarshal(task.body, command)
	//	if err != nil {
	//		log.Warnf("REST: failed to parse /command/update body (error: %s)", err)
	//		return
	//	}

	return
}

// UpdateCommand() function updates the device command.
func (service *Service) UpdateCommand(device *core.Device, command *core.Command, timeout time.Duration) (err error) {
	log.Tracef("REST: updating command %q to %q...", command.Name, device.Id)

	task, err := service.prepareUpdateCommand(device, command)
	if err != nil {
		log.Warnf("REST: failed to prepare /command/update task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /command/update task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processUpdateCommand(task, command)
		if err != nil {
			log.Warnf("REST: failed to process /command/update task (error: %s)", err)
			return
		}
	}

	return
}
