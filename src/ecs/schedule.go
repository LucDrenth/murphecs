package ecs

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/lucdrenth/murphecs/src/utils"
)

// SystemErrorPackageDepth controls how many path segments are shown in system error messages.
var SystemErrorPackageDepth = 3

// Schedule is a named group of systems that run together.
type Schedule string

// System is a function (or [*systemGroupBuilder]) that can be added to a [Schedule].
type System any

type Scheduler struct {
	systems map[Schedule]*ScheduleSystems
	order   []Schedule
}

func newScheduler() Scheduler {
	return Scheduler{
		systems: map[Schedule]*ScheduleSystems{},
		order:   []Schedule{},
	}
}

func (s *Scheduler) addSchedule(schedule Schedule, scheduleSystemsId ScheduleSystemsId, order ScheduleOrder, isPaused bool) (err error) {
	if _, exists := s.systems[schedule]; exists {
		return ErrScheduleAlreadyExists
	}

	scheduleSystems := &ScheduleSystems{id: scheduleSystemsId}
	if isPaused {
		scheduleSystems.isPaused.Store(true)
	}
	s.systems[schedule] = scheduleSystems

	s.order, err = order.insert(schedule, s.order)
	if err != nil {
		return fmt.Errorf("ScheduleOrder failed to insert: %w", err)
	}

	return nil
}

func (s *Scheduler) addSystem(schedule Schedule, system System, source string, world *World, outerWorlds *map[WorldId]*World, logger Logger, eventStorage *EventStorage) error {
	scheduleSystems, exists := s.systems[schedule]
	if !exists {
		return fmt.Errorf("%w: %s", ErrScheduleNotFound, schedule)
	}

	return scheduleSystems.add(system, source, world, outerWorlds, logger, eventStorage)
}

func (s *Scheduler) getScheduleSystems() ([]*ScheduleSystems, error) {
	if len(s.order) != len(s.systems) {
		return nil, fmt.Errorf("order of length %d does not match schedules of length %d", len(s.order), len(s.systems))
	}

	result := make([]*ScheduleSystems, len(s.order))

	for i, schedule := range s.order {
		scheduleSystems, ok := s.systems[schedule]
		if !ok {
			return nil, fmt.Errorf("schedule %s from schedule order does not exist", schedule)
		}

		result[i] = scheduleSystems
	}

	return result, nil
}

func (s *Scheduler) getScheduleSystemsBySchedules(schedules []Schedule) ([]*ScheduleSystems, error) {
	result := make([]*ScheduleSystems, 0, len(schedules))

	for _, schedule := range schedules {
		scheduleSystems, ok := s.systems[schedule]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrScheduleNotFound, schedule)
		}

		result = append(result, scheduleSystems)
	}

	return result, nil
}

func (s *Scheduler) numberOfSystems() uint {
	result := uint(0)
	for _, scheduleSystems := range s.systems {
		for _, systemGroup := range scheduleSystems.systemGroups {
			result += uint(len(systemGroup.systems))
		}
	}
	return result
}

func (s *Scheduler) numberOfSchedules() uint {
	return uint(len(s.systems))
}

type ScheduleOrder interface {
	insert(Schedule, []Schedule) ([]Schedule, error)
}

var (
	_ ScheduleOrder = &ScheduleLast{}
	_ ScheduleOrder = &ScheduleBefore{}
	_ ScheduleOrder = &ScheduleAfter{}
)

type ScheduleLast struct{}

func (scheduleOrder ScheduleLast) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	return append(schedules, schedule), nil
}

type ScheduleBefore struct {
	Other Schedule
}

func (scheduleOrder ScheduleBefore) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	i := slices.Index(schedules, scheduleOrder.Other)
	if i == -1 {
		return schedules, fmt.Errorf("%w: '%s'", ErrScheduleNotFound, scheduleOrder.Other)
	}

	return slices.Insert(schedules, i, schedule), nil
}

type ScheduleAfter struct {
	Other Schedule
}

func (scheduleOrder ScheduleAfter) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	i := slices.Index(schedules, scheduleOrder.Other)
	if i == -1 {
		return schedules, fmt.Errorf("%w: '%s'", ErrScheduleNotFound, scheduleOrder.Other)
	}

	return slices.Insert(schedules, i+1, schedule), nil
}

type ScheduleSystems struct {
	systemGroups []systemGroup
	id           ScheduleSystemsId

	isPaused               atomic.Bool
	isFirstExecSincePaused bool
}

func (s *ScheduleSystems) Id() ScheduleSystemsId {
	return s.id
}

