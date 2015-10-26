package core

import (
	"encoding/json"
	"testing"
)

// Test Command JSON marshaling
func TestJsonCommand(t *testing.T) {
	check := func(command Command, expectedJson string) {
		jsonBytes, err := json.Marshal(command)
		if err != nil {
			t.Errorf("Command: Cannot convert %+v to JSON (error: %s)", command, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("Command: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				command, string(jsonBytes), expectedJson)
		}

		t.Logf("Command: %+v converted to %s", command, string(jsonBytes))
	}

	command := Command{
		Name:       "cmd-name",
		Timestamp:  "2005-10-22",
		Parameters: "hello"}

	check(command, `{"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello"}`)

	command.Result = "custom data"
	command.Status = "done"
	command.Id = 100
	check(command, `{"id":100,"timestamp":"2005-10-22","command":"cmd-name","parameters":"hello","status":"done","result":"custom data"}`)
}
