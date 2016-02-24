package rest

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	taskId = uint64(0)
)

// Asynchronous request/task.
type asyncTask struct {
	Identifier uint64

	request  *http.Request
	response *http.Response
}

// create new async task.
// assigns new unique identifier
func newTask() *asyncTask {
	task := new(asyncTask)
	task.Identifier = atomic.AddUint64(&taskId, 1)
	return task
}

// Do a request/task synchronously.
func (service *Service) doSync(task *asyncTask) error {
	task.log().Debugf("[%s]: sending: %s %s", TAG,
		task.request.Method, task.request.URL)

	var err error
	task.response, err = service.client.Do(task.request)
	if err != nil {
		task.log().Debugf("[%s]: failed to do request: %s", TAG, err)
		return fmt.Errorf("failed to do request: %s", err)
	}

	task.log().Debugf("[%s]: got response status: %s", TAG, task.response.Status)
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
