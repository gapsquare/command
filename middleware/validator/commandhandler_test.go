package validator

import (
	"context"
	"errors"
	"testing"

	"github.com/gapsquare/command"

	"github.com/gapsquare/command/mocks"
)

func TestCommandHandler_WithValidationError(t *testing.T) {
	inner := &mocks.MockCommandHandler{}
	m := NewMiddleware()
	h := command.UseHandlerMiddleware(inner, m)
	cmd := &mocks.Command{
		ID:      1,
		Content: "content",
	}
	e := errors.New("a validation error")
	c := CommandWithValidation(cmd, func() error { return e })
	if err := h.HandleCommand(context.Background(), c); err != e {
		t.Error("there should be an error:", e)
	}
}

func TestCommandHandler_WithValidationNoError(t *testing.T) {
	inner := &mocks.MockCommandHandler{}
	m := NewMiddleware()
	h := command.UseHandlerMiddleware(inner, m)
	cmd := &mocks.Command{
		ID:      1,
		Content: "content",
	}
	c := CommandWithValidation(cmd, func() error { return nil })
	if err := h.HandleCommand(context.Background(), c); err != nil {
		t.Error("there should be no error:", err)
	}

}
