package ecs

type WorldConfigs struct {
	// Decides the capacity of a component storage when a component is spawned or inserted for the first time.
	InitialComponentCapacityStrategy initialComponentCapacityStrategy

	// Decides by how much the component storage capacity grows when inserting a component if the component storage is full
	ComponentCapacityGrowthStrategy componentCapacityGrowthStrategy
}

func DefaultWorldConfigs() WorldConfigs {
	return WorldConfigs{
		InitialComponentCapacityStrategy: &StaticDefaultComponentCapacity{Capacity: 1024},
		ComponentCapacityGrowthStrategy:  &ComponentCapacityGrowthDouble{},
	}
}

type initialComponentCapacityStrategy interface {
	GetDefaultComponentCapacity(ComponentId) uint
}

type StaticDefaultComponentCapacity struct {
	Capacity uint
}

func (s *StaticDefaultComponentCapacity) GetDefaultComponentCapacity(ComponentId) uint {
	return s.Capacity
}

// ComponentSpecificDefaultComponentCapacity lets you specify a default component capacity for
// specific components, or use a default if its not specified.
//
// It is useful to reduce memory usage for components that you do not expect to make a lot of.
// It is also useful to prevent increasing capacity for component when you expect a lot of them to be made.
type ComponentSpecificDefaultComponentCapacity struct {
	ComponentCapacities map[ComponentId]uint
	Default             uint
}

func (s *ComponentSpecificDefaultComponentCapacity) GetDefaultComponentCapacity(component ComponentId) uint {
	if capacity, ok := s.ComponentCapacities[component]; ok {
		return capacity
	}

	return s.Default
}

type componentCapacityGrowthStrategy interface {
	GetExtraCapacity(currentCapacity uint) uint
}

type ComponentCapacityGrowthDouble struct{}

func (s *ComponentCapacityGrowthDouble) GetExtraCapacity(currentCapacity uint) uint {
	return currentCapacity
}
