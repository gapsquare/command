package command

// EntityID uniqueidentifier type for aggregates
type EntityID int

// Entity is an item which is identified by ID
type Entity interface {
	EntityID() EntityID
}

//EntityVersionable is an Entity with require Version check
type EntityVersionable interface {
	Entity
	Versionable
	IncrementVersion()
}
