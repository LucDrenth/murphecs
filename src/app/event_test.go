package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventWriter(t *testing.T) {
	type testEvent struct {
		Event
		id int
	}

	t.Run("Write adds an event with the EventWriter its current ScheduleSystemsId", func(t *testing.T) {
		assert := assert.New(t)

		const numberOfEvents int = 5
		const scheduleSystemsId ScheduleSystemsId = ScheduleSystemsId(10)

		eventWriter := EventWriter[*testEvent]{}
		eventWriter.setScheduleSystemsWriter(scheduleSystemsId)

		for i := range numberOfEvents {
			eventWriter.Write(&testEvent{id: i})
			assert.Len(eventWriter.events, i+1)
			assert.Equal(scheduleSystemsId, eventWriter.events[i].scheduleSystemsWriter)
		}
	})

	t.Run("extractEvents returns all events in EventWriter and clears the events", func(t *testing.T) {
		assert := assert.New(t)

		const numberOfEvents int = 3
		const scheduleSystemsId ScheduleSystemsId = ScheduleSystemsId(11)

		eventWriter := EventWriter[*testEvent]{}
		eventWriter.setScheduleSystemsWriter(scheduleSystemsId)
		for i := range numberOfEvents {
			eventWriter.Write(&testEvent{id: i})
		}

		extractedEvents := eventWriter.extractEvents(0)
		assert.Len(extractedEvents, numberOfEvents)
		assert.Len(eventWriter.events, 0)
		for i, event := range extractedEvents {
			event, ok := reflect.TypeAssert[*testEvent](event)
			assert.True(ok)
			assert.Equal(scheduleSystemsId, event.scheduleSystemsWriter)
			assert.Equal(i, event.id)
		}
	})
}

func TestEventReader(t *testing.T) {
	type testEvent struct {
		Event
		id int
	}

	t.Run("First returns false if there are no elements", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{},
		}

		_, found := eventReader.First()
		assert.False(found)
	})

	t.Run("First returns the first element", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
			},
		}

		element, found := eventReader.First()
		assert.True(found)
		assert.Equal(1, element.id)
	})

	t.Run("First returns the first element that is not marked as removed", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
			},
		}
		eventReader.events[0].Remove()

		element, found := eventReader.First()
		assert.True(found)
		assert.Equal(2, element.id)
	})

	t.Run("Last returns false if there are no elements", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{},
		}

		_, found := eventReader.Last()
		assert.False(found)
	})

	t.Run("Last returns the last element", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
			},
		}

		element, found := eventReader.Last()
		assert.True(found)
		assert.Equal(2, element.id)
	})

	t.Run("Last returns the last element that is not marked as removed", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
			},
		}
		eventReader.events[1].Remove()

		element, found := eventReader.Last()
		assert.True(found)
		assert.Equal(1, element.id)
	})

	t.Run("Len returns the number of elements that are not marked as removed", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
				{id: 3},
			},
		}
		eventReader.events[1].Remove()

		assert.Equal(2, eventReader.Len())
	})

	t.Run("Empty returns true if all elements are marked as removed", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1},
				{id: 2},
			},
		}
		eventReader.events[0].Remove()
		assert.False(eventReader.IsEmpty())
		eventReader.events[1].Remove()
		assert.True(eventReader.IsEmpty())
	})

	t.Run("Empty returns true if there are no elements", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{},
		}
		assert.True(eventReader.IsEmpty())
	})

	t.Run("clearEvents", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				{id: 1, Event: Event{tickAddedToEventReader: 1, scheduleSystemsWriter: 0, remove: true}},
				{id: 2, Event: Event{tickAddedToEventReader: 1, scheduleSystemsWriter: 1}},
				{id: 3, Event: Event{tickAddedToEventReader: 1, scheduleSystemsWriter: 2}},
				{id: 4, Event: Event{tickAddedToEventReader: 1, scheduleSystemsWriter: 3}},
			},
		}

		eventReader.clearEvents(2, 0)
		assert.Len(eventReader.events, 3)
		eventReader.clearEvents(2, 1)
		assert.Len(eventReader.events, 3)

		eventReader.clearEvents(2, 3)
		expectedEventIds := []int{2, 4}
		assert.Len(eventReader.events, len(expectedEventIds))

		for i := range eventReader.events {
			assert.Equal(expectedEventIds[i], eventReader.events[i].id)
		}
	})
}
