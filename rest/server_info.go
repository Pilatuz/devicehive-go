package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetServerInfo() function gets the main server's information.
func (service *Service) GetServerInfo(timeout time.Duration) (info *devicehive.ServerInfo, err error) {
	task := newTask()
	task.log().Debugf("[%s]: getting server info...", TAG)

	// create request
	url := *service.baseUrl
	url.Path += "/info"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		task.log().Warnf("[%s]: failed to create /info request: %s", TAG, err)
		return nil, fmt.Errorf("failed to create /info request: %s", err)
	}

	// authorization
	service.prepareAuthorization(request, nil)

	select {
	case <-time.After(timeout):
		close(request.Cancel) // cancel request
		task.log().Warnf("[%s]: failed to wait %s for /info response", TAG, timeout)
		return nil, fmt.Errorf("failed to get /info response: timed out (%s)", timeout)

	case task = <-service.doAsync(task):
		info = &core.ServerInfo{}
		err = service.processGetServerInfo(task, info)
		if err != nil {
			log.Warnf("REST: failed to process /info task (error: %s)", err)
			return
		}
	}

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
