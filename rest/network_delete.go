package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// DeleteNetwork() function deletes the network.
func (service *Service) DeleteNetwork(network *devicehive.Network, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/network/%d", network.ID)

	// do DELETE and check status is 2xx
	task := newTask("DELETE", &URL, timeout)
	err := service.do2xx(task, "/network/delete", nil, nil)
	if err != nil {
		return err
	}

	return nil // OK
}
