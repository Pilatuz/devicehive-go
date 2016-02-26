package devicehive

import (
	"time"
)

// DeviceService is an abstract DeviceHive /device API.
type DeviceService interface {
	GetServerInfo(timeout time.Duration) (info *ServerInfo, err error)

	RegisterDevice(device *Device, timeout time.Duration) (err error)
	GetDevice(deviceID, deviceKey string, timeout time.Duration) (device *Device, err error)

	GetCommand(device *Device, commandID uint64, timeout time.Duration) (command *Command, err error)
	UpdateCommand(device *Device, command *Command, timeout time.Duration) (err error)
	SubscribeCommands(device *Device, timestamp string, timeout time.Duration) (listener *CommandListener, err error)
	UnsubscribeCommands(device *Device, timeout time.Duration) (err error)

	GetNotification(device *Device, notificationID uint64, timeout time.Duration) (notification *Notification, err error)
	InsertNotification(device *Device, notification *Notification, timeout time.Duration) (err error)
}

// TODO: ClientService
