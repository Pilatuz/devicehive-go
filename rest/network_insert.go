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

// Prepare InsertNetwork task
func (service *Service) prepareInsertNetwork(network *core.Network) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/network", service.baseUrl)

	net_data := *network
	net_data.Id = 0 // do not put ID to the request body
	body, err := json.Marshal(&net_data)
	if err != nil {
		log.Warnf("REST: failed to format /network/insert request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /network/insert request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process InsertNetwork task
func (service *Service) processInsertNetwork(task Task, network *core.Network) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /network/insert status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, network)
	if err != nil {
		log.Warnf("REST: failed to parse /network/insert body (error: %s)", err)
		return
	}

	return
}

// InsertNetwork() function inserts the network.
func (service *Service) InsertNetwork(network *core.Network, timeout time.Duration) (err error) {
	log.Tracef("REST: inserting network %q...", network.Name)

	task, err := service.prepareInsertNetwork(network)
	if err != nil {
		log.Warnf("REST: failed to prepare /network/insert task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /network/insert task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processInsertNetwork(task, network)
		if err != nil {
			log.Warnf("REST: failed to process /network/insert task (error: %s)", err)
			return
		}
	}

	return
}
