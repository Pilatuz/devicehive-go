package rest

import (
	"time"

	"github.com/pilatuz/go-devicehive"
)

// InsertNetwork inserts new network.
func (service *Service) InsertNetwork(network *devicehive.Network, timeout time.Duration) error {
	// build URL
	URL := *service.baseURL
	URL.Path += "/network"

	// request body
	body := *network // copy
	body.ID = 0      // do not put ID to the request body

	// do POST and check status is 2xx
	task := newTask("POST", &URL, timeout)
	err := service.do2xx(task, "/network/insert", &body, network)
	if err != nil {
		return err
	}

	return nil // OK
}
