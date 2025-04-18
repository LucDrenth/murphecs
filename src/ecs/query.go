// Due to yield only being able to return 2 params, it can not be implemented for queries that return
// more than 2 components.
package ecs

import (
	"slices"
)

type QueryResult interface {
	NumberOfResult() uint
}

type query1Result[A IComponent] struct {
	componentsA []*A
	entityIds   []EntityId
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *query1Result[A]) Iter(f func(entityId EntityId, a *A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *query1Result[A]) Range() func(yield func(a *A) bool) {
	return func(yield func(a *A) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i]) {
				return
			}
		}
	}
}

func (q *query1Result[A]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Query gets the given component of all entities that match the options.
// Use Query2, Query3, Query4 etc. to query multiple components.
//
// For filtering, choose from:
//   - ecs.With
//   - ecs.Without
//
// Or use ecs.Or and ecs.And to combine filters.
// If you pass multiple filters in options, all of them must pass for an entity to come up
// in the results.
//
// By default, entities have to have the given component. You can mark the component that
// you query as optional by passing ecs.Optional as an option. This will result in nil
// being returned for that component for the entities that don't have that component.
func Query[A IComponent](world *world, options ...queryOption) query1Result[A] {
	return Query1[A](world, options...)
}

// Query1 gets the given component of all entities that match the options.
// Use Query2, Query3, Query4 etc. to query multiple components.
//
// For filtering, choose from:
//   - ecs.With
//   - ecs.Without
//
// Or use ecs.Or and ecs.And to combine filters.
// If you pass multiple filters in options, all of them must pass for an entity to come up
// in the results.
//
// By default, entities have to have the given component. You can mark the component that
// you query as optional by passing ecs.Optional as an option. This will result in nil
// being returned for that component for the entities that don't have that component.
func Query1[A IComponent](world *world, options ...queryOption) query1Result[A] {
	result := query1Result[A]{}
	queryOptions, err := createCombinedQueryOptions(options)
	if err != nil {
		// TODO log warning but do not return.
	}

	for entityId, entry := range world.entities {
		if ok := validateQueryFilters(entry, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](entry, &queryOptions)
		if !match {
			continue
		}

		result.componentsA = append(result.componentsA, a)
		result.entityIds = append(result.entityIds, entityId)
	}

	return result
}

type query2Result[A, B IComponent] struct {
	componentsA []*A
	componentsB []*B
	entityIds   []EntityId
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *query2Result[A, B]) Iter(f func(entityId EntityId, a *A, b *B) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i]); err != nil {
			return err
		}
	}

	return nil
}

// Query2 gets the given components of all entities that match the options.
// Use Query, Query3, Query4 etc. to get a different number of components.
//
// For filtering, choose from:
//   - ecs.With
//   - ecs.Without
//
// Or use ecs.Or and ecs.And to combine filters.
// If you pass multiple filters in options, all of them must pass for an entity to come up
// in the results.
//
// By default, entities have to have the given components. You can mark components
// as optional by passing 1 or more ecs.Options as an option. This will result in nil
// being returned for that component for the entities that don't have that component.
func Query2[A IComponent, B IComponent](world *world, options ...queryOption) query2Result[A, B] {
	result := query2Result[A, B]{}
	queryOptions, err := createCombinedQueryOptions(options)
	if err != nil {
		// TODO log warning but do not return.
	}

	for entityId, entry := range world.entities {
		if ok := validateQueryFilters(entry, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](entry, &queryOptions)
		if !match {
			continue
		}

		b, match := getQueryComponent[B](entry, &queryOptions)
		if !match {
			continue
		}

		result.componentsA = append(result.componentsA, a)
		result.componentsB = append(result.componentsB, b)
		result.entityIds = append(result.entityIds, entityId)
	}

	return result
}

func getQueryComponent[T IComponent](entry *entry, queryOptions *combinedQueryOptions) (result *T, match bool) {
	result, _, _ = getComponentFromEntry[T](entry)
	match = result != nil || slices.Contains(queryOptions.optionalComponents, getComponentType[T]())
	return result, match
}

func validateQueryFilters(entry *entry, queryOptions *combinedQueryOptions) bool {
	for _, filter := range queryOptions.filters {
		if !filter.validate(entry) {
			return false
		}
	}

	return true
}
