package ws

import (
	"fmt"
	"github.com/devicehive/devicehive-go/devicehive/core"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

// Prepare GetServerInfo task
func (service *Service) prepareGetServerInfo() (task *Task, err error) {
	task = service.newTask()
	task.dataToSend = map[string]interface{}{
		"action":    "server/info",
		"requestId": task.id}
	return
}

// Process GetServerInfo task
func (service *Service) processGetServerInfo(task *Task, info *core.ServerInfo) (err error) {
	// check response status
	err = task.CheckStatus()
	if err != nil {
		log.Warnf("WS: bad /info status (error: %s)", err)
		return
	}

	// parse response
	err = info.AssignJSON(task.dataRecved["info"])
	if err != nil {
		log.Warnf("WS: failed to parse /info body (error: %s)", err)
		return
	}

	return
}

// GetServerInfo() function gets the main server's information.
func (service *Service) GetServerInfo(timeout time.Duration) (info *core.ServerInfo, err error) {
	task, err := service.prepareGetServerInfo()
	if err != nil {
		log.Warnf("WS: failed to prepare /info task (error: %s)", err)
		return
	}

	// add to the TX pipeline
	service.tx <- task

	select {
	case <-time.After(timeout):
		log.Warnf("WS: failed to wait %s for /info task", timeout)
		err = fmt.Errorf("timed out")

	case <-task.done:
		info = &core.ServerInfo{}
		err = service.processGetServerInfo(task, info)
		if err != nil {
			log.Warnf("WS: failed to process /info task (error: %s)", err)
			return
		}
	}

	return
}
