package app

import (
	"reflect"
)

type EventStorage struct {
	eventReaders map[eventId]*reflect.Value
	eventWriters map[eventId]*reflect.Value
}

func newEventStorage() EventStorage {
	return EventStorage{
		eventReaders: map[eventId]*reflect.Value{},
		eventWriters: map[eventId]*reflect.Value{},
	}
}

// getReader gets a reader or creates and stores a new one
func (s *EventStorage) getReader(reader iEventReader) *reflect.Value {
	id := reader.readerEventId()

	result, exists := s.eventReaders[id]
	if !exists {
		concreteReader := reflect.ValueOf(reader)
		result = &concreteReader
		s.eventReaders[id] = &concreteReader
	}

	return result
}

// getWriter gets a writer or creates and stores a new one
func (s *EventStorage) getWriter(writer iEventWriter) *reflect.Value {
	id := writer.writerEventId()

	result, exists := s.eventWriters[id]
	if !exists {
		concreteWriter := reflect.ValueOf(writer)
		result = &concreteWriter
		s.eventWriters[id] = &concreteWriter
	}

	return result
}

// ProcessEvents moves events from writers to readers and cleans up reader events.
func (s *EventStorage) ProcessEvents(systemSetId SystemSetId, currentTick uint) {
	for eventId, reflectWriter := range s.eventWriters {
		writer := reflectWriter.Interface().(iEventWriter)
		writerEvents := writer.extractEvents(currentTick)

		readerEntry, ok := s.eventReaders[eventId]
		if !ok {
			// This event does not have any readers
			continue
		}

		reader := readerEntry.Interface().(iEventReader)
		reader.clearEvents(systemSetId, currentTick)

		for _, reflectEvent := range writerEvents {
			reader.addEvent(reflectEvent)
		}
	}
}

type IEvent interface {
	shouldRemove() bool
	getSystemSetWriter() SystemSetId
	setSystemSetWriter(SystemSetId)
	setTickAddedToEventReader(uint)
	getTickAddedToEventReader() uint
}

type eventId reflect.Type

type Event struct {
	// If true, this event should not be read anymore and should be removed
	remove                 bool
	systemSetWriter        SystemSetId // the SystemSetId of the SystemSet during which this event was written to an EventWriter
	tickAddedToEventReader uint        // the tick number during which this event was added to an EventReader
}

func (e *Event) shouldRemove() bool {
	return e.remove
}

func (e *Event) setSystemSetWriter(id SystemSetId) {
	e.systemSetWriter = id
}

func (e *Event) getSystemSetWriter() SystemSetId {
	return e.systemSetWriter
}

func (e *Event) setTickAddedToEventReader(currentTick uint) {
	e.tickAddedToEventReader = currentTick
}

func (e *Event) getTickAddedToEventReader() uint {
	return e.tickAddedToEventReader
}

// Remove marks the event for cleanup and prevents other systems from reading this event.
func (e *Event) Remove() {
	e.remove = true
}

type EventWriter[E IEvent] struct {
	events      []E
	SystemSetId // this will be updated every systemSet run that this event writer is in
}

// Write adds event if it is not nil.
//
// It will be available for reading from the start of next schedule and until the end the next
// iteration of the current schedule.
//
// For example, see the following scenario with 3 schedules: "pre-update", "update" and "post-update":
//
//   - pre-update 	system 1: 	not readable
//
//   - update 		system 1: 	[write occurs] not readable
//
//   - update 		system 2: 	not readable
//
//   - post-update 	system 1:	readable
//
//     ==== next loop ===
//
//   - pre-update 	system 1: 	readable
//
//   - update 		system 1: 	readable
//
//   - update 		system 2: 	readable
//
//   - post-update 	system 1:	[event cleared] not readable
func (writer *EventWriter[E]) Write(event E) {
	event.setSystemSetWriter(writer.SystemSetId)
	writer.events = append(writer.events, event)
}

func (writer *EventWriter[E]) writerEventId() eventId {
	return reflect.TypeFor[E]()
}

// extractEvents returns the events as reflect.Value's and removes them from the event writer
func (writer *EventWriter[E]) extractEvents(tick uint) []reflect.Value {
	result := make([]reflect.Value, len(writer.events))
	for i, event := range writer.events {
		event.setTickAddedToEventReader(tick)
		result[i] = reflect.ValueOf(event)
	}

	writer.events = []E{}

	return result
}

func (writer *EventWriter[E]) setSystemSetWriter(id SystemSetId) {
	writer.SystemSetId = id
}

type iEventWriter interface {
	writerEventId() eventId
	extractEvents(tick uint) []reflect.Value
	setSystemSetWriter(SystemSetId)
}

var _ iEventReader = &EventReader[*Event]{}

type EventReader[E IEvent] struct {
	events []E
}

// Read ranges over al events that are not yet marked as removed
func (reader *EventReader[E]) Read(yield func(E) bool) {
	for _, event := range reader.events {
		if event.shouldRemove() {
			continue
		}

		if !yield(event) {
			return
		}
	}
}

// First returns the first written event.
// Returns (_, false) if there are no elements.
func (reader *EventReader[E]) First() (E, bool) {
	for _, event := range reader.events {
		if !event.shouldRemove() {
			return event, true
		}
	}

	var result E
	return result, false
}

// Last returns the last written event.
// Returns (_, false) if there are no elements.
func (reader *EventReader[E]) Last() (E, bool) {
	for i := len(reader.events) - 1; i >= 0; i-- {
		if !reader.events[i].shouldRemove() {
			return reader.events[i], true
		}
	}

	var result E
	return result, false
}

// Len returns the number of events in the reader
func (reader *EventReader[E]) Len() int {
	result := 0
	for _, event := range reader.events {
		if !event.shouldRemove() {
			result++
		}
	}

	return result
}

func (writer *EventReader[E]) addEvent(event reflect.Value) {
	writer.events = append(writer.events, event.Interface().(E))
}

// Len returns wether there are any events in the reader
func (reader *EventReader[E]) IsEmpty() bool {
	return reader.Len() == 0
}

func (reader *EventReader[E]) readerEventId() eventId {
	return reflect.TypeFor[E]()
}

// clearEvents removes all events that satisfy one of the following:
//   - marked to be removed
//   - written by [SystemSet] with given [SystemSetId] AND added to reader at least 1 tick back
func (reader *EventReader[E]) clearEvents(systemSetId SystemSetId, currentTick uint) {
	newEvents := []E{}

	for _, event := range reader.events {
		if event.shouldRemove() {
			continue
		}

		if event.getSystemSetWriter() == systemSetId && currentTick > event.getTickAddedToEventReader() {
			continue
		}

		newEvents = append(newEvents, event)
	}

	reader.events = newEvents
}

type iEventReader interface {
	readerEventId() eventId
	addEvent(event reflect.Value)
	clearEvents(systemSetId SystemSetId, currentTick uint)
}

var _ iEventReader = &EventReader[*Event]{}
