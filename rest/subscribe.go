// +build ignore

package rest

import (
	"time"
)

// subscribe for commands
func (service *Service) SubscribeCommands(device *core.Device, timestamp string, timeout time.Duration) (listener *core.CommandListener, err error) {
	if listener, ok := service.commandListeners[device.Id]; ok {
		return listener, nil
	}

	// install new
	listener = core.NewCommandListener()
	service.commandListeners[device.Id] = listener

	go func(deviceId string) {
		log.Debugf("REST: start command polling %q", deviceId)
		for {
			names := ""
			wait := "30"
			cmds, err := service.PollCommands(device, timestamp, names, wait, 60*time.Second)
			if err != nil {
				log.Warnf("REST: failed to poll commands (error: %s)", err)
				// TODO: break? wait and try again?
			}
			if listener, ok := service.commandListeners[deviceId]; ok {
				for _, cmd := range cmds {
					log.Debugf("REST: got command %s received", cmd)
					timestamp = cmd.Timestamp
					listener.C <- &cmd
				}
			} else {
				log.Debugf("REST: stop command polling %q", deviceId)
				return // stop
			}
		}
	}(device.Id)

	return
}

// unsubscribe from commands
func (service *Service) UnsubscribeCommands(device *core.Device, timeout time.Duration) (err error) {
	delete(service.commandListeners, device.Id) // poll loop will be stopped
	return nil
}

// subscribe for notifications
func (service *Service) SubscribeNotifications(device *core.Device, timestamp string, timeout time.Duration) (listener *core.NotificationListener, err error) {
	if listener, ok := service.notificationListeners[device.Id]; ok {
		return listener, nil
	}

	// install new
	listener = core.NewNotificationListener()
	service.notificationListeners[device.Id] = listener

	go func(deviceId string) {
		log.Debugf("REST: start notification polling %q", deviceId)
		for {
			names := ""
			wait := "30"
			ntfs, err := service.PollNotifications(device, timestamp, names, wait, 60*time.Second)
			if err != nil {
				log.Warnf("REST: failed to poll notifications (error: %s)", err)
				// TODO: break? wait and try again?
			}
			if listener, ok := service.notificationListeners[deviceId]; ok {
				for _, ntf := range ntfs {
					log.Debugf("REST: got notification %s received", ntf)
					timestamp = ntf.Timestamp
					listener.C <- &ntf
				}
			} else {
				log.Debugf("REST: stop notification polling %q", deviceId)
				return // stop
			}
		}
	}(device.Id)

	return
}

// unsubscribe from notifications
func (service *Service) UnsubscribeNotifications(device *core.Device, timeout time.Duration) (err error) {
	delete(service.notificationListeners, device.Id) // poll loop will be stopped
	return nil
}
