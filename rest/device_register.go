package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare RegisterDevice task
func (service *Service) prepareRegisterDevice(device *core.Device) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s", service.baseUrl, device.Id)

	dev_data := *device
	dev_data.Id = "" // do not put ID to the request body
	body, err := json.Marshal(&dev_data)
	if err != nil {
		log.Warnf("REST: failed to format /device/register request (error: %s)", err)
		return
	}

	task.request, err = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warnf("REST: failed to create /device/register request (error: %s)", err)
		return
	}
	task.request.Header.Add("Content-Type", "application/json")

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process RegisterDevice task
func (service *Service) processRegisterDevice(task Task) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode < http.StatusOK ||
		task.response.StatusCode > http.StatusPartialContent {
		log.Warnf("REST: unexpected /device/register status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	return
}

// RegisterDevice() function registers the device.
func (service *Service) RegisterDevice(device *core.Device, timeout time.Duration) (err error) {
	log.Tracef("REST: registering device %q...", device.Id)

	task, err := service.prepareRegisterDevice(device)
	if err != nil {
		log.Warnf("REST: failed to prepare /device/register task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /device/register task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processRegisterDevice(task)
		if err != nil {
			log.Warnf("REST: failed to process /device/register task (error: %s)", err)
			return
		}
	}

	return
}
