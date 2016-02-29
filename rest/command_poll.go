package rest

import (
	"fmt"
	"net/url"

	dh "github.com/pilatuz/devicehive-go"
)

// PollCommands polls the commands.
func (service *Service) PollCommands(device *dh.Device, timestamp, names, waitTimeout string) ([]*dh.Command, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/device/%s/command/poll", device.ID)
	query := url.Values{}
	if len(timestamp) != 0 {
		query.Set("timestamp", timestamp)
	}
	if len(names) != 0 {
		query.Set("names", names)
	}
	if len(waitTimeout) != 0 {
		query.Set("waitTimeout", waitTimeout)
	}
	URL.RawQuery = query.Encode()

	// result
	var commands []*dh.Command

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	task.deviceAuth = device
	err := service.do200(task, "/command/poll", nil, &commands)
	if err != nil {
		return nil, err
	}

	// convert map to commands
	//	commands := make([]*dh.Command, 0, len(res))
	//	for _, data := range res {
	//		c := new(dh.Command)
	//		if err := c.FromMap(data); err != nil {
	//			return nil, err
	//		}
	//		commands = append(commands, c)
	//	}

	return commands, nil // OK
}
