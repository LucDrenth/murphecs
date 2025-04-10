package utils

// If there are no duplicates, returns (nil, -1, -1)
// If there are duplicates, return the duplicate element and the array indices in which the duplicates are
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
