package ecs

import (
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type componentRegistry struct {
	currentId atomic.Uint64

	components                     map[reflect.Type]uint
	concurrencySafeComponentsMutex sync.RWMutex
	concurrencySafeComponents      map[reflect.Type]uint
}

func newComponentRegistry() componentRegistry {
	return componentRegistry{
		components:                map[reflect.Type]uint{},
		concurrencySafeComponents: map[reflect.Type]uint{},
	}
}

func (c *componentRegistry) getId(componentType reflect.Type) uint {
	if id, exists := c.components[componentType]; exists {
		return id
	}

	c.concurrencySafeComponentsMutex.RLock()
	id, exists := c.concurrencySafeComponents[componentType]
	c.concurrencySafeComponentsMutex.RUnlock()
	if exists {
		return id
	}

	c.currentId.Add(1)
	newId := c.currentId.Load()
	c.concurrencySafeComponentsMutex.Lock()
	c.concurrencySafeComponents[componentType] = uint(newId)
	c.concurrencySafeComponentsMutex.Unlock()
	return uint(newId)
}

// processComponentIdRegistries moves ids from the mutex protected map
// to the non-mutex protected map.
//
// ! This call is not concurrency safe !
func (c *componentRegistry) processComponentIdRegistries() {
	maps.Copy(c.components, c.concurrencySafeComponents)
}

type IComponent interface {
	RequiredComponents() []IComponent
}

type Component struct{}

func (Component) RequiredComponents() []IComponent {
	return []IComponent{}
}

type ComponentId struct {
	id            uint
	componentType reflect.Type
}

func (c *ComponentId) DebugString() string {
	result, _ := strings.CutPrefix(c.componentType.String(), "*")
	return result
}

func (c *ComponentId) Is(other *ComponentId) bool {
	return other.id == c.id
}

func (c *ComponentId) Id() uint {
	return c.id
}

func (c *ComponentId) ReflectType() reflect.Type {
	return c.componentType
}

// ComponentIdOf returns a unique representation of the component ID
func ComponentIdOf(component IComponent, world *World) ComponentId {
	componentType := reflect.TypeOf(component)

	if componentType.Kind() == reflect.Pointer {
		componentType = componentType.Elem()
	}

	return ComponentId{
		id:            world.componentRegistry.getId(componentType),
		componentType: componentType,
	}
}

// ComponentIdFor returns a unique representation of the component ID
func ComponentIdFor[T IComponent](world *World) ComponentId {
	componentType := reflect.TypeFor[T]()

	if componentType.Kind() == reflect.Pointer {
		componentType = componentType.Elem()
	}

	return ComponentId{
		id:            world.componentRegistry.getId(componentType),
		componentType: componentType,
	}
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

func toComponentIds(components []IComponent, world *World) []ComponentId {
	componentIds := make([]ComponentId, len(components))

	for i, component := range components {
		componentIds[i] = ComponentIdOf(component, world)
	}

	return componentIds
}

// getARequiredComponents non-exhaustively gets required components of `components` adds those components to `result`, and their types to `componentsToExclude`.
// Required components of which their type exists in `componentsToExclude` are skipped.
func getRequiredComponents(componentsToExclude *[]ComponentId, components []IComponent, result *[]IComponent, world *World) (newComponents []IComponent) {
	newComponents = []IComponent{}

	for _, component := range components {
		for _, required_component := range component.RequiredComponents() {
			if required_component == nil {
				continue
			}

			componentId := ComponentIdOf(required_component, world)

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
func getAllRequiredComponents(componentsToExclude *[]ComponentId, components []IComponent, world *World) []IComponent {
	result := []IComponent{}

	for {
		components = getRequiredComponents(componentsToExclude, components, &result, world)

		if len(components) == 0 {
			break
		}
	}

	return result
}
