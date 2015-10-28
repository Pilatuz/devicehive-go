package core

import "fmt"

// Represents notification object - a set of data sent from devices to DeviceHive.
type Notification struct {
	// Unique identifier [do not change].
	Id uint64 `json:"id,omitempty"`

	// Timestamp, UTC.
	Timestamp string `json:"timestamp,omitempty"`

	// Notification name.
	Name string `json:"notification,omitempty"`

	// JSON object with an arbitrary structure [optional].
	Parameters interface{} `json:"parameters,omitempty"`
}

// NewNotification creates a new notification.
func NewNotification(name string, parameters interface{}) *Notification {
	return &Notification{Name: name, Parameters: parameters}
}

// Get Notification string representation
func (notification Notification) String() string {
	body := ""

	// Id [optional]
	if notification.Id != 0 {
		body += fmt.Sprintf("Id:%d, ", notification.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", notification.Name)

	// Timestamp
	if len(notification.Timestamp) != 0 {
		body += fmt.Sprintf(", Timestamp:%q", notification.Timestamp)
	}

	// Parameters [optional]
	if notification.Parameters != nil {
		body += fmt.Sprintf(", Parameters:%v", notification.Parameters)
	}

	return fmt.Sprintf("Notification{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (notification *Notification) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("Notification: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Notification: %v - unexpected data type", rawData)
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case float64:
			notification.Id = uint64(v)
		case uint64:
			notification.Id = v
		default:
			return fmt.Errorf("Notification: %v - unexpected value for id", id)
		}
	}

	// timestamp
	if ts, ok := data["timestamp"]; ok {
		switch v := ts.(type) {
		case string:
			notification.Timestamp = v
		default:
			return fmt.Errorf("Notification: %v - unexpected value for timestamp", ts)
		}
	}

	// name
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			notification.Name = v
		default:
			return fmt.Errorf("Notification: %v - unexpected value for name", name)
		}
	}

	// parameters (as is)
	if p, ok := data["parameters"]; ok {
		notification.Parameters = p
	}

	return nil // OK
}
