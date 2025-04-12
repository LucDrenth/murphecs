package ecs

import "errors"

var ErrEntityNotFound error = errors.New("entity not found")
var ErrComponentNotFound error = errors.New("component not found")
