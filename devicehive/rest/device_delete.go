package rest

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"net/http"
	"time"
)

// Prepare DeleteDevice task
func (service *Service) prepareDeleteDevice(device *core.Device) (task Task, err error) {
	// create request
	url := fmt.Sprintf("%s/device/%s", service.baseUrl, device.Id)
	task.request, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Warnf("REST: failed to create /device/delete request (error: %s)", err)
		return
	}

	// authorization
	service.prepareAuthorization(task.request, device)

	return
}

// Process DeleteDevice task
func (service *Service) processDeleteDevice(task Task) (err error) {
	// check task error first
	if task.err != nil {
		err = task.err
		return
	}

	// check status code
	if task.response.StatusCode != http.StatusOK {
		log.Warnf("REST: unexpected /device/delete status %s",
			task.response.Status)
		err = fmt.Errorf("unexpected status: %s",
			task.response.Status)
		return
	}

	return
}

// DeleteDevice() function deletes the device.
func (service *Service) DeleteDevice(device *core.Device, timeout time.Duration) (err error) {
	log.Tracef("REST: deleting device...")

	task, err := service.prepareDeleteDevice(device)
	if err != nil {
		log.Warnf("REST: failed to prepare /device/delete task (error: %s)", err)
		return
	}

	select {
	case <-time.After(timeout):
		log.Warnf("REST: failed to wait %s for /device/delete task", timeout)
		err = fmt.Errorf("timed out")

	case task = <-service.doAsync(task):
		err = service.processDeleteDevice(task)
		if err != nil {
			log.Warnf("REST: failed to process /device/delete task (error: %s)", err)
			return
		}
	}

	return
}
