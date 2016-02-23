package devicehive

import (
	"bytes"
	"fmt"
)

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
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if network.Id != 0 {
		body.WriteString(fmt.Sprintf("Id:%d, ", network.Id))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", network.Name))

	// Key [optional]
	if len(network.Key) != 0 {
		body.WriteString(fmt.Sprintf(", Key:%q", network.Key))
	}

	// Description [optional]
	if len(network.Description) != 0 {
		body.WriteString(fmt.Sprintf(", Description:%q", network.Description))
	}

	return fmt.Sprintf("Network{%s}", body)
}

// Assign fields from map.
// This method is used to assign already parsed JSON data.
func (network *Network) FromMap(data interface{}) error {
	return fromJsonMap(network, data)
}
