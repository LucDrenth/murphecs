package app

import "errors"

var (
	ErrResourceAlreadyPresent error = errors.New("resource already present")
	ErrResourceNotAPointer    error = errors.New("resource is not a pointer")
	ErrResourceNotFound       error = errors.New("resource not found")
	ErrResourceTypeNotValid   error = errors.New("resource type not valid")
)
