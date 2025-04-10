package ecs

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murphy/src/utils"
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
	duplicate := utils.GetFirstDuplicate(typeIds)
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
