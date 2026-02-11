package app

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murphecs/src/ecs"
)

var (
	queryType         = reflect.TypeFor[ecs.Query]()
	eventReaderType   = reflect.TypeFor[anyEventReader]()
	eventWriterType   = reflect.TypeFor[anyEventWriter]()
	outerResourceType = reflect.TypeFor[ecs.AnyOuterResource]()
)

type queryToOuterWorld struct {
	worldId ecs.WorldId
	query   ecs.Query
}

type systemGroup struct {
	systems                         []systemEntry
	systemParamQueries              []ecs.Query
	systemParamQueriesToOuterWorlds []queryToOuterWorld
	eventWriters                    []anyEventWriter
}

type systemGroupBuilder struct {
	systems []System
	// when true, these systems will be ran one after another, not running them in parallel.
	chain bool
}

func Systems(systems ...System) *systemGroupBuilder {
	return &systemGroupBuilder{systems: systems}
}

// Chain makes the systems run sequential (not in parallel)
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

func (s *systemGroupBuilder) build(source string, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, eventStorage *EventStorage) (systemGroup, error) {
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
			} else if parameterType == reflect.TypeFor[*ecs.World]() {
				params[i] = reflect.ValueOf(world)
			} else if parameterType == reflect.TypeFor[ecs.World]() {
				// ecs.World may not be used by-value because:
				//	1. it is a potentially big object and copying it could give bad performance
				//	2. it is probably unintended and would cause unexpected behavior
				return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ErrSystemParamWorldNotAPointer)
			} else if parameterType.Implements(eventReaderType) {
				eventReader, ok := reflect.TypeAssert[anyEventReader](reflect.New(parameterType.Elem()))
				if !ok {
					panic("failed to type assert iEventReader")
				}
				params[i] = *eventStorage.getReader(eventReader)
			} else if parameterType.Implements(eventWriterType) {
				eventWriter, ok := reflect.TypeAssert[anyEventWriter](reflect.New(parameterType.Elem()))
				if !ok {
					panic("failed to type assert iEventWriter")
				}

				reflectedEventWriter := *eventStorage.getWriter(eventWriter)
				eventWriterParam, ok := reflect.TypeAssert[anyEventWriter](reflectedEventWriter)
				if !ok {
					panic("failed to type assert iEventWriter")
				}
				systemGroup.eventWriters = append(systemGroup.eventWriters, eventWriterParam)
				params[i] = reflectedEventWriter
			} else if parameterType.Implements(outerResourceType) {
				return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ErrSystemParamOuterResourceIsAPointer)
			} else if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(outerResourceType) {
				instance := reflect.New(parameterType)
				outerRes := instance.Interface().(ecs.AnyOuterResource)
				worldId, resType := outerRes.OuterResourceInfo()

				outerWorld, exists := (*outerWorlds)[*worldId]
				if !exists {
					return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), ecs.ErrTargetWorldNotFound)
				}

				resource, err := outerWorld.Resources().GetReflectResource(resType)
				if err != nil {
					return systemGroup, fmt.Errorf("%s: parameter %s: %w", systemToDebugString(sys), systemParameterDebugString(sys, i), err)
				}

				valueField := instance.Elem().FieldByName("Value")
				if resType.Kind() == reflect.Pointer {
					valueField.Set(resource)
				} else {
					valueField.Set(resource.Elem())
				}

				params[i] = instance.Elem()
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
