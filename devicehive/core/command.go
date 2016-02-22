package core

import "fmt"

// Represents command object - a set of data sent from DeviceHive to devices.
type Command struct {
	// Unique identifier [do not change].
	Id uint64 `json:"id,omitempty"`

	// Timestamp, UTC.
	Timestamp string `json:"timestamp,omitempty"`

	// Accociated user identifier.
	UserId uint64 `json:"userId,omitempty"`

	// Command name.
	Name string `json:"command,omitempty"`

	// Lifetime, number of seconds until this command expires.
	Lifetime uint64 `json:"lifetime,omitempty"`

	// JSON object with an arbitrary structure [optional].
	Parameters interface{} `json:"parameters,omitempty"`

	// Status reported by device.
	Status string `json:"status,omitempty"`

	// Execution result reported by device.
	Result interface{} `json:"result,omitempty"`
}

// Command listener is used to listen for asynchronous commands.
type CommandListener struct {
	// channel to receive commands
	C chan *Command
}

// NewCommand creates a new command.
// No lifetime set by default.
func NewCommand(name string, parameters interface{}) *Command {
	return &Command{Name: name, Parameters: parameters}
}

// NewCommandResult creates a new command result.
func NewCommandResult(id uint64, status string, result interface{}) *Command {
	return &Command{Id: id, Status: status, Result: result}
}

// NewCommandListener creates a new command listener.
func NewCommandListener() *CommandListener {
	ch := make(chan *Command) // TODO: buffered?
	return &CommandListener{C: ch}
}

// Get Command string representation
func (command Command) String() string {
	body := ""

	// Id [optional]
	if command.Id != 0 {
		body += fmt.Sprintf("Id:%d, ", command.Id)
	}

	// Name
	body += fmt.Sprintf("Name:%q", command.Name)

	// Timestamp
	if len(command.Timestamp) != 0 {
		body += fmt.Sprintf(", Timestamp:%q", command.Timestamp)
	}

	// UserId
	if command.UserId != 0 {
		body += fmt.Sprintf(", UserId:%d", command.UserId)
	}

	// Lifetime
	if command.Lifetime != 0 {
		body += fmt.Sprintf(", Lifetime:%d", command.Lifetime)
	}

	// Parameters [optional]
	if command.Parameters != nil {
		body += fmt.Sprintf(", Parameters:%v", command.Parameters)
	}

	// Status
	if len(command.Status) != 0 {
		body += fmt.Sprintf(", Status:%q", command.Status)
	}

	// Result [optional]
	if command.Result != nil {
		body += fmt.Sprintf(", Result:%v", command.Result)
	}

	return fmt.Sprintf("Command{%s}", body)
}

// Assign parsed JSON.
// This method is used to assign already parsed JSON data.
func (command *Command) AssignJSON(rawData interface{}) error {
	if rawData == nil {
		return fmt.Errorf("Command: no data")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Command: %v - unexpected data type", rawData)
	}

	// identifier
	if id, ok := data["id"]; ok {
		switch v := id.(type) {
		case float64:
			command.Id = uint64(v)
		case uint64:
			command.Id = v
		default:
			return fmt.Errorf("Command: %v - unexpected value for id", id)
		}
	}

	// timestamp
	if ts, ok := data["timestamp"]; ok {
		switch v := ts.(type) {
		case string:
			command.Timestamp = v
		default:
			return fmt.Errorf("Command: %v - unexpected value for timestamp", ts)
		}
	}

	// user identifier
	if id, ok := data["userId"]; ok {
		switch v := id.(type) {
		case float64:
			command.UserId = uint64(v)
		case uint64:
			command.UserId = v
		default:
			return fmt.Errorf("Command: %v - unexpected value for userId", id)
		}
	}

	// name
	if name, ok := data["command"]; ok {
		switch v := name.(type) {
		case string:
			command.Name = v
		default:
			return fmt.Errorf("Command: %v - unexpected value for name", name)
		}
	}

	// lifetime
	if lt, ok := data["lifetime"]; ok {
		switch v := lt.(type) {
		case float64:
			command.Lifetime = uint64(v)
		case uint64:
			command.Lifetime = v
		case nil:
			command.Lifetime = 0
		default:
			return fmt.Errorf("Command: %v - unexpected value for lifetime", lt)
		}
	}

	// parameters (as is)
	if p, ok := data["parameters"]; ok {
		command.Parameters = p
	}

	// status
	if s, ok := data["status"]; ok {
		switch v := s.(type) {
		case string:
			command.Status = v
		case nil:
			command.Status = ""
		default:
			return fmt.Errorf("Command: %v - unexpected value for status", s)
		}
	}

	// result (as is)
	if r, ok := data["result"]; ok {
		command.Result = r
	}

	return nil // OK
}
