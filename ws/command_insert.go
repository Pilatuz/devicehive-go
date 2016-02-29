package ws

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// InsertCommand insert the new command.
// /client service only.
func (service *Service) InsertCommand(device *dh.Device, command *dh.Command) error {
	const OP = "/command/insert"

	// request data (do not put all fields)
	data := *command // deep copy
	data.ID = 0      // do not put Id inside

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":     "command/insert",
		"requestId":  task.identifier,
		"deviceGuid": device.ID,
		"command":    data,
	}

	// prepare authorization
	task.prepareAuthorization(device)

	err := service.do(task, OP)
	if err != nil {
		return err
	}

	// parse response
	err = command.FromMap(task.dataReceived["command"])
	if err != nil {
		task.log().WithError(err).Warnf("[%s]: failed to parse %s response", TAG, OP)
		return fmt.Errorf("failed to parse %s response: %s", OP, err)
	}

	return nil // OK
}
