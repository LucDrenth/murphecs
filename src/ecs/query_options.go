package ecs

import (
	"fmt"
	"slices"
)

type combinedQueryOptions struct {
	filters            []queryOption
	optionalComponents []componentType
}

// createCombinedQueryOptions combines all queryOption's in to 1 efficient struct.
//
// Returns an ErrDuplicateComponent error if there are any duplicate optionalQueryComponent
// while still combining the rest of the options in to combinedQueryOptions.
// This error is more of a warning to be passed to the user, and should not be considered a failure.
func createCombinedQueryOptions(options []queryOption) (result combinedQueryOptions, err error) {
	for _, option := range options {
		switch option := option.(type) {
		case optionalQueryComponent:
			if slices.Contains(result.optionalComponents, option.componentType) {
				err = fmt.Errorf("%w: %s is specified as optional component multiple times", ErrDuplicateComponent, option.componentType.String())
			} else {
				result.optionalComponents = append(result.optionalComponents, option.componentType)
			}

			continue
		default:
			result.filters = append(result.filters, option)
		}
	}

	return result, err
}

// We have 1 interface for all filters and optional components. This allows the Query function to be easy
// to use, because we can use a variadic of queryOption.
//
// All these queryOption's will be merged in tot a combinedQueryOptions.
type queryOption interface {
	// validate that entry satisfies the filter
	validate(*entry) bool
}

type queryFilterAnd struct {
	a queryOption
	b queryOption
}

func (filter queryFilterAnd) validate(e *entry) bool {
	return filter.a.validate(e) && filter.b.validate(e)
}

type queryFilterOr struct {
	a queryOption
	b queryOption
}

func (filter queryFilterOr) validate(e *entry) bool {
	return filter.a.validate(e) || filter.b.validate(e)
}

type queryFilterWith struct {
	c componentType
}

func (filter queryFilterWith) validate(e *entry) bool {
	return e.containsComponentType(filter.c)
}

type queryFilterWithout struct {
	c componentType
}

func (filter queryFilterWithout) validate(e *entry) bool {
	return !e.containsComponentType(filter.c)
}

type optionalQueryComponent struct {
	componentType componentType
}

func (filter optionalQueryComponent) validate(e *entry) bool {
	return true
}

// And returns a query filter that asserts that both a and b are true
func And(a queryOption, b queryOption) queryOption {
	return queryFilterAnd{a, b}
}

// Or returns a query filter that asserts that either a or b is true
func Or(a queryOption, b queryOption) queryOption {
	return queryFilterOr{a, b}
}

// With returns a query filter that asserts that an entity has component T
func With[T IComponent]() queryOption {
	return queryFilterWith{c: getComponentType[T]()}
}

// Without returns a query filter that asserts that an entity does not have component T
func Without[T IComponent]() queryOption {
	return queryFilterWithout{c: getComponentType[T]()}
}

// Optional makes a queried component optional, resulting in nil if it is not present in an entity that passes all filters.
func Optional[T IComponent]() queryOption {
	return optionalQueryComponent{
		componentType: getComponentType[T](),
	}
}
