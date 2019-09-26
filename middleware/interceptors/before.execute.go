package interceptors

import (
	"context"

	"github.com/gapsquare/command"
)

type Command interface {
	command.Command
	BeforeExecute() error
}

func NewBeforeMiddleware() command.HandlerMiddleware {
	return command.HandlerMiddleware(func(h command.Handler) command.Handler {
		return command.HandlerFunc(func(ctx context.Context, cmd command.Command) error {
			// Call the validation method if it exists
			if c, ok := cmd.(Command); ok {
				err := c.BeforeExecute()
				if err != nil {
					return err
				}
			}

			// Immediate command execution.
			return h.HandleCommand(ctx, cmd)
		})
	})
}

// CommandInterceptBefore returns a wrapped command with before interceptor method
func CommandInterceptBefore(cmd command.Command, b func() error) Command {
	return &commandImp{Command: cmd, before: b}
}

type commandImp struct {
	command.Command
	before func() error
}

//BeforeExecute implements the Command interface
func (c *commandImp) BeforeExecute() error {
	return c.before()
}
