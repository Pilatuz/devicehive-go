package rest

import (
	"fmt"

	dh "github.com/pilatuz/devicehive-go"
)

// GetNetwork gets the network data.
func (service *Service) GetNetwork(networkID uint64) (*dh.Network, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += fmt.Sprintf("/network/%d", networkID)

	// result
	network := &dh.Network{
		ID: networkID,
	}

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	err := service.do200(task, "/network/get", nil, network)
	if err != nil {
		return nil, err
	}

	return network, nil // OK
}
