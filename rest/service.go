package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pilatuz/go-devicehive"
)

var (
	// package logger instance
	log = logrus.New()
	TAG = "DH-REST"
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
	commandListeners      map[string]*devicehive.CommandListener
	notificationListeners map[string]*devicehive.NotificationListener
}

// Get string representation of a service.
func (service *Service) String() string {
	return fmt.Sprintf("RestService{baseUrl:%q}", service.baseUrl)
}

// NewService creates new service.
func NewService(baseUrl, accessKey string) (*Service, error) {
	log.Debugf("[%s]: creating service url: %s", TAG, baseUrl)

	service := new(Service)
	service.accessKey = accessKey

	// remove trailing slashes from URL
	for len(baseUrl) > 1 && strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[0 : len(baseUrl)-1]
	}

	// parse URL
	if u, err := url.Parse(baseUrl); err != nil {
		log.Warnf("[%s]: failed to parse URL: %s", TAG, err)
		return nil, fmt.Errorf("failed to parse URL: %s", err)
	} else {
		service.baseUrl = u
	}

	// initialize HTTP client
	service.client = new(http.Client)
	// TODO: client.Transport
	// TODO: client.CookieJar
	// TODO: client.Timeout

	// create empty set of listeners
	service.commandListeners = make(map[string]*devicehive.CommandListener)
	service.notificationListeners = make(map[string]*devicehive.NotificationListener)

	return service, nil // OK
}

// Adds Authorization header if access key is not empty
// device might be nil.
func (service *Service) prepareAuthorization(request *http.Request, device *devicehive.Device) {
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

// log returns task related log entry.
func (task *asyncTask) log() *logrus.Entry {
	return log.WithField("task", task.Identifier)
}

// Change the package log level
func SetLogLevel(level string) (err error) {
	log.Level, err = logrus.ParseLevel(level)
	return
}
