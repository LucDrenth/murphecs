package utils

import "slices"

// If there are no duplicates, returns (nil, -1, -1).
// If there are duplicates, return the duplicate element and the array indices of the duplicates
func GetFirstDuplicate[T comparable](typeIds []T) (*T, int, int) {
	for i := range len(typeIds) {
		for j := range len(typeIds) {
			if i == j {
				continue
			}

			if typeIds[i] == typeIds[j] {
				return &typeIds[i], i, j
			}
		}
	}

	return nil, -1, -1
}

// RemoveFromSlice removes s[index].
//
// This function does not check if index is out of bounds, and will panic if index >= len(s) or index < 0
//
// This function does not maintain the order of the slice. If you want to maintain the order, use
// RemoveFromSliceAndMaintainOrder instead. But be aware that RemoveFromSliceAndMaintainOrder is significantly slower.
func RemoveFromSlice[T any](s *[]T, index int) {
	(*s)[index] = (*s)[len(*s)-1]
	*s = (*s)[:len((*s))-1]
}

// RemoveFromSliceAndMaintainOrder remove s[index] while maintaining the order of the elements.
//
// This function does not check if index is out of bounds, and will panic if index >= len(s) or index < 0
//
// If you don't care about maintaining the order of the elements, use RemoveFromSlice instead.
// It is significantly faster.
func RemoveFromSliceAndMaintainOrder[T any](s *[]T, index int) {
	*s = slices.Delete((*s), index, index+1)
}
