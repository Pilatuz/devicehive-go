package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare GetDevice task
func (service *Service) prepareGetDevice(device *core.Device) (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "device/get",
		"requestId": task.id}

	// prepare authorization
	task.prepareAuthorization(device)

	return
}

// Process GetDevice task
func (service *Service) processGetDevice(task *Task, device *core.Device) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /device/get status (error: %s)", err)
		return
	}

	// parse response
	err = device.AssignJSON(task.dataRecved["device"])
	if err != nil {
		log.Warnf("WS: failed to parse /device/get body (error: %s)", err)
		return
	}

	return
}

// GetDevice() function gets the device information.
func (service *Service) GetDevice(deviceId, deviceKey string, timeout time.Duration) (device *core.Device, err error) {
	device = &core.Device{Id: deviceId, Key: deviceKey}
	task, err := service.prepareGetDevice(device)
	if err != nil {
		log.Warnf("WS: failed to prepare /device/get task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /device/get task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		err = service.processGetDevice(task, device)
		if err != nil {
			log.Warnf("WS: failed to process /device/get task (error: %s)", err)
			return
		}
	}

	return
}
