package rest

import (
	"fmt"

	dh "github.com/pilatuz/go-devicehive"
)

// UpdateNetwork updates the network.
func (service *Service) UpdateNetwork(network *dh.Network) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/network/%d", network.ID)

	// request body
	body := *network // copy
	body.ID = 0      // do not put ID to the request body

	// do PUT and check status is 200
	task := newTask("PUT", &URL, service.DefaultTimeout)
	err := service.do2xx(task, "/network/update", &body, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
