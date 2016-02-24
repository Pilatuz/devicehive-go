package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetServerInfo() function gets the main server's information.
func (service *Service) GetServerInfo(timeout time.Duration) (*devicehive.ServerInfo, error) {
	task := newTask()
	task.log().Debugf("[%s]: getting /info...", TAG)

	// create request
	var err error
	url := *service.baseUrl
	url.Path += "/info"
	task.request, err = http.NewRequest("GET", url.String(), nil)
	if err != nil {
		task.log().Warnf("[%s]: failed to create /info request: %s", TAG, err)
		return nil, fmt.Errorf("failed to create /info request: %s", err)
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	select {
	case <-time.After(timeout):
		// TODO: task.request.Cancel // cancel request
		task.log().Warnf("[%s]: failed to wait /info response: timed out (%s)", TAG, timeout)
		return nil, fmt.Errorf("failed to wait /info response: timed out (%s)", timeout)

	case err = <-service.doAsync(task):
		if err != nil {
			task.log().Warnf("[%s]: failed to get /info response: %s", TAG, err)
			return nil, fmt.Errorf("failed to get /info response: %s", err)
		}
	}

	// important to close response body!
	defer task.response.Body.Close()

	// check status code
	if task.response.StatusCode != http.StatusOK {
		task.log().Warnf("[%s]: unexpected /info status: %s", TAG, task.response.Status)
		return nil, fmt.Errorf("unexpected /info status: %s", task.response.Status)
	}

	// unmarshal
	info := new(devicehive.ServerInfo)
	dec := json.NewDecoder(task.response.Body)
	err = dec.Decode(info)
	if err != nil {
		task.log().Warnf("[%s]: failed to parse /info body: %s", TAG, err)
		return nil, fmt.Errorf("failed to parse /info body: %s", err)
	}

	task.log().Infof("[%s]: parsed /info body: %s", TAG, info)
	return info, nil // OK
}
