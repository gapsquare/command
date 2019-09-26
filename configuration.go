package command

import (
	"github.com/isuruceanu/goevent"
)

// Configuration configuration option func
type Configuration func(*configureOption)

// BuildConfiguration build command with optional options
// Default: //TODO: specify defaults values
func BuildConfiguration(options ...Configuration) configureOption {
	cfg := configureOption{
		//eventStore:
	}

	for _, option := range options {
		option(&cfg)
	}

	return cfg
}

type configureOption struct {
	eventStore   goevent.EventStore
	commandStore Store
	eventBus     goevent.EventBus
}

// WithEventStore sets specific EventStore
func WithEventStore(evs goevent.EventStore) Configuration {
	return func(c *configureOption) {
		c.eventStore = evs
	}
}

// WithCommandStore sets specific command Store
func WithCommandStore(cmdStore Store) Configuration {
	return func(c *configureOption) {
		c.commandStore = cmdStore
	}
}

// WithEventBus sets specific EventBus
func WithEventBus(eBus goevent.EventBus) Configuration {
	return func(c *configureOption) {
		c.eventBus = eBus
	}
}
