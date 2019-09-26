package command

import "errors"

// ErrVersionMismatched error for version mismatched
var ErrVersionMismatched = errors.New("Version mismatched")

// Versionable is an item that has a version number
type Versionable interface {
	Version() VersionType
}

// VersionType version type
type VersionType uint64
