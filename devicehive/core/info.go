package core

import "fmt"

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

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (info *ServerInfo) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("ServerInfo: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("ServerInfo: %v - unexpected data type")
	}

	// version
	if av, ok := data["apiVersion"]; ok {
		switch v := av.(type) {
		case string:
			info.Version = v
		default:
			return fmt.Errorf("ServerInfo: %v - unexpected value for version", av)
		}
	}

	// timestamp
	if ts, ok := data["serverTimestamp"]; ok {
		switch v := ts.(type) {
		case string:
			info.Timestamp = v
		default:
			return fmt.Errorf("ServerInfo: %v - unexpected value for timestamp", ts)
		}
	}

	// websocket URL
	if ws, ok := data["webSocketServerUrl"]; ok {
		switch v := ws.(type) {
		case string:
			info.WebsocketUrl = v
		case nil:
			// do nothing
		default:
			return fmt.Errorf("ServerInfo: %v - unexpected value for websocketUrl", ws)
		}
	}

	// rest URL
	if rest, ok := data["restServerUrl"]; ok {
		switch v := rest.(type) {
		case string:
			info.RestUrl = v
		case nil:
			// do nothing
		default:
			return fmt.Errorf("ServerInfo: %v - unexpected value for restUrl", rest)
		}
	}

	return nil // OK
}
