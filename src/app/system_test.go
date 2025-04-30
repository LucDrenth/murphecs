package app

import (
	"testing"

	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/log"
	"github.com/stretchr/testify/assert"
)

func TestAddSystem(t *testing.T) {
	t.Run("error if adding a system that is not a function", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add("not a func", &world, &logger, &resourceStorage)
		assert.Error(err)
	})

	t.Run("can use empty function as system", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("can use world as system param", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(world *ecs.World) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("can use logger as system param", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(logger log.Logger) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("returns an error when using non-pointer ecs query as system param", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(_ ecs.Query1[componentA, ecs.NoFilter, ecs.NoOptional, ecs.NoReadOnly]) {}, &world, &logger, &resourceStorage)
		assert.Error(err)
	})

	t.Run("can use an ecs query as system param", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(_ *ecs.Query1[componentA, ecs.NoFilter, ecs.NoOptional, ecs.NoReadOnly]) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("can use a resource as system param by value", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()
		err := resourceStorage.add(&resourceA{})
		assert.NoError(err)

		err = systemSet.add(func(_ resourceA) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("can use a resource as system param by reference", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()
		err := resourceStorage.add(&resourceA{})
		assert.NoError(err)

		err = systemSet.add(func(_ *resourceA) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})
}

func TestExecSystem(t *testing.T) {
	t.Run("runs system without system params", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		didRun := false
		err := systemSet.add(func() { didRun = true }, &world, &logger, &resourceStorage)
		assert.NoError(err)

		systemSet.exec(&logger)
		assert.True(didRun)
	})

	t.Run("system with by-reference resource param can mutate the resource", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		err = systemSet.add(func(r *resourceA) { r.value = 20 }, &world, &logger, &resourceStorage)
		assert.NoError(err)
		systemSet.exec(&logger)
		assert.Equal(20, resource.value)
	})

	t.Run("system with by-value resource param can not mutate the resource", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		err = systemSet.add(func(r resourceA) { r.value = 20 }, &world, &logger, &resourceStorage)
		assert.NoError(err)
		systemSet.exec(&logger)
		assert.Equal(10, resource.value)
	})

	t.Run("system with by-value resource uses a copy of the latest resource value", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		// first system updates the resource
		err = systemSet.add(func(r *resourceA) { r.value = 20 }, &world, &logger, &resourceStorage)
		assert.NoError(err)

		// second system should have the updated resource value
		err = systemSet.add(func(r resourceA) {
			assert.Equal(20, r.value)
			r.value = 30 // should not do anything
		}, &world, &logger, &resourceStorage)
		assert.NoError(err)

		// third system should have the updated resource value from the first system
		err = systemSet.add(func(r resourceA) {
			assert.Equal(20, r.value)
		}, &world, &logger, &resourceStorage)
		assert.NoError(err)

		// fourth system should also have the updated resource value from the first system
		err = systemSet.add(func(r *resourceA) {
			assert.Equal(20, r.value)
			r.value = 40
		}, &world, &logger, &resourceStorage)
		assert.NoError(err)

		// fifth system should have an updated value from the fourth system
		err = systemSet.add(func(r resourceA) {
			assert.Equal(40, r.value)
		}, &world, &logger, &resourceStorage)
		assert.NoError(err)

		systemSet.exec(&logger)
	})
}