func (s *ScheduleSystems) Exec(world *World, outerWorlds *map[WorldId]*World, eventStorage *EventStorage, currentTick uint) []error {
	if s.isPaused.Load() {
		if s.isFirstExecSincePaused {
			// The first exec since the schedule is paused needs to call ProcessEvents
			// to clear the event readers.
			// If we don't do this, the event in the readers will never be cleared and
			// can be infinitely read.
			eventStorage.ProcessEvents(s.id, currentTick)
			s.isFirstExecSincePaused = false
		}

		return []error{}
	}

	for _, systemGroup := range s.systemGroups {
		for _, eventWriter := range systemGroup.eventWriters {
			eventWriter.SetScheduleSystemsWriter(s.id)
		}
	}
	defer eventStorage.ProcessEvents(s.id, currentTick)

	world.Mutex.Lock()
	defer world.Mutex.Unlock()

	if outerWorlds != nil {
		for _, outerWorld := range *outerWorlds {
			outerWorld.Mutex.Lock()
			defer outerWorld.Mutex.Unlock()
		}
	}

	err := s.handleSystemParamQueries(world, outerWorlds)
	if err != nil {
		return []error{
			fmt.Errorf("did not execute system set because query failed: %w", err),
		}
	}

	if outerWorlds != nil {
		err = s.updateNonPointerOuterResources(outerWorlds)
		if err != nil {
			return []error{
				fmt.Errorf("did not execute system set because updating outer resources failed: %w", err),
			}
		}
	}

	world.currentScheduleSystemsId = s.id
	defer func() { world.currentScheduleSystemsId = 0 }()

	return s.execSystems()
}

