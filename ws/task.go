package ws

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"strings"
)

type Task struct {
	id         uint64
	dataToSend map[string]interface{}
	dataRecved map[string]interface{}
	done       chan *Task
}

// prepare device-key authorization
func (task *Task) prepareAuthorization(device *core.Device) {
	if device != nil {
		// DeviceId [optional]
		if len(device.Id) != 0 {
			task.dataToSend["deviceId"] = device.Id
		}

		// DeviceKey [optional]
		if len(device.Key) != 0 {
			task.dataToSend["deviceKey"] = device.Key
		}
	}
}

// Format the JSON data
func (task *Task) Format() (body []byte, err error) {
	body, err = json.Marshal(task.dataToSend)
	return
}

// Check "success" status
func (task *Task) CheckStatus() (err error) {
	status := safeString(task.dataRecved["status"])
	if !strings.EqualFold(status, "success") {
		err = fmt.Errorf("unexpected status: %q [%s %s]", status,
			safeString(task.dataRecved["code"]),
			safeString(task.dataRecved["error"]))
	}
	return
}

// create new empty task and put to active set
func (service *Service) newTask() (task *Task) {
	task = new(Task)
	task.done = make(chan *Task)

	service.taskLock.Lock()
	defer service.taskLock.Unlock()
	service.lastTaskId += 1 // generate unique identifier
	task.id = service.lastTaskId
	service.tasks[task.id] = task // put to "active" set

	//log.Tracef("WS: new task #%d created", task.id)
	return
}

// find task and remove it from active list
// return nil if not found
func (service *Service) takeTask(id uint64) (task *Task) {
	service.taskLock.Lock()
	defer service.taskLock.Unlock()

	task = service.tasks[id]
	delete(service.tasks, id)

	return
}
