package ecs

import (
	"errors"
	"slices"

	"github.com/lucdrenth/murphecs/src/utils"
)

type filterType = int

const (
	filterTypeWith filterType = iota
	filterTypeWithout
	filterTypeAnd
	filterTypeOr
	filterTypeNone
)

type QueryParamFilter interface {
	getComponents(world *World) []ComponentId
	getFilterType() filterType
	getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error)
}
type NoFilter struct{}
type With[A AnyComponent] struct{}
type Without[A AnyComponent] struct{}
type And[A, B QueryParamFilter] struct{}
type Or[A, B QueryParamFilter] struct{}

func (filter NoFilter) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter And[A, B]) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter Or[A, B]) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter With[A]) getComponents(world *World) []ComponentId {
	return []ComponentId{ComponentIdFor[A](world)}
}

func (filter Without[A]) getComponents(world *World) []ComponentId {
	return []ComponentId{ComponentIdFor[A](world)}
}

func (filter NoFilter) getFilterType() filterType {
	return filterTypeNone
}

func (filter With[A]) getFilterType() filterType {
	return filterTypeWith
}

func (filter Without[A]) getFilterType() filterType {
	return filterTypeWithout
}

func (filter And[A, B]) getFilterType() filterType {
	return filterTypeAnd
}

func (filter Or[A, B]) getFilterType() filterType {
	return filterTypeOr
}

func (filter NoFilter) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")
}

func (filter With[A]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")
}

func (filter Without[A]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")
}

func (filter And[A, B]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	a, err = utils.ToConcrete[A]()
	if err != nil {
		return nil, nil, err
	}

	b, err = utils.ToConcrete[B]()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (filter Or[A, B]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	a, err = utils.ToConcrete[A]()
	if err != nil {
		return nil, nil, err
	}

	b, err = utils.ToConcrete[B]()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (filter With[A]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[With[A], NoOptional, NotLazy, DefaultWorld]](world)
}
func (filter Without[A]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[Without[A], NoOptional, NotLazy, DefaultWorld]](world)
}
func (filter And[A, B]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[And[A, B], NoOptional, NotLazy, DefaultWorld]](world)
}
func (filter Or[A, B]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[Or[A, B], NoOptional, NotLazy, DefaultWorld]](world)
}

type QueryFilter interface {
	// EntityMeetsCriteria returns false is the entity is filtered out
	EntityMeetsCriteria(*EntityData) bool

	// EntityMeetsCriteria returns false is the archetype is filtered out
	ArchetypeMeetsCriteria(*Archetype) bool
}
type queryFilterAnd struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterOr struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterWith struct {
	c []ComponentId
}
type queryFilterWithout struct {
	c []ComponentId
}

func (filter *queryFilterAnd) EntityMeetsCriteria(e *EntityData) bool {
	return filter.a.EntityMeetsCriteria(e) && filter.b.EntityMeetsCriteria(e)
}

func (filter *queryFilterOr) EntityMeetsCriteria(e *EntityData) bool {
	return filter.a.EntityMeetsCriteria(e) || filter.b.EntityMeetsCriteria(e)
}

func (filter *queryFilterWith) EntityMeetsCriteria(e *EntityData) bool {
	for _, c := range filter.c {
		if !e.hasComponent(c) {
			return false
		}
	}

	return true
}

func (filter *queryFilterWithout) EntityMeetsCriteria(e *EntityData) bool {
	return !slices.ContainsFunc(filter.c, e.hasComponent)
}

func (filter *queryFilterAnd) ArchetypeMeetsCriteria(archetype *Archetype) bool {
	return filter.a.ArchetypeMeetsCriteria(archetype) && filter.b.ArchetypeMeetsCriteria(archetype)
}

func (filter *queryFilterOr) ArchetypeMeetsCriteria(archetype *Archetype) bool {
	return filter.a.ArchetypeMeetsCriteria(archetype) || filter.b.ArchetypeMeetsCriteria(archetype)
}

func (filter *queryFilterWith) ArchetypeMeetsCriteria(archetype *Archetype) bool {
	for _, c := range filter.c {
		if !archetype.HasComponent(c) {
			return false
		}
	}

	return true
}

func (filter *queryFilterWithout) ArchetypeMeetsCriteria(archetype *Archetype) bool {
	return !slices.ContainsFunc(filter.c, archetype.HasComponent)
}
