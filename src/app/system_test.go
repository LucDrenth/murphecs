package app

import (
	"errors"
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

func TestAddSystem(t *testing.T) {
	t.Run("error if adding a system that is not a function", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add("not a func", world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemNotAFunction)
	})

	t.Run("can use empty function as system", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("can use world as system param", func(t *testing.T) {
		assert := assert.New(t)
		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(world *ecs.World) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns an error when using non-pointer ecs.Query as system param", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(_ ecs.Query1[componentA, ecs.Default]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotAPointer)
	})

	t.Run("returns an error when using Query interface as system parameter", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(_ ecs.Query) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can use an ecs.Query as system param", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(_ *ecs.Query1[componentA, ecs.Default]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("fails when app is not aware of the outer world in the ecs.Query", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		outerWorlds := map[ecs.WorldId]*ecs.World{}
		eventStorage := newEventStorage()

		err := systemSet.add(func(_ *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {}, world, &outerWorlds, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
		assert.ErrorIs(err, ErrTargetWorldNotKnown)
	})

	t.Run("fails when app is not aware of the outer world in the ecs.Query", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err = systemSet.add(func(_ *ecs.Query1[componentA, ecs.QueryOptions2[ecs.TestCustomTargetWorld, ecs.Lazy]]) {}, world, &outerWorlds, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can insert ecs.Query that targets an outer world", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err = systemSet.add(func(_ *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {}, world, &outerWorlds, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns an error if a system parameter is invalid", func(t *testing.T) {
		type resourceA struct{}
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(_ resourceA) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamNotValid)
	})

	t.Run("can use a resource as system param by value", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		err := resourceStorage.add(&resourceA{})
		assert.NoError(err)
		eventStorage := newEventStorage()

		err = systemSet.add(func(_ resourceA) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("can use a resource as system param by reference", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		err := resourceStorage.add(&resourceA{})
		assert.NoError(err)
		eventStorage := newEventStorage()

		err = systemSet.add(func(_ *resourceA) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns err if system param EventReader is not a pointer", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(EventReader[*testEvent]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamEventReaderNotAPointer)
	})

	t.Run("can use system param *EventReader", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(*EventReader[*testEvent]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns err if system param EventWriter is not a pointer", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(EventWriter[*testEvent]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamEventWriterNotAPointer)
	})

	t.Run("can use system param *EventWriter", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func(*EventWriter[*testEvent]) {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns an error if the system returns something that is not an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() int { return 10 }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.ErrorIs(err, ErrSystemInvalidReturnType)
	})

	t.Run("can add a system that returns an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() error { return nil }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
	})
}

func TestExecSystem(t *testing.T) {
	t.Run("runs system without system params", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		didRun := false
		err := systemSet.add(func() { didRun = true }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.True(didRun)
	})

	t.Run("system with by-reference resource param can mutate the resource", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		err = systemSet.add(func(r *resourceA) { r.value = 20 }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.Equal(20, resource.value)
	})

	t.Run("system with by-value resource param can not mutate the resource", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := newEventStorage()
		resourceStorage := newResourceStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		err = systemSet.add(func(r resourceA) { r.value = 20 }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.Equal(10, resource.value)
	})

	t.Run("system with by-value resource uses a copy of the latest resource value", func(t *testing.T) {
		type resourceA struct{ value int }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := newEventStorage()
		resourceStorage := newResourceStorage()

		resource := resourceA{value: 10}

		err := resourceStorage.add(&resource)
		assert.NoError(err)

		// first system updates the resource
		err = systemSet.add(func(r *resourceA) { r.value = 20 }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		// second system should have the updated resource value
		err = systemSet.add(func(r resourceA) {
			assert.Equal(20, r.value)
			r.value = 30 // should not do anything
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		// third system should have the updated resource value from the first system
		err = systemSet.add(func(r resourceA) {
			assert.Equal(20, r.value)
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		// fourth system should also have the updated resource value from the first system
		err = systemSet.add(func(r *resourceA) {
			assert.Equal(20, r.value)
			r.value = 40
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		// fifth system should have an updated value from the fourth system
		err = systemSet.add(func(r resourceA) {
			assert.Equal(40, r.value)
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		systemSet.Exec(world, nil, &eventStorage)
	})

	t.Run("returns no errors if the system does not return anything", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() {}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		errors := systemSet.Exec(world, nil, &eventStorage)

		assert.Empty(errors)
	})

	t.Run("returns no errors if the system does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() error { return nil }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		errors := systemSet.Exec(world, nil, &eventStorage)

		assert.Empty(errors)
	})

	t.Run("returns an error if the system returns an error", func(t *testing.T) {
		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		err := systemSet.add(func() error { return errors.New("oops") }, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		errors := systemSet.Exec(world, nil, &eventStorage)

		assert.Len(errors, 1)
	})

	t.Run("automatically execute non-lazy queries in system params before executing systems", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		_, err := ecs.Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = systemSet.add(func(q *ecs.Query1[componentA, ecs.Default]) {
			numberOfResults = int(q.Result().NumberOfResult())
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)

		assert.Equal(1, numberOfResults)
	})

	t.Run("clear and not execute queries in system params before executing systems", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		_, err := ecs.Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = systemSet.add(func(q *ecs.Query1[componentA, ecs.Lazy]) {
			numberOfResults = int(q.Result().NumberOfResult())
		}, world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)

		assert.Equal(0, numberOfResults)
	})

	t.Run("executes query to outer world", func(t *testing.T) {
		type componentA struct{ ecs.Component }

		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		systemSet := SystemSet{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		_, err = ecs.Spawn(&outerWorld, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = systemSet.add(func(q *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {
			numberOfResults = int(q.Result().NumberOfResult())
		}, world, &outerWorlds, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		errors := systemSet.Exec(world, &outerWorlds, &eventStorage)
		assert.Empty(errors)

		assert.Equal(1, numberOfResults)
	})

	t.Run("event system", func(t *testing.T) {
		type testEvent struct {
			Event
			id int
		}

		const eventId = 3

		assert := assert.New(t)

		systemSet := SystemSet{id: 1}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		resourceStorage := newResourceStorage()
		eventStorage := newEventStorage()

		var events []*testEvent
		var doWriteEvent = true

		err := systemSet.add(
			func(eventWriter *EventWriter[*testEvent]) {
				if doWriteEvent {
					eventWriter.Write(&testEvent{id: 3})
				}
			},
			world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)
		err = systemSet.add(
			func(eventReader *EventReader[*testEvent]) {
				events = []*testEvent{}

				for event := range eventReader.Read {
					events = append(events, event)
				}
			},
			world, nil, &logger, &resourceStorage, &eventStorage)
		assert.NoError(err)

		// first run
		errors := systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.Empty(events)

		eventStorage.ProcessEvents(systemSet.id + 1) // pretend that next systemSet has run
		doWriteEvent = false

		// second run
		errors = systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.Len(events, 1)
		assert.Equal(eventId, events[0].id)

		eventStorage.ProcessEvents(systemSet.id + 1) // pretend that next systemSet has run

		// third run
		errors = systemSet.Exec(world, nil, &eventStorage)
		assert.Empty(errors)
		assert.Empty(events) // no need events read
	})
}
