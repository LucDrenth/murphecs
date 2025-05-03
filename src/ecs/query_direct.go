package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Queries that are created with functions. Compared to queries from [query_generic], these
// are easier to use because you can specify 1 query param at a time, but systems that use them can
// not be ran in parallel.

type iDirectQuery interface {
	getOptions() *combinedQueryOptions
	getComponentTypes() []ComponentType
}

type directQuery struct {
	options        combinedQueryOptions
	componentTypes []ComponentType
}

func (q *directQuery) getOptions() *combinedQueryOptions {
	return &q.options
}

func (q *directQuery) getComponentTypes() []ComponentType {
	return q.componentTypes
}

func QueryWithOptional[C IComponent](query iDirectQuery) error {
	options := query.getOptions()
	componentTypes := query.getComponentTypes()

	componentType := GetComponentType[C]()
	if !slices.Contains(componentTypes, componentType) {
		return fmt.Errorf("query does not include component %s", componentType.String())
	}

	options.OptionalComponents = append(options.OptionalComponents, componentType)

	return nil
}

func QueryWithReadOnly[C IComponent](query iDirectQuery) error {
	options := query.getOptions()
	componentTypes := query.getComponentTypes()

	componentType := GetComponentType[C]()
	if !slices.Contains(componentTypes, componentType) {
		return fmt.Errorf("query does not include component %s", componentType.String())
	}

	options.ReadOnlyComponents.ComponentTypes = append(options.ReadOnlyComponents.ComponentTypes, componentType)

	return nil
}

func QueryWithAllReadOnly(query iDirectQuery) {
	query.getOptions().ReadOnlyComponents.IsAllReadOnly = true
}

func QueryWith[C IComponent](query iDirectQuery) {
	componentType := GetComponentType[C]()
	options := query.getOptions()
	options.Filters = append(options.Filters, queryFilterWith{c: []ComponentType{componentType}})
}

func QueryWithout[C IComponent](query iDirectQuery) {
	componentType := GetComponentType[C]()
	options := query.getOptions()
	options.Filters = append(options.Filters, queryFilterWithout{c: []ComponentType{componentType}})
}

func QueryWithFilters[Filters QueryParamFilter](query iDirectQuery) error {
	concreteFilters, err := utils.ToConcrete[Filters]()
	if err != nil {
		return fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	filter, err := getFilterFromConcreteQueryParamFilter(concreteFilters)
	if err != nil {
		return fmt.Errorf("failed to create filter: %w", err)
	}

	if filter == nil {
		return nil
	}

	options := query.getOptions()
	options.Filters = append(options.Filters, filter)
	return nil
}

type directQuery1[ComponentA IComponent] struct {
	directQuery
}
type directQuery2[ComponentA, ComponentB IComponent] struct {
	directQuery
}
type directQuery3[ComponentA, ComponentB, ComponentC IComponent] struct {
	directQuery
}
type directQuery4[ComponentA, ComponentB, ComponentC, ComponentD IComponent] struct {
	directQuery
}

func NewQuery1[ComponentA IComponent]() directQuery1[ComponentA] {
	return directQuery1[ComponentA]{directQuery: directQuery{componentTypes: []ComponentType{
		GetComponentType[ComponentA](),
	}}}
}
func NewQuery2[ComponentA, ComponentB IComponent]() directQuery2[ComponentA, ComponentB] {
	return directQuery2[ComponentA, ComponentB]{directQuery: directQuery{componentTypes: []ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
	}}}
}
func NewQuery3[ComponentA, ComponentB, ComponentC IComponent]() directQuery3[ComponentA, ComponentB, ComponentC] {
	return directQuery3[ComponentA, ComponentB, ComponentC]{directQuery: directQuery{componentTypes: []ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
		GetComponentType[ComponentC](),
	}}}
}
func NewQuery4[ComponentA, ComponentB, ComponentC, ComponentD IComponent]() directQuery4[ComponentA, ComponentB, ComponentC, ComponentD] {
	return directQuery4[ComponentA, ComponentB, ComponentC, ComponentD]{directQuery: directQuery{componentTypes: []ComponentType{
		GetComponentType[ComponentA](),
		GetComponentType[ComponentB](),
		GetComponentType[ComponentC](),
		GetComponentType[ComponentD](),
	}}}
}

func (q *directQuery1[ComponentA]) Exec(world *World) Query1Result[ComponentA] {
	query := Query1[ComponentA, DefaultQueryOptions]{
		options: q.options,
	}
	query.Exec(world)
	return query.results
}
func (q *directQuery2[ComponentA, ComponentB]) Exec(world *World) Query2Result[ComponentA, ComponentB] {
	query := Query2[ComponentA, ComponentB, DefaultQueryOptions]{
		options: q.options,
	}
	query.Exec(world)
	return query.results
}
func (q *directQuery3[ComponentA, ComponentB, ComponentC]) Exec(world *World) Query3Result[ComponentA, ComponentB, ComponentC] {
	query := Query3[ComponentA, ComponentB, ComponentC, DefaultQueryOptions]{
		options: q.options,
	}
	query.Exec(world)
	return query.results
}
func (q *directQuery4[ComponentA, ComponentB, ComponentC, ComponentD]) Exec(world *World) Query4Result[ComponentA, ComponentB, ComponentC, ComponentD] {
	query := Query4[ComponentA, ComponentB, ComponentC, ComponentD, DefaultQueryOptions]{
		options: q.options,
	}
	query.Exec(world)
	return query.results
}
