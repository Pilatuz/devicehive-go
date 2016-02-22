// DeviceHive abstract service interface.
package devicehive

import (
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"github.com/devicehive/devicehive-go/devicehive/rest"
	"github.com/devicehive/devicehive-go/devicehive/ws"
	"strings"
	"time"
)

const (
	// Datetime layout used for timestamps
	DateTimeLayout = core.DateTimeLayout
)

// Abstract DeviceHive /device API.
type Service interface {
	GetServerInfo(timeout time.Duration) (info *core.ServerInfo, err error)

	RegisterDevice(device *core.Device, timeout time.Duration) (err error)
	GetDevice(deviceId, deviceKey string, timeout time.Duration) (device *core.Device, err error)

	GetCommand(device *core.Device, commandId uint64, timeout time.Duration) (command *core.Command, err error)
	UpdateCommand(device *core.Device, command *core.Command, timeout time.Duration) (err error)
	SubscribeCommands(device *core.Device, timestamp string, timeout time.Duration) (listener *core.CommandListener, err error)
	UnsubscribeCommands(device *core.Device, timeout time.Duration) (err error)

	GetNotification(device *core.Device, notificationId uint64, timeout time.Duration) (notification *core.Notification, err error)
	InsertNotification(device *core.Device, notification *core.Notification, timeout time.Duration) (err error)
}

// NewRestService creates a new REST service.
// Base REST URL should be provided.
// Access key is optional, might be empty.
func NewRestService(baseUrl, accessKey string) (service Service, err error) {
	return rest.NewService(baseUrl, accessKey)
}

// NewWebsocketService creates a new Websocket service.
// Base Websocket URL should be provided.
// Access key is optional, might be empty.
func NewWebsocketService(baseUrl, accessKey string) (service Service, err error) {
	return ws.NewService(baseUrl, accessKey)
}

// NewService creates a new service (either REST or Websocket).
// Base URL should be provided.
// If protocol is "ws://" or "wss://" Websocket service will be created,
// otherwise REST service will be used as a fallback.
// Access key is optional, might be empty.
func NewService(baseUrl, accessKey string) (service Service, err error) {
	url := strings.ToLower(baseUrl)
	if strings.HasPrefix(url, `ws://`) || strings.HasPrefix(url, `wss://`) {
		return NewWebsocketService(baseUrl, accessKey)
	}

	// use REST service as a fallback
	return NewRestService(baseUrl, accessKey)
}

// NewDevice creates a new device without network.
// No user data by default.
func NewDevice(id, name string, class *core.DeviceClass) *core.Device {
	return core.NewDevice(id, name, class)
}

// NewDeviceWithNetwork creates a new device with network initialized.
// No user data by default.
func NewDeviceWithNetwork(id, name string, class *core.DeviceClass, network *core.Network) *core.Device {
	return core.NewDeviceWithNetwork(id, name, class, network)
}

// NewDeviceClass creates a new device class.
// No user data and no equipment by default.
func NewDeviceClass(name, version string) *core.DeviceClass {
	return core.NewDeviceClass(name, version)
}

// NewEquipment creates a new equipment.
// No user data by default.
func NewEquipment(name, code, type_ string) *core.Equipment {
	return core.NewEquipment(name, code, type_)
}

// NewNetwork creates a new network.
// Network Description is empty.
func NewNetwork(name, key string) *core.Network {
	return core.NewNetwork(name, key)
}

// NewCommand creates a new command.
// No lifetime set by default.
func NewCommand(name string, parameters interface{}) *core.Command {
	return core.NewCommand(name, parameters)
}

// NewCommandResult creates a new command result.
func NewCommandResult(id uint64, status string, result interface{}) *core.Command {
	return core.NewCommandResult(id, status, result)
}

// NewNotification creates a new notification.
func NewNotification(name string, parameters interface{}) *core.Notification {
	return core.NewNotification(name, parameters)
}

// SetLogLevel changes the global logging level:
// "WARN" "INFO" "DEBUG" "TRACE"
func SetLogLevel(level string) {
	log.SetLevelByName(level)
}
