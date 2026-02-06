package utils

// PointerTo returns a pointer to value. Useful for pointing to a literal value.
//
// For example: `stringPointer := PointerTo("my string")`, which can not be done
// like this: `stringPointer := &"my string"`
func PointerTo[T any](value T) *T {
	return &value
}
