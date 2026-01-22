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

	query := ecs.Query1[*componentA, ecs.Default]{}
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

// TestMixingComponentTypes tests that it does not matter whether components are passed
// by value or by reference when spawning entities, as long as the types are the same.
func TestMixingComponentTypes(t *testing.T) {
	type componentA struct{ ecs.Component }
	type componentB struct{ ecs.Component }

	t.Run("spawn", func(t *testing.T) {
		assert := assert.New(t)
		world := ecs.NewDefaultWorld()

		entity1, err := ecs.Spawn(world, componentA{}, &componentB{})
		assert.NoError(err)
		entity2, err := ecs.Spawn(world, &componentA{}, componentB{})
		assert.NoError(err)
		entity3, err := ecs.Spawn(world, componentA{}, componentB{})
		assert.NoError(err)
		entity4, err := ecs.Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)

		assert.Equal(4, world.CountEntities())
		assert.Equal(8, world.CountComponents())

		query := ecs.Query2[*componentA, *componentB, ecs.Default]{}
		err = query.Prepare(world, nil)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(4), query.NumberOfResult())
		queriedEntities := []ecs.EntityId{}
		query.Iter(func(entityId ecs.EntityId, _ *componentA, _ *componentB) {
			queriedEntities = append(queriedEntities, entityId)
		})
		assert.Contains(queriedEntities, entity1)
		assert.Contains(queriedEntities, entity2)
		assert.Contains(queriedEntities, entity3)
		assert.Contains(queriedEntities, entity4)
	})

	t.Run("insert", func(t *testing.T) {
		assert := assert.New(t)
		world := ecs.NewDefaultWorld()

		entity1, err := ecs.Spawn(world)
		assert.NoError(err)
		err = ecs.Insert(world, entity1, componentA{}, &componentB{})
		assert.NoError(err)
		entity2, err := ecs.Spawn(world)
		assert.NoError(err)
		err = ecs.Insert(world, entity2, &componentA{}, componentB{})
		assert.NoError(err)
		entity3, err := ecs.Spawn(world)
		assert.NoError(err)
		err = ecs.Insert(world, entity3, componentA{}, componentB{})
		assert.NoError(err)
		entity4, err := ecs.Spawn(world)
		assert.NoError(err)
		err = ecs.Insert(world, entity4, &componentA{}, &componentB{})
		assert.NoError(err)

		assert.Equal(4, world.CountEntities())
		assert.Equal(8, world.CountComponents())

		query := ecs.Query2[*componentA, *componentB, ecs.Default]{}
		err = query.Prepare(world, nil)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(4), query.NumberOfResult())
		queriedEntities := []ecs.EntityId{}
		query.Iter(func(entityId ecs.EntityId, _ *componentA, _ *componentB) {
			queriedEntities = append(queriedEntities, entityId)
		})
		assert.Contains(queriedEntities, entity1)
		assert.Contains(queriedEntities, entity2)
		assert.Contains(queriedEntities, entity3)
		assert.Contains(queriedEntities, entity4)
	})
}
