package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	dh "github.com/pilatuz/go-devicehive"
)

var (
	// package logger instance
	log = logrus.New()

	// TAG is a log prefix
	TAG = "DH-WS"

	// indicates stop
	errorStopped = fmt.Errorf("stopped")
)

// Action is a structure to hold some common action fields.
type Action struct {
	RequestID uint32 `json:"requestId,omitempty"`
	Action    string `json:"action,omitempty"`

	// command/insert
	DeviceID string `json:"deviceGuid,omitempty"`

	// status + error
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	ErrMsg string `json:"error,omitempty"`
}

// Error checks if action contains error
func (action *Action) Error() error {
	if !strings.EqualFold(action.Status, "success") {
		return fmt.Errorf("unexpected status: %q [%s %s]",
			action.Status, action.Code, action.ErrMsg)
	}

	return nil // OK
}

// Websocket service representation.
type Service struct {
	// Base URL.
	baseURL *url.URL

	// Access key, might be empty - means no access key authorizathion used.
	accessKey string

	// Websocket connection
	conn *websocket.Conn

	// active task set
	taskLock   sync.Mutex
	lastTaskId uint32
	tasks      map[uint32]*Task

	// command listeners
	commandListenerLock sync.Mutex
	commandListeners    map[string]*dh.CommandListener

	// notification listeners
	notificationListenerLock sync.Mutex
	notificationListeners    map[string]*dh.NotificationListener

	// transmitter
	tx chan *Task

	stopped uint32
	stop    chan interface{}
	wg      sync.WaitGroup

	// default operation timeout
	DefaultTimeout time.Duration
}

// Get string representation of a Websocket service.
func (service *Service) String() string {
	return fmt.Sprintf("WebsocketService{url:%q}", service.baseURL)
}

// NewDeviceService creates new Websocket /device service.
func NewDeviceService(baseUrl, accessKey string) (*Service, error) {
	return newService(baseUrl, accessKey, "/device")
}

// NewClientService creates new Websocket service.
func NewClientService(baseUrl, accessKey string) (*Service, error) {
	return newService(baseUrl, accessKey, "/client")
}

// newService creates new Websocket service.
func newService(baseUrl, accessKey string, path string) (*Service, error) {
	log.WithField("url", baseUrl).Debugf("[%s]: creating %s service", TAG, path)

	var err error
	service := new(Service)
	service.accessKey = accessKey

	// remove trailing slashes from URL
	for len(baseUrl) > 1 && strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[0 : len(baseUrl)-1] // strings.TrimSuffix(baseURL, "/")
	}

	// parse URL
	if service.baseURL, err = url.Parse(baseUrl); err != nil {
		log.WithError(err).Warnf("[%s]: failed to parse URL", TAG)
		return nil, fmt.Errorf("failed to parse URL: %s", err)
	}

	// connect to /device or /client endpoint
	service.baseURL.Path += path
	headers := http.Header{}
	headers.Add("Origin", "http://localhost/")
	if len(service.accessKey) != 0 {
		headers.Add("Authorization", "Bearer "+service.accessKey)
	}
	service.conn, _, err = websocket.DefaultDialer.Dial(service.baseURL.String(), headers)
	if err != nil {
		log.WithError(err).Warnf("[%s]: failed to dial", TAG)
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	// set of active tasks
	service.tasks = make(map[uint32]*Task)

	// create empty set of listeners
	service.notificationListeners = make(map[string]*dh.NotificationListener)
	service.commandListeners = make(map[string]*dh.CommandListener)

	// create stop channel
	service.stop = make(chan interface{})

	// default timeout
	service.DefaultTimeout = 60 * time.Second

	// create TX channel
	service.tx = make(chan *Task, 64) // TODO: dedicated contant for buffer size

	// and start RX/TX threads
	service.wg.Add(2)
	go service.doRX()
	go service.doTX()

	return service, nil // OK
}

// Stop stops all active requests and working goroutines
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

		// close connection
		// TODO: send Close frame?
		service.conn.Close()

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

