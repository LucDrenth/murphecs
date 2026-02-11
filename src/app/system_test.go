package app

import (
	"errors"
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

func TestAddSystem(t *testing.T) {
	type componentA struct{ ecs.Component }

	t.Run("error if adding an invalid system", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem("not a func")
		assert.ErrorIs(err, ErrSystemTypeNotValid)
	})

	t.Run("can use empty function as system", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func() {})
		assert.NoError(err)
	})

	t.Run("can use world as system param", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(world *ecs.World) {})
		assert.NoError(err)
	})

	t.Run("returns an error when using non-pointer ecs.Query as system param", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ ecs.Query1[componentA, ecs.Default]) {})
		assert.ErrorIs(err, ErrSystemParamQueryNotAPointer)
	})

	t.Run("returns an error when using Query interface as system parameter", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ ecs.Query) {})
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can use an ecs.Query as system param", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ *ecs.Query1[componentA, ecs.Default]) {})
		assert.NoError(err)
	})

	t.Run("fails when app is not aware of the outer world in the ecs.Query", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		outerWorlds := map[ecs.WorldId]*ecs.World{}
		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func(_ *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
		assert.ErrorIs(err, ecs.ErrTargetWorldNotFound)
	})

	t.Run("fails when app is not aware of the outer world in the ecs.Query", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ *ecs.Query1[componentA, ecs.QueryOptions2[ecs.TestCustomTargetWorld, ecs.Lazy]]) {}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can insert ecs.Query that targets an outer world", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns an error if a system parameter is invalid", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ resourceA) {})
		assert.ErrorIs(err, ErrSystemParamNotValid)
	})

	t.Run("can use a resource as system param by value", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		err := world.Resources().Add(&resourceA{})
		assert.NoError(err)
		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ resourceA) {}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
	})

	t.Run("can use a resource as system param by reference", func(t *testing.T) {
		type resourceA struct{}

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		err := world.Resources().Add(&resourceA{})
		assert.NoError(err)
		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ *resourceA) {}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
	})

	t.Run("returns err if system param EventReader is not a pointer", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)
		err := simpleTestAddSystem(func(EventReader[*testEvent]) {})
		assert.ErrorIs(err, ErrSystemParamEventReaderNotAPointer)
	})

	t.Run("can use system param *EventReader", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)
		err := simpleTestAddSystem(func(*EventReader[*testEvent]) {})
		assert.NoError(err)
	})

	t.Run("returns err if system param EventWriter is not a pointer", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)
		err := simpleTestAddSystem(func(EventWriter[*testEvent]) {})
		assert.ErrorIs(err, ErrSystemParamEventWriterNotAPointer)
	})

	t.Run("can use system param *EventWriter", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)
		err := simpleTestAddSystem(func(*EventWriter[*testEvent]) {})
		assert.NoError(err)
	})

	t.Run("returns an error if the system returns something that is not an error", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func() int { return 10 })
		assert.ErrorIs(err, ErrSystemInvalidReturnType)
	})

	t.Run("can add a system that returns an error", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func() error { return nil })
		assert.NoError(err)
	})

	t.Run("can add system with Systems function", func(t *testing.T) {
		assert := assert.New(t)

		func1 := func() {}
		func2 := func() {}

		err := simpleTestAddSystem(Systems(func1))
		assert.NoError(err)

		err = simpleTestAddSystem(Systems(func1, func2))
		assert.NoError(err)

		err = simpleTestAddSystem(Systems(func1, func2, "not a function"))
		assert.ErrorIs(err, ErrSystemNotAFunction)
	})
}

func simpleTestAddSystem(system System) error {
	scheduleSystems := ScheduleSystems{}
	world := ecs.NewDefaultWorld()
	logger := NoOpLogger{}

	eventStorage := NewEventStorage()

	return scheduleSystems.add(system, "", world, nil, &logger, &eventStorage)
}

