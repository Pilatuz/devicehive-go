package rest

import (
	"github.com/devicehive/devicehive-go/devicehive/core"
	"testing"
)

// TestInsertCommand() unit test for /command/insert POST method,
// /command/update PUT method, /command/get GET method
// test device should be already registered!
func TestInsertCommand(t *testing.T) {
	TestRegisterDevice(t)
	if t.Failed() {
		return // nothing to test without device
	}

	s, err := NewService(testServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create service (error: %s)", err)
		return
	}

	device := &core.Device{Id: testDeviceId, Key: testDeviceKey}
	command := &core.Command{Name: "cmd-test", Parameters: 123, Lifetime: 600}
	err = s.InsertCommand(device, command, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to insert command (error: %s)", err)
		return
	}
	t.Logf("command: %s", command)

	command.Status = "Done"
	command.Result = 12345
	err = s.UpdateCommand(device, command, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to update command (error: %s)", err)
		return
	}

	*command, err = s.GetCommand(device, command.Id, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get command (error: %s)", err)
		return
	}
	t.Logf("command: %s", command)
}

// TODO: TestPollCommand
