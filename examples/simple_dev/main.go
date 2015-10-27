package main

import (
	"github.com/devicehive/devicehive-go/devicehive"
	"github.com/devicehive/devicehive-go/devicehive/log"
	"github.com/devicehive/devicehive-go/devicehive/rest"
	"time"
	"os/signal"
	"os"
)

func main() {
	log.SetLevel(log.TRACE)

	sig_ch := make(chan os.Signal, 1)
	signal.Notify(sig_ch, os.Interrupt)

	waitTimeout := 300 * time.Second
	deviceId := "go-test-device-id123"
	deviceKey := "go-test-device-key"

	url := "ws://playground.devicehive.com/api/websocket"
	key := "<access key should be here>"
	s, err := rest.NewService(url, key)
	if err != nil {
		log.Fatalf("Failed to create service (error: %s)", err)
	}
	log.Alwaysf("service created: %s", s)

	info, err := s.GetServerInfo(waitTimeout)
	if err != nil {
		log.Fatalf("Failed to get server info (error: %s)", err)
	}
	log.Alwaysf("server info: %s", info)

	devices, err := s.GetDeviceList(0, 0, waitTimeout)
	if err != nil {
		log.Fatalf("Failed to get device list (error: %s)", err)
	}
	for _, d := range devices {
		log.Alwaysf("device %s", d)
	}
	log.Alwaysf("done %d", len(devices))
	return

	device := devicehive.NewDevice(deviceId, "test-name",
		devicehive.NewDeviceClass("go-dev-class", "1.2.3"))
	device.DeviceClass.AddEquipment(devicehive.NewEquipment("n1", "c1", "t1"))
	device.DeviceClass.AddEquipment(devicehive.NewEquipment("n2", "c2", "t2"))
	//device.Network = devicehive.NewNetwork("dev-net", "net-key")
	device.Key = deviceKey

//	err = s.Authenticate(device, waitTimeout)
//	if err != nil {
//		log.Fatalf("Failed to authenticate (error: %s)", err)
//	}

	err = s.RegisterDevice(device, waitTimeout)
	if err != nil {
		log.Fatalf("Failed to register device (error: %s)", err)
	}

//	*device, err = s.GetDevice(deviceId, deviceKey, waitTimeout)
//	if err != nil {
//		log.Fatalf("Failed to get device (error: %s)", err)
//	}
	log.Alwaysf("device: %s", device)

	notification := devicehive.NewNotification("hello", 12345)
	err = s.InsertNotification(device, notification, waitTimeout)
	if err != nil {
		log.Fatalf("Failed to insert notification (error: %s)", err)
	}
	log.Alwaysf("notification: %s", notification)

	cmd_listener, err := s.SubscribeCommands(device, info.Timestamp, waitTimeout)
	if err != nil {
		log.Fatalf("failed to subscribe commands (error: %s)", err)
	}
	for {
		select {
			case command := <- cmd_listener.C:
				log.Alwaysf("command received: %s", command)
				command.Status = "Success"
				command.Result = "No"
				err = s.UpdateCommand(device, command, waitTimeout)
				if err != nil {
					log.Fatalf("failed to update command status (error: %s)", err)
				}
			case <- time.After(time.Second):
				log.Alwaysf(".")
			case <- sig_ch:
				log.Alwaysf("Exiting...")
				return
		}
	}
	return

	command := devicehive.NewCommand("hello", 12345)
	err = s.InsertCommand(device, command, waitTimeout)
	if err != nil {
		log.Fatalf("Failed to insert command (error: %s)", err)
	}
	log.Alwaysf("command: %s", command)

	command.Status = "Done"
	command.Result = "No result"
	err = s.UpdateCommand(device, command, waitTimeout)
	if err != nil {
		log.Fatalf("Failed to update command (error: %s)", err)
	}
}
