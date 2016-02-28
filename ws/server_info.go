package ws

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// GetServerInfo gets the main server's information.
func (service *Service) GetServerInfo() (*dh.ServerInfo, error) {
	const OP = "/info"

	task := service.newTask(service.DefaultTimeout)
	task.dataToSend = map[string]interface{}{
		"action":    "server/info",
		"requestId": task.identifier,
	}

	err := service.do(task, "/info")
	if err != nil {
		return nil, err
	}

	info := new(dh.ServerInfo)
	err = info.FromMap(task.dataReceived["info"])
	if err != nil {
		log.WithError(err).Warnf("[%s]: failed to parse %s response", TAG, OP)
		return nil, fmt.Errorf("failed to parse %s response: %s", OP, err)
	}

	return info, nil // OK
}