func (s *ScheduleSystems) handleSystemParamQueries(world *World, outerWorlds *map[WorldId]*World) error {
	for _, systemGroup := range s.systemGroups {
		for _, query := range systemGroup.systemParamQueries {
			if query.IsLazy() {
				query.ClearResults()
			} else {
				err := query.Exec(world)
				if err != nil {
					return err
				}
			}
		}

		for _, outerWorldQuery := range systemGroup.systemParamQueriesToOuterWorlds {
			outerWorld := (*outerWorlds)[outerWorldQuery.worldId]
			err := outerWorldQuery.query.Exec(outerWorld)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ScheduleSystems) prepare(outerWorlds *map[WorldId]*World) error {
	for _, systemGroup := range s.systemGroups {
		for _, orp := range systemGroup.outerResources {
			outerWorld, exists := (*outerWorlds)[orp.worldId]
			if !exists {
				return fmt.Errorf("%w: world id %d", ErrTargetWorldNotFound, orp.worldId)
			}

			resource, err := outerWorld.Resources().GetReflectResource(orp.resourceType)
			if err != nil {
				return err
			}

			instance := reflect.New(orp.outerResourceType)
			valueField := instance.Elem().FieldByName("Value")
			if orp.resourceType.Kind() == reflect.Pointer {
				valueField.Set(resource)
			} else {
				valueField.Set(resource.Elem())
			}

			systemGroup.systems[orp.systemIndex].params[orp.paramIndex] = instance.Elem()
		}
	}

	return nil
}

// updateNonPointerOuterResources refreshes the value of non-pointer outer resources
// so that systems see the latest resource state. Pointer outer resources don't need
// this because they reference the resource memory directly.
func (s *ScheduleSystems) updateNonPointerOuterResources(outerWorlds *map[WorldId]*World) error {
	for _, systemGroup := range s.systemGroups {
		for _, orp := range systemGroup.outerResources {
			if orp.resourceType.Kind() == reflect.Pointer {
				continue
			}

			outerWorld := (*outerWorlds)[orp.worldId]
			resource, err := outerWorld.Resources().GetReflectResource(orp.resourceType)
			if err != nil {
				return err
			}

			instance := reflect.New(orp.outerResourceType)
			valueField := instance.Elem().FieldByName("Value")
			valueField.Set(resource.Elem())

			systemGroup.systems[orp.systemIndex].params[orp.paramIndex] = instance.Elem()
		}
	}

	return nil
}

func (s *ScheduleSystems) execSystems() []error {
	errors := []error{}

	for _, systemGroup := range s.systemGroups {
		for i := range systemGroup.systems {
			err := systemGroup.systems[i].exec()
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}

func (s *ScheduleSystems) add(sys System, source string, world *World, outerWorlds *map[WorldId]*World, logger Logger, eventStorage *EventStorage) error {
	systemValue := reflect.ValueOf(sys)
	systemGroupBuilderType2 := reflect.TypeFor[*systemGroupBuilder]()

	if systemValue.Kind() == reflect.Func {
		systemGroup := Systems(sys)
		systemValue = reflect.ValueOf(systemGroup)
	} else if systemValue.Type() != systemGroupBuilderType2 {
		return ErrSystemTypeNotValid
	}

	systemGroupBuilder, ok := reflect.TypeAssert[*systemGroupBuilder](systemValue)
	if !ok {
		panic("failed to type assert *systemGroup")
	}
	if err := systemGroupBuilder.validate(); err != nil {
		return err
	}

	systemGroup, err := systemGroupBuilder.build(source, world, outerWorlds, logger, eventStorage)
	if err != nil {
		return err
	}
	s.systemGroups = append(s.systemGroups, systemGroup)

	return nil
}

type systemEntry struct {
	system     reflect.Value
	params     []reflect.Value
	sourcePath string
}

func (s *systemEntry) exec() error {
	result := s.system.Call(s.params)

	if len(result) == 1 {
		returnedError, isErr := reflect.TypeAssert[error](result[0])
		if isErr {
			return fmt.Errorf("%s: %w", s.sourcePath, returnedError)
		}
	}

	return nil
}

var (
	queryType         = reflect.TypeFor[Query]()
	eventReaderType   = reflect.TypeFor[AnyEventReader]()
	eventWriterType   = reflect.TypeFor[AnyEventWriter]()
	outerResourceType = reflect.TypeFor[AnyOuterResource]()
)

type queryToOuterWorld struct {
	worldId WorldId
	query   Query
}

type outerResourceParam struct {
	systemIndex       int
	paramIndex        int
	worldId           WorldId
	resourceType      reflect.Type
	outerResourceType reflect.Type // the full OuterResource[R, T] struct type
}

type systemGroup struct {
	systems                         []systemEntry
	systemParamQueries              []Query
	systemParamQueriesToOuterWorlds []queryToOuterWorld
	outerResources                  []outerResourceParam
	eventWriters                    []AnyEventWriter
}

type systemGroupBuilder struct {
	systems []System
	// when true, these systems will be run one after another, not in parallel.
	chain bool
}

// Systems creates a group of systems that will be run together.
func Systems(systems ...System) *systemGroupBuilder {
	return &systemGroupBuilder{systems: systems}
}

// Chain makes the systems run sequentially (not in parallel).
func (s *systemGroupBuilder) Chain() *systemGroupBuilder {
	s.chain = true
	return s
}

func (s *systemGroupBuilder) validate() error {
	for _, system := range s.systems {
		systemValue := reflect.ValueOf(system)

		if systemValue.Kind() != reflect.Func {
			return fmt.Errorf("system  %s: %w", systemToDebugString(system), ErrSystemNotAFunction)
		}

		if err := validateSystemReturnTypes(systemValue); err != nil {
			return fmt.Errorf("system %s: %w: %w", systemToDebugString(system), ErrSystemInvalidReturnType, err)
		}
	}

	return nil
}

func (s *systemGroupBuilder) build(source string, world *World, outerWorlds *map[WorldId]*World, logger Logger, eventStorage *EventStorage) (systemGroup, error) {
	systemGroup := systemGroup{}

	for _, sys := range s.systems {
		systemValue := reflect.ValueOf(sys)

		numberOfParams := systemValue.Type().NumIn()
		params := make([]reflect.Value, numberOfParams)

		for i := range numberOfParams {
			parameterType := systemValue.Type().In(i)

			if parameterType.Implements(queryType) {
				query, err := parseQueryParam(parameterType, world, logger, outerWorlds)
				if err != nil {
					return systemGroup, fmt.Errorf("%s: parameter %s: %w: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ErrSystemParamQueryNotValid, err)
				}

				if query.TargetWorld() != nil {
					systemGroup.systemParamQueriesToOuterWorlds = append(systemGroup.systemParamQueriesToOuterWorlds, queryToOuterWorld{
						worldId: *query.TargetWorld(),
						query:   query,
					})
				} else {
					systemGroup.systemParamQueries = append(systemGroup.systemParamQueries, query)
				}

				params[i] = reflect.ValueOf(query)
			} else if parameterType == reflect.TypeFor[*World]() {
				params[i] = reflect.ValueOf(world)
			} else if parameterType == reflect.TypeFor[World]() {
				// World may not be used by-value because:
				//	1. it is a potentially big object and copying it could give bad performance
				//	2. it is probably unintended and would cause unexpected behavior
				return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ErrSystemParamWorldNotAPointer)
			} else if parameterType.Implements(eventReaderType) {
				eventReader, ok := reflect.TypeAssert[AnyEventReader](reflect.New(parameterType.Elem()))
				if !ok {
					panic("failed to type assert AnyEventReader")
				}
				params[i] = *eventStorage.GetReader(eventReader)
			} else if parameterType.Implements(eventWriterType) {
				eventWriter, ok := reflect.TypeAssert[AnyEventWriter](reflect.New(parameterType.Elem()))
				if !ok {
					panic("failed to type assert AnyEventWriter")
				}

				reflectedEventWriter := *eventStorage.GetWriter(eventWriter)
				eventWriterParam, ok := reflect.TypeAssert[AnyEventWriter](reflectedEventWriter)
				if !ok {
					panic("failed to type assert AnyEventWriter")
				}
				systemGroup.eventWriters = append(systemGroup.eventWriters, eventWriterParam)
				params[i] = reflectedEventWriter
			} else if parameterType.Implements(outerResourceType) {
				return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ErrSystemParamOuterResourceIsAPointer)
			} else if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(outerResourceType) {
				instance := reflect.New(parameterType)
				outerRes := instance.Interface().(AnyOuterResource)
				worldId, resType := outerRes.OuterResourceInfo()

				systemGroup.outerResources = append(systemGroup.outerResources, outerResourceParam{
					systemIndex:       len(systemGroup.systems),
					paramIndex:        i,
					worldId:           *worldId,
					resourceType:      resType,
					outerResourceType: parameterType,
				})

				params[i] = reflect.Zero(parameterType)
			} else {
				// check if its a resource
				resource, err := world.Resources().GetReflectResource(parameterType)
				if err != nil {
					// err just means its not a resource, no need to return this specific error.

					err = handleInvalidSystemParam(parameterType)
					return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), err)
				}

				if parameterType.Kind() == reflect.Pointer {
					params[i] = resource
				} else {
					params[i] = resource.Elem()
				}
			}
		}

		entry := systemEntry{
			system:     systemValue,
			params:     params,
			sourcePath: source,
		}
		systemGroup.systems = append(systemGroup.systems, entry)
	}

	return systemGroup, nil
}

func handleInvalidSystemParam(parameterType reflect.Type) error {
	if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(reflect.TypeFor[Query]()) {
		return ErrSystemParamQueryNotAPointer
	}

	if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(reflect.TypeFor[AnyEventReader]()) {
		return fmt.Errorf("EventReader: %w", ErrSystemParamEventReaderNotAPointer)
	}

	if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(reflect.TypeFor[AnyEventWriter]()) {
		return fmt.Errorf("EventWriter: %w", ErrSystemParamEventWriterNotAPointer)
	}

	return ErrSystemParamNotValid
}

func parseQueryParam(parameterType reflect.Type, world *World, logger Logger, outerWorlds *map[WorldId]*World) (Query, error) {
	if parameterType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("can not be an interface")
	}

	query, ok := reflect.TypeAssert[Query](reflect.New(parameterType.Elem()))
	if !ok {
		return nil, fmt.Errorf("failed to cast param to query")
	}

	err := query.Prepare(world, outerWorlds)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query param: %w", err)
	}

	if query.TargetWorld() != nil && query.IsLazy() {
		// We have this limitation because there would be no way to Execute such a query from inside
		// the system param because there is no way to use another world as a system param.
		//
		// We could easily implement this by having some OuterWorld[ecs.TargetWorld] struct that we can
		// use as a system parameter. For now there is no valid use case for it.
		return nil, fmt.Errorf("query cannot target an outer world if its lazy")
	}

	warning := query.Validate()
	if warning != nil {
		logger.Warn("query %s is not optimized: %v", parameterType.String(), warning)
	}

	return query, nil
}

