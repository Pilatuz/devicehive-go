package devicehive

import (
	"bytes"
	"fmt"
)

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
	// Channel to receive commands.
	C chan *Command
}

// NewEmptyCommand creates a new empty command.
func NewEmptyCommand() *Command {
	command := new(Command)
	return command
}

// NewCommand creates a new command.
// No lifetime set by default.
func NewCommand(name string, parameters interface{}) *Command {
	command := new(Command)
	command.Name = name
	command.Parameters = parameters
	return command
}

// NewCommandResult creates a new command result.
func NewCommandResult(id uint64, status string, result interface{}) *Command {
	command := new(Command)
	command.Id = id
	command.Status = status
	command.Result = result
	return command
}

// NewCommandListener creates a new command listener.
func NewCommandListener(buffered int) *CommandListener {
	listener := new(CommandListener)
	listener.C = make(chan *Command, buffered)
	return listener
}

// Get Command string representation.
func (command Command) String() string {
	// NOTE all errors are ignored!
	body := new(bytes.Buffer)

	// Id [optional]
	if command.Id != 0 {
		body.WriteString(fmt.Sprintf("Id:%d, ", command.Id))
	}

	// Name
	body.WriteString(fmt.Sprintf("Name:%q", command.Name))

	// Timestamp
	if len(command.Timestamp) != 0 {
		body.WriteString(fmt.Sprintf(", Timestamp:%q", command.Timestamp))
	}

	// UserId
	if command.UserId != 0 {
		body.WriteString(fmt.Sprintf(", UserId:%d", command.UserId))
	}

	// Lifetime
	if command.Lifetime != 0 {
		body.WriteString(fmt.Sprintf(", Lifetime:%d", command.Lifetime))
	}

	// Parameters [optional]
	if command.Parameters != nil {
		body.WriteString(fmt.Sprintf(", Parameters:%v", command.Parameters))
	}

	// Status
	if len(command.Status) != 0 {
		body.WriteString(fmt.Sprintf(", Status:%q", command.Status))
	}

	// Result [optional]
	if command.Result != nil {
		body.WriteString(fmt.Sprintf(", Result:%v", command.Result))
	}

	return fmt.Sprintf("Command{%s}", body)
}

// Assign fields from map.
// This method is used to assign already parsed JSON data.
func (command *Command) FromMap(data interface{}) error {
	return fromJsonMap(command, data)
}
