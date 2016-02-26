package devicehive

import (
	"bytes"
	"fmt"
)

// DeviceClass represents device class which holds meta-information related to device.
type DeviceClass struct {
	// Unique identifier [do not change].
	ID uint64 `json:"id,omitempty"`

	// Display name.
	Name string `json:"name"`

	// Version string.
	Version string `json:"version"`

	// Indicates whether device class is permanent.
	// Permanent device classes could not be modified during device registration.
	IsPermanent bool `json:"isPermanent,omitempty"`

	// If set, specifies inactivity timeout in seconds before the framework
	// changes device status to "Offline". Device considered inactive when
	// it's not persistently connected and does not send any notifications.
	OfflineTimeout int `json:"offlineTimeout,omitempty"`

	// JSON object with an arbitrary structure [optional].
	Data interface{} `json:"data,omitempty"`

	// Associated equipment objects.
	Equipment []*Equipment `json:"equipment,omitempty"`
}

// NewDeviceClass creates a new device class.
// No user data and no equipment by default.
func NewDeviceClass(name, version string) *DeviceClass {
	deviceClass := new(DeviceClass)
	deviceClass.Name = name
	deviceClass.Version = version
	return deviceClass
}

// AddEquipment adds a new equipment to the device class
func (deviceClass *DeviceClass) AddEquipment(equipment ...*Equipment) {
	deviceClass.Equipment = append(deviceClass.Equipment, equipment...)
}

// Get DeviceClass string representation
func (deviceClass DeviceClass) String() string {
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if deviceClass.ID != 0 {
		body.WriteString(fmt.Sprintf("ID:%d, ", deviceClass.ID))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", deviceClass.Name))

	// Version [optional]
	if len(deviceClass.Version) != 0 {
		body.WriteString(fmt.Sprintf(", Version:%q", deviceClass.Version))
	}

	// IsPermanent [optional]
	if deviceClass.IsPermanent {
		body.WriteString(fmt.Sprintf(", Permanent:%t", deviceClass.IsPermanent))
	}

	// OfflineTimeout [optional]
	if deviceClass.OfflineTimeout != 0 {
		body.WriteString(fmt.Sprintf(", OfflineTimeout:%d", deviceClass.OfflineTimeout))
	}

	// Data [optional]
	if deviceClass.Data != nil {
		body.WriteString(fmt.Sprintf(", Data:%v", deviceClass.Data))
	}

	// Equipment [optional]
	if len(deviceClass.Equipment) != 0 {
		body.WriteString(fmt.Sprintf(", Equipment:%v", deviceClass.Equipment))
	}

	return fmt.Sprintf("DeviceClass{%s}", body)
}

// FromMap assigns fields from map.
// This method is used to assign already parsed JSON data.
func (deviceClass *DeviceClass) FromMap(data interface{}) error {
	return fromJSON(deviceClass, data)
}
