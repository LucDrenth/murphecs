package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Queries that are created using generics. They can be used as system parameters which allows systems
// to be executed in parallel.

// Query1 queries 1 component.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query1[ComponentA IComponent, _ QueryOption] struct {
	options combinedQueryOptions
	results Query1Result[ComponentA]
}

// Query2 queries 2 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query2[ComponentA, ComponentB IComponent, _ QueryOption] struct {
	options combinedQueryOptions
	results Query2Result[ComponentA, ComponentB]
}

// Query3 queries 3 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query3[ComponentA, ComponentB, ComponentC IComponent, _ QueryOption] struct {
	options combinedQueryOptions
	results Query3Result[ComponentA, ComponentB, ComponentC]
}

// Query4 queries 4 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query4[ComponentA, ComponentB, ComponentC, ComponentD IComponent, _ QueryOption] struct {
	options combinedQueryOptions
	results Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
}

func (q *Query1[ComponentA, QueryOptions]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match := getQueryComponent[ComponentA](world, entityData, &q.options)
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match := getQueryComponent[ComponentA](world, entityData, &q.options)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, &q.options)
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match := getQueryComponent[ComponentA](world, entityData, &q.options)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, &q.options)
		if !match {
			continue
		}

		c, match := getQueryComponent[ComponentC](world, entityData, &q.options)
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.componentsC = append(q.results.componentsC, c)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match := getQueryComponent[ComponentA](world, entityData, &q.options)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, &q.options)
		if !match {
			continue
		}

		c, match := getQueryComponent[ComponentC](world, entityData, &q.options)
		if !match {
			continue
		}

		d, match := getQueryComponent[ComponentD](world, entityData, &q.options)
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.componentsC = append(q.results.componentsC, c)
		q.results.componentsD = append(q.results.componentsD, d)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}
}

func (q *Query1[A, QueryOptions]) Prepare() (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions]()
	return err
}
func (q *Query2[A, B, QueryOptions]) Prepare() (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions]()
	return err
}
func (q *Query3[A, B, C, QueryOptions]) Prepare() (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions]()
	return err
}
func (q *Query4[A, B, C, D, QueryOptions]) Prepare() (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions]()
	return err
}

func (q *Query1[ComponentA, _]) Validate() error {
	return q.options.validateOptions([]ComponentType{
		GetComponentType[ComponentA](),
	})
}
func (q *Query2[ComponentA, ComponentB, _]) Validate() error {
	return q.options.validateOptions([]ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
	})
}
func (q *Query3[ComponentA, ComponentB, ComponentC, _]) Validate() error {
	return q.options.validateOptions([]ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
		GetComponentType[ComponentC](),
	})
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, _]) Validate() error {
	return q.options.validateOptions([]ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
		GetComponentType[ComponentC](),
		GetComponentType[ComponentD](),
	})
}

func (q *Query1[ComponentA, QueryOptions]) Result() *Query1Result[ComponentA] {
	return &q.results
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) Result() *Query2Result[ComponentA, ComponentB] {
	return &q.results
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) Result() *Query3Result[ComponentA, ComponentB, ComponentC] {
	return &q.results
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) Result() *Query4Result[ComponentA, ComponentB, ComponentC, ComponentD] {
	return &q.results
}

// getQueryComponent returns a pointer to T if the component is found on the entity.
//
// match is true when the entity has the component or if the component is marked marked as optional.
// When match is true, the entity should be present in the query results.
func getQueryComponent[T IComponent](world *World, entityData *EntityData, queryOptions *combinedQueryOptions) (result *T, match bool) {
	componentType := GetComponentType[T]()

	componentRegistryIndex, entityHasComponent := entityData.components[componentType]
	if !entityHasComponent {
		return nil, slices.Contains(queryOptions.OptionalComponents, componentType)
	}

	result, err := getComponentFromComponentRegistry[T](world.components[componentType], componentRegistryIndex)
	if err != nil {
		world.logger.Error(fmt.Sprintf("getQueryComponent encountered unexpected error: %v", err))
		return nil, false
	}

	if result != nil && (queryOptions.ReadOnlyComponents.IsAllReadOnly || slices.Contains(queryOptions.ReadOnlyComponents.ComponentTypes, componentType)) {
		result = utils.ClonePointerValue(result)
	}

	return result, true
}
