package command

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gapsquare/goevent"
)

// Executer an interface that is responsable to execute a command
type Executer interface {
	// Execute execute the command
	Execute(context.Context, Command) error
}

type commandExecuter struct {
	repository ReadWriteRepository
	store      Store
	bus        goevent.EventBus
}

// NewExecuter creates an instance of Executer
func NewExecuter(config configureOption, repository ReadWriteRepository) (Executer, error) {
	if repository == nil {
		return nil, errors.New("repository is nil")
	}

	if config.commandStore == nil {
		return nil, errors.New("store is nil")
	}

	if config.eventBus == nil {
		return nil, errors.New("event bus is nil")
	}
	return &commandExecuter{
		repository: repository,
		store:      config.commandStore,
		bus:        config.eventBus,
	}, nil
}

func (ce *commandExecuter) Execute(ctx context.Context, cmd Command) error {

	handler, err := GetCommandHandler(cmd.CommandType())
	if err != nil {
		return fmt.Errorf("Can not find command handler for command %s, Error: %v", cmd.CommandType(), err)
	}

	if dest, ok := handler.(*DestructiveHandler); ok {
		if e := ce.executeDestructiveHandler(ctx, cmd, dest); e != nil {
			return e
		}
	} else {
		if e := ce.executeConstructiveHandler(ctx, cmd, handler); e != nil {
			return e
		}
	}

	return ce.publishEvents(ctx, cmd)
}

func (ce *commandExecuter) executeDestructiveHandler(ctx context.Context, cmd Command, handler Handler) error {
	entity := cmd.Entity()

	if entity == nil {
		return fmt.Errorf("Can not run destructive action on a nil Entity for command type %s", cmd.CommandType())
	}

	err := ce.repository.Find(entity)
	if err != nil {
		return err
	}

	if e := ce.store.Save(cmd, ce.repository); e != nil {
		return e
	}

	if e := handler.HandleCommand(ctx, cmd); e != nil {
		return e
	}

	return nil
}

func (ce *commandExecuter) executeConstructiveHandler(ctx context.Context, cmd Command, handler Handler) error {

	entity := cmd.Entity()
	if entity != nil {

		// load aggregate and send it to commandhandler
		err := ce.repository.Find(entity)
		if err != nil {
			return err
		}

		// if cmd is versionable check version, entity also should be versionable
		cmdVersionable, cmdOk := cmd.(Versionable)
		entityVersionable, eOk := entity.(EntityVersionable)

		if cmdOk != eOk {
			return fmt.Errorf("version check fails. cmd.(Versionbale): %v, entity.(Versionable): %v", cmdOk, eOk)
		}
		if eOk && cmdOk {
			if !reflect.DeepEqual(cmdVersionable.Version(), entityVersionable.Version()) {
				return ErrVersionMismatched
			}
		}
	}

	if e := handler.HandleCommand(ctx, cmd); e != nil {
		return e
	}

	if e := ce.store.Save(cmd, ce.repository); e != nil {
		return e
	}

	if entity != nil {
		if entityVersionable, ok := entity.(EntityVersionable); ok {
			entityVersionable.IncrementVersion()

		}

		if e := ce.repository.Save(entity); e != nil {
			return e
		}
	}

	return nil
}

func (ce *commandExecuter) publishEvents(ctx context.Context, cmd Command) error {
	if c, ok := cmd.(WithEvents); ok {
		events := c.Events(ctx)
		for _, ev := range events {
			if err := ce.bus.Publish(ctx, ev); err != nil {
				return err
			}
		}
	}

	return nil
}
