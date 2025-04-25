package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type filterType = int

const (
	filterTypeWith filterType = iota
	filterTypeWithout
	filterTypeAnd
	filterTypeOr
	filterTypeNone
)

type CombinedQueryOptions struct {
	Filters            []QueryFilter
	OptionalComponents []ComponentType
}

type QueryParamFilter interface {
	getComponents() []ComponentType
	getFilterType() filterType
}
type NoFilter struct{}
type With[A IComponent] struct{}
type Without[A IComponent] struct{}
type Or[A, B QueryParamFilter] struct{}
type And[A, B QueryParamFilter] struct{}

func (filter NoFilter) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter And[A, B]) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter Or[A, B]) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter With[A]) getComponents() []ComponentType {
	return []ComponentType{GetComponentType[A]()}
}

func (filter Without[A]) getComponents() []ComponentType {
	return []ComponentType{GetComponentType[A]()}
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

type QueryFilter interface {
	// Validate that entityData satisfies the filter
	Validate(*EntityData) bool
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
	c []ComponentType
}
type queryFilterWithout struct {
	c []ComponentType
}

func (filter queryFilterAnd) Validate(e *EntityData) bool {
	return filter.a.Validate(e) && filter.b.Validate(e)
}

func (filter queryFilterOr) Validate(e *EntityData) bool {
	return filter.a.Validate(e) || filter.b.Validate(e)
}

func (filter queryFilterWith) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if !e.hasComponent(c) {
			return false
		}
	}

	return true
}

func (filter queryFilterWithout) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if e.hasComponent(c) {
			return false
		}
	}

	return true
}

type OptionalComponents interface {
	getOptionalComponentTypes() []ComponentType
}
type AllRequired struct{}
type Optional1[A IComponent] struct{}
type Optional2[A, B IComponent] struct{}
type Optional3[A, B, C IComponent] struct{}
type Optional4[A, B, C, D IComponent] struct{}

func (o AllRequired) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{}
}

func (o Optional1[A]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
	}
}

func (o Optional2[A, B]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
	}
}

func (o Optional3[A, B, C]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
	}
}

func (o Optional4[A, B, C, D]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
		GetComponentType[D](),
	}
}

func getCombinedQueryOptions[filters QueryParamFilter, optionals OptionalComponents]() (CombinedQueryOptions, error) {
	result := CombinedQueryOptions{}

	concreteFilters, err := utils.ToConcrete[filters]()
	if err != nil {
		return result, fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	switch concreteFilterType := concreteFilters.getFilterType(); concreteFilterType {
	case filterTypeWith:
		result.Filters = append(result.Filters, queryFilterWith{c: concreteFilters.getComponents()})
	case filterTypeWithout:
		result.Filters = append(result.Filters, queryFilterWithout{c: concreteFilters.getComponents()})
	case filterTypeNone:
		break
	case filterTypeAnd:
		// TODO
		return result, fmt.Errorf("filterTypeAnd not yet implemented")
	case filterTypeOr:
		// TODO
		return result, fmt.Errorf("filterTypeOr not yet implemented")
	default:
		return result, fmt.Errorf("unhandled filter type: %d", concreteFilterType)
	}

	concreteOptionals, err := utils.ToConcrete[optionals]()
	if err != nil {
		return result, fmt.Errorf("failed to cast optionals to concrete type: %w", err)
	}
	result.OptionalComponents = concreteOptionals.getOptionalComponentTypes()

	return result, nil
}
