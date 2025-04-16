package ecs

import (
	"slices"
)

type query1Result[A IComponent] struct {
	componentsA []*A
	entityIds   []entityId
}

func (q query1Result[A]) Iter() func(yield func(entityId entityId, a *A) bool) {
	return func(yield func(entityId entityId, a *A) bool) {
		for i := range q.entityIds {
			if !yield(q.entityIds[i], q.componentsA[i]) {
				return
			}
		}
	}
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
	entityIds   []entityId
}

func (q query2Result[A, B]) Iter() func(yield func(entityId entityId, a *A, b *B) bool) {
	return func(yield func(entityId entityId, a *A, b *B) bool) {
		for i := range q.entityIds {
			if !yield(q.entityIds[i], q.componentsA[i], q.componentsB[i]) {
				return
			}
		}
	}
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
func Query2[A IComponent, B IComponent](world *world, options ...queryOption) (query2Result[A, B], error) {
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

	return result, nil
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
