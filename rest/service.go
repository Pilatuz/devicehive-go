package rest

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// REST service.
type Service struct {
	// Base URL.
	baseUrl *url.URL

	// Access key, might be empty - means no access key authorizathion used.
	accessKey string

	// HTTP client is used to perform all requests
	client *http.Client

	// set of command/notification listeners
	commandListeners map[string]*core.CommandListener
	notificationListeners map[string]*core.NotificationListener
}

// Get string representation of a service.
func (service *Service) String() string {
	return fmt.Sprintf("RestService{baseUrl:%q, accessKey:%q}",
		service.baseUrl, service.accessKey)
}

// NewService creates new service.
func NewService(baseUrl, accessKey string) (service *Service, err error) {
	log.Tracef("REST: creating service (url:%q)", baseUrl)
	service = &Service{accessKey: accessKey}

	// remove trailing slashes from URL
	for len(baseUrl) > 1 && strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[0 : len(baseUrl)-1]
	}

	// parse URL
	service.baseUrl, err = url.Parse(baseUrl)
	if err != nil {
		log.Warnf("REST: failed to parse URL (error: %s)", err)
		return
	}

	// initialize HTTP client
	service.client = &http.Client{}
	// TODO: client.Transport
	// TODO: client.CookieJar
	// TODO: client.Timeout

	service.commandListeners = make(map[string]*core.CommandListener)
	service.notificationListeners = make(map[string]*core.NotificationListener)
	return
}

// Adds Authorization header if access key is not empty
func (service *Service) prepareAuthorization(request *http.Request, device *core.Device) {
	// access key
	if len(service.accessKey) != 0 {
		request.Header.Add("Authorization", "Bearer "+service.accessKey)
	}

	// device id+key
	if device != nil && (len(device.Id) != 0 || len(device.Key) != 0) {
		request.Header.Add("Auth-DeviceID", device.Id)
		request.Header.Add("Auth-DeviceKey", device.Key)
	}
}

// Asynchronous request/task
type Task struct {
	request  *http.Request
	response *http.Response
	body     []byte
	err      error
}

// Do a request/task asynchronously
func (service *Service) doAsync(task Task) <-chan Task {
	ch := make(chan Task, 1)

	go func() {
		defer func() { ch <- task }()

		log.Tracef("REST: sending: %+v", task.request)
		task.response, task.err = service.client.Do(task.request)
		if task.err != nil {
			log.Warnf("REST: failed to do %s %s request (error: %s)",
				task.request.Method, task.request.URL, task.err)
			return
		}
		log.Tracef("REST: got %s %s response: %+v",
			task.request.Method, task.request.URL, task.response)

		// read body
		defer task.response.Body.Close()
		task.body, task.err = ioutil.ReadAll(task.response.Body)
		if task.err != nil {
			log.Warnf("REST: failed to read %s %s response body (error: %s)",
				task.request.Method, task.request.URL, task.err)
			return
		}

		log.Debugf("REST: got %s %s body: %s",
			task.request.Method, task.request.URL, string(task.body))
	}()

	return ch
}

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
