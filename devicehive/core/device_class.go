package core

import "fmt"

// Represents device class which holds meta-information related to device.
type DeviceClass struct {
	// Unique identifier [do not change].
	Id uint64 `json:"id,omitempty"`

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

// Get DeviceClass string representation
func (deviceClass DeviceClass) String() string {
	body := ""

	// Id [optional]
	if deviceClass.Id != 0 {
		body += fmt.Sprintf("Id:%d, ", deviceClass.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", deviceClass.Name)

	// Version [optional]
	if len(deviceClass.Version) != 0 {
		body += fmt.Sprintf(", Version:%q", deviceClass.Version)
	}

	// IsPermanent [optional]
	if deviceClass.IsPermanent {
		body += fmt.Sprintf(", Permanent:%t", deviceClass.IsPermanent)
	}

	// OfflineTimeout [optional]
	if deviceClass.OfflineTimeout != 0 {
		body += fmt.Sprintf(", OfflineTimeout:%d", deviceClass.OfflineTimeout)
	}

	// Data [optional]
	body += fmt.Sprintf(", Data:%v", deviceClass.Data)

	// Equipment [optional]
	if len(deviceClass.Equipment) != 0 {
		body += fmt.Sprintf(", Equipment:%v", deviceClass.Equipment)
	}

	return fmt.Sprintf("DeviceClass{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (deviceClass *DeviceClass) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("DeviceClass: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("DeviceClass: %v - unexpected data type")
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case float64:
			deviceClass.Id = uint64(v)
		case uint64:
			deviceClass.Id = v
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for id", id)
		}
	}

	// name
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			deviceClass.Name = v
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for name", name)
		}
	}

	// version
	if vs, ok := data["version"]; ok {
		switch v := vs.(type) {
		case string:
			deviceClass.Version = v
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for version", vs)
		}
	}

	// is permanent flag
	if p, ok := data["isPermanent"]; ok {
		switch v := p.(type) {
		case bool:
			deviceClass.IsPermanent = v
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for isPermanent", p)
		}
	}

	// offline timeout
	if ot, ok := data["offlineTimeout"]; ok {
		switch v := ot.(type) {
		case float64:
			deviceClass.OfflineTimeout = int(v)
		case nil:
			// do nothing
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for offlineTimeout", ot)
		}
	}

	// data (as is)
	if d, ok := data["data"]; ok {
		deviceClass.Data = d
	}

	// equipment
	if e, ok := data["equipment"]; ok {
		switch v := e.(type) {
		case []interface{}:
			for i, x := range v {
				if i < len(deviceClass.Equipment) {
					// update existing
					err := deviceClass.Equipment[i].AssignJSON(x)
					if err != nil {
						return err
					}
				} else {
					// create new
					tmp := &Equipment{}
					err := tmp.AssignJSON(x)
					if err != nil {
						return err
					}
					deviceClass.AddEquipment(tmp)
				}
			}
			deviceClass.Equipment = deviceClass.Equipment[0:len(v)]
		default:
			return fmt.Errorf("DeviceClass: %v - unexpected value for equipment", e)
		}
	}

	return nil // OK
}

// AddEquipment() adds a new equipment
func (deviceClass *DeviceClass) AddEquipment(equipment *Equipment) {
	deviceClass.Equipment = append(deviceClass.Equipment, equipment)
}
