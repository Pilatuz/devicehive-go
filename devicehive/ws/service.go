package ws

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"github.com/gorilla/websocket"

	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TODO: support /client websocket endpoint

// Websocket service representation.
type Service struct {
	// Base URL.
	baseUrl *url.URL

	// Access key, might be empty - means no access key authorizathion used.
	accessKey string

	// Websocket connection
	conn *websocket.Conn

	// active task set
	taskLock   sync.Mutex
	lastTaskId uint64
	tasks      map[uint64]*Task

	// command listeners
	commandListenerLock sync.Mutex
	commandListeners    map[string]*core.CommandListener

	// transmitter
	tx chan *Task
}

// Get string representation of a Websocket service.
func (s Service) String() string {
	return fmt.Sprintf("WebsocketService{baseUrl:%q, accessKey:%q}", s.baseUrl, s.accessKey)
}

func (service *Service) GetCommand(device *core.Device, commandId uint64, timeout time.Duration) (command *core.Command, err error) {
	return &core.Command{}, nil
}

func (service *Service) InsertCommand(device *core.Device, command *core.Command, timeout time.Duration) (err error) {
	return nil
}

func (service *Service) GetNotification(device *core.Device, notificationId uint64, timeout time.Duration) (notification *core.Notification, err error) {
	return &core.Notification{}, nil
}

// find command listener
func (service *Service) findCommandListener(deviceId string) *core.CommandListener {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	listener := service.commandListeners[deviceId]
	return listener
}

// insert new command listener
func (service *Service) insertCommandListener(deviceId string, listener *core.CommandListener) {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	service.commandListeners[deviceId] = listener
}

// remove command listener
func (service *Service) removeCommandListener(deviceId string) {
	service.commandListenerLock.Lock()
	defer service.commandListenerLock.Unlock()
	delete(service.commandListeners, deviceId)
}

// NewService creates new Websocket /device service.
func NewService(baseUrl, accessKey string) (service *Service, err error) {
	log.Tracef("WS: creating service (url:%q)", baseUrl)
	service = &Service{accessKey: accessKey}

	// remove trailing slashes from URL
	for len(baseUrl) > 1 && strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[0 : len(baseUrl)-1]
	}

	// parse URL
	service.baseUrl, err = url.Parse(baseUrl)
	if err != nil {
		log.Warnf("WS: failed to parse URL (error: %s)", err)
		return
	}

	// connect to /device endpoint
	ws_url := fmt.Sprintf("%s/device", service.baseUrl)
	origin := "http://localhost/"
	service.conn, _, err = websocket.DefaultDialer.Dial(ws_url,
		http.Header{"Origin": []string{origin}})
	if err != nil {
		log.Warnf("WS: failed to dial (error: %s)", err)
		return
	}

	// set of active tasks
	service.tasks = make(map[uint64]*Task)

	// command listeners
	service.commandListeners = make(map[string]*core.CommandListener)

	// create TX channel
	service.tx = make(chan *Task)

	// and start RX/TX threads
	go service.doRX()
	go service.doTX()

	return
}

// TX thread
func (service *Service) doTX() {
	for {
		select {
		case task, ok := <-service.tx:
			if !ok || task == nil {
				log.Infof("WS: TX thread stopped")
				service.conn.Close() // TODO: send Close frame?
				return
			}

			body, err := task.Format()
			if err != nil {
				log.Warnf("WS: failed to format message (error: %s)", err)
				continue // TODO: return?
			}

			log.Tracef("WS: sending message: %s", string(body))
			err = service.conn.WriteMessage(websocket.TextMessage, body)
			if err != nil {
				log.Warnf("WS: failed to send message (error: %s)", err)
				continue // TODO: return?
			}

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
	for {
		_, body, err := service.conn.ReadMessage()
		if err != nil {
			log.Warnf("WS: failed to receive message (error: %s)", err)
			return
		}
		log.Tracef("WS: received message: %s", string(body))

		// parse JSON
		var msg map[string]interface{}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Warnf("WS: failed to parse JSON (error: %s), ignored", err)
			continue
		}

		// log.Tracef("WS: parsed JSON %+v", msg)
		if v, ok := msg["requestId"]; ok {
			id := safeUint64(v)
			task := service.takeTask(id)
			if task != nil {
				task.dataRecved = msg
				task.done <- task
			}
		}

		if v, ok := msg["action"]; ok {
			action := safeString(v)
			service.handleAction(action, msg)
		}
	}
}

// handle asynchronous actions
func (service *Service) handleAction(action string, data map[string]interface{}) {
	switch action {
	case "command/insert":
		if v, ok := data["deviceGuid"]; ok {
			deviceId := safeString(v)
			listener := service.findCommandListener(deviceId)
			if listener != nil {
				command := &core.Command{}
				err := command.AssignJSON(data["command"])
				if err != nil {
					log.Warnf("WS: failed to parse commnad/insert body (error: %s)", err)
					return
				}
				listener.C <- command
			} else {
				log.Warnf("WS: no command listener installed, %v ignored", data)
			}
		} else {
			log.Warnf("WS: no deviceId provided for command/insert, %v ignored", data)
		}
	default:
		log.Warnf("WS: unexpected action received: %v, ignored", data)
	}
}

// get uint64
func safeUint64(v interface{}) uint64 {
	switch x := v.(type) {
	case float64:
		return uint64(x)

	case uint64:
		return x

	// TODO: add other types

	default:
		log.Warnf("WS: unable to convert %v to uint64", v)
		return 0
	}
}

// get string
func safeString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x

	case float64:
		return fmt.Sprintf("%g", x)

	// TODO: add other types

	default:
		log.Warnf("WS: unable to convert %v to string", v)
		return ""
	}
}
