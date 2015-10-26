package core

import "fmt"

// Represents equipment object - a peripheral or sensor hosted on device.
type Equipment struct {
	// Unique identifier [do not change].
	Id uint64 `json:"id,omitempty"`

	// Display name.
	Name string `json:"name"`

	// Unique code. It's used to reference particular equipment
	// and it should be unique within a device class.
	Code string `json:"code"`

	// An arbitrary string representing equipment capabilities.
	Type string `json:"type"`

	// JSON object with an arbitrary structure [optional].
	Data interface{} `json:"data,omitempty"`
}

// Get Equipment string representation
func (equipment Equipment) String() string {
	body := ""

	// Id [optional]
	if equipment.Id != 0 {
		body += fmt.Sprintf("Id:%d, ", equipment.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", equipment.Name)

	// Code
	body += fmt.Sprintf(", Code:%q", equipment.Code)

	// Type
	body += fmt.Sprintf(", Type:%q", equipment.Type)

	// Data [optional]
	body += fmt.Sprintf(", Data:%v", equipment.Data)

	return fmt.Sprintf("Equipment{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (equipment *Equipment) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("Equipment: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Equipment: %v - unexpected data type")
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case float64:
			equipment.Id = uint64(v)
		case uint64:
			equipment.Id = v
		default:
			return fmt.Errorf("Equipment: %v - unexpected value for id", id)
		}
	}

	// name
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			equipment.Name = v
		default:
			return fmt.Errorf("Equipment: %v - unexpected value for name", name)
		}
	}

	// code
	if code, ok := data["code"]; ok {
		switch v := code.(type) {
		case string:
			equipment.Code = v
		default:
			return fmt.Errorf("Equipment: %v - unexpected value for code", code)
		}
	}

	// type
	if t, ok := data["type"]; ok {
		switch v := t.(type) {
		case string:
			equipment.Type = v
		default:
			return fmt.Errorf("Equipment: %v - unexpected value for type", t)
		}
	}

	// data (as is)
	if d, ok := data["data"]; ok {
		equipment.Data = d
	}

	return nil // OK
}
