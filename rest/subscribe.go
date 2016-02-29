package rest

import (
	"time"

	dh "github.com/pilatuz/devicehive-go"
)

// SubscribeCommands adds a new command listener
func (service *Service) SubscribeCommands(device *dh.Device, timestamp string) (*dh.CommandListener, error) {
	if listener := service.findCommandListener(device.ID); listener != nil {
		return listener, nil // already exists
	}

	// install new
	listener := dh.NewCommandListener(64) // TODO: dedicated variable for buffer size
	service.insertCommandListener(device.ID, listener)

	service.wg.Add(1)
	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start command polling", TAG)
		defer func() {
			log.WithField("ID", deviceID).Debugf("[%s]: stop command polling", TAG)
			service.wg.Done()
		}()

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
			if listener := service.findCommandListener(deviceID); listener != nil {
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
	// poll loop will be stopped once listener is removed from the map
	service.removeCommandListener(device.ID)
	return nil // OK
}

// SubscribeNotifications adds a new notification listener
func (service *Service) SubscribeNotifications(device *dh.Device, timestamp string) (*dh.NotificationListener, error) {
	if listener := service.findNotificationListener(device.ID); listener != nil {
		return listener, nil // already exists
	}

	// install new
	listener := dh.NewNotificationListener(64) // TODO: dedicated variable for buffer size
	service.insertNotificationListener(device.ID, listener)

	service.wg.Add(1)
	go func(deviceID string, timestamp string) {
		log.WithField("ID", deviceID).Debugf("[%s]: start notification polling", TAG)
		defer func() {
			log.WithField("ID", deviceID).Debugf("[%s]: stop notification polling", TAG)
			service.wg.Done()
		}()

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
			if listener := service.findNotificationListener(deviceID); listener != nil {
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
	// poll loop will be stopped once listener is removed from the map
	service.removeNotificationListener(device.ID)
	return nil // OK
}
