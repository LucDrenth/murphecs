package ecs

import (
	"reflect"
	"slices"
	"strings"
)

type IComponent interface {
	RequiredComponents() []IComponent
}

type Component struct{}

func (Component) RequiredComponents() []IComponent {
	return []IComponent{}
}

type ComponentId = reflect.Type

// ComponentIdOf returns a unique representation of the component ID
func ComponentIdOf(component IComponent) ComponentId {
	result := reflect.TypeOf(component)

	if result.Kind() == reflect.Pointer {
		return result.Elem()
	}

	return result
}

// ComponentIdFor returns a unique representation of the component ID
func ComponentIdFor[T IComponent]() ComponentId {
	result := reflect.TypeFor[T]()

	if result.Kind() == reflect.Pointer {
		return result.Elem()
	}

	return result
}

// ComponentDebugStringOf returns a string reflection of a component id such as "ecs.Entity"
func ComponentDebugStringOf(component IComponent) string {
	return reflect.TypeOf(component).String()
}

// ComponentDebugStringFor returns a string reflection of a component id such as "ecs.Entity"
func ComponentDebugStringFor[T IComponent]() string {
	result := reflect.TypeFor[T]().String()
	result, _ = strings.CutPrefix(result, "*")
	return result
}

func toComponentIds(components []IComponent) []ComponentId {
	componentIds := make([]ComponentId, len(components))

	for i, component := range components {
		componentIds[i] = ComponentIdOf(component)
	}

	return componentIds
}

// getARequiredComponents non-exhaustively gets required components of `components` adds those components to `result`, and their types to `componentsToExclude`.
// Required components of which their type exists in `componentsToExclude` are skipped.
func getRequiredComponents(componentsToExclude *[]ComponentId, components []IComponent, result *[]IComponent) (newComponents []IComponent) {
	newComponents = []IComponent{}

	for _, component := range components {
		for _, required_component := range component.RequiredComponents() {
			componentId := ComponentIdOf(required_component)

			if slices.Contains(*componentsToExclude, componentId) {
				continue
			}

			*componentsToExclude = append(*componentsToExclude, componentId)
			*result = append(*result, required_component)
			newComponents = append(newComponents, required_component)
		}
	}

	return newComponents
}

// getAllRequiredComponents exhaustively gets all required components of `components`.
//
// `componentsToExclude` gets updated with the types from the result.
func getAllRequiredComponents(componentsToExclude *[]ComponentId, components []IComponent) []IComponent {
	result := []IComponent{}

	for {
		components = getRequiredComponents(componentsToExclude, components, &result)

		if len(components) == 0 {
			break
		}
	}

	return result
}
