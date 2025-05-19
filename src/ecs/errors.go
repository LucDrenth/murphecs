package ecs

import "errors"

var (
	ErrEntityNotFound error = errors.New("entity not found")

	ErrComponentNotFound       error = errors.New("component not found")
	ErrDuplicateComponent      error = errors.New("duplicate component")
	ErrComponentAlreadyPresent error = errors.New("component is already present")
	ErrComponentIsNotAPointer  error = errors.New("component is not a pointer")

	ErrInvalidComponentStorageCapacity   error = errors.New("invalid component storage capacity")
	ErrComponentRegistryIndexOutOfBounds error = errors.New("component registry index is out of bounds")
)
