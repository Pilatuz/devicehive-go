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

// Prepare InsertCommand task
func (service *Service) prepareInsertCommand(device *core.Device, command *core.Command) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s/command", service.baseUrl, device.Id)

	// do not put some fields to the request body
	command = &core.Command{Name: command.Name,
		Parameters: command.Parameters,
		Lifetime:   command.Lifetime}

	body, err := json.Marshal(command)
	if err != nil {
		log.Warnf("REST: failed to format /command/insert request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /command/insert request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process InsertCommand task
func (service *Service) processInsertCommand(task Task, command *core.Command) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /command/insert status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, command)
	if err != nil {
		log.Warnf("REST: failed to parse /command/insert body (error: %s)", err)
		return
	}

	return
}

// InsertCommand() function inserts the device command.
func (service *Service) InsertCommand(device *core.Device, command *core.Command, timeout time.Duration) (err error) {
	log.Tracef("REST: inserting command %q to %q...", command.Name, device.Id)

	task, err := service.prepareInsertCommand(device, command)
	if err != nil {
		log.Warnf("REST: failed to prepare /command/insert task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /command/insert task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processInsertCommand(task, command)
		if err != nil {
			log.Warnf("REST: failed to process /command/insert task (error: %s)", err)
			return
		}
	}

	return
}
