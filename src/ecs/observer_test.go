package ecs_test

import (
	"testing"

	. "github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

// TestGlobalObserver tests observers in the [ecs.World] global observer registry
func TestGlobalObserver(t *testing.T) {
	type myComponent1 struct{ Component }
	type myComponent2 struct{ Component }

	type observer1 struct{ Observer }
	type observer2 struct{ Observer }

	t.Run("Spawn nil observer does nothing", func(t *testing.T) {
		world := NewDefaultWorld()
		assert.NoError(t, On[observer1](world, nil))
	})

	t.Run("Spawn observer pointer panics", func(t *testing.T) {
		world := NewDefaultWorld()
		assert.Panics(t, func() {
			_ = On[*observer1](world, func(world *World, observer *observer1) {})
		})
	})

	t.Run("Custom observer can be registered and triggered", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		var observed1, observed2 uint

		assert.NoError(On[observer1](world, func(world *World, observer observer1) {
			observed1++
		}))
		assert.NoError(On[observer2](world, func(world *World, observer observer2) {
			observed2++
		}))

		Trigger(world, observer1{})
		Trigger(world, observer2{})
		Trigger(world, observer1{})
		Trigger(world, observer1{})

		assert.Equal(uint(3), observed1)
		assert.Equal(uint(1), observed2)
	})

	t.Run("Triggering observer without registering an observer does nothing", func(t *testing.T) {
		world := NewDefaultWorld()
		Trigger(world, observer1{})
	})

	t.Run("OnSpawn", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		nrObserved := 0
		var expectedEntityId EntityId

		assert.NoError(On[OnSpawn[myComponent1]](world, func(world *World, observed OnSpawn[myComponent1]) {
			nrObserved++
			assert.Equal(expectedEntityId, observed.Entity)
		}))

		assert.NoError(On[OnDespawn[myComponent1]](world, func(world *World, observed OnDespawn[myComponent1]) {
			assert.FailNow("did not expect OnDespawn to trigger")
		}))

		expectedEntityId = 1
		_, err := Spawn(world, myComponent1{}) // triggers
		assert.NoError(err)
		_, err = Spawn(world, myComponent2{}) // does not trigger
		assert.NoError(err)
		expectedEntityId = 3
		_, err = Spawn(world, myComponent1{}, myComponent2{}) // triggers
		assert.NoError(err)
		expectedEntityId = 4
		_, err = Spawn(world, myComponent2{}, myComponent1{}) // triggers
		assert.NoError(err)

		assert.Equal(3, nrObserved)
	})

	t.Run("OnDespawn", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		nrObserved := 0
		var expectedEntityId EntityId

		assert.NoError(On[OnDespawn[myComponent1]](world, func(world *World, observed OnDespawn[myComponent1]) {
			nrObserved++
			assert.Equal(expectedEntityId, observed.Entity)
		}))

		id1, err := Spawn(world, myComponent1{}) // despawn will trigger
		assert.NoError(err)
		_, err = Spawn(world, myComponent2{}) // despawn won't trigger
		assert.NoError(err)
		id3, err := Spawn(world, myComponent1{}, myComponent2{}) // despawn will trigger
		assert.NoError(err)
		id4, err := Spawn(world, myComponent2{}, myComponent1{}) // despawn will trigger
		assert.NoError(err)

		assert.Equal(0, nrObserved)

		expectedEntityId = id3
		err = Despawn(world, id3)
		assert.NoError(err)
		expectedEntityId = id1
		err = Despawn(world, id1)
		assert.NoError(err)
		expectedEntityId = id4
		err = Despawn(world, id4)
		assert.NoError(err)

		assert.Equal(3, nrObserved)
	})
}

func TestEntityObserver(t *testing.T) {
	type myComponent1 struct{ Component }
	type myComponent2 struct{ Component }

	type observer1 struct{ Observer }
	type observer2 struct{ Observer }

	t.Run("err when entity not found", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := Observe[observer1](world, EntityId(5), func(world *World, o observer1) {})
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("can trigger non-observed", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world)
		assert.NoError(err)
		err = TriggerEntity(world, entity, observer1{})
		assert.NoError(err)
	})

	t.Run("Custom observer", func(t *testing.T) {
		t.Run("TriggerEntity only triggers for the given entity", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()
			e1, err := Spawn(world)
			assert.NoError(err)
			e2, err := Spawn(world)
			assert.NoError(err)

			var triggersForEntity1, triggersForEntity2 uint

			err = Observe[observer1](world, e1, func(world *World, o observer1) { triggersForEntity1++ })
			assert.NoError(err)
			err = Observe[observer1](world, e2, func(world *World, o observer1) { triggersForEntity2++ })
			assert.NoError(err)

			err = TriggerEntity(world, e1, observer1{})
			assert.NoError(err)
			assert.Equal(uint(1), triggersForEntity1)
			assert.Equal(uint(0), triggersForEntity2)
		})

		t.Run("TriggerEntity only triggers for the given observer", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()
			entity, err := Spawn(world)
			assert.NoError(err)

			var triggersForObserve1, triggersForObserve2 uint
			err = Observe[observer1](world, entity, func(world *World, o observer1) { triggersForObserve1++ })
			assert.NoError(err)
			err = Observe[observer2](world, entity, func(world *World, o observer2) { triggersForObserve2++ })
			assert.NoError(err)

			err = TriggerEntity(world, entity, observer1{})
			assert.NoError(err)
			assert.Equal(uint(1), triggersForObserve1)
			assert.Equal(uint(0), triggersForObserve2)
		})
	})

	t.Run("OnDespawn", func(t *testing.T) {
		t.Run("gets triggered when removing component", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()

			numberOfTriggers := 0

			entity, err := Spawn(world, myComponent1{})
			assert.NoError(err)
			assert.NoError(Observe[OnDespawn[myComponent1]](world, entity, func(world *World, o OnDespawn[myComponent1]) { numberOfTriggers++ }))
			assert.NoError(Remove1[myComponent1](world, entity))

			assert.Equal(1, numberOfTriggers)
		})

		t.Run("gets triggered when despawning entity", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()

			numberOfTriggers := 0

			entity, err := Spawn(world, myComponent1{})
			assert.NoError(err)
			assert.NoError(Observe[OnDespawn[myComponent1]](world, entity, func(world *World, o OnDespawn[myComponent1]) { numberOfTriggers++ }))
			assert.NoError(Despawn(world, entity))

			assert.Equal(1, numberOfTriggers)
		})
	})

	t.Run("OnSpawn", func(t *testing.T) {
		t.Run("gets triggered on Insert", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()

			numberOfTriggers := 0

			entity, err := Spawn(world, myComponent1{})
			assert.NoError(err)
			assert.NoError(Observe[OnSpawn[myComponent2]](world, entity, func(world *World, o OnSpawn[myComponent2]) { numberOfTriggers++ }))
			assert.NoError(Insert(world, entity, myComponent2{}))

			assert.Equal(1, numberOfTriggers)
		})

		t.Run("gets triggered on InsertOrOverwrite", func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()

			numberOfTriggers := 0

			entity, err := Spawn(world, myComponent1{})
			assert.NoError(err)
			assert.NoError(Observe[OnSpawn[myComponent2]](world, entity, func(world *World, o OnSpawn[myComponent2]) { numberOfTriggers++ }))
			assert.NoError(InsertOrOverwrite(world, entity, myComponent2{}))

			assert.Equal(1, numberOfTriggers)
		})
	})
}
