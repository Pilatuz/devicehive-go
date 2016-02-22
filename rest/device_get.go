package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare GetDevice task
func (service *Service) prepareGetDevice(deviceId, deviceKey string) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s", service.baseUrl, deviceId)

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /device/get request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request,
		&core.Device{Id: deviceId, Key: deviceKey})

	return
}

// Process GetDevice task
func (service *Service) processGetDevice(task Task, device *core.Device) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /device/get status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, device)
	if err != nil {
		log.Warnf("REST: failed to parse /device/get body (error: %s)", err)
		return
	}

	return
}

// GetDevice() function get the device data.
func (service *Service) GetDevice(deviceId, deviceKey string, timeout time.Duration) (device *core.Device, err error) {
	log.Tracef("REST: getting device %q...", deviceId)

	task, err := service.prepareGetDevice(deviceId, deviceKey)
	if err != nil {
		log.Warnf("REST: failed to prepare /device/get task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /device/get task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		device = &core.Device{Id: deviceId, Key: deviceKey}
		err = service.processGetDevice(task, device)
		if err != nil {
			log.Warnf("REST: failed to process /device/get task (error: %s)", err)
			return
		}
	}

	return
}
