package app

import (
	"reflect"
	"time"
)

type EventStorage struct {
	eventReaders map[eventId]*reflect.Value
	eventWriters map[eventId]*reflect.Value
}

func NewEventStorage() EventStorage {
	return EventStorage{
		eventReaders: map[eventId]*reflect.Value{},
		eventWriters: map[eventId]*reflect.Value{},
	}
}

// getReader gets a reader or creates and stores a new one
func (s *EventStorage) getReader(reader anyEventReader) *reflect.Value {
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
func (s *EventStorage) getWriter(writer anyEventWriter) *reflect.Value {
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
func (s *EventStorage) ProcessEvents(scheduleSystemsId ScheduleSystemsId, currentTick uint) {
	for eventId, reflectWriter := range s.eventWriters {
		writer, ok := reflect.TypeAssert[anyEventWriter](*reflectWriter)
		if !ok {
			panic("failed to type assert iEventWriter")
		}
		writerEvents := writer.extractEvents(currentTick)

		readerEntry, ok := s.eventReaders[eventId]
		if !ok {
			// This event does not have any readers
			continue
		}

		reader, ok := reflect.TypeAssert[anyEventReader](*readerEntry)
		if !ok {
			panic("failed to type assert iEventReader")
		}
		reader.clearEvents(scheduleSystemsId, currentTick)

		for _, reflectEvent := range writerEvents {
			reader.addEvent(reflectEvent)
		}
	}
}

type IEvent interface {
	shouldRemove() bool
	getScheduleSystemsWriter() ScheduleSystemsId
	setScheduleSystemsWriter(ScheduleSystemsId)
	setTimeWritten(time.Time)
	setTickAddedToEventReader(uint)
	getTickAddedToEventReader() uint
}

type eventId reflect.Type

type Event struct {
	// If true, this event should not be read anymore and should be removed
	remove                 bool
	scheduleSystemsWriter  ScheduleSystemsId // the [ScheduleSystemsId] of the [ScheduleSystems] during which this event was written to an EventWriter
	tickAddedToEventReader uint              // the tick number during which this event was added to an EventReader
	timeWritten            time.Time         // the time at which the event has been written to the EventWriter
}

func (e *Event) shouldRemove() bool {
	return e.remove
}

func (e *Event) setScheduleSystemsWriter(id ScheduleSystemsId) {
	e.scheduleSystemsWriter = id
}

func (e *Event) getScheduleSystemsWriter() ScheduleSystemsId {
	return e.scheduleSystemsWriter
}

func (e *Event) setTickAddedToEventReader(currentTick uint) {
	e.tickAddedToEventReader = currentTick
}

func (e *Event) getTickAddedToEventReader() uint {
	return e.tickAddedToEventReader
}

func (e *Event) setTimeWritten(t time.Time) {
	e.timeWritten = t
}

// Remove marks the event for cleanup and prevents other systems from reading this event.
func (e *Event) Remove() {
	e.remove = true
}

// TimeWritten returns the time during which the event was written to the EventWriter
func (e *Event) TimeWritten() time.Time {
	return e.timeWritten
}

type EventWriter[E IEvent] struct {
	events            []E
	ScheduleSystemsId // this will be updated every [ScheduleSystems] run that this event writer is in
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
	event.setTimeWritten(time.Now())
	event.setScheduleSystemsWriter(writer.ScheduleSystemsId)
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

func (writer *EventWriter[E]) setScheduleSystemsWriter(id ScheduleSystemsId) {
	writer.ScheduleSystemsId = id
}

type anyEventWriter interface {
	writerEventId() eventId
	extractEvents(tick uint) []reflect.Value
	setScheduleSystemsWriter(ScheduleSystemsId)
}

var _ anyEventReader = &EventReader[*Event]{}

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
	element, ok := reflect.TypeAssert[E](event)
	if !ok {
		panic("failed to type assert event")
	}
	writer.events = append(writer.events, element)
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
//   - written by [ScheduleSystems] with given [ScheduleSystemsId] AND added to reader at least 1 tick back
func (reader *EventReader[E]) clearEvents(scheduleSystemsId ScheduleSystemsId, currentTick uint) {
	newEvents := []E{}

	for _, event := range reader.events {
		if event.shouldRemove() {
			continue
		}

		if event.getScheduleSystemsWriter() == scheduleSystemsId && currentTick > event.getTickAddedToEventReader() {
			continue
		}

		newEvents = append(newEvents, event)
	}

	reader.events = newEvents
}

type anyEventReader interface {
	readerEventId() eventId
	addEvent(event reflect.Value)
	clearEvents(scheduleSystemsId ScheduleSystemsId, currentTick uint)
}

var _ anyEventReader = &EventReader[*Event]{}
