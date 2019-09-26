package command

// Store is an interface for storing commands
type Store interface {

	// Save command
	Save(Command, WriteRepository) error
}
