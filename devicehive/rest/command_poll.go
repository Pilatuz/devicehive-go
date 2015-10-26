package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"net/url"
	"time"
)

// Prepare PollCommand task
func (service *Service) preparePollCommand(device *core.Device, timestamp, names, waitTimeout string) (task Task, err error) {
	// create request
	query := url.Values{}
	if len(timestamp) != 0 {
		query.Set("timestamp", timestamp)
	}
	if len(names) != 0 {
		query.Set("names", names)
	}
	if len(waitTimeout) != 0 {
		query.Set("waitTimeout", waitTimeout)
	}
	var url string
	if len(query) != 0 {
		url = fmt.Sprintf("%s/device/%s/command?%s", service.baseUrl, device.Id, query.Encode())
	} else {
		url = fmt.Sprintf("%s/device/%s/command", service.baseUrl, device.Id)
	}

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /command/poll request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process PollCommand task
func (service *Service) processPollCommand(task Task) (commands []core.Command, err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /command/poll status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, &commands)
	if err != nil {
		log.Warnf("REST: failed to parse /command/poll body (error: %s)", err)
		return
	}

	return
}

// GetCommand() function poll the commands.
func (service *Service) PollCommands(device *core.Device, timestamp, names, waitTimeout string, timeout time.Duration) (commands []core.Command, err error) {
	log.Tracef("REST: polling commands %q...", device.Id)

	task, err := service.preparePollCommand(device, timestamp, names, waitTimeout)
	if err != nil {
		log.Warnf("REST: failed to prepare /command/poll task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /command/poll task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		commands, err = service.processPollCommand(task)
		if err != nil {
			log.Warnf("REST: failed to process /command/poll task (error: %s)", err)
			return
		}
	}

	return
}
