package core

import (
	"encoding/json"
	"testing"
)

// Test ServiceInfo JSON marshaling
func TestJsonServiceInfo(t *testing.T) {
	check := func(info ServerInfo, expectedJson string) {
		jsonBytes, err := json.Marshal(info)
		if err != nil {
			t.Errorf("ServerInfo: Cannot convert %+v to JSON (error: %s)", info, err)
		}
		if expectedJson != string(jsonBytes) {
			t.Errorf("ServerInfo: Cannot convert %+v to JSON\n\t   found: %s,\n\texpected:%s",
				info, string(jsonBytes), expectedJson)
		}

		t.Logf("ServerInfo: %+v converted to %s", info, string(jsonBytes))
	}

	info := ServerInfo{
		Version:      "1.2.3",
		Timestamp:    "2015-10-22T14:15:16.999",
		WebsocketUrl: "ws://devicehive.com"}

	check(info, `{"apiVersion":"1.2.3","serverTimestamp":"2015-10-22T14:15:16.999","webSocketServerUrl":"ws://devicehive.com"}`)

	info.WebsocketUrl = ""
	info.RestUrl = "https://devicehive.com"
	check(info, `{"apiVersion":"1.2.3","serverTimestamp":"2015-10-22T14:15:16.999","restServerUrl":"https://devicehive.com"}`)
}
