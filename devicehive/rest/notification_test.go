package rest

import (
	"github.com/devicehive/devicehive-go/devicehive/core"
	"testing"
)

// TestInsertNotification() unit test for /notification/insert POST method
// and /notification/get GET method
func TestInsertNotification(t *testing.T) {
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
	notification := &core.Notification{Name: "ntf-test", Parameters: 12345}
	err = s.InsertNotification(device, notification, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to insert notification (error: %s)", err)
		return
	}
	t.Logf("notification: %s", notification)

	*notification, err = s.GetNotification(device, notification.Id, testWaitTimeout)
	if err != nil {
		t.Errorf("Failed to get notification (error: %s)", err)
		return
	}
	t.Logf("notification: %s", notification)
}

// TODO: TestPollNotification
