package ws

import (
	"testing"

	dh "github.com/pilatuz/devicehive-go"
	"github.com/stretchr/testify/assert"
)

// Test InsertCommand and UpdateCommand methods
func TestCommandInsertAndUpdate(t *testing.T) {
	service := testNewWsDevice(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	client := testNewWsClient(t)
	if client == nil {
		return // nothing to test
	}
	defer client.Stop()

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-ws"
	device.Name += "-ws"

	err := client.Authenticate(device)
	assert.NoError(t, err, "Failed to authenticate device")

	err = service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		i := 123 // t.Logf("device: %s", device)

		a := dh.NewCommand("go-test-command", i)
		err := client.InsertCommand(device, a)
		assert.NoError(t, err, "Failed to insert command")
		t.Logf("command-A: %s", a)

		c := dh.NewCommandResult(a.ID, "OK", i)
		service.UpdateCommand(device, c)
		assert.NoError(t, err, "Failed to update command")
		t.Logf("command-C: %s", c)
	}
}

// Test InsertCommand and SubscribeCommands methods
func TestCommandInsertAndSubscribe(t *testing.T) {
	service := testNewWsDevice(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	client := testNewWsClient(t)
	if client == nil {
		return // nothing to test
	}
	defer client.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-ws"
	device.Name += "-ws"

	err = client.Authenticate(device)
	assert.NoError(t, err, "Failed to authenticate device")

	err = service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		i := 12345 // t.Logf("device: %s", device)

		listener, err := service.SubscribeCommands(device, info.Timestamp)
		assert.NoError(t, err, "Failed to subscribe commands")
		assert.NotNil(t, listener, "No command listener available")
		defer func() {
			err := service.UnsubscribeCommands(device)
			assert.NoError(t, err, "Failed to unsubscribe commands")
		}()

		a := dh.NewCommand("go-test-command", i)
		err = client.InsertCommand(device, a)
		assert.NoError(t, err, "Failed to insert command")

		b := <-listener.C // wait for command polled
		assert.NotNil(t, b, "No any commands polled")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "unexpected command polled")
	}
}
