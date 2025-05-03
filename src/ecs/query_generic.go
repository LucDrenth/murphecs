package ecs

// Thus
//

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type Query interface {
	Exec(world *World)

	// PrepareOptions extracts the query options and puts it in CombinedQueryOptions. This should be called
	// once, after which the query is ready to be used (e.g. Exec can be called).
	PrepareOptions() error
}
type Query1[ComponentA IComponent, _ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents] struct {
	options CombinedQueryOptions
	results Query1Result[ComponentA]
}
type Query2[ComponentA, ComponentB IComponent, _ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents] struct {
	options CombinedQueryOptions
	results Query2Result[ComponentA, ComponentB]
}
type Query3[ComponentA, ComponentB, ComponentC IComponent, _ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents] struct {
	options CombinedQueryOptions
	results Query3Result[ComponentA, ComponentB, ComponentC]
}
type Query4[ComponentA, ComponentB, ComponentC, ComponentD IComponent, _ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents] struct {
	options CombinedQueryOptions
	results Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
}

func (q *Query1[ComponentA, Filters, OptionalComponents, ReadOnlyComponent]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if !q.options.validateFilters(entityData) {
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
func (q *Query2[ComponentA, ComponentB, Filters, OptionalComponents, ReadOnlyComponent]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if !q.options.validateFilters(entityData) {
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
func (q *Query3[ComponentA, ComponentB, ComponentC, Filters, OptionalComponents, ReadOnlyComponent]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if !q.options.validateFilters(entityData) {
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
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, Filters, OptionalComponents, ReadOnlyComponent]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if !q.options.validateFilters(entityData) {
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

func (q *Query1[A, Filters, OptionalComponents, ReadOnlyComponent]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents, ReadOnlyComponent]()
	return err
}
func (q *Query2[A, B, Filters, OptionalComponents, ReadOnlyComponent]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents, ReadOnlyComponent]()
	return err
}
func (q *Query3[A, B, C, Filters, OptionalComponents, ReadOnlyComponent]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents, ReadOnlyComponent]()
	return err
}
func (q *Query4[A, B, C, D, Filters, OptionalComponents, ReadOnlyComponent]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents, ReadOnlyComponent]()
	return err
}

func (q *Query1[ComponentA, Filters, OptionalComponents, ReadOnlyComponent]) Result() *Query1Result[ComponentA] {
	return &q.results
}
func (q *Query2[ComponentA, ComponentB, Filters, OptionalComponents, ReadOnlyComponent]) Result() *Query2Result[ComponentA, ComponentB] {
	return &q.results
}
func (q *Query3[ComponentA, ComponentB, ComponentC, Filters, OptionalComponents, ReadOnlyComponent]) Result() *Query3Result[ComponentA, ComponentB, ComponentC] {
	return &q.results
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, Filters, OptionalComponents, ReadOnlyComponent]) Result() *Query4Result[ComponentA, ComponentB, ComponentC, ComponentD] {
	return &q.results
}

// getQueryComponent returns a pointer to T if the component is found on the entity.
//
// match is true when the entity has the component or if the component is marked marked as optional.
// When match is true, the entity should be present in the query results.
func getQueryComponent[T IComponent](world *World, entityData *EntityData, queryOptions *CombinedQueryOptions) (result *T, match bool) {
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

	if result != nil && queryOptions.ReadOnlyComponents.IsAllReadOnly || slices.Contains(queryOptions.ReadOnlyComponents.ComponentTypes, componentType) {
		result = utils.ClonePointerValue(result)
	}

	return result, true
}
