package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetServerInfo function gets the main server's information.
func (service *Service) GetServerInfo(timeout time.Duration) (*devicehive.ServerInfo, error) {
	var err error
	task := newTask()
	task.log().Debugf("[%s]: getting /info...", TAG)

	// create request
	url := *service.baseURL
	url.Path += "/info"
	task.request, err = http.NewRequest("GET", url.String(), nil)
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to create /info request", TAG)
		return nil, fmt.Errorf("failed to create /info request: %s", err)
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	select {
	case <-time.After(timeout):
		// TODO: task.request.Cancel // cancel request
		task.log().WithField("timeout", timeout).Warnf("[%s]: failed to wait /info response: timed out", TAG)
		return nil, fmt.Errorf("failed to wait /info response: timed out (%s)", timeout)

	case err = <-service.doAsync(task):
		if err != nil {
			task.log().WithError(err).Warnf("[%s]: failed to get /info response", TAG)
			return nil, fmt.Errorf("failed to get /info response: %s", err)
		}
	}

	// important to close response body!
	defer task.response.Body.Close()

	// check status code
	if task.response.StatusCode != http.StatusOK {
		task.log().WithField("status", task.response.Status).Warnf("[%s]: unexpected /info status", TAG)
		return nil, fmt.Errorf("unexpected /info status: %s", task.response.Status)
	}

	// unmarshal
	info := new(devicehive.ServerInfo)
	dec := json.NewDecoder(task.response.Body)
	err = dec.Decode(info)
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to parse /info body", TAG)
		return nil, fmt.Errorf("failed to parse /info body: %s", err)
	}

	task.log().WithField("info", info).Infof("[%s]: parsed /info body", TAG)
	return info, nil // OK
}