// TX thread
func (service *Service) doTX() {
	defer func() {
		log.Debugf("[%s]: TX thread stopped", TAG)
		service.wg.Done()
	}()

	for {
		select {
		case <-service.stop:
			return

		case task := <-service.tx:
			task.log().WithField("msg", task.dataToSend).Debugf("[%s]: sending message", TAG)
			body, err := json.Marshal(task.dataToSend)
			if err != nil {
				task.log().WithError(err).Warnf("[%s]: failed to format message", TAG)
				service.takeTask(task.identifier) // remove from "active" set
				task.ReportDone(nil, fmt.Errorf("failed to format message: %s", err))
				continue
			}

			// task.log().WithField("msg", string(body)).Debugf("[%s]: sending message", TAG)
			err = service.conn.WriteMessage(websocket.TextMessage, body)
			if err != nil {
				task.log().WithError(err).Warnf("[%s]: failed to send message", TAG)
				service.takeTask(task.identifier) // remove from "active" set
				task.ReportDone(nil, fmt.Errorf("failed to send message: %s", err))
				continue
			}

			// TODO: ping/pong messages
			//			case <-service.pingTimer.C:
			//			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
			//				log.Warnf("WS: could not write ping message (error: %s)", err)
			//				return
			//			}
		}
	}
}

// RX thread
func (service *Service) doRX() {
	defer func() {
		log.Debugf("[%s]: RX thread stopped", TAG)
		service.wg.Done()
	}()

	rx_chan := make(chan []byte, 16)
	var rx_err error

	// receive message in dedicated goroutine
	service.wg.Add(1)
	go func() {
		defer func() {
			log.Debugf("[%s]: *RX thread stopped", TAG)
			service.wg.Done()
		}()

		for {
			// read a message...
			_, body, err := service.conn.ReadMessage()
			if err != nil {
				rx_err = err
				close(rx_chan)
				return
			}

			// ...and pass it for processing
			rx_chan <- body
		}
	}()

	for {
		select {
		case <-service.stop:
			return

		case msg_body := <-rx_chan:
			if rx_err != nil {
				log.WithError(rx_err).Warnf("[%s]: failed to receive message", TAG)
				return // stop processing
			}

			// parse JSON to the map
			var msg map[string]interface{}
			err := json.Unmarshal(msg_body, &msg)
			if err != nil {
				log.WithField("msg", string(msg_body)).Debugf("[%s]: new message received", TAG)
				log.WithError(err).Warnf("[%s]: failed to parse JSON, ignored", TAG)
				continue
			}

			log.WithField("msg", msg).Debugf("[%s]: new message received", TAG)
			service.handleMessage(msg)
		}
	}
}

// handle received messages
func (service *Service) handleMessage(data map[string]interface{}) {
	// decode common fields
	msg := new(Action)
	if err := dh.FromJSON(msg, data); err != nil {
		log.WithError(err).Warnf("[%s]: failed to assign JSON, ignored", TAG)
		return
	}

	// handle pending requests first
	if task := service.takeTask(msg.RequestID); task != nil {
		task.ReportDone(data, msg.Error())
		return
	}

	// asynchronous actions
	switch msg.Action {
	case "command/insert", "command/update":
		if listener := service.findCommandListener(msg.DeviceID); listener != nil {
			command := new(dh.Command)
			err := command.FromMap(data["command"])
			if err != nil {
				log.WithError(err).Warnf("[%s]: failed to parse %q body, ignored", TAG, msg.Action)
				return
			}
			listener.C <- command
		} else {
			log.WithField("deviceId", msg.DeviceID).Warnf("[%s]: no command listener installed, ignored", TAG)
		}
	case "notification/insert":
		if listener := service.findNotificationListener(msg.DeviceID); listener != nil {
			notification := new(dh.Notification)
			err := notification.FromMap(data["notification"])
			if err != nil {
				log.WithError(err).Warnf("[%s]: failed to parse %q body, ignored", TAG, msg.Action)
				return
			}
			listener.C <- notification
		} else {
			log.WithField("deviceId", msg.DeviceID).Warnf("[%s]: no notification listener installed, ignored", TAG)
		}
	default:
		log.WithField("action", msg.Action).Warnf("[%s]: unknown action received, ignored", TAG)
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
