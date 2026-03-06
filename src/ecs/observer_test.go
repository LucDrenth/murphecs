package ecs_test

import (
	"testing"

	. "github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

func TestObserver(t *testing.T) {
	type myComponent1 struct{ Component }
	type myComponent2 struct{ Component }

	type observer1 struct{ Observer }
	type observer2 struct{ Observer }

	t.Run("Spawn nil observer does nothing", func(t *testing.T) {
		world := NewDefaultWorld()
		On[observer1](world, nil)
	})

	t.Run("Spawn observer pointer panics", func(t *testing.T) {
		world := NewDefaultWorld()
		assert.Panics(t, func() {
			On(world, func(world *World, observer *observer1) {})
		})
	})

	t.Run("Custom observer can be registered and triggered", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		var observed1, observed2 uint

		On(world, func(world *World, observer observer1) {
			observed1++
		})
		On(world, func(world *World, observer observer2) {
			observed2++
		})

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

		On(world, func(world *World, observed OnSpawn[myComponent1]) {
			nrObserved++
			assert.Equal(expectedEntityId, observed.Entity)
		})

		On(world, func(world *World, observed OnDespawn[myComponent1]) {
			assert.FailNow("did not expect OnDespawn to trigger")
		})

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

		On(world, func(world *World, observed OnDespawn[myComponent1]) {
			nrObserved++
			assert.Equal(expectedEntityId, observed.Entity)
		})

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
