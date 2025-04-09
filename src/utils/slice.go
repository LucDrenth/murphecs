package utils

func GetFirstDuplicate[T comparable](typeIds []T) *T {
	for i := range len(typeIds) {
		for j := range len(typeIds) {
			if i == j {
				continue
			}

			if typeIds[i] == typeIds[j] {
				return &typeIds[i]
			}
		}
	}

	return nil
}
