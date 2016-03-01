package devicehive

import (
	"time"
)

// DeviceService is an abstract DeviceHive /device API.
type DeviceService interface {
	Stop()

	SetTimeout(timeout time.Duration)
	GetServerInfo() (info *ServerInfo, err error)

	RegisterDevice(device *Device) (err error)
	GetDevice(deviceID, deviceKey string) (device *Device, err error)

	// GetCommand(device *Device, commandID uint64) (command *Command, err error)
	UpdateCommand(device *Device, command *Command) (err error)
	SubscribeCommands(device *Device, timestamp string) (listener *CommandListener, err error)
	UnsubscribeCommands(device *Device) (err error)

	// GetNotification(device *Device, notificationID uint64) (notification *Notification, err error)
	InsertNotification(device *Device, notification *Notification) (err error)
}

// ClientService is an abstract DeviceHive /client API.
type ClientService interface {
	Stop()

	SetTimeout(timeout time.Duration)
	GetServerInfo() (info *ServerInfo, err error)

	// GetDevice(deviceID, deviceKey string) (device *Device, err error)

	// GetCommand(device *Device, commandID uint64) (command *Command, err error)
	UpdateCommand(device *Device, command *Command) (err error)
	SubscribeCommands(device *Device, timestamp string) (listener *CommandListener, err error)
	UnsubscribeCommands(device *Device) (err error)

	// GetNotification(device *Device, notificationID uint64) (notification *Notification, err error)
	InsertNotification(device *Device, notification *Notification) (err error)
	SubscribeNotifications(device *Device, timestamp string) (listener *NotificationListener, err error)
	UnsubscribeNotifications(device *Device) (err error)
}
