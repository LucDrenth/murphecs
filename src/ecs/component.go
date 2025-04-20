package ecs

import (
	"reflect"
	"slices"
	"strings"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type IComponent interface {
	RequiredComponents() []IComponent
}

type Component struct{}

func (c Component) RequiredComponents() []IComponent {
	return []IComponent{}
}

type componentType = reflect.Type

// toComponentType returns a unique representation of the component type
func toComponentType(component IComponent) componentType {
	result := reflect.TypeOf(component)

	if result.Kind() == reflect.Pointer {
		return result.Elem()
	}

	return result
}

// getComponentType returns a unique representation of the component type
func getComponentType[T IComponent]() componentType {
	result := utils.TypeOf[T]()

	if result.Kind() == reflect.Pointer {
		return result.Elem()
	}

	return result
}

// toComponentDebugType returns a string reflection of the component type such as "ecs.Entity"
func toComponentDebugType(component IComponent) string {
	return reflect.TypeOf(component).String()
}

// getComponentDebugType returns a string reflection of the component type such as "ecs.Entity"
func getComponentDebugType[T IComponent]() string {
	result := utils.TypeOf[T]().String()
	result, _ = strings.CutPrefix(result, "*")
	return result
}

func toComponentTypes(components []IComponent) []componentType {
	componentTypes := make([]componentType, len(components))

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
		for _, required_component := range component.RequiredComponents() {
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
