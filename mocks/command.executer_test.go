package mocks

import (
	"context"
	"errors"
	"testing"

	"github.com/gapsquare/command"

	"github.com/stretchr/testify/assert"
)

var mockCommandStore = &MockCommandStore{}
var eventBus = &EventBus{}
var defaultConfiguration = command.BuildConfiguration(command.WithCommandStore(mockCommandStore))
var defaultRepository = &MockRepository{}
var commandHandler = &MockCommandHandler{}

func TestNewCommandExecuterHandleErrors(t *testing.T) {
	resetConfiguration()
	// Test Repository as nil
	_, err := command.NewExecuter(defaultConfiguration, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "repository is nil", err.Error())

	defaultConfiguration = command.BuildConfiguration(
		command.WithCommandStore(nil),
	)
	_, err = command.NewExecuter(defaultConfiguration, defaultRepository)
	assert.NotNil(t, err)
	assert.Equal(t, "store is nil", err.Error())

	defaultConfiguration = command.BuildConfiguration(
		command.WithCommandStore(mockCommandStore),
		command.WithEventBus(nil),
	)

	_, err = command.NewExecuter(defaultConfiguration, defaultRepository)
	assert.NotNil(t, err)
	assert.Equal(t, "event bus is nil", err.Error())
}

func TestCanNewCommandExecuter(t *testing.T) {
	resetConfiguration()
	ce, err := command.NewExecuter(defaultConfiguration, defaultRepository)
	assert.Nil(t, err)
	assert.NotNil(t, ce)
}

func TestCommandexecuterHandleErrors(t *testing.T) {
	resetConfiguration()
	ce, err := command.NewExecuter(defaultConfiguration, defaultRepository)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	assertThrowsError := func(cmd *MockSimpleCommand, expectedError string) {
		err := ce.Execute(context.Background(), cmd)
		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err.Error())
	}

	cmd := &MockSimpleCommand{ID: 1, Name: "mock", entity: &SimpleModel{}}

	// Test Find entity throws error
	command.RegisterCommandHandler(MockSimpleCommandType, commandHandler)
	defaultRepository.FindErr = errors.New("Model not found")
	assertThrowsError(cmd, "Model not found")

	// Test CommandHandler return error
	defaultRepository.FindErr = nil
	defaultRepository.Entity = &SimpleModel{ID: 1, Content: "some content"}
	commandHandler.Err = errors.New("Simulate CommandHandler error")
	assertThrowsError(cmd, "Simulate CommandHandler error")

	// Test CommandStore throws error
	commandHandler.Err = nil
	mockCommandStore.Err = errors.New("Save command failed")
	assertThrowsError(cmd, "Save command failed")
	mockCommandStore.Err = nil

	// Test Save Entity throws error
	defaultRepository.SaveErr = errors.New("Save entity failed")
	assertThrowsError(cmd, "Save entity failed")
	defaultRepository.SaveErr = nil
}

func TestVersionableModels(t *testing.T) {
	resetConfiguration()
	ce, err := command.NewExecuter(defaultConfiguration, defaultRepository)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	assert.NotNil(t, ce)
	command.RegisterCommandHandler(MockSimpleCommandType, commandHandler)
	command.RegisterCommandHandler(MockVersionableCommandType, commandHandler)
	defaultRepository.Entity = &MockVersionableModel{ID: 1, VersionInt: 1, Content: "some content"}

	err = ce.Execute(context.Background(), &MockSimpleCommand{ID: 1, Name: "mock", entity: &MockVersionableModel{}})
	assert.NotNil(t, err)
	assert.Equal(t, "version check fails. cmd.(Versionbale): false, entity.(Versionable): true", err.Error())

	defaultRepository.Entity = &SimpleModel{ID: 1, Content: "some content"}
	err = ce.Execute(context.Background(),
		&MockVersionableCommand{MockSimpleCommand: MockSimpleCommand{entity: &SimpleModel{}}})
	assert.NotNil(t, err)
	assert.Equal(t, "version check fails. cmd.(Versionbale): true, entity.(Versionable): false", err.Error())

	//test Versionable Version differs throws error

}

func TestVersionableEntityThrowsErrorOnDifVersion(t *testing.T) {
	resetConfiguration()
	ce, err := command.NewExecuter(defaultConfiguration, defaultRepository)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	assert.NotNil(t, ce)
	command.RegisterCommandHandler(MockVersionableCommandType, commandHandler)
	defaultRepository.Entity = &MockVersionableModel{ID: 1, VersionInt: 2, Content: "some content"}

	cmd := &MockVersionableCommand{MockSimpleCommand: MockSimpleCommand{ID: 1, Name: "some content", entity: &MockVersionableModel{}}, Ver: 1}
	err = ce.Execute(context.Background(), cmd)
	assert.NotNil(t, err)
	assert.Equal(t, command.ErrVersionMismatched, err)
}

func TestVersionableEntityIncrementVersion(t *testing.T) {
	resetConfiguration()
	ce, err := command.NewExecuter(defaultConfiguration, defaultRepository)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	assert.NotNil(t, ce)
	command.RegisterCommandHandler(MockVersionableCommandType, commandHandler)
	defaultRepository.Entity = &MockVersionableModel{ID: 1, VersionInt: 1, Content: "some content"}

	cmd := &MockVersionableCommand{
		MockSimpleCommand: MockSimpleCommand{ID: 1, Name: "some content", entity: &MockVersionableModel{}},
		Ver:               1}
	err = ce.Execute(context.Background(), cmd)
	assert.Nil(t, err)

	changedEntity, ok := defaultRepository.Entity.(*MockVersionableModel)
	assert.True(t, ok)
	assert.Equal(t, 2, changedEntity.VersionInt)

}

func resetConfiguration() {
	defaultConfiguration = command.BuildConfiguration(
		command.WithCommandStore(mockCommandStore),
		command.WithEventBus(eventBus),
	)
}

var _ = command.Command(&MockSimpleCommand{})
var MockSimpleCommandType = command.Type("mock.simple.command")
var MockVersionableCommandType = command.Type("mock.versionable.command")

type MockSimpleCommand struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	entity command.Entity
}

func (c *MockSimpleCommand) CommandType() command.Type {
	return MockSimpleCommandType
}

func (c *MockSimpleCommand) Entity() command.Entity {
	return c.entity
}

var _ = command.Command(&MockVersionableCommand{})
var _ = command.Versionable(&MockVersionableCommand{})

type MockVersionableCommand struct {
	MockSimpleCommand
	Ver int `json: "version"`
}

func (c *MockVersionableCommand) CommandType() command.Type {
	return MockVersionableCommandType
}

func (c *MockVersionableCommand) Version() command.VersionType {
	return command.VersionType(c.Ver)
}
