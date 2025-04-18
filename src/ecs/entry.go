package ecs

type entry struct {
	components []IComponent
}

func (entry *entry) containsComponentType(componentTypeToCheck componentType) bool {
	for i := range entry.components {
		if toComponentType(entry.components[i]) == componentTypeToCheck {
			return true
		}
	}

	return false
}

func (entry *entry) countComponents() int {
	return len(entry.components)
}

// getComponentFromEntry returns a pointer to the component, the index of the component and nil if entry contains the component.
// returns nil, -1, error if entry does not contain the component.
func getComponentFromEntry[T IComponent](entry *entry) (*T, int, error) {
	for i := range entry.components {
		if result, ok := entry.components[i].(T); ok {
			return &result, i, nil
		}
	}

	return nil, -1, ErrComponentNotFound
}
