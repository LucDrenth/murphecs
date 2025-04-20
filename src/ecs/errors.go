package ecs

import "errors"

var ErrEntityNotFound error = errors.New("entity not found")
var ErrComponentNotFound error = errors.New("component not found")
var ErrDuplicateComponent error = errors.New("duplicate component")
var ErrComponentAlreadyPresent error = errors.New("component is already present")
var ErrComponentIsNotAPointer error = errors.New("component is not a pointer")
