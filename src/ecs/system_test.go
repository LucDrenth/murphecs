package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSystem(t *testing.T) {
	type componentA struct{ Component }

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
		err := simpleTestAddSystem(func(world *World) {})
		assert.NoError(err)
	})

	t.Run("returns an error when using non-pointer Query as system param", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ Query1[componentA, Default]) {})
		assert.ErrorIs(err, ErrSystemParamQueryNotAPointer)
	})

	t.Run("returns an error when using Query interface as system parameter", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ Query) {})
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can use a Query as system param", func(t *testing.T) {
		assert := assert.New(t)
		err := simpleTestAddSystem(func(_ *Query1[componentA, Default]) {})
		assert.NoError(err)
	})

	t.Run("fails when app is not aware of the outer world in the Query", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}

		outerWorlds := map[WorldId]*World{}
		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func(_ *Query1[componentA, TestCustomTargetWorld]) {}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
		assert.ErrorIs(err, ErrTargetWorldNotFound)
	})

	t.Run("fails when query targets outer world with lazy option", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ *Query1[componentA, QueryOptions2[TestCustomTargetWorld, Lazy]]) {}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.ErrorIs(err, ErrSystemParamQueryNotValid)
	})

	t.Run("can insert Query that targets an outer world", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		err = scheduleSystems.add(func(_ *Query1[componentA, TestCustomTargetWorld]) {}, "", world, &outerWorlds, &logger, &eventStorage)
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
		world := NewDefaultWorld()
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
		world := NewDefaultWorld()
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

	t.Run("returns err if system param OuterResource is a pointer", func(t *testing.T) {
		type resource struct{}
		assert := assert.New(t)
		err := simpleTestAddSystem(func(*OuterResource[*resource, TestCustomTargetWorld]) {})
		assert.ErrorIs(err, ErrSystemParamOuterResourceIsAPointer)
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
	world := NewDefaultWorld()
	logger := NoOpLogger{}

	eventStorage := NewEventStorage()

	return scheduleSystems.add(system, "", world, nil, &logger, &eventStorage)
}

func TestExecSystem(t *testing.T) {
	type componentA struct{ Component }
	type resourceA struct{ value int }

	t.Run("runs system without system params", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		didRun := false
		err := scheduleSystems.add(func() { didRun = true }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)
		assert.True(didRun)
	})

	t.Run("system with by-reference resource param can mutate the resource", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(r *resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)
		assert.Equal(20, resource.value)
	})

	t.Run("system with by-value resource param can not mutate the resource", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(r resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)
		assert.Equal(10, resource.value)
	})

	t.Run("system with by-value resource uses a copy of the latest resource value", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 10}

		err := world.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(r *resourceA) { r.value = 20 }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(20, r.value)
			r.value = 30
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(20, r.value)
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(r *resourceA) {
			assert.Equal(20, r.value)
			r.value = 40
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(r resourceA) {
			assert.Equal(40, r.value)
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		scheduleSystems.Exec(world, nil, &eventStorage, 1)
	})

	t.Run("OuterResource with resource pointer can mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res OuterResource[*resourceA, TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res OuterResource[*resourceA, TestCustomTargetWorld]) {
			assert.Equal(77, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		assert.NoError(scheduleSystems.prepare(&outerWorlds))

		errs := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errs)
	})

	t.Run("OuterResource without resource pointer can not mutate", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		err = scheduleSystems.add(func(res OuterResource[*resourceA, TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
			res.Value.value = 77
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		err = scheduleSystems.add(func(res OuterResource[resourceA, TestCustomTargetWorld]) {
			assert.Equal(7, res.Value.value)
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		assert.NoError(scheduleSystems.prepare(&outerWorlds))

		errs := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errs)
	})

	t.Run("non-pointer OuterResource is updated between executions", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		resource := resourceA{value: 7}
		err = outerWorld.Resources().Add(&resource)
		assert.NoError(err)

		execCount := 0
		err = scheduleSystems.add(func(res OuterResource[resourceA, TestCustomTargetWorld]) {
			if execCount == 0 {
				assert.Equal(7, res.Value.value)
			} else {
				assert.Equal(77, res.Value.value)
			}
			execCount++
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		assert.NoError(scheduleSystems.prepare(&outerWorlds))

		errs := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errs)

		resource.value = 77

		errs = scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 2)
		assert.Empty(errs)
		assert.Equal(2, execCount)
	})

	t.Run("returns no errors if the system does not return anything", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() {}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Empty(errs)
	})

	t.Run("returns no errors if the system does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() error { return nil }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Empty(errs)
	})

	t.Run("returns an error if the system returns an error", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		err := scheduleSystems.add(func() error { return errors.New("oops") }, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)
		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)

		assert.Len(errs, 1)
	})

	t.Run("automatically execute non-lazy queries in system params before executing systems", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		_, err := Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *Query1[componentA, Default]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)

		assert.Equal(1, numberOfResults)
	})

	t.Run("clear and not execute queries in system params before executing systems", func(t *testing.T) {
		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		_, err := Spawn(world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *Query1[componentA, Lazy]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)

		assert.Equal(0, numberOfResults)
	})

	t.Run("executes query to outer world", func(t *testing.T) {
		assert := assert.New(t)

		outerWorldConfigs := DefaultWorldConfigs()
		outerWorldConfigs.Id = &TestCustomTargetWorldId
		outerWorld, err := NewWorld(outerWorldConfigs)
		assert.NoError(err)
		outerWorlds := map[WorldId]*World{
			TestCustomTargetWorldId: &outerWorld,
		}

		scheduleSystems := ScheduleSystems{}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()

		_, err = Spawn(&outerWorld, &componentA{})
		assert.NoError(err)

		numberOfResults := 0
		err = scheduleSystems.add(func(q *Query1[componentA, TestCustomTargetWorld]) {
			numberOfResults = int(q.NumberOfResult())
		}, "", world, &outerWorlds, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, &outerWorlds, &eventStorage, 1)
		assert.Empty(errs)

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
		world := NewDefaultWorld()
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

		executeAll := func(tick uint) {
			for _, ss := range []*ScheduleSystems{scheduleSystems1, scheduleSystems2, scheduleSystems3} {
				errs := ss.Exec(world, nil, &eventStorage, tick)
				assert.Empty(errs)
			}
		}

		executeAll(1)

		assert.Empty(eventsFromScheduleSystems1)
		assert.Empty(eventsFromScheduleSystems2)
		assert.Len(eventsFromScheduleSystems3, 1)
		assert.Equal(eventId, eventsFromScheduleSystems3[0].id)

		doWriteEvent = false

		executeAll(2)

		assert.Len(eventsFromScheduleSystems1, 1)
		assert.Equal(eventId, eventsFromScheduleSystems1[0].id)
		assert.Len(eventsFromScheduleSystems2, 1)
		assert.Equal(eventId, eventsFromScheduleSystems2[0].id)
		assert.Empty(eventsFromScheduleSystems3)

		executeAll(3)
		assert.Empty(eventsFromScheduleSystems1)
		assert.Empty(eventsFromScheduleSystems2)
		assert.Empty(eventsFromScheduleSystems3)
	})

	t.Run("EventReader without an EventWriter", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{id: 1}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()
		err := scheduleSystems.add(
			func(eventReader *EventReader[*testEvent]) {
				for range eventReader.Read {
				}
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)
		errs = scheduleSystems.Exec(world, nil, &eventStorage, 2)
		assert.Empty(errs)
	})

	t.Run("EventWriter without an EventReader", func(t *testing.T) {
		type testEvent struct{ Event }

		assert := assert.New(t)

		scheduleSystems := ScheduleSystems{id: 1}
		world := NewDefaultWorld()
		logger := NoOpLogger{}
		eventStorage := NewEventStorage()
		err := scheduleSystems.add(
			func(eventWriter *EventWriter[*testEvent]) {
				eventWriter.Write(&testEvent{})
			},
			"", world, nil, &logger, &eventStorage)
		assert.NoError(err)

		errs := scheduleSystems.Exec(world, nil, &eventStorage, 1)
		assert.Empty(errs)
		errs = scheduleSystems.Exec(world, nil, &eventStorage, 2)
		assert.Empty(errs)
	})
}
