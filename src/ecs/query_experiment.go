package ecs

import (
	"fmt"
	"reflect"
)

// ==============================
// 		    QUERY RESULTS
// ==============================

type Query1ResultExperimental[A IComponent] struct {
	componentsA []A
	entityIds   []EntityId
}

func (q *Query1ResultExperimental[A]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}

func (q *Query1ResultExperimental[A]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointer.
func (q *Query1ResultExperimental[A]) Iter(f func(entityId EntityId, a A)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i])
	}
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointer.
func (q *Query1ResultExperimental[A]) IterUntil(f func(entityId EntityId, a A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query1ResultExperimental[A]) Single() (EntityId, A, error) {
	if q.NumberOfResult() != 1 {
		var a A
		return nonExistingEntity, a, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], nil
}

// ==============================
//			   QUERY
// ==============================

type Query1Experimental[ComponentA IComponent, _ QueryOption] struct {
	Query1ResultExperimental[ComponentA]
	queryOptions

	componentIdA ComponentId
}

func (q *Query1Experimental[ComponentA, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResultExperimental[ComponentA](q.componentIdA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query1Experimental[A, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
	}
	q.options.optimize(q.components)
	return nil
}

func (q *Query1Experimental[A, Options]) ClearResults() { q.Clear() }

type exp1 struct{ Component }

var (
	_       = Query1Experimental[exp1, Default]{}
	_       = Query1Experimental[*exp1, Default]{}
	_ Query = &Query1Experimental[exp1, Default]{}
)

func fetchComponentForQueryResultExperimental[T IComponent](componentId ComponentId, entityRow uint, archetype *Archetype) (result T, err error) {
	storage := archetype.components[componentId]
	result, err = getComponentFromComponentStorage[T](storage, entityRow, reflect.TypeFor[T]().Kind() == reflect.Pointer)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve component %s from storage: %v", componentId.DebugString(), err)
	}

	return result, nil
}
