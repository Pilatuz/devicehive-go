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

// Prepare UpdateNetwork task
func (service *Service) prepareUpdateNetwork(network *core.Network) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/network/%d", service.baseUrl, network.Id)

	net_data := *network
	net_data.Id = 0 // do not put ID to the request body
	body, err := json.Marshal(&net_data)
	if err != nil {
		log.Warnf("REST: failed to format /network/update request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /network/update request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process UpdateNetwork task
func (service *Service) processUpdateNetwork(task Task, network *core.Network) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /network/update status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	_ = network // unused
	// unmarshal
	//	err = json.Unmarshal(task.body, network)
	//	if err != nil {
	//		log.Warnf("REST: failed to parse /network/update body (error: %s)", err)
	//		return
	//	}

	return
}

// UpdateNetwork() function updates the network.
func (service *Service) UpdateNetwork(network *core.Network, timeout time.Duration) (err error) {
	log.Tracef("REST: updating network %q...", network.Id)

	task, err := service.prepareUpdateNetwork(network)
	if err != nil {
		log.Warnf("REST: failed to prepare /network/update task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /network/update task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processUpdateNetwork(task, network)
		if err != nil {
			log.Warnf("REST: failed to process /network/update task (error: %s)", err)
			return
		}
	}

	return
}
