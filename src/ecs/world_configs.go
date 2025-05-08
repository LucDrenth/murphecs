package ecs

type WorldConfigs struct {
	// Decides the capacity of a component storage when a component is spawned or inserted for the first time.
	ComponentCapacityStrategy initialComponentCapacityStrategy
}

func DefaultWorldConfigs() WorldConfigs {
	return WorldConfigs{
		ComponentCapacityStrategy: &StaticDefaultComponentCapacity{Capacity: 1024},
	}
}

type initialComponentCapacityStrategy interface {
	GetDefaultComponentCapacity(ComponentType) uint
}

type StaticDefaultComponentCapacity struct {
	Capacity uint
}

func (s *StaticDefaultComponentCapacity) GetDefaultComponentCapacity(ComponentType) uint {
	return s.Capacity
}

// ComponentSpecificDefaultComponentCapacity lets you specify a default component capacity for
// specific components, or use a default if its not specified.
//
// It is useful to reduce memory usage for components that you do not expect to make a lot of.
// It is also useful to prevent increasing capacity for component when you expect a lot of them to be made.
type ComponentSpecificDefaultComponentCapacity struct {
	ComponentCapacities map[ComponentType]uint
	Default             uint
}

func (s *ComponentSpecificDefaultComponentCapacity) GetDefaultComponentCapacity(component ComponentType) uint {
	if capacity, ok := s.ComponentCapacities[component]; ok {
		return capacity
	}

	return s.Default
}
