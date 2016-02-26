package devicehive

import (
	"fmt"
)

// ServerInfo represents DeviceHive server information
type ServerInfo struct {
	// Server's API version.
	Version string `json:"apiVersion"`

	// Current server time, UTC.
	Timestamp string `json:"serverTimestamp"`

	// Alternative websocket URL. Empty for Websocket service.
	WebsocketURL string `json:"webSocketServerUrl,omitempty"`

	// Alternative REST URL. Empty for REST service.
	RestURL string `json:"restServerUrl,omitempty"`
}

// Get ServerInfo string representation
func (info ServerInfo) String() string {
	switch {
	// both URLs provided, should be impossible
	case len(info.WebsocketURL) != 0 && len(info.RestURL) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, WS:%q, REST:%q}",
			info.Version, info.Timestamp, info.WebsocketURL, info.RestURL)

	// Websocket URL
	case len(info.WebsocketURL) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, WS:%q}",
			info.Version, info.Timestamp, info.WebsocketURL)

	// REST URL
	case len(info.RestURL) != 0:
		return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q, REST:%q}",
			info.Version, info.Timestamp, info.RestURL)
	}

	// default, no alternative URLs
	return fmt.Sprintf("ServerInfo{Version:%q, Timestamp:%q}",
		info.Version, info.Timestamp)
}

// FromMap assigns fields from map.
// This method is used to assign already parsed JSON data.
func (info *ServerInfo) FromMap(data interface{}) error {
	return fromJSON(info, data)
}
