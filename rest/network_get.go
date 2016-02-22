package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare GetNetwork task
func (service *Service) prepareGetNetwork(networkId uint64) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/network/%d", service.baseUrl, networkId)

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /network/get request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process GetNetwork task
func (service *Service) processGetNetwork(task Task, network *core.Network) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /network/get status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, network)
	if err != nil {
		log.Warnf("REST: failed to parse /network/get body (error: %s)", err)
		return
	}

	return
}

// GetNetwork() function get the network data.
func (service *Service) GetNetwork(networkId uint64, timeout time.Duration) (network *core.Network, err error) {
	log.Tracef("REST: getting network %d...", networkId)

	task, err := service.prepareGetNetwork(networkId)
	if err != nil {
		log.Warnf("REST: failed to prepare /network/get task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /network/get task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		network = &core.Network{Id: networkId}
		err = service.processGetNetwork(task, network)
		if err != nil {
			log.Warnf("REST: failed to process /network/get task (error: %s)", err)
			return
		}
	}

	return
}
