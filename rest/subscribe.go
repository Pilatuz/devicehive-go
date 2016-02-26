package rest

import (
	"time"

	"github.com/pilatuz/go-devicehive"
)

// SubscribeCommands adds a new command listener
func (service *Service) SubscribeCommands(device *devicehive.Device, timestamp string, timeout time.Duration) (*devicehive.CommandListener, error) {
	if listener, ok := service.commandListeners[device.ID]; ok {
		return listener, nil // already exists
	}

	// install new
	listener := devicehive.NewCommandListener(64)  // TODO: dedicated variable for buffer size
	service.commandListeners[device.ID] = listener // TODO: use mutex here!

	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start command polling", TAG)
		defer log.WithField("ID", deviceID).Debugf("[%s]: stop command polling", TAG)

		for {
			const names = ""
			const wait = "30"
			commands, err := service.PollCommands(device, timestamp, names, wait, 60*time.Second)
			if err != nil {
				log.WithError(err).Warnf("[%s]: failed to poll commands", TAG)
				// TODO: break? wait and try again?
				time.Sleep(1 * time.Second) // sleep a while
			}
			if listener, ok := service.commandListeners[deviceID]; ok {
				for _, command := range commands {
					log.WithField("command", command).Infof("[%s]: new command received", TAG)
					timestamp = command.Timestamp // continue with the latest command timestamp!
					listener.C <- command
				}
			} else {
				return // stop
			}
		}
	}(device.ID, timestamp)

	return listener, nil // OK
}

// UnsubscribeCommands removes the command listener
func (service *Service) UnsubscribeCommands(device *devicehive.Device, timeout time.Duration) error {
	// TODO: use mutex here!
	delete(service.commandListeners, device.ID) // poll loop will be stopped

	return nil // OK
}

// SubscribeNotifications adds a new notification listener
func (service *Service) SubscribeNotifications(device *devicehive.Device, timestamp string, timeout time.Duration) (*devicehive.NotificationListener, error) {
	if listener, ok := service.notificationListeners[device.ID]; ok {
		return listener, nil // already exists
	}

	// install new
	listener := devicehive.NewNotificationListener(64)  // TODO: dedicated variable for buffer size
	service.notificationListeners[device.ID] = listener // TODO: use mutex here!

	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start notification polling", TAG)
		defer log.WithField("ID", deviceID).Debugf("[%s]: stop notification polling", TAG)

		for {
			const names = ""
			const wait = "30"
			notifications, err := service.PollNotifications(device, timestamp, names, wait, 60*time.Second)
			if err != nil {
				log.WithError(err).Warnf("[%s]: failed to poll notifications", TAG)
				// TODO: break? wait and try again?
				time.Sleep(1 * time.Second) // sleep a while
			}
			if listener, ok := service.notificationListeners[deviceID]; ok {
				for _, notification := range notifications {
					log.WithField("notification", notification).Infof("[%s]: new notification received", TAG)
					timestamp = notification.Timestamp // continue with the latest notification timestamp!
					listener.C <- notification
				}
			} else {
				return // stop
			}
		}
	}(device.ID, timestamp)

	return listener, nil // OK
}

// UnsubscribeNotifications removes the notification listener
func (service *Service) UnsubscribeNotifications(device *devicehive.Device, timeout time.Duration) error {
	// TODO: use mutex here!
	delete(service.notificationListeners, device.ID) // poll loop will be stopped

	return nil // OK
}
