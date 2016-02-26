package devicehive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test JSON to string method.
func TestServerInfoString(t *testing.T) {
	info := new(ServerInfo)
	info.Version = "1.2.3"
	info.Timestamp = "2005-10-22"
	assert.Equal(t, info.String(), `ServerInfo{Version:"1.2.3", Timestamp:"2005-10-22"}`)

	info.RestURL, info.WebsocketURL = "http://localhost/rest", ""
	assert.Equal(t, info.String(), `ServerInfo{Version:"1.2.3", Timestamp:"2005-10-22", REST:"http://localhost/rest"}`)

	info.RestURL, info.WebsocketURL = "", "http://localhost/ws"
	assert.Equal(t, info.String(), `ServerInfo{Version:"1.2.3", Timestamp:"2005-10-22", WS:"http://localhost/ws"}`)
}

// Test ServerInfo JSON marshaling.
func TestServerInfoJson(t *testing.T) {
	info := new(ServerInfo)
	info.Version = "1.2.3"
	info.Timestamp = "2005-10-22"
	assert.JSONEq(t, toJsonStr(t, info), `{"apiVersion":"1.2.3", "serverTimestamp":"2005-10-22"}`)

	info.RestURL, info.WebsocketURL = "http://localhost/rest", ""
	assert.JSONEq(t, toJsonStr(t, info), `{"apiVersion":"1.2.3", "serverTimestamp":"2005-10-22", "restServerUrl":"http://localhost/rest"}`)

	info.RestURL, info.WebsocketURL = "", "http://localhost/ws"
	assert.JSONEq(t, toJsonStr(t, info), `{"apiVersion":"1.2.3", "serverTimestamp":"2005-10-22", "webSocketServerUrl":"http://localhost/ws"}`)
}

// Test ServerInfo assignment.
func TestServerInfoAssign(t *testing.T) {
	info := new(ServerInfo)
	assert.NoError(t, info.FromMap(map[string]interface{}{
		"apiVersion":      123,
		"serverTimestamp": "2005-10-22",
	}))
	assert.JSONEq(t, toJsonStr(t, info), `{"apiVersion":"123", "serverTimestamp":"2005-10-22"}`)

	assert.NoError(t, info.FromMap(map[string]interface{}{
		"restServerUrl":      "http://localhost/rest",
		"webSocketServerUrl": "http://localhost/ws",
	}))
	assert.JSONEq(t, toJsonStr(t, info), `{"apiVersion":"123", "serverTimestamp":"2005-10-22", "restServerUrl":"http://localhost/rest", "webSocketServerUrl":"http://localhost/ws"}`)
}
