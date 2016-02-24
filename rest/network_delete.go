// +build ignore

package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// Prepare DeleteNetwork task
func (service *Service) prepareDeleteNetwork(network *core.Network) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/network/%d", service.baseUrl, network.Id)
	task.request, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /network/delete request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process DeleteNetwork task
func (service *Service) processDeleteNetwork(task Task) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /network/delete status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	return
}

// DeleteNetwork() function deletes the network.
func (service *Service) DeleteNetwork(network *core.Network, timeout time.Duration) (err error) {
	log.Tracef("REST: deleting network %d...", network.Id)

	task, err := service.prepareDeleteNetwork(network)
	if err != nil {
		log.Warnf("REST: failed to prepare /network/delete task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /network/delete task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processDeleteNetwork(task)
		if err != nil {
			log.Warnf("REST: failed to process /network/delete task (error: %s)", err)
			return
		}
	}

	return
}
