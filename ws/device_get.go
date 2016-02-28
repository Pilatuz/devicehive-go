package ws

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// GetDevice gets the device information.
func (service *Service) GetDevice(deviceID, deviceKey string) (*dh.Device, error) {
	const OP = "/device/get"

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "device/get",
		"requestId": task.identifier,
	}

	device := new(dh.Device)
	device.ID = deviceID
	device.Key = deviceKey

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, OP)
	if err != nil {
		return nil, err
	}

	// parse response
	err = device.FromMap(task.dataReceived["device"])
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to parse %s response", TAG, OP)
		return nil, fmt.Errorf("failed to parse %s response: %s", OP, err)
	}

	return device, nil // OK
}
