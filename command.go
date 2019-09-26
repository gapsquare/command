package command

import (
	"context"

	"github.com/gapsquare/goevent"
)

// Command base interface for a command
type Command interface {
	CommandType() Type
	Entity() Entity
}

// Type command type
type Type string

//Encoder interface for command encode and decode
type Encoder interface {
	Marshal(*Command) []byte
	Unmarshal([]byte, Command) error
}

// WithEvents an interface for command which publish events
type WithEvents interface {
	Command
	// BuildEvents build a list of events need to be published after Command is handled
	Events(context.Context) goevent.Events
}
