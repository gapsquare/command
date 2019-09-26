package command

import (
	"errors"
	"fmt"
	"sync"
)

// ErrCommandHandlerNotRegistered command not registered
var ErrCommandHandlerNotRegistered = errors.New("Command handler not registered")
var errNilCommandHandler = errors.New("Created command is nil")
var errEmptyCommandType = errors.New("Can not register a command handler for empty type")
var errRegisterDuplicateCommand = func(cmdType Type) error {
	return fmt.Errorf("Attempt to register duplicate command handler for type %s", cmdType)
}

var errUnregisterNotRegisteredCommandHandler = func(cmdType Type) error {
	return fmt.Errorf("Can not un-register not registered command handler for type %s", cmdType)
}

var handlerMu sync.RWMutex
var handlers = make(map[Type]Handler)

// RegisterCommandHandler register a command Handler
func RegisterCommandHandler(cmdType Type, h Handler) error {

	if cmdType == Type("") {
		return errEmptyCommandType
	}

	handlerMu.Lock()
	defer handlerMu.Unlock()

	if _, ok := handlers[cmdType]; ok {
		return errRegisterDuplicateCommand(cmdType)
	}
	handlers[cmdType] = h
	return nil
}

// UnRegisterCommandHandler un register command Handler
func UnRegisterCommandHandler(cmdType Type) error {
	if cmdType == Type("") {
		return errEmptyCommandType
	}

	handlerMu.Lock()
	defer handlerMu.Unlock()

	if _, ok := handlers[cmdType]; !ok {
		return errUnregisterNotRegisteredCommandHandler(cmdType)
	}
	delete(handlers, cmdType)
	return nil
}

// GetCommandHandler returns a registered command Handler
func GetCommandHandler(cmdType Type) (Handler, error) {
	handlerMu.Lock()
	defer handlerMu.Unlock()

	if handler, ok := handlers[cmdType]; ok {
		return handler, nil
	}
	return nil, ErrCommandHandlerNotRegistered
}
