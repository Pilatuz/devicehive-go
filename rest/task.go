package rest

import (
	"net/http"
	"sync/atomic"
)

var (
	taskID = uint64(0) // unique identifier generator
)

// Asynchronous request/task.
type asyncTask struct {
	identifier uint64

	request  *http.Request
	response *http.Response
}

// create new async task.
// assigns new unique identifier
func newTask() *asyncTask {
	task := new(asyncTask)
	task.identifier = atomic.AddUint64(&taskID, 1)
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
