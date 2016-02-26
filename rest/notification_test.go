package rest

import (
	"testing"

	"github.com/pilatuz/go-devicehive"
	"github.com/stretchr/testify/assert"
)

// Test InsertNotification and GetNotification methods
func TestNotificationInsertAndGet(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	devices, err := service.GetDeviceList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	for i, device := range devices {
		t.Logf("device-%d: %s", i, device)

		notification := devicehive.NewNotification("go-test-notification", i)
		err := service.InsertNotification(device, notification, testWaitTimeout)
		assert.NoError(t, err, "Failed to insert notification")
		t.Logf("notification-A: %s", notification)

		notification, err = service.GetNotification(device, notification.ID, testWaitTimeout)
		assert.NoError(t, err, "Failed to get notification")
		t.Logf("notification-B: %s", notification)
	}
}

// Test InsertNotification and PollNotification methods
func TestNotificationInsertAndPoll(t *testing.T) {
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	info, err := service.GetServerInfo(testWaitTimeout)
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	devices, err := service.GetDeviceList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	// TODO: register and delete dedicated device!

	for i, device := range devices {
		t.Logf("device-%d: %s", i, device)

		notification := devicehive.NewNotification("go-test-notification", i)
		err := service.InsertNotification(device, notification, testWaitTimeout)
		assert.NoError(t, err, "Failed to insert notification")
		t.Logf("sent notification: %s", notification)

		notifications, err := service.PollNotifications(device, info.Timestamp, "", "", testWaitTimeout)
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
	service := testNewRest(t)
	if service == nil {
		return // nothing to test
	}

	info, err := service.GetServerInfo(testWaitTimeout)
	assert.NoError(t, err, "Failed to get server info")
	assert.NotEmpty(t, info.Timestamp, "No server timestamp avaialble")

	devices, err := service.GetDeviceList(0, 0, testWaitTimeout)
	assert.NoError(t, err, "Failed to get list of devices")
	assert.NotEmpty(t, devices, "No any device available")

	// TODO: register and delete dedicated device!

	for i, device := range devices {
		// t.Logf("device-%d: %s", i, device)

		listener, err := service.SubscribeNotifications(device, info.Timestamp, testWaitTimeout)
		assert.NoError(t, err, "Failed to subscribe notifications")
		assert.NotNil(t, listener, "No notification listener available")
		defer func() {
			err := service.UnsubscribeNotifications(device, testWaitTimeout)
			assert.NoError(t, err, "Failed to unsubscribe notifications")
		}()

		a := devicehive.NewNotification("go-test-notification", i)
		err = service.InsertNotification(device, a, testWaitTimeout)
		assert.NoError(t, err, "Failed to insert notification")

		b := <-listener.C // wait for notification polled
		assert.NotNil(t, b, "No any notification polled")
		assert.JSONEq(t, toJsonStr(a), toJsonStr(b), "unexpected notification polled")
		return
	}
}
