package devicehive

import (
	"bytes"
	"fmt"
)

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

// NewEquipment creates a new equipment.
// No user data by default.
func NewEquipment(name, code, type_ string) *Equipment {
	equipment := new(Equipment)
	equipment.Name = name
	equipment.Code = code
	equipment.Type = type_
	return equipment
}

// Get Equipment string representation
func (equipment Equipment) String() string {
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if equipment.Id != 0 {
		body.WriteString(fmt.Sprintf("Id:%d, ", equipment.Id))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", equipment.Name))

	// Code
	body.WriteString(fmt.Sprintf(", Code:%q", equipment.Code))

	// Type
	body.WriteString(fmt.Sprintf(", Type:%q", equipment.Type))

	// Data [optional]
	if equipment.Data != nil {
		body.WriteString(fmt.Sprintf(", Data:%v", equipment.Data))
	}

	return fmt.Sprintf("Equipment{%s}", body)
}

// Assign fields from map.
// This method is used to assign already parsed JSON data.
func (equipment *Equipment) FromMap(data interface{}) error {
	return fromJsonMap(equipment, data)
}
