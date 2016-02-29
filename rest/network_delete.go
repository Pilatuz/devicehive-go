package rest

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// DeleteNetwork() function deletes the network.
func (service *Service) DeleteNetwork(network *dh.Network) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/network/%d", network.ID)

	// do DELETE and check status is 2xx
	task := newTask("DELETE", &URL, service.DefaultTimeout)
	err := service.do2xx(task, "/network/delete", nil, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
