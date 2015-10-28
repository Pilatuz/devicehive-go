package core

import "fmt"

// Represents a device - unit that communicate with DeviceHive.
type Device struct {
	// Unique identifier.
	Id string `json:"id,omitempty"`

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
	return &Device{Id: id, Name: name, Status: "Online",
		DeviceClass: class, Network: nil}
}

// NewDeviceWithNetwork creates a new device with network initialized.
// No user data by default.
func NewDeviceWithNetwork(id, name string, class *DeviceClass, network *Network) *Device {
	return &Device{Id: id, Name: name, Status: "Online",
		DeviceClass: class, Network: network}
}

// Get Device string representation
func (device Device) String() string {
	body := ""

	// Id [optional]
	if len(device.Id) != 0 {
		body += fmt.Sprintf("Id:%q, ", device.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", device.Name)

	// Key [optional]
	if len(device.Key) != 0 {
		body += fmt.Sprintf(", Key:%q", device.Key)
	}

	// Status [optional]
	if len(device.Status) != 0 {
		body += fmt.Sprintf(", Status:%q", device.Status)
	}

	// Data [optional]
	body += fmt.Sprintf(", Data:%v", device.Data)

	// Network [optional]
	if device.Network != nil {
		body += fmt.Sprintf(", %v", *device.Network)
	}

	// DeviceCLass [optional]
	if device.DeviceClass != nil {
		body += fmt.Sprintf(", %v", *device.DeviceClass)
	}

	return fmt.Sprintf("Device{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (device *Device) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("Device: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Device: %v - unexpected data type", rawData)
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case string:
			device.Id = v
		default:
			return fmt.Errorf("Device: %v - unexpected value for id", id)
		}
	}

	// name
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			device.Name = v
		default:
			return fmt.Errorf("Device: %v - unexpected value for name", name)
		}
	}

	// key
	if key, ok := data["key"]; ok {
		switch v := key.(type) {
		case string:
			device.Key = v
		default:
			return fmt.Errorf("Device: %v - unexpected value for key", key)
		}
	}

	// status
	if s, ok := data["status"]; ok {
		switch v := s.(type) {
		case string:
			device.Status = v
		default:
			return fmt.Errorf("Device: %v - unexpected value for status", s)
		}
	}

	// data (as is)
	if d, ok := data["data"]; ok {
		device.Data = d
	}

	// network
	if net, ok := data["network"]; ok {
		if net != nil {
			if device.Network == nil {
				device.Network = &Network{}
			}
			err := device.Network.AssignJSON(net)
			if err != nil {
				return err
			}
		} else {
			device.Network = nil
		}
	}

	// deviceClass
	if dc, ok := data["deviceClass"]; ok {
		if dc != nil {
			if device.DeviceClass == nil {
				device.DeviceClass = &DeviceClass{}
			}
			err := device.DeviceClass.AssignJSON(dc)
			if err != nil {
				return err
			}
		} else {
			device.DeviceClass = nil
		}
	}

	return nil // OK
}
