package rest

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"net/url"
	"time"
)

// Prepare ListDevices task
func (service *Service) prepareGetDeviceList(take, skip int) (task Task, err error) {
	// create request
	query := url.Values{}
	if take > 0 {
		query.Set("take", fmt.Sprintf("%d", take))
	}
	if skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", skip))
	}
	url := fmt.Sprintf("%s/device", service.baseUrl)
	if len(query) != 0 {
		url += "?" + query.Encode()
	}

	task.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /device/list request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, nil)

	return
}

// Process ListDevices task
func (service *Service) processGetDeviceList(task Task) (devices []core.Device, err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /device/list status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	// unmarshal
	err = json.Unmarshal(task.body, &devices)
	if err != nil {
		log.Warnf("REST: failed to parse /device/list body (error: %s)", err)
		return
	}

	return
}

// ListDevices() function get the device list.
func (service *Service) GetDeviceList(take, skip int, timeout time.Duration) (devices []core.Device, err error) {
	log.Tracef("REST: getting device list (take:%d, skip:%d)...", take, skip)

	task, err := service.prepareGetDeviceList(take, skip)
	if err != nil {
		log.Warnf("REST: failed to prepare /device/list task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /device/list task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		devices, err = service.processGetDeviceList(task)
		if err != nil {
			log.Warnf("REST: failed to process /device/list task (error: %s)", err)
			return
		}
	}

	return
}
