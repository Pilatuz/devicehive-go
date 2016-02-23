package devicehive

import (
	"fmt"
)

// Represents DeviceHive server information
type ServerInfo struct {
	// Server's API version.
	Version string `json:"apiVersion"`

	// Current server time, UTC.
	Timestamp string `json:"serverTimestamp"`

	// Alternative websocket URL. Empty for Websocket service.
	WebsocketUrl string `json:"webSocketServerUrl,omitempty"`

	// Alternative REST URL. Empty for REST service.
	RestUrl string `json:"restServerUrl,omitempty"`
}

// Get ServerInfo string representation
func (info ServerInfo) String() string {
	switch {
	// both URLs provided, should be impossible
	case len(info.WebsocketUrl) != 0 && len(info.RestUrl) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, WS:%q, REST:%q}",
			info.Version, info.Timestamp, info.WebsocketUrl, info.RestUrl)

	// Websocket URL
	case len(info.WebsocketUrl) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, WS:%q}",
			info.Version, info.Timestamp, info.WebsocketUrl)

	// REST URL
	case len(info.RestUrl) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, REST:%q}",
			info.Version, info.Timestamp, info.RestUrl)
	}

	// default, no alternative URLs
	return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q}",
		info.Version, info.Timestamp)
}

// Assign fields from map.
// This method is used to assign already parsed JSON data.
func (info *ServerInfo) FromMap(data interface{}) error {
	return fromJsonMap(info, data)
}
