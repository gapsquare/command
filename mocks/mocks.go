package mocks

import (
	"context"
	"reflect"
	"time"

	"github.com/gapsquare/command"

	"github.com/gapsquare/goevent"
)

func init() {
	goevent.RegisterEventData(Topic, func() goevent.EventData { return &EventData{} })
}

const (
	// Topic main mock event
	Topic goevent.EventTopic = "Event"

	// TopicOther - other event
	TopicOther goevent.EventTopic = "OtherEvent"

	// CommandType is the type for Command.
	CommandType command.Type = "Command"
	// CommandOtherType is the type for CommandOther.
	CommandOtherType command.Type = "CommandOther"
)

// EventData is a mocked event data, useful in testing :)
type EventData struct {
	Content string
}

// MockVersionableModel is a mocked read model, useful in testing
type MockVersionableModel struct {
	ID         int       `json:"id" 	bson:"_id"`
	VersionInt int       `json:"version" bson:"version"`
	Content    string    `json:"content" bson:"content"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

var _ = command.Entity(&MockVersionableModel{})
var _ = command.Versionable(&MockVersionableModel{})
var _ = command.EntityVersionable(&MockVersionableModel{})

// EntityID implements the EntityID method of the command.Entity
func (m *MockVersionableModel) EntityID() command.EntityID { return command.EntityID(m.ID) }

// Version implements Version method of the command.Versionable interface
func (m *MockVersionableModel) Version() command.VersionType { return command.VersionType(m.VersionInt) }

func (m *MockVersionableModel) IncrementVersion() {
	m.VersionInt++
}

// SimpleModel is a mocked read model for a simple model without version
type SimpleModel struct {
	ID      int    `json:"id" 	bson:"_id"`
	Content string `json:"content" bson:"content"`
}

var _ = command.Entity(&SimpleModel{})

// EntityID implements the EntityID method of the command.Entity
func (m *SimpleModel) EntityID() command.EntityID { return command.EntityID(m.ID) }

// EventHandler is a mocked command.EventHandler, useful in testing.
type EventHandler struct {
	Type   string
	Events goevent.Events
	Time   time.Time
	Recv   chan goevent.Event
	// Used to simulate errors when publishing.
	Err error
}

var _ = goevent.EventHandler(&EventHandler{})

// NewEventHandler creates a new EventHandler.
func NewEventHandler(handlerType string) *EventHandler {
	return &EventHandler{
		Type:   handlerType,
		Events: goevent.Events{},
		Recv:   make(chan goevent.Event, 10),
	}
}

func (m *EventHandler) HandlerType() goevent.EventHandlerType {
	return goevent.EventHandlerType(m.Type)
}

func (m *EventHandler) HandleEvent(ctx context.Context, event goevent.Event) error {
	if m.Err != nil {
		return m.Err
	}

	m.Events = append(m.Events, event)
	m.Time = time.Now()
	m.Recv <- event
	return nil
}

// Reset resets the mock data
func (m *EventHandler) Reset() {
	m.Events = goevent.Events{}
	m.Time = time.Time{}
}

// Wait is a helper to wait some duration until for an event to be handled
func (m *EventHandler) Wait(d time.Duration) bool {
	select {
	case <-m.Recv:
		return true
	case <-time.After(d):
		return false
	}
}

var _ = goevent.EventBus(&EventBus{})

// EventBus is a mocked command.EventBus, useful in testing
type EventBus struct {
	Events goevent.Events
	Err    error
}

// Publish implements PublishEvent method of the command.EventBus interface
func (b *EventBus) Publish(ctx context.Context, event goevent.Event) error {
	if b.Err != nil {
		return b.Err
	}

	b.Events = append(b.Events, event)
	return nil
}

// AddHandler implements AddHandler method of the command.EventBus interface
func (b *EventBus) AddHandler(m goevent.EventMatcher, h goevent.EventHandler) error {
	return b.Err
}

// Errors implements Errors method of the command.EventBus interface
func (b *EventBus) Errors() <-chan goevent.EventBusError {
	return make(chan goevent.EventBusError)
}

//MockRepository is a mock repository for command.ReadWriteRepository, useful in tests
type MockRepository struct {
	Entity                               command.Entity
	LoadErr, SaveErr, FindErr            error
	FindCalled, SaveCalled, RemoveCalled bool
}

var _ = command.ReadWriteRepository(&MockRepository{})

// Find implements the Find method of command.ReadRepository
func (r *MockRepository) Find(entity command.Entity) error {
	r.FindCalled = true
	if r.FindErr != nil {
		return r.FindErr
	}

	if r.Entity == nil {
		entity = nil
		return nil
	}

	refPtrVal := reflect.ValueOf(entity)
	refPtrVal.Elem().Set(reflect.ValueOf(r.Entity).Elem())

	return nil
}

// Remove implements the Remove method of command.WriteRepository
func (r *MockRepository) Remove(entity command.Entity) error {
	r.RemoveCalled = true
	if r.SaveErr != nil {
		return r.SaveErr
	}

	r.Entity = nil
	return nil
}

// Save implements the Save method of command.WriteRepository
func (r *MockRepository) Save(entity command.Entity) error {
	r.SaveCalled = true
	if r.SaveErr != nil {
		return r.SaveErr
	}
	r.Entity = entity
	return nil
}

var _ = command.Command(Command{})

//Command is a mocked Command usefull in testing
type Command struct {
	ID      int
	Content string
}

func (t Command) CommandType() command.Type { return CommandType }
func (t Command) Entity() command.Entity    { return &SimpleModel{} }
func (t Command) Handler() command.Handler  { return &MockCommandHandler{} }

var _ = command.Handler(&MockCommandHandler{})

type MockCommandHandler struct {
	// use to simulate error
	Err error
	// Func of business logics
	BuFn func(command.Command, command.Entity) error
}

func (h *MockCommandHandler) HandleCommand(ctx context.Context, c command.Command) error {
	if h.Err != nil {
		return h.Err
	}

	// do business logics
	if h.BuFn == nil {
		return nil
	}

	return h.BuFn(c, c.Entity())
}

type MockCommandStore struct {
	SaveCall bool
	// Use to simulate errors
	Err error
}

func (cs MockCommandStore) Save(command.Command, command.WriteRepository) error {
	if cs.Err != nil {
		return cs.Err
	}

	cs.SaveCall = true
	return nil

}

type contextKey int

const (
	contextKeyOne contextKey = iota
)

func WithContextOne(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, contextKeyOne, val)
}

func ContextOne(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(contextKeyOne).(string)
	return val, ok
}
