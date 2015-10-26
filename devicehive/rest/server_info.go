package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare GetServerInfo task
func (service *Service) prepareGetServerInfo() (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/info", service.baseUrl)
	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /info request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process GetServerInfo task
func (service *Service) processGetServerInfo(task Task, info *core.ServerInfo) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /info status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, info)
	if err != nil {
		log.Warnf("REST: failed to parse /info body (error: %s)", err)
		return
	}

	return
}

// GetServerInfo() function gets the main server's information.
func (service *Service) GetServerInfo(timeout time.Duration) (info *core.ServerInfo, err error) {
	log.Tracef("REST: getting server info...")

	task, err := service.prepareGetServerInfo()
	if err != nil {
		log.Warnf("REST: failed to prepare /info task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /info task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		info = &core.ServerInfo{}
		err = service.processGetServerInfo(task, info)
		if err != nil {
			log.Warnf("REST: failed to process /info task (error: %s)", err)
			return
		}
	}

	return
}
