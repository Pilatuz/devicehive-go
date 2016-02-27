package rest

import (
	"time"

	dh "github.com/pilatuz/go-devicehive"
)

// SubscribeCommands adds a new command listener
func (service *Service) SubscribeCommands(device *dh.Device, timestamp string) (*dh.CommandListener, error) {
	if listener, ok := service.commandListeners[device.ID]; ok {
		return listener, nil // already exists
	}

	// install new
	listener := dh.NewCommandListener(64)          // TODO: dedicated variable for buffer size
	service.commandListeners[device.ID] = listener // TODO: use mutex here!

	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start command polling", TAG)
		defer log.WithField("ID", deviceID).Debugf("[%s]: stop command polling", TAG)

		for {
			const names = ""
			const wait = "30" // TODO: check wait < service.DefaultTimeout!
			commands, err := service.PollCommands(device, timestamp, names, wait)
			if err != nil {
				if err == errorStopped {
					// if service is stopped
					// just stop polling
					return
				}

				// sleep a while...
				log.WithError(err).Warnf("[%s]: failed to poll commands", TAG)
				select {
				case <-time.After(service.PollRetryTimeout):
					continue // ...and try again

				case <-service.stop:
					return // stop
				}
			}
			if listener, ok := service.commandListeners[deviceID]; ok {
				for _, command := range commands {
					log.WithField("command", command).Infof("[%s]: new command received", TAG)
					timestamp = command.Timestamp // continue with the latest command timestamp!
					// TODO: check timestamp consistency!
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
func (service *Service) UnsubscribeCommands(device *dh.Device) error {
	// TODO: use mutex here!
	if listener, ok := service.commandListeners[device.ID]; ok {
		delete(service.commandListeners, device.ID) // poll loop will be stopped
		close(listener.C)
	}

	return nil // OK
}

// SubscribeNotifications adds a new notification listener
func (service *Service) SubscribeNotifications(device *dh.Device, timestamp string) (*dh.NotificationListener, error) {
	if listener, ok := service.notificationListeners[device.ID]; ok {
		return listener, nil // already exists
	}

	// install new
	listener := dh.NewNotificationListener(64)          // TODO: dedicated variable for buffer size
	service.notificationListeners[device.ID] = listener // TODO: use mutex here!

	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start notification polling", TAG)
		defer log.WithField("ID", deviceID).Debugf("[%s]: stop notification polling", TAG)

		for {
			const names = ""
			const wait = "30" // TODO: check wait < service.DefaultTimeout!
			notifications, err := service.PollNotifications(device, timestamp, names, wait)
			if err != nil {
				if err == errorStopped {
					// if service is stopped
					// just stop polling
					return
				}

				// sleep a while...
				log.WithError(err).Warnf("[%s]: failed to poll notifications", TAG)
				select {
				case <-time.After(service.PollRetryTimeout):
					continue // ...and try again

				case <-service.stop:
					return // stop
				}
			}
			if listener, ok := service.notificationListeners[deviceID]; ok {
				for _, notification := range notifications {
					log.WithField("notification", notification).Infof("[%s]: new notification received", TAG)
					timestamp = notification.Timestamp // continue with the latest notification timestamp!
					// TODO: check timestamp consistency!
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
func (service *Service) UnsubscribeNotifications(device *dh.Device) error {
	// TODO: use mutex here!
	if listener, ok := service.notificationListeners[device.ID]; ok {
		delete(service.notificationListeners, device.ID) // poll loop will be stopped
		close(listener.C)
	}

	return nil // OK
}
