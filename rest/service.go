package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	dh "github.com/pilatuz/go-devicehive"
)

var (
	// package logger instance
	log = logrus.New()

	// TAG is a log prefix
	TAG = "DH-REST"

	// indicates stop
	errorStopped = fmt.Errorf("stopped")
)

// Service is a REST service for devices (TODO: and for clients).
type Service struct {
	// Base URL.
	baseURL *url.URL

	// Access key, might be empty - means no access key authorizathion used.
	accessKey string

	// HTTP client is used to perform all requests
	client *http.Client

	// set of command/notification listeners
	commandListenerLock      sync.Mutex
	commandListeners         map[string]*dh.CommandListener
	notificationListenerLock sync.Mutex
	notificationListeners    map[string]*dh.NotificationListener
	PollRetryTimeout         time.Duration

	stopped uint32
	stop    chan interface{}
	wg      sync.WaitGroup

	// default operation timeout
	DefaultTimeout time.Duration
}

// Get string representation of a service.
func (service *Service) String() string {
	return fmt.Sprintf("RestService{url:%q}", service.baseURL)
}

// NewService creates new service.
func NewService(baseURL, accessKey string) (*Service, error) {
	log.WithField("url", baseURL).Debugf("[%s]: creating service", TAG)

	var err error
	service := new(Service)
	service.accessKey = accessKey

	// remove trailing slashes from URL
	for len(baseURL) > 1 && strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1] // strings.TrimSuffix(baseURL, "/")
	}

	// parse URL
	if service.baseURL, err = url.Parse(baseURL); err != nil {
		log.WithError(err).Warnf("[%s]: failed to parse URL", TAG)
		return nil, fmt.Errorf("failed to parse URL: %s", err)
	}

	// initialize HTTP client
	service.client = new(http.Client)
	// TODO: client.Transport
	// TODO: client.CookieJar
	// TODO: client.Timeout

	// create empty set of listeners
	service.PollRetryTimeout = 1 * time.Second
	service.commandListeners = make(map[string]*dh.CommandListener)
	service.notificationListeners = make(map[string]*dh.NotificationListener)

	// create stop channel
	service.stop = make(chan interface{})

	// default timeout
	service.DefaultTimeout = 60 * time.Second

	return service, nil // OK
}

// Stop stops all active requests and polling loops
func (service *Service) Stop() {
	if atomic.CompareAndSwapUint32(&service.stopped, 0, 1) {
		log.Debugf("[%s]: stopping service...", TAG)

		// close channel
		close(service.stop)

		// clear all command listeners
		func() {
			service.commandListenerLock.Lock()
			defer service.commandListenerLock.Unlock()
			for ID := range service.commandListeners {
				service.removeCommandListenerUnsafe(ID)
			}
		}()

		// clear all notification listeners
		func() {
			service.notificationListenerLock.Lock()
			defer service.notificationListenerLock.Unlock()
			for ID := range service.notificationListeners {
				service.removeNotificationListenerUnsafe(ID)
			}
		}()

		log.Debugf("[%s]: waiting...", TAG)
		service.wg.Wait()
		log.Debugf("[%s]: service stopped", TAG)
	}
}

// SetTimeout sets the default timeout
func (service *Service) SetTimeout(timeout time.Duration) {
	log.WithField("timeout", timeout).Infof("[%s]: default timeout changed", TAG)
	service.DefaultTimeout = timeout
}

// check is the service stopped?
func (service *Service) isStopped() bool {
	return atomic.LoadUint32(&service.stopped) > 0
}

// Adds Authorization header if access key is not empty, device might be nil.
func (service *Service) prepareAuthorization(request *http.Request, device *dh.Device) {
	// access key
	if len(service.accessKey) != 0 {
		request.Header.Add("Authorization", "Bearer "+service.accessKey)
	}

	// device id+key
	if device != nil && (len(device.ID) != 0 || len(device.Key) != 0) {
		request.Header.Add("Auth-DeviceID", device.ID)
		request.Header.Add("Auth-DeviceKey", device.Key)
	}
}

// find command listener
func (service *Service) findCommandListener(deviceID string) *dh.CommandListener {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	if listener, ok := service.commandListeners[deviceID]; ok {
		return listener
	}
	return nil // not found
}

// insert new command listener
func (service *Service) insertCommandListener(deviceID string, listener *dh.CommandListener) {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	service.commandListeners[deviceID] = listener
}

// remove command listener
func (service *Service) removeCommandListener(deviceID string) {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	service.removeCommandListenerUnsafe(deviceID)
}

// remove command listener (without mutex lock)
func (service *Service) removeCommandListenerUnsafe(deviceID string) {
	if listener, ok := service.commandListeners[deviceID]; ok {
		delete(service.commandListeners, deviceID)
		close(listener.C)
	}
}

// find notification listener
func (service *Service) findNotificationListener(deviceID string) *dh.NotificationListener {
	service.notificationListenerLock.Lock()
	defer service.notificationListenerLock.Unlock()
	if listener, ok := service.notificationListeners[deviceID]; ok {
		return listener
	}
	return nil // not found
}

// insert new notification listener
func (service *Service) insertNotificationListener(deviceID string, listener *dh.NotificationListener) {
	service.notificationListenerLock.Lock()
	defer service.notificationListenerLock.Unlock()
	service.notificationListeners[deviceID] = listener
}

// remove notification listener
func (service *Service) removeNotificationListener(deviceID string) {
	service.notificationListenerLock.Lock()
	defer service.notificationListenerLock.Unlock()
	service.removeNotificationListenerUnsafe(deviceID)
}

// remove notification listener (without mutex lock)
func (service *Service) removeNotificationListenerUnsafe(deviceID string) {
	if listener, ok := service.notificationListeners[deviceID]; ok {
		delete(service.notificationListeners, deviceID)
		close(listener.C)
	}
}

// log returns task related log entry.
func (task *Task) log() *logrus.Entry {
	return log.WithField("task", task.identifier)
}

// SetLogLevel changes the package log level.
func SetLogLevel(level string) (err error) {
	log.Level, err = logrus.ParseLevel(level)
	return
}
