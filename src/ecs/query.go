package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type Query interface {
	Exec(world *World)

	// Prepare extracts the query options and puts it in CombinedQueryOptions. This should be called
	// once, after which the query is ready to be used (e.g. Exec can be called).
	//
	// This step is necessary because the way that queries are created is optimized for users, and not
	// for the computer. This method closes that gap.
	Prepare() error

	// Validate checks if there are any unexpected or unoptimized things. It returns an error if there
	// is something that can be optimized. The error should be treated as a warning, and does not mean
	// that the query can not be executed.
	//
	// Prepare must be called without returning any errors before calling this method.
	Validate() error
}

// QueryComponent represents a component that can be queried
type QueryComponent interface {
	isReadOnly() bool
	isOptional() bool
	innerQueryComponent() QueryComponent
	IComponent
}

type Query1[ComponentA QueryComponent, _ QueryOption] struct {
	a       ComponentA
	options combinedQueryOptions
	results Query1Result[ComponentA]
}
type Query2[ComponentA, ComponentB QueryComponent, _ QueryOption] struct {
	a       ComponentA
	b       ComponentB
	options combinedQueryOptions
	results Query2Result[ComponentA, ComponentB]
}
type Query3[ComponentA, ComponentB, ComponentC QueryComponent, _ QueryOption] struct {
	a       ComponentA
	b       ComponentB
	c       ComponentC
	options combinedQueryOptions
	results Query3Result[ComponentA, ComponentB, ComponentC]
}
type Query4[ComponentA, ComponentB, ComponentC, ComponentD QueryComponent, _ QueryOption] struct {
	a       ComponentA
	b       ComponentB
	c       ComponentC
	d       ComponentD
	options combinedQueryOptions
	results Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
}

func (q *Query1[ComponentA, QueryOptions]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match := getQueryComponent[ComponentA](world, entityData, q.a, q.options.IsAllReadOnly)
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

		a, match := getQueryComponent[ComponentA](world, entityData, q.a, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, q.b, q.options.IsAllReadOnly)
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

		a, match := getQueryComponent[ComponentA](world, entityData, q.a, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, q.b, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		c, match := getQueryComponent[ComponentC](world, entityData, q.c, q.options.IsAllReadOnly)
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

		a, match := getQueryComponent[ComponentA](world, entityData, q.a, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		b, match := getQueryComponent[ComponentB](world, entityData, q.b, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		c, match := getQueryComponent[ComponentC](world, entityData, q.c, q.options.IsAllReadOnly)
		if !match {
			continue
		}

		d, match := getQueryComponent[ComponentD](world, entityData, q.d, q.options.IsAllReadOnly)
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
	return q.options.validateOptions([]QueryComponent{
		q.a,
	})
}
func (q *Query2[ComponentA, ComponentB, _]) Validate() error {
	return q.options.validateOptions([]QueryComponent{
		q.a, q.b,
	})
}
func (q *Query3[ComponentA, ComponentB, ComponentC, _]) Validate() error {
	return q.options.validateOptions([]QueryComponent{
		q.a, q.b, q.c,
	})
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, _]) Validate() error {
	return q.options.validateOptions([]QueryComponent{
		q.a, q.b, q.c, q.d,
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
func getQueryComponent[T QueryComponent](world *World, entityData *EntityData, queryComponent QueryComponent, isAllReadOnly bool) (result *T, match bool) {
	concrete, err := utils.ToConcrete[T]()
	if err != nil {
		panic(err)
	}
	componentType := toComponentType(getInnerQueryComponent(concrete))

	componentRegistryIndex, entityHasComponent := entityData.components[componentType]
	if !entityHasComponent {
		return nil, queryComponent.isOptional()
	}

	// TODO we do not want to cast to T, but instead to the underlying component type (componentType)
	result, err = getComponentFromComponentRegistry[T](world.components[componentType], componentRegistryIndex)
	if err != nil {
		world.logger.Error(fmt.Sprintf("getQueryComponent encountered unexpected error: %v", err))
		return nil, false
	}

	if result != nil && (isAllReadOnly || queryComponent.isReadOnly()) {
		result = utils.ClonePointerValue(result)
	}

	return result, true
}
