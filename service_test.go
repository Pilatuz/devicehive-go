// +build ignore

package devicehive

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"
)

var (
	testRestServerUrl      = "http://playground.devicehive.com/api/rest"
	testWsServerUrl        = "ws://playground.devicehive.com/api/websocket"
	testAccessKey          = ""
	testDeviceId           = "go-test-dev-id"
	testDeviceName         = "go-test-dev-name"
	testDeviceKey          = "go-test-dev-key"
	testDeviceClassName    = "go-device-class"
	testDeviceClassVersion = "1.2.3"
	testNetworkName        = ""
	testNetworkKey         = ""
	testNetworkDesc        = ""
	testWaitTimeout        = 60 * time.Second
	testLogLevel           = "NOLOG"

	testGapMs    = 1000
	testBatchLen = 100
)

// initialize test environment
func init() {
	flag.StringVar(&testRestServerUrl, "rest-url", testRestServerUrl, "REST service URL")
	flag.StringVar(&testWsServerUrl, "ws-url", testWsServerUrl, "Websocket service URL")
	flag.StringVar(&testAccessKey, "access-key", testAccessKey, "key to access playground")

	flag.StringVar(&testDeviceId, "device-id", testDeviceId, "test Device identifier")
	flag.StringVar(&testDeviceName, "device-name", testDeviceName, "test Device name")
	flag.StringVar(&testDeviceKey, "device-key", testDeviceKey, "test Device key")

	flag.StringVar(&testDeviceClassName, "device-class-name", testDeviceClassName, "test Device class name")
	flag.StringVar(&testDeviceClassVersion, "device-class-version", testDeviceClassVersion, "test Device class version")

	flag.StringVar(&testNetworkName, "network-name", testNetworkName, "test Network name")
	flag.StringVar(&testNetworkKey, "network-key", testNetworkKey, "test Network key")
	flag.StringVar(&testNetworkDesc, "network-desc", testNetworkDesc, "test Network description")

	flag.IntVar(&testGapMs, "gap", testGapMs, "gap interval, milliseconds")
	flag.IntVar(&testBatchLen, "batch-len", testBatchLen, "batch length")

	flag.StringVar(&testLogLevel, "log-level", testLogLevel, "Logging level: WARN INFO DEBUG TRACE or NOLOG")
	flag.Parse()

	log.SetLevelByName(testLogLevel)
}

// creates new REST service
func testNewRest(t *testing.T) (service *rest.Service) {
	if len(testRestServerUrl) == 0 {
		return
	}

	service, err := rest.NewService(testRestServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create REST service (error: %s)", err)
	}
	return
}

// creates new Websocket service
func testNewWs(t *testing.T) (service *ws.Service) {
	if len(testWsServerUrl) == 0 {
		return
	}

	service, err := ws.NewService(testWsServerUrl, testAccessKey)
	if err != nil {
		t.Errorf("Failed to create WS service (error: %s)", err)
	}
	return
}

// creates new REST service (abstract interface)
func testNewRestService(t *testing.T) (service Service) {
	if rs := testNewRest(t); rs != nil {
		service = rs
	}
	return
}

// creates new Websocket service (abstract interface)
func testNewWsService(t *testing.T) (service Service) {
	if wss := testNewWs(t); wss != nil {
		service = wss
	}
	return
}

// creates new test Device with device class initialized
func testNewDevice() (device *core.Device) {
	dc := core.NewDeviceClass(testDeviceClassName, testDeviceClassVersion)
	device = core.NewDevice(testDeviceId, testDeviceName, dc)
	device.Key = testDeviceKey
	return
}

// creates new test Network
func testNewNetwork() (network *core.Network) {
	if len(testNetworkName) != 0 {
		network = core.NewNetwork(testNetworkName, testNetworkKey)
		network.Description = testNetworkDesc
	}
	return
}

type TimeStat struct {
	min   time.Duration
	max   time.Duration
	sum   time.Duration
	count int
}

// Add interval to the statistics
func (s *TimeStat) Add(t time.Duration) {
	if s.count == 0 { // first time
		s.min = t
		s.max = t
	} else {
		if t < s.min {
			s.min = t
		}
		if t > s.max {
			s.max = t
		}
	}

	s.sum += t
	s.count += 1
}

