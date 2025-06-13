package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(element.id, 1)
	})

	t.Run("First returns the first element that is not marked as removed", func(t *testing.T) {
		assert := assert.New(t)

		eventReader := EventReader[*testEvent]{
			events: []*testEvent{
				&testEvent{id: 1},
				&testEvent{id: 2},
			},
		}
		eventReader.events[0].Remove()

		element, found := eventReader.First()
		assert.True(found)
		assert.Equal(element.id, 2)
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
		assert.Equal(element.id, 2)
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
		assert.Equal(element.id, 1)
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

		assert.Equal(eventReader.Len(), 2)
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
}
