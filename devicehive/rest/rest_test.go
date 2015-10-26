package rest

import (
	"github.com/devicehive/devicehive-go/devicehive/log"
	"time"
)

var (
	testServerUrl   = "http://playground.devicehive.com/api/rest"
	testAccessKey   = "<place key here>"
	testDeviceId    = "go-test-device-id"
	testDeviceKey   = "go-test-device-key"
	testWaitTimeout = 60 * time.Second
)

func init() {
	// TODO: load test parameters from command line arguments
	// or from environment variables

	// TODO: check the test parameters

	log.SetLevel(log.TRACE)
}
