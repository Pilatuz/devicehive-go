package ws

/*

import (
	"testing"

	dh "github.com/pilatuz/go-devicehive"
	"github.com/stretchr/testify/assert"
)

// Test InsertCommand and GetCommand and UpdateCommand methods
func TestCommandInsertAndUpdate(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	devices, err := service.GetDeviceList(0, 0)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	for i, device := range devices {
		t.Logf("device-%d: %s", i, device)

		command := dh.NewCommand("go-test-command", i)
		err := service.InsertCommand(device, command)
		assert.NoError(t, err, "Failed to insert command")
		t.Logf("command-A: %s", command)

		command, err = service.GetCommand(device, command.ID)
		assert.NoError(t, err, "Failed to get command")
		t.Logf("command-B: %s", command)

		command = dh.NewCommandResult(command.ID, "OK", i)
		service.UpdateCommand(device, command)
		assert.NoError(t, err, "Failed to update command")
		t.Logf("command-C: %s", command)

		command, err = service.GetCommand(device, command.ID)
		assert.NoError(t, err, "Failed to get command")
		t.Logf("command-D: %s", command)
	}
}

// Test InsertCommand and PollCommand methods
func TestCommandInsertAndPoll(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	devices, err := service.GetDeviceList(0, 0)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	// TODO: register and delete dedicated device!

	for i, device := range devices {
		t.Logf("device-%d: %s", i, device)

		command := dh.NewCommand("go-test-command", i)
		err := service.InsertCommand(device, command)
		assert.NoError(t, err, "Failed to insert command")
		t.Logf("sent command: %s", command)

		commands, err := service.PollCommands(device, info.Timestamp, "", "")
		assert.NoError(t, err, "Failed to poll commands")
		assert.NotEmpty(t, commands, "No any commands polled")

		for _, c := range commands {
			t.Logf("check command: %s", c)
			if c.ID == command.ID {
				return // OK
			}
		}

		assert.Fail(t, "Failed to poll command")
	}
}

// Test InsertCommand and SubscribeCommands methods
func TestCommandInsertAndSubscribe(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	devices, err := service.GetDeviceList(0, 0)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	// TODO: register and delete dedicated device!

	for i, device := range devices {
		// t.Logf("device-%d: %s", i, device)

		listener, err := service.SubscribeCommands(device, info.Timestamp)
		assert.NoError(t, err, "Failed to subscribe commands")
		assert.NotNil(t, listener, "No command listener available")
		defer func() {
			err := service.UnsubscribeCommands(device)
			assert.NoError(t, err, "Failed to unsubscribe commands")
		}()

		a := dh.NewCommand("go-test-command", i)
		err = service.InsertCommand(device, a)
		assert.NoError(t, err, "Failed to insert command")

		b := <-listener.C // wait for command polled
		assert.NotNil(t, b, "No any commands polled")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "unexpected command polled")
		return
	}
}
*/
