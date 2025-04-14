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

func (c Component) RequiredComponents() []IComponent {
	return []IComponent{}
}

type componentType = string

func toComponentType(component IComponent) componentType {
	// TODO reflect.Type.String is not safe because it does not guarantee uniqueness. Instead, it recommends
	// to compare the types directly. We probably don't want use reflect.Type as componentType because it is
	// way bigger than we need. Instead, we could make a map[reflect.Type]componentType and look it up in there.
	return reflect.TypeOf(component).String()
}

func getComponentType[T IComponent]() componentType {
	// TODO reflect.Type.String is not safe because it does not guarantee uniqueness. Instead, it recommends
	// to compare the types directly. We probably don't want use reflect.Type as componentType because it is
	// way bigger than we need. Instead, we could make a map[reflect.Type]componentType and look it up in there.
	return reflect.TypeOf((*T)(nil)).String()
}

// toComponentDebugType returns a string reflection of the component type such as "ecs.Entity"
func toComponentDebugType(component IComponent) string {
	return reflect.TypeOf(component).String()
}

// getComponentDebugType returns a string reflection of the component type such as "ecs.Entity"
func getComponentDebugType[T IComponent]() string {
	result := reflect.TypeOf((*T)(nil)).String()
	result, _ = strings.CutPrefix(result, "*")
	return result
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
