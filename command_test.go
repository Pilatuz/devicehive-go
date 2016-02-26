package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test JSON to string method.
func TestCommandString(t *testing.T) {
	command := NewCommand("cmd-name", "hello")
	command.Timestamp = "2005-10-22"
	assert.Equal(t, command.String(), `Command{Name:"cmd-name", Timestamp:"2005-10-22", Parameters:hello}`)

	command.ID = 100
	command.Result = "data"
	command.Status = "done"
	assert.Equal(t, command.String(), `Command{ID:100, Name:"cmd-name", Timestamp:"2005-10-22", Parameters:hello, Status:"done", Result:data}`)
}

// Test Command JSON marshaling.
func TestCommandJson(t *testing.T) {
	command := NewCommand("cmd-name", "hello")
	command.Timestamp = "2005-10-22"
	assert.JSONEq(t, toJsonStr(t, command), `{"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello"}`)

	command.ID = 100
	command.Result = "data"
	command.Status = "done"
	assert.JSONEq(t, toJsonStr(t, command), `{"id":100,"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello","status":"done","result":"data"}`)
}

// Test Command assignment.
func TestCommandAssign(t *testing.T) {
	command := NewEmptyCommand()
	assert.NoError(t, command.FromMap(map[string]interface{}{
		"command":    "cmd-name",
		"parameters": "hello",
		"timestamp":  "2005-10-22",
	}))
	assert.JSONEq(t, toJsonStr(t, command), `{"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello"}`)

	assert.NoError(t, command.FromMap(map[string]interface{}{
		"id":     "100",
		"result": "data",
		"status": 123,
	}))
	assert.JSONEq(t, toJsonStr(t, command), `{"id":100,"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello","status":"123","result":"data"}`)
}
