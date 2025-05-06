package app

import (
	"errors"
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
		assert.ErrorIs(err, ErrSystemNotAFunction)
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

	t.Run("returns an error when using non-pointer world as system param", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(world ecs.World) {}, &world, &logger, &resourceStorage)
		assert.ErrorIs(err, ErrSystemParamWorldNotAPointer)
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

		err := systemSet.add(func(_ ecs.Query1[componentA, ecs.Default]) {}, &world, &logger, &resourceStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotAPointer)
	})

	t.Run("can use an ecs query as system param", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(_ *ecs.Query1[componentA, ecs.Default]) {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
	})

	t.Run("returns an error if a system parameter is invalid", func(t *testing.T) {
		type resourceA struct{}
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func(_ resourceA) {}, &world, &logger, &resourceStorage)
		assert.ErrorIs(err, ErrSystemParamNotValid)
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

	t.Run("returns an error if the system returns something that is not an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() int { return 10 }, &world, &logger, &resourceStorage)
		assert.ErrorIs(err, ErrSystemInvalidReturnType)
	})

	t.Run("can add a system that returns an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() error { return nil }, &world, &logger, &resourceStorage)
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

		systemSet.exec(&world)
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
		systemSet.exec(&world)
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
		systemSet.exec(&world)
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

		systemSet.exec(&world)
	})

	t.Run("returns no errors if the system does not return anything", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() {}, &world, &logger, &resourceStorage)
		assert.NoError(err)
		errors := systemSet.exec(&world)

		assert.Empty(errors)
	})

	t.Run("returns no errors if the system does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() error { return nil }, &world, &logger, &resourceStorage)
		assert.NoError(err)
		errors := systemSet.exec(&world)

		assert.Empty(errors)
	})

	t.Run("returns an error if the system returns an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		err := systemSet.add(func() error { return errors.New("oops") }, &world, &logger, &resourceStorage)
		assert.NoError(err)
		errors := systemSet.exec(&world)

		assert.Len(errors, 1)
	})

	t.Run("automatically system param queries before executing systems", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewWorld()
		logger := log.NoOp()
		resourceStorage := newResourceStorage()

		_, err := ecs.Spawn(&world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = systemSet.add(func(q *ecs.Query1[componentA, ecs.QueryOptions[ecs.NoFilter, ecs.AllReadOnly]]) {
			numberOfResults = int(q.Result().NumberOfResult())
		}, &world, &logger, &resourceStorage)
		assert.NoError(err)

		errors := systemSet.exec(&world)
		assert.Empty(errors)

		assert.Equal(1, numberOfResults)
	})
}
