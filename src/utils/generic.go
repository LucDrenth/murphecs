package utils

import (
	"errors"
	"reflect"
)

// ToConcrete converts T into the default concrete value of that generic.
// Returns an error when passing an interface.
func ToConcrete[T any]() (T, error) {
	var v T

	if reflect.TypeFor[T]().Kind() == reflect.Interface {
		return v, errors.New("interface can not be converted in to a concrete type")
	}

	return v, nil
}
