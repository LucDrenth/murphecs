package app

import "errors"

var (
	ErrResourceAlreadyPresent error = errors.New("resource already present")
	ErrResourceNotAPointer    error = errors.New("resource is not a pointer")
	ErrResourceIsNil          error = errors.New("resource is nil")
	ErrResourceNotAStruct     error = errors.New("resource is not a struct")
	ErrResourceNotFound       error = errors.New("resource not found")
	ErrResourceTypeNotValid   error = errors.New("resource type not valid")
	ErrResourceTypeNotAllowed error = errors.New("resource type not allowed")

	ErrSystemNotAFunction          error = errors.New("not a function")
	ErrSystemInvalidReturnType     error = errors.New("invalid return type(s)")
	ErrSystemParamQueryNotAPointer error = errors.New("query must be a pointer")
	ErrSystemParamQueryNotValid    error = errors.New("query param not valid")
	ErrSystemParamWorldNotAPointer error = errors.New("world must be a pointer")
	ErrSystemParamNotValid         error = errors.New("not valid")

	ErrTargetWorldNotKnown error = errors.New("target world not known")
)
