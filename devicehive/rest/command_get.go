package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare GetCommand task
func (service *Service) prepareGetCommand(device *core.Device, commandId uint64) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s/command/%d", service.baseUrl, device.Id, commandId)

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /command/get request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process GetCommand task
func (service *Service) processGetCommand(task Task, command *core.Command) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /command/get status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, command)
	if err != nil {
		log.Warnf("REST: failed to parse /command/get body (error: %s)", err)
		return
	}

	return
}

// GetCommand() function get the command data.
func (service *Service) GetCommand(device *core.Device, commandId uint64, timeout time.Duration) (command *core.Command, err error) {
	log.Tracef("REST: getting command %q/%d...", device.Id, commandId)

	task, err := service.prepareGetCommand(device, commandId)
	if err != nil {
		log.Warnf("REST: failed to prepare /command/get task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /command/get task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		command = &core.Command{Id: commandId}
		err = service.processGetCommand(task, command)
		if err != nil {
			log.Warnf("REST: failed to process /command/get task (error: %s)", err)
			return
		}
	}

	return
}
