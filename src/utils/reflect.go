package utils

import (
	"fmt"
	"reflect"
)

// AlignedSize returns the aligned size of t.
//
// Some systems already return the aligned size when calling reflect.Type.Size, but some don't.
// Use this function to ensure getting the aligned size. This is necessary when manually managing
// blocks of memory.
func AlignedSize(t reflect.Type) uintptr {
	size := t.Size()
	align := uintptr(t.Align())
	return (size + (align - 1)) / align * align
}

// MethodHasPointerReceiver checks wether the implementation of a method uses a pointer receiver.
func MethodHasPointerReceiver(i any, methodName string) (bool, error) {
	if i == nil {
		return false, fmt.Errorf("interface value is nil")
	}

	if methodName == "" {
		return false, fmt.Errorf("methodName can not be empty")
	}

	if !StringStartsWithUpper(methodName) {
		return false, fmt.Errorf("method must be exported")
	}

	concreteType := reflect.TypeOf(i)
	if concreteType.Kind() == reflect.Pointer {
		concreteType = concreteType.Elem()
	}

	if concreteType.Kind() != reflect.Struct {
		return false, fmt.Errorf("input has invalid type")
	}

	method, found := concreteType.MethodByName(methodName)
	if !found {
		pointerType := reflect.PointerTo(concreteType)
		method, found = pointerType.MethodByName(methodName)
		if !found {
			return false, fmt.Errorf("method '%s' not found on type %v or its pointer equivalent", methodName, concreteType)
		}
	}

	receiverType := method.Type.In(0)
	return receiverType.Kind() == reflect.Pointer, nil
}
