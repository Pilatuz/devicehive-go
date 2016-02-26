package rest

import (
	"time"

	"github.com/pilatuz/go-devicehive"
)

// GetServerInfo gets the main server information.
func (service *Service) GetServerInfo(timeout time.Duration) (*devicehive.ServerInfo, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += "/info"

	// do GET and check status is 200
	task := newTask("GET", &URL, timeout)
	info := new(devicehive.ServerInfo)
	err := service.do200(task, "/info", nil, info)
	if err != nil {
		return nil, err
	}

	return info, nil // OK
}