func TestExecSystem(t *testing.T) {
	type componentA struct{ ecs.Component }
	type resourceA struct{ value int }

	t.Run("runs system without system params", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		didRun := false
		err := scheduleSystems.add(func() { didRun = true }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)
		assert.True(didRun)
	})

	t.Run("system with by-reference resource param can mutate the resource", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(r *resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)
		assert.Equal(20, resource.value)
	})

	t.Run("system with by-value resource param can not mutate the resource", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(r resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)
		assert.Equal(10, resource.value)
	})

	t.Run("system with by-value resource uses a copy of the latest resource value", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		// first system updates the resource
		err = scheduleSystems.add(func(r *resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		// second system should have the updated resource value
		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(20, r.value)
			r.value = 30 // should not do anything
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		// third system should have the updated resource value from the first system
		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(20, r.value)
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		// fourth system should also have the updated resource value from the first system
		err = scheduleSystems.add(func(r *resourceA) {
			assert.Equal(20, r.value)
			r.value = 40
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		// fifth system should have an updated value from the fourth system
		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(40, r.value)
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		scheduleSystems.Exec(world, nil, &eventStorage, 1)
	})

	t.Run("OuterResource pointer with resource pointer can mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res *ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res *ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(77, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errors)
	})

	t.Run("OuterResource pointer without resource pointer can not mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res *ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res *ecs.OuterResource[resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errors)
	})

	t.Run("OuterResource with resource pointer can mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(77, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errors)
	})

	t.Run("OuterResource without resource pointer can not mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res ecs.OuterResource[*resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res ecs.OuterResource[resourceA, ecs.TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errors)
	})

	t.Run("returns no errors if the system does not return anything", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() {}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Empty(errors)
	})

	t.Run("returns no errors if the system does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() error { return nil }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Empty(errors)
	})

	t.Run("returns an error if the system returns an error", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() error { return errors.New("oops") }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Len(errors, 1)
	})

	t.Run("automatically execute non-lazy queries in system params before executing systems", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		_, err := ecs.Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *ecs.Query1[componentA, ecs.Default]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)

		assert.Equal(1, numberOfResults)
	})

	t.Run("clear and not execute queries in system params before executing systems", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		_, err := ecs.Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *ecs.Query1[componentA, ecs.Lazy]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)

		assert.Equal(0, numberOfResults)
	})

	t.Run("executes query to outer world", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := ecs.DefaultWorldConfigs()
		outerWorldConfigs.Id = &ecs.TestCustomTargetWorldId
		outerWorld, err := ecs.NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[ecs.WorldId]*ecs.World{
			ecs.TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		_, err = ecs.Spawn(&outerWorld, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *ecs.Query1[componentA, ecs.TestCustomTargetWorld]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
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

		scheduleSystems1 := &ScheduleSystems{id: 1}
		scheduleSystems2 := &ScheduleSystems{id: 2}
		scheduleSystems3 := &ScheduleSystems{id: 3}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()

		eventsFromScheduleSystems1 := []*testEvent{}
		eventsFromScheduleSystems2 := []*testEvent{}
		eventsFromScheduleSystems3 := []*testEvent{}
		doWriteEvent := true

		err := scheduleSystems1.add(
			func(eventReader *EventReader[*testEvent]) {
				eventsFromScheduleSystems1 = []*testEvent{}

				for event := range eventReader.Read {
					eventsFromScheduleSystems1 = append(eventsFromScheduleSystems1, event)
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems2.add(
			func(eventWriter *EventWriter[*testEvent]) {
				if doWriteEvent {
					eventWriter.Write(&testEvent{id: 3})
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		err = scheduleSystems2.add(
			func(eventReader *EventReader[*testEvent]) {
				eventsFromScheduleSystems2 = []*testEvent{}

				for event := range eventReader.Read {
					eventsFromScheduleSystems2 = append(eventsFromScheduleSystems2, event)
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems3.add(
			func(eventReader *EventReader[*testEvent]) {
				eventsFromScheduleSystems3 = []*testEvent{}

				for event := range eventReader.Read {
					eventsFromScheduleSystems3 = append(eventsFromScheduleSystems3, event)
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		executeScheduleSystemss := func(tick uint) {
			for _, scheduleSystems := range []*ScheduleSystems{scheduleSystems1, scheduleSystems2, scheduleSystems3} {
				errors := scheduleSystems.Exec(world, nil, &eventStorage, tick)
				assert.Empty(errors)
			}
		}

		// first run (event written by scheduleSystems1)
		executeScheduleSystemss(1)

		assert.Empty(eventsFromScheduleSystems1)
		assert.Empty(eventsFromScheduleSystems2)
		assert.Len(eventsFromScheduleSystems3, 1)
		assert.Equal(eventId, eventsFromScheduleSystems3[0].id)

		doWriteEvent = false

		// second run
		executeScheduleSystemss(2)

		assert.Len(eventsFromScheduleSystems1, 1)
		assert.Equal(eventId, eventsFromScheduleSystems1[0].id)
		assert.Len(eventsFromScheduleSystems2, 1)
		assert.Equal(eventId, eventsFromScheduleSystems2[0].id)
		assert.Empty(eventsFromScheduleSystems3)

		// third run
		executeScheduleSystemss(3)
		assert.Empty(eventsFromScheduleSystems1)
		assert.Empty(eventsFromScheduleSystems2)
		assert.Empty(eventsFromScheduleSystems3)
	})

	t.Run("EventReader without an EventWriter", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{id: 1}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()
		err := scheduleSystems.add(
			func(eventReader *EventReader[*testEvent]) {
				for range eventReader.Read {
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)
		errors = scheduleSystems.Exec(world, nil, &eventStorage, 2)
		assert.Empty(errors)
	})

	t.Run("EventWriter without an EventReader", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{id: 1}
		world := ecs.NewDefaultWorld()
		logger := NoOpLogger{}

		eventStorage := NewEventStorage()
		err := scheduleSystems.add(
			func(eventWriter *EventWriter[*testEvent]) {
				eventWriter.Write(&testEvent{})
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errors := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errors)
		errors = scheduleSystems.Exec(world, nil, &eventStorage, 2)
		assert.Empty(errors)
	})
}
