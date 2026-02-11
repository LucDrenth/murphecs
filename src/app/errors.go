package app

import "errors"

var (
	ErrSystemTypeNotValid          error = errors.New("system type is not valid")
	ErrSystemNotAFunction          error = errors.New("not a function")
	ErrSystemInvalidReturnType     error = errors.New("invalid return type(s)")
	ErrSystemParamQueryNotAPointer error = errors.New("query must be a pointer")
	ErrSystemParamQueryNotValid    error = errors.New("query param not valid")
	ErrSystemParamWorldNotAPointer error = errors.New("world must be a pointer")
	ErrSystemParamNotValid         error = errors.New("not valid")

	ErrSystemParamEventReaderNotAPointer error = errors.New("must be a pointer")
	ErrSystemParamEventWriterNotAPointer error = errors.New("must be a pointer")
	ErrSystemParamOuterResourceIsAPointer error = errors.New("OuterResource must not be a pointer")

	ErrScheduleAlreadyExists error = errors.New("schedule already exists")
	ErrScheduleNotFound      error = errors.New("schedule not found")
	ErrScheduleTypeNotFound  error = errors.New("schedule type not found")
)
