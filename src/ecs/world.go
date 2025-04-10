package ecs

import (
	"fmt"
	"reflect"
)

type entityId = uint

type entry struct {
	components []any
}

type world struct {
	entityIdCounter uint
	entities        map[entityId]entry
}

func NewWorld() world {
	return world{
		entities: map[entityId]entry{},
	}
}

func (world *world) Spawn(components ...any) (entityId, error) {
	typeIds := make([]string, len(components))
	for i, component := range components {
		typeIds[i] = reflect.TypeOf(component).String()
	}

	// check for duplicates
	duplicate := getFirstDuplicate(typeIds)
	if duplicate != nil {
		return 0, fmt.Errorf("found duplicate component: %s", *duplicate)
	}

	world.entityIdCounter++
	entityId := world.entityIdCounter
	world.entities[entityId] = entry{
		components: components,
	}

	return entityId, nil
}

func getFirstDuplicate(typeIds []string) *string {
	for i := range len(typeIds) {
		for j := range len(typeIds) {
			if i == j {
				continue
			}

			if typeIds[i] == typeIds[j] {
				return &typeIds[i]
			}
		}
	}

	return nil
}
