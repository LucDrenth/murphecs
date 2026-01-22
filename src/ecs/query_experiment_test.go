package ecs_test

import (
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

func TestQueryExperiment(t *testing.T) {
	type componentA struct {
		ecs.Component
		Value int
	}

	type componentB struct {
		ecs.Component
		Value string
	}

	t.Run("by value", func(t *testing.T) {
		assert := assert.New(t)

		world := ecs.NewDefaultWorld()

		_, err := ecs.Spawn(world, componentA{Value: 10})
		assert.NoError(err)

		query := ecs.Query1Experimental[componentA, ecs.Default]{}
		err = query.Prepare(world, nil)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId ecs.EntityId, a componentA) {
			a.Value += 15
		})
		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId ecs.EntityId, a componentA) {
			assert.Equal(10, a.Value)
		})
	})

	t.Run("by reference", func(t *testing.T) {
		assert := assert.New(t)

		world := ecs.NewDefaultWorld()

		_, err := ecs.Spawn(world, componentA{Value: 10})
		assert.NoError(err)

		query := ecs.Query1Experimental[*componentA, ecs.Default]{}
		err = query.Prepare(world, nil)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId ecs.EntityId, a *componentA) {
			a.Value += 15
		})
		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId ecs.EntityId, a *componentA) {
			assert.Equal(25, a.Value)
		})
	})
}
