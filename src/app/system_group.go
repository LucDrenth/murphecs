package app

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type queryToOuterWorld struct {
	worldId ecs.WorldId
	query   ecs.Query
}

type systemGroup struct {
	systems                         []systemEntry
	systemParamQueries              []ecs.Query
	systemParamQueriesToOuterWorlds []queryToOuterWorld
	eventWriters                    []iEventWriter
}

type systemGroupBuilder struct {
	systems []System
}

func Systems(systems ...System) *systemGroupBuilder {
	return &systemGroupBuilder{systems: systems}
}

func (s *systemGroupBuilder) validate() error {
	for i, system := range s.systems {
		systemValue := reflect.ValueOf(system)

		if systemValue.Kind() != reflect.Func {
			return fmt.Errorf("system at index %d: %w", i, ErrSystemNotAFunction)
		}

		if err := validateSystemReturnTypes(systemValue); err != nil {
			return fmt.Errorf("system at index %d: %w: %w", i, ErrSystemInvalidReturnType, err)
		}
	}

	return nil
}

func (s *systemGroupBuilder) build(world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, resources *resourceStorage, eventStorage *EventStorage) (systemGroup, error) {
	systemGroup := systemGroup{}

	queryType := reflect.TypeOf((*ecs.Query)(nil)).Elem()
	eventReaderType := reflect.TypeOf((*iEventReader)(nil)).Elem()
	eventWriterType := reflect.TypeOf((*iEventWriter)(nil)).Elem()

	for _, sys := range s.systems {
		systemValue := reflect.ValueOf(sys)

		numberOfParams := systemValue.Type().NumIn()
		params := make([]reflect.Value, numberOfParams)

		for i := range numberOfParams {
			parameterType := systemValue.Type().In(i)

			if parameterType.Implements(queryType) {
				query, err := parseQueryParam(parameterType, world, logger, outerWorlds)
				if err != nil {
					return systemGroup, fmt.Errorf("%w: %w", ErrSystemParamQueryNotValid, err)
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
			} else if parameterType == reflect.TypeFor[*ecs.World]() {
				params[i] = reflect.ValueOf(world)
			} else if parameterType == reflect.TypeFor[ecs.World]() {
				// ecs.World may not be used by-value because:
				//	1. it is a potentially big object and copying it could give bad performance
				//	2. it is probably unintended and would cause unexpected behavior
				return systemGroup, fmt.Errorf("system parameter %d: %w", i+1, ErrSystemParamWorldNotAPointer)
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
				systemGroup.eventWriters = append(systemGroup.eventWriters, eventWriterParam)
				params[i] = reflectedEventWriter
			} else {
				// check if its a resource
				resource, err := resources.getReflectResource(parameterType)
				if err != nil {
					// err just means its not a resource, no need to return this specific error.

					err = handleInvalidSystemParam(parameterType)
					return systemGroup, fmt.Errorf("system parameter %d: %w", i+1, err)
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
		systemGroup.systems = append(systemGroup.systems, entry)
	}

	return systemGroup, nil
}
