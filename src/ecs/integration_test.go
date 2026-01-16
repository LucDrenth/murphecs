package ecs_test

import (
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

// This file contains tests that implement multiple components of the ECS

func TestArchetypeMoveInsideQUery(t *testing.T) {
	type componentA struct {
		ecs.Component

		value int
	}
	type componentB struct{ ecs.Component }

	assert := assert.New(t)

	world := ecs.NewDefaultWorld()
	entityId, err := ecs.Spawn(world, &componentA{value: 10})
	assert.NoError(err)

	query := ecs.Query1[componentA, ecs.Default]{}
	err = query.Prepare(world, nil)
	assert.NoError(err)
	err = query.Exec(world)
	assert.NoError(err)

	// update value before archetype change
	{
		query.Iter(func(entityId ecs.EntityId, component *componentA) {
			component.value += 10

			err = ecs.Insert(world, entityId, &componentB{})
			assert.NoError(err)
		})
		comp, err := ecs.Get1[componentA](world, entityId)
		assert.NoError(err)
		assert.Equal(20, comp.value)
	}

	// reset
	err = ecs.Delete(world, entityId)
	assert.NoError(err)
	entityId, err = ecs.Spawn(world, &componentA{value: 10})
	assert.NoError(err)
	err = query.Exec(world)
	assert.NoError(err)

	// update value after archetype change
	{
		query.Iter(func(entityId ecs.EntityId, component *componentA) {
			err = ecs.Insert(world, entityId, &componentB{})
			assert.NoError(err)

			component.value += 10 // <--- `component` now still points to the old location
		})
		comp, err := ecs.Get1[componentA](world, entityId)
		assert.NoError(err)
		assert.NotEqual(20, comp.value)
	}
}
