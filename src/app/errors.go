package app

import "errors"

var (
	ErrResourceAlreadyPresent error = errors.New("resource already present")
	ErrResourceNotAPointer    error = errors.New("resource is not a pointer")
	ErrResourceNotFound       error = errors.New("resource not found")
	ErrResourceTypeNotValid   error = errors.New("resource type not valid")

	ErrSystemNotAFunction          error = errors.New("not a function")
	ErrSystemInvalidReturnType     error = errors.New("invalid return type(s)")
	ErrSystemParamQueryNotAPointer error = errors.New("query must be a pointer")
	ErrSystemParamQueryNotValid    error = errors.New("query not valid")
	ErrSystemParamNotValid         error = errors.New("not valid")
)
