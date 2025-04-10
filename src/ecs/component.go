package ecs

import (
	"reflect"
	"slices"
)

type IComponent interface {
	requiredComponents() []IComponent
}

type Component struct{}

func (c Component) requiredComponents() []IComponent {
	return []IComponent{}
}

type componentType = string

func toComponentType(component IComponent) componentType {
	return reflect.TypeOf(component).String()
}

func toComponentTypes(components []IComponent) []componentType {
	componentTypes := make([]string, len(components))

	for i, component := range components {
		componentTypes[i] = toComponentType(component)
	}

	return componentTypes
}

// getARequiredComponents non-exhaustively gets required components of `components` adds those components to `result`, and their types to `typesToExclude`.
// Required components of which their type exists in `typesToExclude` are skipped.
func getRequiredComponents(typesToExclude *[]componentType, components []IComponent, result *[]IComponent) (newComponents []IComponent) {
	newComponents = []IComponent{}

	for _, component := range components {
		for _, required_component := range component.requiredComponents() {
			componentType := toComponentType(required_component)

			if slices.Contains(*typesToExclude, componentType) {
				continue
			}

			*typesToExclude = append(*typesToExclude, componentType)
			*result = append(*result, required_component)
			newComponents = append(newComponents, required_component)
		}
	}

	return newComponents
}

// getAllRequiredComponents exhaustively gets all required components of `components`.
//
// `typesToExclude` gets updated with the types from the result.
func getAllRequiredComponents(typesToExclude *[]componentType, components []IComponent) []IComponent {
	result := []IComponent{}

	for {
		components = getRequiredComponents(typesToExclude, components, &result)

		if len(components) == 0 {
			break
		}
	}

	return result
}
