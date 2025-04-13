package ecs

type entry struct {
	components []IComponent
}

func (entry *entry) containsComponentType(componentTypeToCheck componentType) bool {
	for _, component := range entry.components {
		if toComponentType(component) == componentTypeToCheck {
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
	for i, component := range entry.components {
		if maybeTarget, ok := component.(T); ok {
			return &maybeTarget, i, nil
		}
	}

	return nil, -1, ErrComponentNotFound
}
