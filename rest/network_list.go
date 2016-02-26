package rest

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetNetworkList gets the list of networks.
func (service *Service) GetNetworkList(take, skip int, timeout time.Duration) ([]*devicehive.Network, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += "/network"
	query := url.Values{}
	if take > 0 {
		query.Set("take", fmt.Sprintf("%d", take))
	}
	if skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", skip))
	}
	URL.RawQuery = query.Encode()

	// result
	var networks []*devicehive.Network

	// do GET and check status is 200
	task := newTask("GET", &URL, timeout)
	err := service.do200(task, "/network/list", nil, &networks)
	if err != nil {
		return nil, err
	}

	return networks, nil // OK
}
