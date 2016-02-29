package ws

import (
	"testing"

	dh "github.com/pilatuz/devicehive-go"
	"github.com/stretchr/testify/assert"
)

// Test InsertNotification method
func TestNotificationInsert(t *testing.T) {
	service := testNewWsDevice(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-ws"
	device.Name += "-ws"

	err := service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		notification := dh.NewNotification("go-test-notification", 123)
		err := service.InsertNotification(device, notification)
		assert.NoError(t, err, "Failed to insert notification")
		t.Logf("notification: %s", notification)
	}
}

// Test InsertNotification and SubscribeNotifications methods
func TestNotificationInsertAndSubscribe(t *testing.T) {
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

		listener, err := client.SubscribeNotifications(device, info.Timestamp)
		assert.NoError(t, err, "Failed to subscribe notifications")
		assert.NotNil(t, listener, "No notification listener available")
		defer func() {
			err := client.UnsubscribeNotifications(device)
			assert.NoError(t, err, "Failed to unsubscribe notifications")
		}()

		a := dh.NewNotification("go-test-notification", i)
		err = service.InsertNotification(device, a)
		assert.NoError(t, err, "Failed to insert notification")

		b := <-listener.C // wait for notification polled
		assert.NotNil(t, b, "No any notification polled")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "unexpected notification polled")
	}
}
