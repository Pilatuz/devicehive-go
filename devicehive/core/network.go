package core

import "fmt"

// Represents a network - a custom set of devices.
type Network struct {
	// Unique identifier [do not change].
	Id uint64 `json:"id,omitempty"`

	// Display name.
	Name string `json:"name"`

	// Authentication key [optional]. This key is used to protect the network
	// from unauthorized device registrations. When defined, devices will
	// need to pass the key in order to register to the network.
	Key string `json:"key,omitempty"`

	// Text description [optional].
	Description string `json:"description,omitempty"`
}

// NewNetwork creates a new network.
// Network Description is empty.
func NewNetwork(name, key string) *Network {
	return &Network{Name: name, Key: key}
}

// Get Network string representation
func (network Network) String() string {
	body := ""

	// Id [optional]
	if network.Id != 0 {
		body += fmt.Sprintf("Id:%d, ", network.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", network.Name)

	// Key [optional]
	if len(network.Key) != 0 {
		body += fmt.Sprintf(", Key:%q", network.Key)
	}

	// Description [optional]
	if len(network.Description) != 0 {
		body += fmt.Sprintf(", Description:%q", network.Description)
	}

	return fmt.Sprintf("Network{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (network *Network) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("Network: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Network: %v - unexpected data type")
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case float64:
			network.Id = uint64(v)
		case uint64:
			network.Id = v
		default:
			return fmt.Errorf("Network: %v - unexpected value for id", id)
		}
	}

	// name
	if name, ok := data["name"]; ok {
		switch v := name.(type) {
		case string:
			network.Name = v
		default:
			return fmt.Errorf("Network: %v - unexpected value for name", name)
		}
	}

	// key
	if key, ok := data["key"]; ok {
		switch v := key.(type) {
		case string:
			network.Key = v
		default:
			return fmt.Errorf("Network: %v - unexpected value for key", key)
		}
	}

	// description
	if d, ok := data["description"]; ok {
		switch v := d.(type) {
		case string:
			network.Description = v
		case nil:
			// do nothing
		default:
			return fmt.Errorf("Network: %v - unexpected value for description", d)
		}
	}

	return nil // OK
}
