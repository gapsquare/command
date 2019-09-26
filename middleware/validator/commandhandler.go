package validator

import (
	"context"

	"github.com/gapsquare/command"
)

// NewMiddleware returns a new async handling middleware that validate commands
// with its own validation method.
func NewMiddleware() command.HandlerMiddleware {
	return command.HandlerMiddleware(func(h command.Handler) command.Handler {
		return command.HandlerFunc(func(ctx context.Context, cmd command.Command) error {
			// Call the validation method if it exists
			if c, ok := cmd.(Command); ok {
				err := c.Validate()
				if err != nil {
					return err
				}
			}

			// Immediate command execution.
			return h.HandleCommand(ctx, cmd)
		})
	})
}

// Command is a command with its own validation method
type Command interface {
	command.Command
	// Validate returns the error when validating the command
	Validate() error
}

// CommandWithValidation returns a wrapped command with a validation method
func CommandWithValidation(cmd command.Command, v func() error) Command {
	return &commandImp{Command: cmd, validate: v}
}

type commandImp struct {
	command.Command
	validate func() error
}

func (c *commandImp) Validate() error {
	return c.validate()
}
