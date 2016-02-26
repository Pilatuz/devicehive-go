package rest

import (
	"fmt"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetNetwork gets the network data.
func (service *Service) GetNetwork(networkID uint64, timeout time.Duration) (*devicehive.Network, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/network/%d", networkID)

	// result
	network := new(devicehive.Network)
	network.ID = networkID

	// do GET and check status is 200
	task := newTask("GET", &URL, timeout)
	err := service.do200(task, "/network/get", nil, network)
	if err != nil {
		return nil, err
	}

	return network, nil // OK
}
