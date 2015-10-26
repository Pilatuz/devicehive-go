package core

import (
	"encoding/json"
	"testing"
)

// Test Network JSON marshaling
func TestJsonNetwork(t *testing.T) {
	check := func(network Network, expectedJson string) {
		jsonBytes, err := json.Marshal(network)
		if err != nil {
			t.Errorf("Network: Cannot convert %+v to JSON (error: %s)", network, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("Network: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				network, string(jsonBytes), expectedJson)
		}

		t.Logf("Network: %+v converted to %s", network, string(jsonBytes))
	}

	network := Network{
		Name:        "net-name",
		Key:         "net-key",
		Description: "custom description"}

	check(network, `{"name":"net-name","key":"net-key","description":"custom description"}`)

	network.Description = ""
	network.Id = 100
	check(network, `{"id":100,"name":"net-name","key":"net-key"}`)
}
