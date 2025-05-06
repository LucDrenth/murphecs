package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type ReadOnly[C QueryComponent] struct {
	Inner C
	Component
}

type Optional[C QueryComponent] struct {
	Inner C
	Component
}

func (c Optional[_]) isOptional() bool {
	return true
}

func (c ReadOnly[_]) isReadOnly() bool {
	return true
}

func (c Optional[C]) innerQueryComponent() QueryComponent {
	return c.Inner
}

func (c ReadOnly[_]) innerQueryComponent() QueryComponent {
	return c.Inner
}

func getInnerQueryComponent(c QueryComponent) IComponent {
	inner := c.innerQueryComponent()

	if inner == nil {
		return c
	} else {
		return getInnerQueryComponent(inner)
	}
}

// combinedQueryOptions is a combination of all possible all query options, parsed for easy use within queries.
type combinedQueryOptions struct {
	Filters       []QueryFilter
	IsAllReadOnly bool
}

func (o *combinedQueryOptions) isFilteredOut(entityData *EntityData) bool {
	for i := range o.Filters {
		if !o.Filters[i].Validate(entityData) {
			return true
		}
	}

	return false
}

// validateOptions returns an error if there are any invalid or non-logical options.
// If an error is returned, it does not mean that the combinedQueryOptions can not be
// used in a query, thus the error should be treated as a warning.
func (o *combinedQueryOptions) validateOptions(queryComponents []QueryComponent) error {
	if o.IsAllReadOnly {
		for _, q := range queryComponents {
			if q.isReadOnly() {
				return fmt.Errorf("query has both read-only component and IsAllReadOnly")
			}
		}
	}

	return nil
}

type QueryOption interface {
	getCombinedQueryOptions() (combinedQueryOptions, error)
}

// Default query options
type Default struct{}
type QueryOptions[_ QueryParamFilter, _ ReadOnlyComponents] struct{}

func (o Default) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return combinedQueryOptions{}, nil
}

func (o QueryOptions[QueryParamFilter, ReadOnlyComponents]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	result := combinedQueryOptions{}

	concreteFilters, err := utils.ToConcrete[QueryParamFilter]()
	if err != nil {
		return result, fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	filter, err := getFilterFromConcreteQueryParamFilter(concreteFilters)
	if err != nil {
		return result, fmt.Errorf("failed to create filter: %w", err)
	}
	if filter != nil {
		result.Filters = append(result.Filters, filter)
	}

	readOnlyComponents, err := utils.ToConcrete[ReadOnlyComponents]()
	if err != nil {
		return result, fmt.Errorf("failed to cast read only components to concrete type: %w", err)
	}
	result.IsAllReadOnly = readOnlyComponents.isAllReadOnly()

	return result, nil
}

func getFilterFromConcreteQueryParamFilter(filters QueryParamFilter) (QueryFilter, error) {
	switch concreteFilterType := filters.getFilterType(); concreteFilterType {
	case filterTypeWith:
		return queryFilterWith{c: filters.getComponents()}, nil
	case filterTypeWithout:
		return queryFilterWithout{c: filters.getComponents()}, nil
	case filterTypeNone:
		return nil, nil
	case filterTypeAnd:
		{
			filterParamA, filterParamB, err := filters.getNestedFilters()
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter: %w", err)
			}

			filterA, err := getFilterFromConcreteQueryParamFilter(filterParamA)
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter for a: %w", err)
			}

			filterB, err := getFilterFromConcreteQueryParamFilter(filterParamB)
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter for b: %w", err)
			}

			return queryFilterAnd{a: filterA, b: filterB}, nil
		}
	case filterTypeOr:
		{
			filterParamA, filterParamB, err := filters.getNestedFilters()
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter: %w", err)
			}

			filterA, err := getFilterFromConcreteQueryParamFilter(filterParamA)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for a: %w", err)
			}

			filterB, err := getFilterFromConcreteQueryParamFilter(filterParamB)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for b: %w", err)
			}

			return queryFilterOr{a: filterA, b: filterB}, nil
		}
	}

	return nil, fmt.Errorf("unhandled filter type: %d", filters.getFilterType())
}

func toCombinedQueryOptions[QueryOptions QueryOption]() (combinedQueryOptions, error) {
	result := combinedQueryOptions{}

	concreteQueryOptions, err := utils.ToConcrete[QueryOptions]()
	if err != nil {
		return result, fmt.Errorf("failed to cast query options to concrete type: %w", err)
	}

	return concreteQueryOptions.getCombinedQueryOptions()
}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions2[A, B QueryOption] struct{}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions3[A, B, C QueryOption] struct{}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions4[A, B, C, D QueryOption] struct{}

func (o QueryOptions2[A, B]) getCombinedQueryOptions() (result combinedQueryOptions, err error) {
	a, err := utils.ToConcrete[A]()
	if err != nil {
		return result, err
	}

	b, err := utils.ToConcrete[B]()
	if err != nil {
		return result, err
	}

	return mergeQueryOptions([]QueryOption{a, b})
}
func (o QueryOptions3[A, B, C]) getCombinedQueryOptions() (result combinedQueryOptions, err error) {
	a, err := utils.ToConcrete[A]()
	if err != nil {
		return result, err
	}

	b, err := utils.ToConcrete[B]()
	if err != nil {
		return result, err
	}

	c, err := utils.ToConcrete[C]()
	if err != nil {
		return result, err
	}

	return mergeQueryOptions([]QueryOption{a, b, c})
}
func (o QueryOptions4[A, B, C, D]) getCombinedQueryOptions() (result combinedQueryOptions, err error) {
	a, err := utils.ToConcrete[A]()
	if err != nil {
		return result, err
	}

	b, err := utils.ToConcrete[B]()
	if err != nil {
		return result, err
	}

	c, err := utils.ToConcrete[C]()
	if err != nil {
		return result, err
	}

	d, err := utils.ToConcrete[D]()
	if err != nil {
		return result, err
	}

	return mergeQueryOptions([]QueryOption{a, b, c, d})
}

func mergeQueryOptions(queryOptions []QueryOption) (result combinedQueryOptions, err error) {
	for _, queryOption := range queryOptions {
		options, err := queryOption.getCombinedQueryOptions()
		if err != nil {
			return result, err
		}

		result.Filters = append(result.Filters, options.Filters...)

		if !result.IsAllReadOnly && options.IsAllReadOnly {
			result.IsAllReadOnly = true
		}
	}

	return result, nil
}