func validateSystemReturnTypes(systemValue reflect.Value) error {
	numberOfSystemReturnValues := systemValue.Type().NumOut()

	if numberOfSystemReturnValues == 0 {
		return nil
	}

	if numberOfSystemReturnValues == 1 {
		returnType := systemValue.Type().Out(0)

		if returnType == reflect.TypeFor[error]() {
			return nil
		}

		return fmt.Errorf("return type is %s but must be either error or nothing", returnType.String())
	}

	return fmt.Errorf("has %d return values but must have either 1 (error) or 0", numberOfSystemReturnValues)
}

// systemToDebugString returns a reflection string of the system but with shortened paths.
func systemToDebugString(system System) string {
	result := reflect.TypeOf(system).String()
	result = applyDebugTypeReplacements(result)
	result = strings.ReplaceAll(result, ",", ", ")
	return result
}

// systemParameterDebugString returns a reflection string of the system parameter but with shortened paths.
func systemParameterDebugString(system System, index int) string {
	systemType := reflect.TypeOf(system)
	parameterType := systemType.In(index)
	result := parameterType.String()
	result = applyDebugTypeReplacements(result)
	return result
}

func applyDebugTypeReplacements(s string) string {
	result := s
	for k, v := range DebugTypeReplacements {
		result = strings.ReplaceAll(result, k, v)
	}
	return result
}

var DebugTypeReplacements = map[string]string{
	"github.com/lucdrenth/murphecs/src/": "murphecs/",
}

func callerSource(skip int) string {
	return utils.Caller(skip+1, SystemErrorPackageDepth)
}
