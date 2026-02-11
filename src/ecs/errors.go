package ecs

import "errors"

var (
	ErrEntityNotFound error = errors.New("entity not found")

	ErrComponentNotFound       error = errors.New("component not found")
	ErrComponentDuplicate      error = errors.New("duplicate component")
	ErrComponentAlreadyPresent error = errors.New("component is already present")
	ErrComponentIsNil          error = errors.New("component is nil")

	ErrResourceAlreadyPresent error = errors.New("resource already present")
	ErrResourceIsNil          error = errors.New("resource is nil")
	ErrResourceNotFound       error = errors.New("resource not found")
	ErrResourceTypeNotValid   error = errors.New("resource type not valid")
	ErrResourceTypeNotAllowed error = errors.New("resource type not allowed")

	ErrInvalidComponentStorageCapacity  error = errors.New("invalid component storage capacity")
	ErrComponentStorageIndexOutOfBounds error = errors.New("component storage index is out of bounds")

	ErrUnexpectedNumberOfQueryResults error = errors.New("unexpected number of query results")

	ErrTargetWorldNotFound error = errors.New("target world not found")
)
