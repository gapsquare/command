package command

import "context"

// Handler interface that all command handler should implement
type Handler interface {

	// HandleCommand handle the command
	HandleCommand(context.Context, Command) error
}

// HandlerMiddleware is a function that middlewares can implement to be able to chain
type HandlerMiddleware func(Handler) Handler

// UseHandlerMiddleware wraps a Command in one or more middlewares.
func UseHandlerMiddleware(h Handler, middleware ...HandlerMiddleware) Handler {
	// Apply in reverse order
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}

	return h
}

// DestructiveHandler destructive command handler
type DestructiveHandler struct {
	handler  Handler
	onDelete func(ctx context.Context, cmd Command) error
}

// HandleCommand implements command.CommandHandler interface
func (h *DestructiveHandler) HandleCommand(ctx context.Context, cmd Command) error {
	if err := h.handler.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return h.onDelete(ctx, cmd)
}

//NewDestructiveHandler creates new NewDestructiveCommandHandler
func NewDestructiveHandler(handler Handler, fn func(ctx context.Context, cmd Command) error) *DestructiveHandler {
	return &DestructiveHandler{handler: handler, onDelete: fn}
}

//HandlerFunc a function that can handle commands
type HandlerFunc func(context.Context, Command) error

// HandleCommand HandlerFunc implementation of the CommandHandler
func (h HandlerFunc) HandleCommand(ctx context.Context, cmd Command) error {
	return h(ctx, cmd)
}
