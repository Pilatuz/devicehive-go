package rest

import (
	"fmt"
	"net/url"

	dh "github.com/pilatuz/devicehive-go"
)

// GetDeviceList get the device list.
func (service *Service) GetDeviceList(take, skip int) ([]*dh.Device, error) {
	// build URL
	URL := *service.baseURL
	URL.Path += "/device"
	query := url.Values{}
	if take > 0 {
		query.Set("take", fmt.Sprintf("%d", take))
	}
	if skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", skip))
	}
	URL.RawQuery = query.Encode()

	// result
	var res []interface{}

	// do GET and check status is 200
	task := newTask("GET", &URL, service.DefaultTimeout)
	err := service.do200(task, "/device/list", nil, &res)
	if err != nil {
		return nil, err
	}

	// convert map to devices
	devices := make([]*dh.Device, 0, len(res))
	for _, data := range res {
		d := new(dh.Device)
		if err := d.FromMap(data); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil // OK
}
