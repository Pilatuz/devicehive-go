package rest

import (
	"testing"

	dh "github.com/pilatuz/go-devicehive"
	"github.com/stretchr/testify/assert"
)

// Test InsertNotification and GetNotification methods
func TestNotificationInsertAndGet(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-rest"
	device.Name += "-rest"

	err := service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		i := 123 // t.Logf("device: %s", device)

		notification := dh.NewNotification("go-test-notification", i)
		err := service.InsertNotification(device, notification)
		assert.NoError(t, err, "Failed to insert notification")
		t.Logf("notification-A: %s", notification)

		notification, err = service.GetNotification(device, notification.ID)
		assert.NoError(t, err, "Failed to get notification")
		t.Logf("notification-B: %s", notification)
	}
}

// Test InsertNotification and PollNotification methods
func TestNotificationInsertAndPoll(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-rest"
	device.Name += "-rest"

	err = service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		i := 1234 // t.Logf("device: %s", device)

		notification := dh.NewNotification("go-test-notification", i)
		err := service.InsertNotification(device, notification)
		assert.NoError(t, err, "Failed to insert notification")
		t.Logf("sent notification: %s", notification)

		notifications, err := service.PollNotifications(device, info.Timestamp, "", "")
		assert.NoError(t, err, "Failed to poll notifications")
		assert.NotEmpty(t, notifications, "No any notifications polled")

		for _, c := range notifications {
			t.Logf("check notification: %s", c)
			if c.ID == notification.ID {
				return // OK
			}
		}

		assert.Fail(t, "Failed to poll notification")
	}
}

// Test InsertNotification and SubscribeNotifications methods
func TestNotificationInsertAndSubscribe(t *testing.T) {
	service := testNewREST(t)
	if service == nil {
		return // nothing to test
	}
	defer service.Stop()

	info, err := service.GetServerInfo()
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	device := testNewDevice()
	device.Network = testNewNetwork()
	device.ID += "-rest"
	device.Name += "-rest"

	err = service.RegisterDevice(device)
	if assert.NoError(t, err, "Failed to register device") {
		i := 12345 // t.Logf("device: %s", device)

		listener, err := service.SubscribeNotifications(device, info.Timestamp)
		assert.NoError(t, err, "Failed to subscribe notifications")
		assert.NotNil(t, listener, "No notification listener available")
		defer func() {
			err := service.UnsubscribeNotifications(device)
			assert.NoError(t, err, "Failed to unsubscribe notifications")
		}()

		a := dh.NewNotification("go-test-notification", i)
		err = service.InsertNotification(device, a)
		assert.NoError(t, err, "Failed to insert notification")

		b := <-listener.C // wait for notification polled
		assert.NotNil(t, b, "No any notification polled")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "unexpected notification polled")
		return
	}
}
