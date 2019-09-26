package command

// WriteRepository is a write repository for entities
type WriteRepository interface {
	Save(Entity) error
	Remove(Entity) error
}

// ReadRepository is a read repository for entities
type ReadRepository interface {
	Find(Entity) error
}

// ReadWriteRepository is a combined read and write repository
type ReadWriteRepository interface {
	ReadRepository
	WriteRepository
}