// get statistics as string
func (s TimeStat) String() string {
	if s.count > 0 {
		return fmt.Sprintf("{min:%s, max:%s, mean:%s}",
			s.min, s.max, s.sum/time.Duration(s.count))
	} else {
		return fmt.Sprintf("{no stats}")
	}
}

func TestBatchCommandInsert(t *testing.T) {
	device := testNewDevice()
	device.Network = testNewNetwork()

	s := testNewRest(t)
	s2 := testNewWs(t)

	device.Id += "-batch-cmd"
	testCheckRegisterDevice1(t, s, *device, true)

	count := testBatchLen
	gap := time.Duration(testGapMs) * time.Millisecond

	tx_cmds := make([]*core.Command, 0, count)
	rx_cmds := make([]*core.Command, 0, count)

	type Stat struct {
		tx_beg time.Time
		tx_end time.Time
		rx_end time.Time
	}
	stat := make(map[string]*Stat)

	// transmitter
	go func() {
		time.Sleep(2 * time.Second) // small delay before start
		log.Infof("TEST/TX: started")
		for i := 0; i < count; i++ {
			p := fmt.Sprintf("%d", i)
			cmd := core.NewCommand("batch-command", p)
			stat[p] = &Stat{}
			stat[p].tx_beg = time.Now()
			err := s.InsertCommand(device, cmd, testWaitTimeout)
			stat[p].tx_end = time.Now()
			if err != nil {
				t.Errorf("failed to insert batch command: %s", err)
				break
			}
			log.Infof("TEST/TX: %s", cmd)
			tx_cmds = append(tx_cmds, cmd)
			time.Sleep(gap)
		}
		log.Infof("TEST/TX: stopped")
	}()

	// receiver
	listener, err := s2.SubscribeCommands(device, "", testWaitTimeout)
	if err != nil {
		t.Errorf("failed to subscribe commands: %s", err)
		return
	}

	log.Infof("TEST/RX: started")
	for len(rx_cmds) < count && !t.Failed() {
		select {
		case cmd := <-listener.C:
			p := cmd.Parameters.(string)
			stat[p].rx_end = time.Now()
			log.Infof("TEST/RX: %s", cmd)
			rx_cmds = append(rx_cmds, cmd)
		case <-time.After(30 * time.Second):
			t.Errorf("failed to wait command: %s", "timed out")
			break
		}
	}
	log.Infof("TEST/RX: stopped")

	err = s2.UnsubscribeCommands(device, testWaitTimeout)
	if err != nil {
		t.Errorf("failed to unsubscribe commands: %s", err)
		return
	}

	// compare tx_cmd == rx_cmd
	if len(tx_cmds) != count || len(rx_cmds) != count {
		t.Errorf("TX:%d != RX:%d commands length mismatch", len(tx_cmds), len(rx_cmds))
		return
	}

	// time statistics:
	// ins - insertion time (tx_end - tx_beg)
	// rtt - round trip (rx_end - tx_beg)
	var ins, rtt TimeStat
	for i, tx := range tx_cmds {
		rx := rx_cmds[i]
		//		log.Infof("%d:\tTX:%q at %q\t\tRX:%q at %q", i,
		//				tx.Parameters, tx.Timestamp,
		//				rx.Parameters, rx.Timestamp)
		if tx.Name != rx.Name {
			t.Errorf("TX:%q != RX:%q command name mismatch", tx.Name, rx.Name)
			continue
		}
		tx_p := tx.Parameters.(string)
		rx_p := rx.Parameters.(string)
		if tx_p != rx_p {
			t.Errorf("TX:%q != RX:%q command parameter mismatch", tx_p, rx_p)
			continue
		}

		t := stat[tx_p]
		ins.Add(t.tx_end.Sub(t.tx_beg))
		rtt.Add(t.rx_end.Sub(t.tx_beg))
	}

	log.Infof("insert time: %s", ins)
	log.Infof(" round trip: %s", rtt)
}

