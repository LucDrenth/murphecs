package ecs

import (
	"fmt"
	"reflect"
)

type observerOuterResourceParam struct {
	paramIndex        int
	worldId           WorldId
	resourceType      reflect.Type // type of the actual resource (e.g. *R or R)
	outerResourceType reflect.Type // type of the OuterResource[R, W] struct
}

type observerEntry struct {
	systemEntry
	observerParamIndex int                        // -1 if observer is not used as a param
	queries            []Query                    // queries to execute before running
	eventWriters       []AnyEventWriter           // event writers to process after running
	outerResources     []observerOuterResourceParam // outer resource params to update before running
}

func (e *observerEntry) execWithObserver(world *World, observerValue reflect.Value) error {
	for _, orp := range e.outerResources {
		if orp.resourceType.Kind() == reflect.Pointer {
			continue // pointer outer resources reference memory directly; no refresh needed
		}
		outerWorld := (*world.OuterWorlds())[orp.worldId]
		resource, err := outerWorld.Resources().GetReflectResource(orp.resourceType)
		if err != nil {
			return err
		}
		instance := reflect.New(orp.outerResourceType)
		instance.Elem().FieldByName("Value").Set(resource.Elem())
		e.params[orp.paramIndex] = instance.Elem()
	}

	for _, q := range e.queries {
		if q.IsLazy() {
			q.ClearResults()
		} else {
			if err := q.Exec(world); err != nil {
				return err
			}
		}
	}

	for _, ew := range e.eventWriters {
		ew.SetScheduleSystemsWriter(world.currentScheduleSystemsId)
	}

	if e.observerParamIndex >= 0 {
		e.params[e.observerParamIndex] = observerValue
	}

	return e.exec()
}

type spawnDespawnObserverEntry struct {
	observerEntry
	observerValue    reflect.Value // pre-allocated OnSpawn/OnDespawn value
	entityFieldIndex int           // index of Entity field in the observer struct
}

type observerRegistry struct {
	observers        map[observerId][]observerEntry
	spawnObservers   map[ComponentId][]spawnDespawnObserverEntry
	despawnObservers map[ComponentId][]spawnDespawnObserverEntry
}

func newObserverRegistry() observerRegistry {
	return observerRegistry{
		observers:        map[observerId][]observerEntry{},
		spawnObservers:   map[ComponentId][]spawnDespawnObserverEntry{},
		despawnObservers: map[ComponentId][]spawnDespawnObserverEntry{},
	}
}

func (registry *observerRegistry) triggerDespawnObservers(world *World, componentIds []ComponentId, entity EntityId) {
	for _, componentId := range componentIds {
		entries, exists := registry.despawnObservers[componentId]
		if !exists {
			continue
		}

		for i := range entries {
			entry := &entries[i]
			entry.observerValue.Field(entry.entityFieldIndex).Set(reflect.ValueOf(entity))
			err := entry.execWithObserver(world, entry.observerValue)
			if err != nil {
				world.logger.Error("exec observer failed: %v", err)
			}
		}
	}
}

func (registry *observerRegistry) triggerSpawnObservers(world *World, componentIds []ComponentId, entity EntityId) {
	for _, componentId := range componentIds {
		entries, exists := registry.spawnObservers[componentId]
		if !exists {
			continue
		}

		for i := range entries {
			entry := &entries[i]
			entry.observerValue.Field(entry.entityFieldIndex).Set(reflect.ValueOf(entity))
			err := entry.execWithObserver(world, entry.observerValue)
			if err != nil {
				world.logger.Error("exec observer failed: %v", err)
			}
		}
	}
}

func registerObserver[O AnyObserver](registry *observerRegistry, world *World, action System, source string) error {
	if action == nil {
		return nil
	}

	var zeroObserver O
	observerReflectType := reflect.TypeFor[O]()
	if observerReflectType.Kind() == reflect.Pointer {
		panic("can not register observer pointer")
	}

	entry, err := buildObserverEntry[O](action, world, source)
	if err != nil {
		return err
	}

	switch zeroObserver.getObserverType() {
	case customObserver:
		registry.observers[observerReflectType] = append(registry.observers[observerReflectType], entry)

	case spawnObserver:
		componentId := zeroObserver.componentId(world)
		obsType := reflect.TypeOf(zeroObserver)
		obsValue := reflect.New(obsType).Elem()
		entityField, _ := obsType.FieldByName("Entity")
		entityFieldIndex := entityField.Index[0]
		registry.spawnObservers[componentId] = append(registry.spawnObservers[componentId], spawnDespawnObserverEntry{
			observerEntry:    entry,
			observerValue:    obsValue,
			entityFieldIndex: entityFieldIndex,
		})

	case despawnObserver:
		componentId := zeroObserver.componentId(world)
		obsType := reflect.TypeOf(zeroObserver)
		obsValue := reflect.New(obsType).Elem()
		entityField, _ := obsType.FieldByName("Entity")
		entityFieldIndex := entityField.Index[0]
		registry.despawnObservers[componentId] = append(registry.despawnObservers[componentId], spawnDespawnObserverEntry{
			observerEntry:    entry,
			observerValue:    obsValue,
			entityFieldIndex: entityFieldIndex,
		})

	default:
		panic("unhandled observer type")
	}

	return nil
}

