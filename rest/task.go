package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	dh "github.com/pilatuz/go-devicehive"
)

var (
	taskID = uint64(0) // unique identifier generator
)

// Asynchronous request/task.
type asyncTask struct {
	identifier uint64
	timeout    time.Duration
	deviceAuth *dh.Device // is used for device authentication

	method   string
	URL      *url.URL
	request  *http.Request
	response *http.Response
}

// create new async task.
// assigns new unique identifier
func newTask(method string, URL *url.URL, timeout time.Duration) *asyncTask {
	task := new(asyncTask)
	task.identifier = atomic.AddUint64(&taskID, 1)
	task.timeout = timeout
	task.method = method
	task.URL = URL
	return task
}

// Do a request/task synchronously.
func (service *Service) doSync(task *asyncTask) error {
	task.log().WithField("url", task.request.URL).
		Debugf("[%s]: sending %s request", TAG, task.request.Method)

	var err error
	task.response, err = service.client.Do(task.request)
	if err != nil {
		return err
	}

	task.log().WithField("status", task.response.Status).Debugf("[%s]: got response", TAG)
	return nil // OK
}

// Do a request/task asynchronously.
func (service *Service) doAsync(task *asyncTask) <-chan error {
	ch := make(chan error, 1)

	go func() {
		ch <- service.doSync(task)
	}()

	return ch
}

// send request and parse response (expectes 200 status code)
func (service *Service) do200(task *asyncTask, OP string, body interface{}, result interface{}) error {
	return service.do(task, OP, body, result, func(status int) bool { return status == http.StatusOK })
}

// send request and parse response (expectes 200..206 status code)
func (service *Service) do2xx(task *asyncTask, OP string, body interface{}, result interface{}) error {
	return service.do(task, OP, body, result, func(status int) bool { return status >= http.StatusOK && status <= 299 })
}

// send request and parse response
func (service *Service) do(task *asyncTask, OP string, body interface{}, result interface{},
	checkStatus func(status int) bool) (err error) {
	// do nothing if service is stopped
	if service.isStopped() {
		return errorStopped
	}

	task.log().WithField("url", task.URL).Debugf("[%s]: doing %s...", TAG, OP)

	// build request body
	var requestBody io.Reader
	if body != nil {
		task.log().WithField("request", body).Infof("[%s]: prepared request", TAG)

		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			task.log().WithError(err).Warnf("[%s]: failed to build %s body", TAG, OP)
			return fmt.Errorf("failed to build %s body: %s", OP, err)
		}

		requestBody = buf
	}

	// create request
	task.request, err = http.NewRequest(task.method, task.URL.String(), requestBody)
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to create %s request", TAG, OP)
		return fmt.Errorf("failed to create %s request: %s", OP, err)
	}

	// content type
	if requestBody != nil {
		task.request.Header.Add("Content-Type", "application/json")
	}

	// authorization
	service.prepareAuthorization(task.request, task.deviceAuth)

	select {
	case <-time.After(task.timeout):
		// TODO: task.request.Cancel // cancel request
		task.log().WithField("timeout", task.timeout).Warnf("[%s]: failed to wait %s response: timed out", TAG, OP)
		return fmt.Errorf("failed to wait %s response: timed out (%s)", OP, task.timeout)

	case <-service.stop:
		task.log().Warnf("[%s]: %s stopped", TAG, OP)
		return errorStopped

	case err = <-service.doAsync(task):
		if err != nil {
			task.log().WithError(err).Warnf("[%s]: failed to get %s response", TAG, OP)
			return fmt.Errorf("failed to get %s response: %s", OP, err)
		}
	}

	// important to close response body!
	defer task.response.Body.Close()

	// check status code
	if !checkStatus(task.response.StatusCode) {
		task.log().WithField("status", task.response.Status).Warnf("[%s]: unexpected %s status", TAG, OP)
		return fmt.Errorf("unexpected %s status: %s", OP, task.response.Status)
	}

	// unmarshal
	if result != nil && task.response.StatusCode != http.StatusNoContent {
		dec := json.NewDecoder(task.response.Body)
		err = dec.Decode(result)
		if err != nil {
			task.log().WithError(err).Warnf("[%s]: failed to parse %s body", TAG, OP)
			return fmt.Errorf("failed to parse %s body: %s", OP, err)
		}

		task.log().WithField("response", result).Infof("[%s]: parsed response", TAG)
	}

	return nil // OK
}
