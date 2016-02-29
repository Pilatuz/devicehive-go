package ws

import (
	"fmt"
	"time"

	dh "github.com/pilatuz/devicehive-go"
)

// Task is one request/response pair.
type Task struct {
	// identifer is 32-bits wide because
	// go uses float64 to store all JSON numbers
	// so uint64 cannot be precisely stored!
	// (float64 uses about 52 bits per mantissa)
	identifier uint32

	dataToSend   map[string]interface{}
	dataReceived map[string]interface{}

	doneCh  chan error
	timeout time.Duration
}

// prepare device-key authorization
func (task *Task) prepareAuthorization(device *dh.Device) {
	if device == nil {
		return // do nothing
	}

	// optional fields
	if len(device.ID) != 0 || len(device.Key) != 0 {
		task.dataToSend["deviceId"] = device.ID
		task.dataToSend["deviceKey"] = device.Key
	}
}

// ReportDone reports task is done
func (task *Task) ReportDone(data map[string]interface{}, err error) {
	task.dataReceived = data
	task.doneCh <- err
}

// create new empty task and put to active set
func (service *Service) newTask(timeout time.Duration) *Task {
	task := new(Task)
	task.timeout = timeout
	task.doneCh = make(chan error, 1)

	service.taskLock.Lock()
	defer service.taskLock.Unlock()
	service.lastTaskId += 1 // generate unique identifier
	if service.lastTaskId == 0 {
		// avoid zero identifier
		service.lastTaskId += 1
	}
	task.identifier = service.lastTaskId
	service.tasks[task.identifier] = task // put to "active" set

	// task.log().Debugf("[%s]: new task created", TAG)
	return task
}

// find task and remove it from active list
// return nil if not found
func (service *Service) takeTask(id uint32) *Task {
	service.taskLock.Lock()
	defer service.taskLock.Unlock()

	if task, ok := service.tasks[id]; ok {
		delete(service.tasks, id)
		return task
	}

	return nil // not found
}

// send request and parse response
func (service *Service) do(task *Task, OP string) (err error) {
	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(task.timeout):
		// TODO: task.request.Cancel // cancel request
		task.log().WithField("timeout", task.timeout).Warnf("[%s]: failed to wait %s response: timed out", TAG, OP)
		return fmt.Errorf("failed to wait %s response: timed out (%s)", OP, task.timeout)

	case <-service.stop:
		task.log().Warnf("[%s]: %s stopped", TAG, OP)
		return errorStopped

	case err := <-task.doneCh:
		if err != nil {
			task.log().WithError(err).Warnf("[%s]: failed to get %s response", TAG, OP)
			return fmt.Errorf("failed to get %s response: %s", OP, err)
		}
	}

	return nil // OK
}
