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

func TestObserverSystemParams(t *testing.T) {
	type myObserver struct{ Observer }

	t.Run("Resource", func(t *testing.T) {
		type myResource struct{ value int }

		t.Run("can register observer with resource by pointer", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.Resources().Add(&myResource{}))
			assert.NoError(On[myObserver](world, func(_ myObserver, _ *myResource) {}))
		})

		t.Run("can register observer with resource by value", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.Resources().Add(&myResource{}))
			assert.NoError(On[myObserver](world, func(_ myObserver, _ myResource) {}))
		})

		t.Run("resource by pointer can be mutated", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.Resources().Add(&myResource{value: 1}))
			assert.NoError(On[myObserver](world, func(_ myObserver, res *myResource) { res.value++ }))
			Trigger(world, myObserver{})
			res, err := GetResource[myResource](world)
			assert.NoError(err)
			assert.Equal(2, res.value)
		})

		t.Run("resource by value cannot be mutated", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.Resources().Add(&myResource{value: 1}))
			assert.NoError(On[myObserver](world, func(_ myObserver, res myResource) { res.value++ }))
			Trigger(world, myObserver{})
			res, err := GetResource[myResource](world)
			assert.NoError(err)
			assert.Equal(1, res.value)
		})

		t.Run("fails if resource is not added", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.Error(On[myObserver](world, func(_ myObserver, _ *myResource) {}))
		})
	})

	t.Run("EventReader", func(t *testing.T) {
		type myEvent struct{ Event }

		t.Run("can register observer with EventReader", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(On[myObserver](world, func(_ myObserver, _ *EventReader[*myEvent]) {}))
		})

		t.Run("fails if EventReader is not a pointer", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.ErrorIs(On[myObserver](world, func(_ myObserver, _ EventReader[*myEvent]) {}), ErrSystemParamEventReaderNotAPointer)
		})

		t.Run("observer can read events written by a schedule system in the same tick", func(t *testing.T) {
			// The observer reads events that were written by an earlier system in the same schedule run.
			// ProcessEvents runs after all systems, so the event is in the reader by the next schedule run.
			// To read within the same tick, the observer must be triggered AFTER the writing system runs.
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.AddSchedule("update", ScheduleLast{}, false))

			var readEvents []*myEvent
			assert.NoError(On[myObserver](world, func(_ myObserver, reader *EventReader[*myEvent]) {
				readEvents = nil
				for e := range reader.Read {
					readEvents = append(readEvents, e)
				}
			}))

			// Writing system runs first in the schedule
			assert.NoError(world.AddSystem("update", func(w *EventWriter[*myEvent]) {
				w.Write(&myEvent{})
			}))
			assert.NoError(world.PrepareSystems())

			schedules, err := world.GetScheduleSystemsBySchedules([]Schedule{"update"})
			assert.NoError(err)
			eventStorage := world.Events()

			// Run once: writing system writes to writer; ProcessEvents moves it to reader
			schedules[0].Exec(world, nil, eventStorage, 1)
			// Observer triggered after schedule: events are now in the reader
			Trigger(world, myObserver{})

			assert.Len(readEvents, 1)
		})
	})

	t.Run("EventWriter", func(t *testing.T) {
		type myEvent struct {
			Event
			value int
		}

		t.Run("can register observer with EventWriter", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(On[myObserver](world, func(_ myObserver, _ *EventWriter[*myEvent]) {}))
		})

		t.Run("fails if EventWriter is not a pointer", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.ErrorIs(On[myObserver](world, func(_ myObserver, _ EventWriter[*myEvent]) {}), ErrSystemParamEventWriterNotAPointer)
		})

		t.Run("events written by observer inside a schedule are readable in subsequent schedule run", func(t *testing.T) {
			// Observer is triggered by Spawn inside a system. The observer's event writer uses
			// the same schedule ID as the running schedule (same lifecycle as regular event writers).
			type spawnedComponent struct{ Component }

			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.AddSchedule("update", ScheduleLast{}, false))

			assert.NoError(On[OnSpawn[spawnedComponent]](world, func(_ OnSpawn[spawnedComponent], w *EventWriter[*myEvent]) {
				w.Write(&myEvent{value: 42})
			}))

			var readEvents []*myEvent
			assert.NoError(world.AddSystem("update", func(reader *EventReader[*myEvent]) {
				readEvents = nil
				for e := range reader.Read {
					readEvents = append(readEvents, e)
				}
			}))
			// This system triggers the observer by spawning
			assert.NoError(world.AddSystem("update", func(w *World) {
				_, _ = Spawn(w, spawnedComponent{})
			}))
			assert.NoError(world.PrepareSystems())

			schedules, err := world.GetScheduleSystemsBySchedules([]Schedule{"update"})
			assert.NoError(err)
			eventStorage := world.Events()

			// Run 1: Spawn triggers observer which writes event; ProcessEvents moves it to reader
			schedules[0].Exec(world, nil, eventStorage, 1)
			assert.Empty(readEvents) // reader system ran before the spawning system

			// Run 2: reader system sees the event
			schedules[0].Exec(world, nil, eventStorage, 2)
			assert.Len(readEvents, 1)
			assert.Equal(42, readEvents[0].value)
		})

		t.Run("observer events follow the same clearing lifecycle as schedule events", func(t *testing.T) {
			// Events written by an observer inside schedule S are cleared when schedule S runs again.
			type spawnedComponent struct{ Component }

			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.NoError(world.AddSchedule("update", ScheduleLast{}, false))

			spawnOnce := true
			assert.NoError(world.AddSystem("update", func(w *World) {
				if spawnOnce {
					spawnOnce = false
					_, _ = Spawn(w, spawnedComponent{})
				}
			}))

			assert.NoError(On[OnSpawn[spawnedComponent]](world, func(_ OnSpawn[spawnedComponent], w *EventWriter[*myEvent]) {
				w.Write(&myEvent{value: 1})
			}))

			var readEvents []*myEvent
			assert.NoError(world.AddSystem("update", func(reader *EventReader[*myEvent]) {
				readEvents = nil
				for e := range reader.Read {
					readEvents = append(readEvents, e)
				}
			}))
			assert.NoError(world.PrepareSystems())
			schedules, err := world.GetScheduleSystemsBySchedules([]Schedule{"update"})
			assert.NoError(err)
			eventStorage := world.Events()

			schedules[0].Exec(world, nil, eventStorage, 1) // spawn triggers observer; event goes to writer; ProcessEvents moves to reader
			assert.Empty(readEvents)                       // reader system ran before the spawning system

			schedules[0].Exec(world, nil, eventStorage, 2) // reader system reads the event; ProcessEvents clears it
			assert.Len(readEvents, 1)

			schedules[0].Exec(world, nil, eventStorage, 3) // event cleared by previous ProcessEvents
			assert.Empty(readEvents)
		})
	})

	t.Run("OuterResource", func(t *testing.T) {
		type myResource struct{ value int }

		t.Run("fails if OuterResource is a pointer", func(t *testing.T) {
			assert := assert.New(t)
			world := NewDefaultWorld()
			assert.ErrorIs(
				On[myObserver](world, func(_ myObserver, _ *OuterResource[*myResource, TestCustomTargetWorld]) {}),
				ErrSystemParamOuterResourceIsAPointer,
			)
		})

		t.Run("can register observer with OuterResource pointer resource", func(t *testing.T) {
			assert := assert.New(t)

			outerWorldConfigs := DefaultWorldConfigs()
			outerWorldConfigs.Id = &TestCustomTargetWorldId
			outerWorld, err := NewWorld(outerWorldConfigs)
			assert.NoError(err)
			assert.NoError(outerWorld.Resources().Add(&myResource{value: 10}))

			world := NewDefaultWorld()
			assert.NoError(world.RegisterOuterWorld(TestCustomTargetWorldId, &outerWorld))

			assert.NoError(On[myObserver](world, func(_ myObserver, _ OuterResource[*myResource, TestCustomTargetWorld]) {}))
		})

		t.Run("outer resource pointer value is accessible when observer triggers", func(t *testing.T) {
			assert := assert.New(t)

			outerWorldConfigs := DefaultWorldConfigs()
			outerWorldConfigs.Id = &TestCustomTargetWorldId
			outerWorld, err := NewWorld(outerWorldConfigs)
			assert.NoError(err)
			assert.NoError(outerWorld.Resources().Add(&myResource{value: 10}))

			world := NewDefaultWorld()
			assert.NoError(world.RegisterOuterWorld(TestCustomTargetWorldId, &outerWorld))

			var gotValue int
			assert.NoError(On[myObserver](world, func(_ myObserver, res OuterResource[*myResource, TestCustomTargetWorld]) {
				gotValue = res.Value.value
			}))

			Trigger(world, myObserver{})
			assert.Equal(10, gotValue)
		})

		t.Run("non-pointer outer resource is refreshed between observer triggers", func(t *testing.T) {
			assert := assert.New(t)

			outerWorldConfigs := DefaultWorldConfigs()
			outerWorldConfigs.Id = &TestCustomTargetWorldId
			outerWorld, err := NewWorld(outerWorldConfigs)
			assert.NoError(err)
			res := &myResource{value: 10}
			assert.NoError(outerWorld.Resources().Add(res))

			world := NewDefaultWorld()
			assert.NoError(world.RegisterOuterWorld(TestCustomTargetWorldId, &outerWorld))

			var gotValue int
			assert.NoError(On[myObserver](world, func(_ myObserver, r OuterResource[myResource, TestCustomTargetWorld]) {
				gotValue = r.Value.value
			}))

			Trigger(world, myObserver{})
			assert.Equal(10, gotValue)

			res.value = 99
			Trigger(world, myObserver{})
			assert.Equal(99, gotValue)
		})
	})
}
