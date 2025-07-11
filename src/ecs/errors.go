package ecs

import "errors"

var (
	ErrEntityNotFound error = errors.New("entity not found")

	ErrComponentNotFound       error = errors.New("component not found")
	ErrDuplicateComponent      error = errors.New("duplicate component")
	ErrComponentAlreadyPresent error = errors.New("component is already present")
	ErrComponentIsNotAPointer  error = errors.New("component is not a pointer")
	ErrComponentIsNil          error = errors.New("component is nil")

	ErrInvalidComponentStorageCapacity  error = errors.New("invalid component storage capacity")
	ErrComponentStorageIndexOutOfBounds error = errors.New("component storage index is out of bounds")

	ErrUnexpectedNumberOfQueryResults error = errors.New("unexpected number of query results")
)
