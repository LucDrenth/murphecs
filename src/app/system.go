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

type SystemSetId int

type SystemSet struct {
	systems                         []systemEntry
	systemParamQueries              []ecs.Query
	systemParamQueriesToOuterWorlds []queryToOuterWorld
	id                              SystemSetId
	eventWriters                    []iEventWriter
}

type queryToOuterWorld struct {
	worldId ecs.WorldId
	query   ecs.Query
}

func (s *SystemSet) Exec(world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, eventStorage *EventStorage, currentTick uint) []error {
	for _, eventWriters := range s.eventWriters {
		eventWriters.setSystemSetWriter(s.id)
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

func (s *SystemSet) handleSystemParamQueries(world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World) error {
	for _, query := range s.systemParamQueries {
		if query.IsLazy() {
			query.ClearResults()
		} else {
			err := query.Exec(world)
			if err != nil {
				return err
			}
		}
	}

	for _, outerWorldQuery := range s.systemParamQueriesToOuterWorlds {
		outerWorld := (*outerWorlds)[outerWorldQuery.worldId]
		err := outerWorldQuery.query.Exec(outerWorld)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SystemSet) execSystems() []error {
	errors := []error{}

	for i := range s.systems {
		err := s.systems[i].exec()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *SystemSet) add(sys System, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, resources *resourceStorage, eventStorage *EventStorage) error {
	systemValue := reflect.ValueOf(sys)
	queryType := reflect.TypeOf((*ecs.Query)(nil)).Elem()
	eventReaderType := reflect.TypeOf((*iEventReader)(nil)).Elem()
	eventWriterType := reflect.TypeOf((*iEventWriter)(nil)).Elem()

	if err := validateSystem(systemValue); err != nil {
		return fmt.Errorf("system is not valid: %w", err)
	}

	numberOfParams := systemValue.Type().NumIn()
	params := make([]reflect.Value, numberOfParams)

	for i := range numberOfParams {
		parameterType := systemValue.Type().In(i)

		if parameterType.Implements(queryType) {
			query, err := parseQueryParam(parameterType, world, logger, outerWorlds)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrSystemParamQueryNotValid, err)
			}

			if query.TargetWorld() != nil {
				s.systemParamQueriesToOuterWorlds = append(s.systemParamQueriesToOuterWorlds, queryToOuterWorld{
					worldId: *query.TargetWorld(),
					query:   query,
				})
			} else {
				s.systemParamQueries = append(s.systemParamQueries, query)
			}

			params[i] = reflect.ValueOf(query)
		} else if parameterType == reflect.TypeFor[*ecs.World]() {
			params[i] = reflect.ValueOf(world)
		} else if parameterType == reflect.TypeFor[ecs.World]() {
			// ecs.World may not be used by-value because:
			//	1. it is a potentially big object and copying it could give bad performance
			//	2. it is probably unintended and would cause unexpected behavior
			return fmt.Errorf("system parameter %d: %w", i+1, ErrSystemParamWorldNotAPointer)
		} else if parameterType.Implements(eventReaderType) {
			eventReader, ok := reflect.TypeAssert[iEventReader](reflect.New(parameterType.Elem()))
			if !ok {
				panic("failed to type assert iEventReader")
			}
			params[i] = *eventStorage.getReader(eventReader)
		} else if parameterType.Implements(eventWriterType) {
			eventWriter, ok := reflect.TypeAssert[iEventWriter](reflect.New(parameterType.Elem()))
			if !ok {
				panic("failed to type assert iEventWriter")
			}

			reflectedEventWriter := *eventStorage.getWriter(eventWriter)
			eventWriterParam, ok := reflect.TypeAssert[iEventWriter](reflectedEventWriter)
			if !ok {
				panic("failed to type assert iEventWriter")
			}
			s.eventWriters = append(s.eventWriters, eventWriterParam)
			params[i] = reflectedEventWriter
		} else {
			// check if its a resource
			resource, err := resources.getReflectResource(parameterType)
			if err != nil {
				// err just means its not a resource, no need to return this specific error.

				err = handleInvalidSystemParam(parameterType)
				return fmt.Errorf("system parameter %d: %w", i+1, err)
			}

			if parameterType.Kind() == reflect.Pointer {
				params[i] = resource
			} else {
				params[i] = resource.Elem()
			}
		}
	}

	entry := systemEntry{
		system: systemValue,
		params: params,
	}
	s.systems = append(s.systems, entry)
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

func validateSystem(sys reflect.Value) error {
	if sys.Kind() != reflect.Func {
		return ErrSystemNotAFunction
	}

	if err := validateSystemReturnTypes(sys); err != nil {
		return fmt.Errorf("%w: %w", ErrSystemInvalidReturnType, err)
	}

	return nil
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
	result = strings.ReplaceAll(result, "github.com/lucdrenth/murphecs/src/", "murphecs/")
	return result
}
