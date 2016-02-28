package devicehive

import (
	"bytes"
	"fmt"
)

// Device represents a device - unit that communicate with DeviceHive.
type Device struct {
	// Unique identifier.
	ID string `json:"id,omitempty"`

	// Display name.
	Name string `json:"name,omitempty"`

	// Authentication key [optional].
	// The key is set during device registration and it has to be provided
	// for all subsequent calls initiated by device.
	// Maximum length is 64 characters.
	Key string `json:"key,omitempty"`

	// Operation status [optional]. The status can be set to any arbitrary value.
	// If device status monitoring feature is enabled, the framework will
	// set status value to "Offline" after defined period of inactivity.
	Status string `json:"status,omitempty"`

	// JSON object with an arbitrary structure [optional].
	Data interface{} `json:"data,omitempty"`

	// Associated network object [optional].
	Network *Network `json:"network,omitempty"`

	// Associated device class object [optional].
	DeviceClass *DeviceClass `json:"deviceClass,omitempty"`
}

// NewDevice creates a new device without network.
// No user data by default.
func NewDevice(id, name string, class *DeviceClass) *Device {
	return NewDeviceWithNetwork(id, name, class, nil)
}

// NewDeviceWithNetwork creates a new device with network initialized.
// No user data by default.
func NewDeviceWithNetwork(id, name string, class *DeviceClass, network *Network) *Device {
	device := new(Device)
	device.ID = id
	device.Name = name
	device.DeviceClass = class
	device.Network = network
	device.Status = "Online"
	return device
}

// Get Device string representation
func (device Device) String() string {
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if len(device.ID) != 0 {
		body.WriteString(fmt.Sprintf("ID:%q, ", device.ID))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", device.Name))

	// Key [optional]
	if len(device.Key) != 0 {
		body.WriteString(fmt.Sprintf(", Key:%q", device.Key))
	}

	// Status [optional]
	if len(device.Status) != 0 {
		body.WriteString(fmt.Sprintf(", Status:%q", device.Status))
	}

	// Data [optional]
	if device.Data != nil {
		body.WriteString(fmt.Sprintf(", Data:%v", device.Data))
	}

	// Network [optional]
	if device.Network != nil {
		body.WriteString(fmt.Sprintf(", %v", *device.Network))
	}

	// DeviceCLass [optional]
	if device.DeviceClass != nil {
		body.WriteString(fmt.Sprintf(", %v", *device.DeviceClass))
	}

	return fmt.Sprintf("Device{%s}", body)
}

// FromMap assigns fields from map.
// This method is used to assign already parsed JSON data.
func (device *Device) FromMap(data interface{}) error {
	return FromJSON(device, data)
}
