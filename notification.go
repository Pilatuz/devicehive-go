package devicehive

import (
	"bytes"
	"fmt"
)

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

// Notification listener is used to listen for asynchronous notifications.
type NotificationListener struct {
	// Channel to receive notifications
	C chan *Notification
}

// NewEmptyNotification creates an empty notification.
func NewEmptyNotification() *Notification {
	notification := new(Notification)
	return notification
}

// NewNotification creates a new notification.
func NewNotification(name string, parameters interface{}) *Notification {
	notification := new(Notification)
	notification.Name = name
	notification.Parameters = parameters
	return notification
}

// NewNotificationListener creates a new notification listener.
func NewNotificationListener(buffered int) *NotificationListener {
	listener := new(NotificationListener)
	listener.C = make(chan *Notification, buffered)
	return listener
}

// Get Notification string representation
func (notification Notification) String() string {
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if notification.Id != 0 {
		body.WriteString(fmt.Sprintf("Id:%d, ", notification.Id))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", notification.Name))

	// Timestamp
	if len(notification.Timestamp) != 0 {
		body.WriteString(fmt.Sprintf(", Timestamp:%q", notification.Timestamp))
	}

	// Parameters [optional]
	if notification.Parameters != nil {
		body.WriteString(fmt.Sprintf(", Parameters:%v", notification.Parameters))
	}

	return fmt.Sprintf("Notification{%s}", body)
}

// Assign fields from map.
// This method is used to assign already parsed JSON data.
func (notification *Notification) FromMap(data interface{}) error {
	return fromJsonMap(notification, data)
}