// Test InsertNotification method
func TestInsertNotification(t *testing.T) {
	// create device (REST)
	device := testNewDevice()
	device.Network = testNewNetwork()
	testCheckRegisterDevice1(t, testNewRestService(t), *device, true)
	if t.Failed() {
		return // nothing to test without device
	}

	notification := core.NewNotification("ntf-test", 12345)
	testCheckInsertNotification(t, testNewRestService(t), device, *notification)
	testCheckInsertNotification(t, testNewWsService(t), device, *notification)
}

func TestBatchNotificationInsert(t *testing.T) {
	device := testNewDevice()
	device.Network = testNewNetwork()

	s := testNewRest(t)
	s2 := s //testNewWs(t)

	device.Id += "-batch-ntf"
	testCheckRegisterDevice1(t, s, *device, true)

	count := testBatchLen
	gap := time.Duration(testGapMs) * time.Millisecond

	tx_ntfs := make([]*core.Notification, 0, count)
	rx_ntfs := make([]*core.Notification, 0, count)

	type Stat struct {
		tx_beg time.Time
		tx_end time.Time
		rx_end time.Time
	}
	stat := make(map[string]*Stat)

	// transmitter
	go func() {
		time.Sleep(2 * time.Second) // small delay before start
		log.Infof("TEST/TX: started")
		for i := 0; i < count; i++ {
			p := fmt.Sprintf("%d", i)
			ntf := core.NewNotification("batch-notification", p)
			stat[p] = &Stat{}
			stat[p].tx_beg = time.Now()
			err := s.InsertNotification(device, ntf, testWaitTimeout)
			stat[p].tx_end = time.Now()
			if err != nil {
				t.Errorf("failed to insert batch notification: %s", err)
				break
			}
			log.Infof("TEST/TX: %s", ntf)
			tx_ntfs = append(tx_ntfs, ntf)
			time.Sleep(gap)
		}
		log.Infof("TEST/TX: stopped")
	}()

	// receiver
	listener, err := s2.SubscribeNotifications(device, "", testWaitTimeout)
	if err != nil {
		t.Errorf("failed to subscribe notifications: %s", err)
		return
	}

	log.Infof("TEST/RX: started")
	for len(rx_ntfs) < count && !t.Failed() {
		select {
		case ntf := <-listener.C:
			if ntf.Name != "batch-notification" {
				// notification ignored
				continue
			}
			p := ntf.Parameters.(string)
			stat[p].rx_end = time.Now()
			log.Infof("TEST/RX: %s", ntf)
			rx_ntfs = append(rx_ntfs, ntf)
		case <-time.After(30 * time.Second):
			t.Errorf("failed to wait notification: %s", "timed out")
			break
		}
	}
	log.Infof("TEST/RX: stopped")

	err = s2.UnsubscribeNotifications(device, testWaitTimeout)
	if err != nil {
		t.Errorf("failed to unsubscribe notifications: %s", err)
		return
	}

	// compare tx_ntfs == rx_ntfs
	if len(tx_ntfs) != count || len(rx_ntfs) != count {
		t.Errorf("TX:%d != RX:%d notifications length mismatch", len(tx_ntfs), len(rx_ntfs))
		return
	}

	// time statistics:
	// ins - insertion time (tx_end - tx_beg)
	// rtt - round trip (rx_end - tx_beg)
	var ins, rtt TimeStat
	for i, tx := range tx_ntfs {
		rx := rx_ntfs[i]
		//		log.Infof("%d:\tTX:%q at %q\t\tRX:%q at %q", i,
		//				tx.Parameters, tx.Timestamp,
		//				rx.Parameters, rx.Timestamp)
		//		if tx.Name != rx.Name {
		//			t.Errorf("TX:%q != RX:%q notification name mismatch", tx.Name, rx.Name)
		//			continue
		//		}
		tx_p := tx.Parameters.(string)
		rx_p := rx.Parameters.(string)
		if tx_p != rx_p {
			t.Errorf("TX:%q != RX:%q notification parameter mismatch", tx_p, rx_p)
			continue
		}

		t := stat[tx_p]
		ins.Add(t.tx_end.Sub(t.tx_beg))
		rtt.Add(t.rx_end.Sub(t.tx_beg))
	}

	log.Infof("insert time: %s", ins)
	log.Infof(" round trip: %s", rtt)
}
