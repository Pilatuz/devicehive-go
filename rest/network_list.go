// +build ignore

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// Prepare GetNetworkList task
func (service *Service) prepareGetNetworkList(take, skip int) (task Task, err error) {
	// create request
	query := url.Values{}
	if take > 0 {
		query.Set("take", fmt.Sprintf("%d", take))
	}
	if skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", skip))
	}
	url := fmt.Sprintf("%s/network", service.baseUrl)
	if len(query) != 0 {
		url += "?" + query.Encode()
	}

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /network/list request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process GetNetworkList task
func (service *Service) processGetNetworkList(task Task) (networks []core.Network, err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /network/list status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, &networks)
	if err != nil {
		log.Warnf("REST: failed to parse /network/list body (error: %s)", err)
		return
	}

	return
}

// GetNetworkList() function get the network list.
func (service *Service) GetNetworkList(take, skip int, timeout time.Duration) (networks []core.Network, err error) {
	log.Tracef("REST: getting network list (take:%d, skip:%d)...", take, skip)

	task, err := service.prepareGetNetworkList(take, skip)
	if err != nil {
		log.Warnf("REST: failed to prepare /network/list task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /network/list task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		networks, err = service.processGetNetworkList(task)
		if err != nil {
			log.Warnf("REST: failed to process /network/list task (error: %s)", err)
			return
		}
	}

	return
}
