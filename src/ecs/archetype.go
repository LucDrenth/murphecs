package ecs

import (
	"fmt"
	"sort"
	"strconv"
)

type archetypeStorage struct {
	componentsHashToArchetype map[string]Archetype
	entityIdToArchetype       map[EntityId]*Archetype
	componentIdToArchetypes   map[ComponentId]*[]*Archetype
	idCounter                 uint
}

func newArchetypeStorage() archetypeStorage {
	return archetypeStorage{
		componentsHashToArchetype: map[string]Archetype{},
		entityIdToArchetype:       map[EntityId]*Archetype{},
		componentIdToArchetypes:   map[ComponentId]*[]*Archetype{},
	}
}

// getArchetype either returns an existing archetype or creates a new one if it doesn't exist yet.
func (s archetypeStorage) getArchetype(world *World, componentIds []ComponentId) (*Archetype, error) {
	hash := hashComponentIds(componentIds)
	archetype, exists := s.componentsHashToArchetype[hash]
	if exists {
		return &archetype, nil
	}

	archetype, err := newArchetype(world, componentIds)
	if err != nil {
		return nil, err
	}

	s.componentsHashToArchetype[hash] = archetype

	for i := range componentIds {
		archetypeList, exists := s.componentIdToArchetypes[componentIds[i]]
		if exists {
			*archetypeList = append(*archetypeList, &archetype)
		} else {
			s.componentIdToArchetypes[componentIds[i]] = &[]*Archetype{&archetype}
		}
	}

	return &archetype, nil
}

// countComponents returns the number of living components
func (storage *archetypeStorage) countComponents() uint {
	count := uint(0)
	for key := range storage.componentsHashToArchetype {
		archetype := storage.componentsHashToArchetype[key]
		count += archetype.CountComponents()
	}

	return count
}

type Archetype struct {
	id                 uint
	componentTypesHash string
	components         map[ComponentId]*componentStorage
	componentIds       []ComponentId
}

// newArchetype returns a new archetype for the given componentIds.
// The componentIds will get sorted, so the order of the given componentIds does not matter.
func newArchetype(world *World, componentIds []ComponentId) (Archetype, error) {
	sortComponentIds(componentIds)

	components := map[ComponentId]*componentStorage{}
	for i := range componentIds {
		storage, err := createComponentStorage(
			world.initialComponentCapacityStrategy.GetDefaultComponentCapacity(componentIds[i]),
			componentIds[i],
		)
		if err != nil {
			return Archetype{}, fmt.Errorf("failed to create component storage for component %s: %w", componentIds[i].DebugString(), err)
		}

		components[componentIds[i]] = &storage
	}

	world.archetypeStorage.idCounter++

	return Archetype{
		id:                 world.archetypeStorage.idCounter,
		componentTypesHash: hashComponentIds(componentIds),
		components:         components,
		componentIds:       componentIds,
	}, nil
}

func (archetype *Archetype) IsFromComponents(componentIds []ComponentId) bool {
	sortComponentIds(componentIds)
	return archetype.componentTypesHash == hashComponentIds(componentIds)
}

func (archetype *Archetype) HasComponent(componentId ComponentId) bool {
	_, hasComponent := archetype.components[componentId]
	return hasComponent
}

func (archetype *Archetype) CountComponents() uint {
	count := uint(0)
	for _, storage := range archetype.components {
		count += storage.numberOfComponents
	}

	return count
}

func sortComponentIds(componentIds []ComponentId) {
	sort.Slice(componentIds, func(i, j int) bool {
		return componentIds[i].id > componentIds[j].id
	})
}

// hashComponentIds returns a unique hash for every different combination of component id's.
// This function is deterministic, meaning the same input results in the same output.
func hashComponentIds(componentIds []ComponentId) string {
	if len(componentIds) == 0 {
		return ""
	}

	// Max uint64 is 20 digits
	// Max uint32 is 10 digits
	//
	// Which one is used is platform dependent, but we'll assume 64 bits just to be safe.
	const maxDigitsPerComponentId = 20

	// For the ',' character. We need it because else, cases like (1, 3) and (13) would give
	// the same result.
	const delimiterSize = 1

	buf := make([]byte, 0, len(componentIds)*(maxDigitsPerComponentId+delimiterSize))

	buf = strconv.AppendUint(buf, uint64(componentIds[0].id), 10)
	for i := 1; i < len(componentIds); i++ {
		buf = append(buf, ',')
		buf = strconv.AppendUint(buf, uint64(componentIds[i].id), 10)
	}

	return string(buf)
}