func triggerObserver[O AnyObserver](world *World, registry *observerRegistry, observed O) {
	observerId := reflect.TypeFor[O]()
	entries, exists := registry.observers[observerId]
	if !exists {
		return
	}

	observedValue := reflect.ValueOf(observed)
	for i := range entries {
		err := entries[i].execWithObserver(world, observedValue)
		if err != nil {
			world.logger.Error("exec observer failed: %v", err)
		}
	}
}

func buildObserverEntry[O AnyObserver](action System, world *World, source string) (observerEntry, error) {
	actionValue := reflect.ValueOf(action)
	if actionValue.Kind() != reflect.Func {
		return observerEntry{}, ErrSystemNotAFunction
	}

	if err := validateSystemReturnTypes(actionValue); err != nil {
		return observerEntry{}, fmt.Errorf("%w: %w", ErrSystemInvalidReturnType, err)
	}

	observerParamType := reflect.TypeFor[O]()
	observerParamIdx := -1
	numberOfParams := actionValue.Type().NumIn()
	params := make([]reflect.Value, numberOfParams)
	var queries []Query
	var eventWriters []AnyEventWriter
	var outerResources []observerOuterResourceParam

	for i := range numberOfParams {
		paramType := actionValue.Type().In(i)

		if paramType == observerParamType {
			observerParamIdx = i
			params[i] = reflect.Zero(observerParamType)
			continue
		}

		if paramType == reflect.TypeFor[*World]() {
			params[i] = reflect.ValueOf(world)
		} else if paramType == reflect.TypeFor[World]() {
			return observerEntry{}, fmt.Errorf("parameter %d: %w", i, ErrSystemParamWorldNotAPointer)
		} else if paramType.Implements(queryType) {
			query, err := parseQueryParam(paramType, world, world.logger, &world.outerWorlds)
			if err != nil {
				return observerEntry{}, fmt.Errorf("parameter %d: %w: %w", i, ErrSystemParamQueryNotValid, err)
			}
			queries = append(queries, query)
			params[i] = reflect.ValueOf(query)
		} else if paramType.Implements(eventReaderType) {
			eventReader, ok := reflect.TypeAssert[AnyEventReader](reflect.New(paramType.Elem()))
			if !ok {
				panic("failed to type assert AnyEventReader")
			}
			params[i] = *world.Events().GetReader(eventReader)
		} else if paramType.Implements(eventWriterType) {
			eventWriter, ok := reflect.TypeAssert[AnyEventWriter](reflect.New(paramType.Elem()))
			if !ok {
				panic("failed to type assert AnyEventWriter")
			}
			reflectedEventWriter := *world.Events().GetWriter(eventWriter)
			eventWriterInstance, ok := reflect.TypeAssert[AnyEventWriter](reflectedEventWriter)
			if !ok {
				panic("failed to type assert AnyEventWriter")
			}
			eventWriters = append(eventWriters, eventWriterInstance)
			params[i] = reflectedEventWriter
		} else if paramType.Implements(outerResourceType) {
			return observerEntry{}, fmt.Errorf("parameter %d: %w", i, ErrSystemParamOuterResourceIsAPointer)
		} else if paramType.Kind() != reflect.Pointer && reflect.PointerTo(paramType).Implements(outerResourceType) {
			instance := reflect.New(paramType)
			outerRes := instance.Interface().(AnyOuterResource)
			worldId, resType := outerRes.OuterResourceInfo()
			outerResources = append(outerResources, observerOuterResourceParam{
				paramIndex:        i,
				worldId:           *worldId,
				resourceType:      resType,
				outerResourceType: paramType,
			})
			if resType.Kind() == reflect.Pointer {
				// Pointer outer resources reference the same memory across all invocations,
				// so we initialize them once here (mirroring what prepare() does for schedule systems).
				outerWorld, exists := (*world.OuterWorlds())[*worldId]
				if !exists {
					return observerEntry{}, fmt.Errorf("parameter %d: %w: world id %d", i, ErrTargetWorldNotFound, *worldId)
				}
				resource, err := outerWorld.Resources().GetReflectResource(resType)
				if err != nil {
					return observerEntry{}, err
				}
				initInstance := reflect.New(paramType)
				initInstance.Elem().FieldByName("Value").Set(resource)
				params[i] = initInstance.Elem()
			} else {
				params[i] = reflect.Zero(paramType)
			}
		} else {
			resource, err := world.Resources().GetReflectResource(paramType)
			if err != nil {
				return observerEntry{}, fmt.Errorf("parameter %d: %w", i, handleInvalidSystemParam(paramType))
			}
			if paramType.Kind() == reflect.Pointer {
				params[i] = resource
			} else {
				params[i] = resource.Elem()
			}
		}
	}

	return observerEntry{
		systemEntry: systemEntry{
			system:     actionValue,
			params:     params,
			sourcePath: source,
		},
		observerParamIndex: observerParamIdx,
		queries:            queries,
		eventWriters:       eventWriters,
		outerResources:     outerResources,
	}, nil
}
