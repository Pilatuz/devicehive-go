package rest

import (
	dh "github.com/pilatuz/go-devicehive"
)

// GetServerInfo gets the main server information.
func (service *Service) GetServerInfo() (*dh.ServerInfo, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += "/info"

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	info := new(dh.ServerInfo)
	err := service.do200(task, "/info", nil, info)
	if err != nil {
		return nil, err
	}

	return info, nil // OK
}
