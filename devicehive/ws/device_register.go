package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare RegisterDevice task
func (service *Service) prepareRegisterDevice(device *core.Device) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "device/save",
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	dev_data := *device // deep copy
	dev_data.Id = ""    // do not put Id inside
	task.dataToSend["device"] = dev_data

	return
}

// Process RegisterDevice task
func (service *Service) processRegisterDevice(task *Task) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /device/register status (error: %s)", err)
		return
	}

	return
}

// RegisterDevice() function registers the device.
func (service *Service) RegisterDevice(device *core.Device, timeout time.Duration) (err error) {
	task, err := service.prepareRegisterDevice(device)
	if err != nil {
		log.Warnf("WS: failed to prepare /device/register task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /device/register task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processRegisterDevice(task)
		if err != nil {
			log.Warnf("WS: failed to process /device/register task (error: %s)", err)
			return
		}
	}

	return
}
