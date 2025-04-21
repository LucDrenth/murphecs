// Due to yield only being able to return 2 params, it can not be implemented for queries that return
// more than 2 components.
package ecs

import (
	"fmt"
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
func (q *query1Result[A]) Range() func(yield func(*A) bool) {
	return func(yield func(*A) bool) {
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

// Query1 gets the given component of all entities that match the options.
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
		world.logger.Warn(fmt.Sprintf("Query1 encountered issue with query options: %v", err))
	}

	for entityId, entityData := range world.entities {
		if ok := validateQueryFilters(entityData, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](world, entityData, &queryOptions)
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

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *query2Result[A, B]) Range() func(yield func(*A, *B) bool) {
	return func(yield func(*A, *B) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i], q.componentsB[i]) {
				return
			}
		}
	}
}

func (q *query2Result[A, B]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Query2 gets the given components of all entities that match the options.
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
func Query2[A, B IComponent](world *world, options ...queryOption) query2Result[A, B] {
	result := query2Result[A, B]{}
	queryOptions, err := createCombinedQueryOptions(options)
	if err != nil {
		world.logger.Warn(fmt.Sprintf("Query2 encountered issue with query options: %v", err))
	}

	for entityId, entityData := range world.entities {
		if ok := validateQueryFilters(entityData, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](world, entityData, &queryOptions)
		if !match {
			continue
		}

		b, match := getQueryComponent[B](world, entityData, &queryOptions)
		if !match {
			continue
		}

		result.componentsA = append(result.componentsA, a)
		result.componentsB = append(result.componentsB, b)
		result.entityIds = append(result.entityIds, entityId)
	}

	return result
}

type query3Result[A, B, C IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	entityIds   []EntityId
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *query3Result[A, B, C]) Iter(f func(entityId EntityId, a *A, b *B, c *C) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i]); err != nil {
			return err
		}
	}

	return nil
}

func (q *query3Result[A, B, C]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Query3 gets the given components of all entities that match the options.
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
func Query3[A, B, C IComponent](world *world, options ...queryOption) query3Result[A, B, C] {
	result := query3Result[A, B, C]{}
	queryOptions, err := createCombinedQueryOptions(options)
	if err != nil {
		world.logger.Warn(fmt.Sprintf("Query3 encountered issue with query options: %v", err))
	}

	for entityId, entityData := range world.entities {
		if ok := validateQueryFilters(entityData, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](world, entityData, &queryOptions)
		if !match {
			continue
		}

		b, match := getQueryComponent[B](world, entityData, &queryOptions)
		if !match {
			continue
		}

		c, match := getQueryComponent[C](world, entityData, &queryOptions)
		if !match {
			continue
		}

		result.componentsA = append(result.componentsA, a)
		result.componentsB = append(result.componentsB, b)
		result.componentsC = append(result.componentsC, c)
		result.entityIds = append(result.entityIds, entityId)
	}

	return result
}

type query4Result[A, B, C, D IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	componentsD []*D
	entityIds   []EntityId
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *query4Result[A, B, C, D]) Iter(f func(entityId EntityId, a *A, b *B, c *C, d *D) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i]); err != nil {
			return err
		}
	}

	return nil
}

func (q *query4Result[A, B, C, D]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Query4 gets the given components of all entities that match the options.
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
func Query4[A, B, C, D IComponent](world *world, options ...queryOption) query4Result[A, B, C, D] {
	result := query4Result[A, B, C, D]{}
	queryOptions, err := createCombinedQueryOptions(options)
	if err != nil {
		world.logger.Warn(fmt.Sprintf("Query4 encountered issue with query options: %v", err))
	}

	for entityId, entityData := range world.entities {
		if ok := validateQueryFilters(entityData, &queryOptions); !ok {
			continue
		}

		a, match := getQueryComponent[A](world, entityData, &queryOptions)
		if !match {
			continue
		}

		b, match := getQueryComponent[B](world, entityData, &queryOptions)
		if !match {
			continue
		}

		c, match := getQueryComponent[C](world, entityData, &queryOptions)
		if !match {
			continue
		}

		d, match := getQueryComponent[D](world, entityData, &queryOptions)
		if !match {
			continue
		}

		result.componentsA = append(result.componentsA, a)
		result.componentsB = append(result.componentsB, b)
		result.componentsC = append(result.componentsC, c)
		result.componentsD = append(result.componentsD, d)
		result.entityIds = append(result.entityIds, entityId)
	}

	return result
}

// getQueryComponent returns a pointer to T if the component is found on the entity.
//
// match is true when the entity has the component or if the component is marked marked as optional.
// When match is true, the entity should be present in the query results.
func getQueryComponent[T IComponent](world *world, entityData *entityData, queryOptions *combinedQueryOptions) (result *T, match bool) {
	componentType := getComponentType[T]()

	componentRegistryIndex, entityHasComponent := entityData.components[componentType]
	if !entityHasComponent {
		return nil, slices.Contains(queryOptions.optionalComponents, getComponentType[T]())
	}

	result, err := getComponentFromComponentRegistry[T](world.components[componentType], componentRegistryIndex)
	if err != nil {
		world.logger.Error(fmt.Sprintf("getQueryComponent encountered unexpected error: %v", err))
		return nil, false
	}

	return result, true
}

func validateQueryFilters(entityData *entityData, queryOptions *combinedQueryOptions) bool {
	for _, filter := range queryOptions.filters {
		if !filter.validate(entityData) {
			return false
		}
	}

	return true
}
