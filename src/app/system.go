package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type System any

type systemEntry struct {
	system reflect.Value
	params []reflect.Value
}

func (s *systemEntry) exec() error {
	result := s.system.Call(s.params)

	if len(result) == 1 {
		returnedError, isErr := reflect.TypeAssert[error](result[0])
		if isErr {
			return returnedError
		}
	}

	return nil
}

type ScheduleSystemsId int

type ScheduleSystems struct {
	systemGroups []systemGroup
	id           ScheduleSystemsId
}

func (s *ScheduleSystems) Exec(world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, eventStorage *EventStorage, currentTick uint) []error {
	for _, systemGroup := range s.systemGroups {
		for _, eventWriters := range systemGroup.eventWriters {
			eventWriters.setScheduleSystemsWriter(s.id)
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

	return s.execSystems()
}

func (s *ScheduleSystems) handleSystemParamQueries(world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World) error {
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

func (s *ScheduleSystems) add(sys System, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, resources *resourceStorage, eventStorage *EventStorage) error {
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

	systemGroup, err := systemGroupBuilder.build(world, outerWorlds, logger, resources, eventStorage)
	if err != nil {
		return err
	}
	s.systemGroups = append(s.systemGroups, systemGroup)

	return nil
}

func handleInvalidSystemParam(parameterType reflect.Type) error {
	if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(reflect.TypeFor[ecs.Query]()) {
		return ErrSystemParamQueryNotAPointer
	}

	// TODO do a proper check for the parameter being an EventReader
	if strings.HasPrefix(parameterType.Name(), "EventReader") {
		return fmt.Errorf("EventReader: %w", ErrSystemParamEventReaderNotAPointer)
	}

	// TODO do a proper check for the parameter being an EventWriter
	if strings.HasPrefix(parameterType.Name(), "EventWriter") {
		return fmt.Errorf("EventWriter: %w", ErrSystemParamEventWriterNotAPointer)
	}

	return ErrSystemParamNotValid
}

func parseQueryParam(parameterType reflect.Type, world *ecs.World, logger Logger, outerWorlds *map[ecs.WorldId]*ecs.World) (ecs.Query, error) {
	if parameterType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("can not be an interface")
	}

	query, ok := reflect.TypeAssert[ecs.Query](reflect.New(parameterType.Elem()))
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

// systemParameterDebugString returns a reflection string of the system but with shortened paths.
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
