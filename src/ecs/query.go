// Due to yield only being able to return 2 params, it can not be implemented for queries that return
// more than 2 components.
package ecs

import (
	"fmt"
	"slices"
)

type Query interface {
	Exec(world *World)

	// PrepareOptions extracts the query options and puts it in CombinedQueryOptions. This should be called
	// once, after which the query is ready to be used (e.g. Exec can be called).
	PrepareOptions() error
}
type Query1[ComponentA IComponent, _ QueryParamFilter, _ OptionalComponents] struct {
	options CombinedQueryOptions
	results Query1Result[ComponentA]
}
type Query2[ComponentA, ComponentB IComponent, _ QueryParamFilter, _ OptionalComponents] struct {
	options CombinedQueryOptions
	results Query2Result[ComponentA, ComponentB]
}
type Query3[ComponentA, ComponentB, ComponentC IComponent, _ QueryParamFilter, _ OptionalComponents] struct {
	options CombinedQueryOptions
	results Query3Result[ComponentA, ComponentB, ComponentC]
}
type Query4[ComponentA, ComponentB, ComponentC, ComponentD IComponent, _ QueryParamFilter, _ OptionalComponents] struct {
	options CombinedQueryOptions
	results Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
}

func (q *Query1[ComponentA, Filters, OptionalComponents]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
		a, match := getQueryComponent[ComponentA](world, entityData, &q.options)
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}
}
func (q *Query2[ComponentA, ComponentB, Filters, OptionalComponents]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
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
func (q *Query3[ComponentA, ComponentB, ComponentC, Filters, OptionalComponents]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
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
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, Filters, OptionalComponents]) Exec(world *World) {
	q.results.Clear()

	for entityId, entityData := range world.entities {
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

func (q *Query1[A, Filters, OptionalComponents]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents]()
	return err
}
func (q *Query2[A, B, Filters, OptionalComponents]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents]()
	return err
}
func (q *Query3[A, B, C, Filters, OptionalComponents]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents]()
	return err
}
func (q *Query4[A, B, C, D, Filters, OptionalComponents]) PrepareOptions() (err error) {
	q.options, err = getCombinedQueryOptions[Filters, OptionalComponents]()
	return err
}

func (q *Query1[ComponentA, Filters, OptionalComponents]) Result() *Query1Result[ComponentA] {
	return &q.results
}
func (q *Query2[ComponentA, ComponentB, Filters, OptionalComponents]) Result() *Query2Result[ComponentA, ComponentB] {
	return &q.results
}
func (q *Query3[ComponentA, ComponentB, ComponentC, Filters, OptionalComponents]) Result() *Query3Result[ComponentA, ComponentB, ComponentC] {
	return &q.results
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, Filters, OptionalComponents]) Result() *Query4Result[ComponentA, ComponentB, ComponentC, ComponentD] {
	return &q.results
}

type QueryResult interface {
	NumberOfResult() uint
	Clear()
}
type Query1Result[A IComponent] struct {
	componentsA []*A
	entityIds   []EntityId
}
type Query2Result[A, B IComponent] struct {
	componentsA []*A
	componentsB []*B
	entityIds   []EntityId
}
type Query3Result[A, B, C IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	entityIds   []EntityId
}
type Query4Result[A, B, C, D IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	componentsD []*D
	entityIds   []EntityId
}

func (q *Query1Result[A]) Clear() {
	q.componentsA = []*A{}
	q.entityIds = []EntityId{}
}
func (q *Query2Result[A, B]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.entityIds = []EntityId{}
}
func (q *Query3Result[A, B, C]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.componentsC = []*C{}
	q.entityIds = []EntityId{}
}
func (q *Query4Result[A, B, C, D]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.componentsC = []*C{}
	q.componentsD = []*D{}
	q.entityIds = []EntityId{}
}

func (q *Query1Result[A]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query2Result[A, B]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query3Result[A, B, C]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query4Result[A, B, C, D]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *Query1Result[A]) Range() func(yield func(*A) bool) {
	return func(yield func(*A) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i]) {
				return
			}
		}
	}
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *Query2Result[A, B]) Range() func(yield func(*A, *B) bool) {
	return func(yield func(*A, *B) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i], q.componentsB[i]) {
				return
			}
		}
	}
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query1Result[A]) Iter(f func(entityId EntityId, a *A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query2Result[A, B]) Iter(f func(entityId EntityId, a *A, b *B) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query3Result[A, B, C]) Iter(f func(entityId EntityId, a *A, b *B, c *C) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query4Result[A, B, C, D]) Iter(f func(entityId EntityId, a *A, b *B, c *C, d *D) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i]); err != nil {
			return err
		}
	}

	return nil
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

	return result, true
}
