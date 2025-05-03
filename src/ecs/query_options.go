package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// combinedQueryOptions is a combination of all possible all query options, parsed for easy use within queries.
type combinedQueryOptions struct {
	Filters            []QueryFilter
	OptionalComponents []ComponentType
	ReadOnlyComponents combinedReadOnlyComponent
}

func (o *combinedQueryOptions) validateFilters(entityData *EntityData) bool {
	for i := range o.Filters {
		if !o.Filters[i].Validate(entityData) {
			return false
		}
	}

	return true
}

type combinedReadOnlyComponent struct {
	ComponentTypes []ComponentType
	IsAllReadOnly  bool
}

type iQueryOptions interface {
	getCombinedQueryOptions() (combinedQueryOptions, error)
}

type DefaultQueryOptions struct{}
type QueryOptionsAllReadOnly struct{}
type QueryOptions[_ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents] struct{}

func (o DefaultQueryOptions) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return combinedQueryOptions{}, nil
}

func (o QueryOptionsAllReadOnly) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return combinedQueryOptions{
		ReadOnlyComponents: combinedReadOnlyComponent{IsAllReadOnly: true},
	}, nil
}

func (o QueryOptions[QueryParamFilter, OptionalComponents, ReadOnlyComponents]) getCombinedQueryOptions() (combinedQueryOptions, error) {
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

	concreteOptionals, err := utils.ToConcrete[OptionalComponents]()
	if err != nil {
		return result, fmt.Errorf("failed to cast optional components to concrete type: %w", err)
	}
	result.OptionalComponents = concreteOptionals.getOptionalComponentTypes()

	readOnlyComponents, err := utils.ToConcrete[ReadOnlyComponents]()
	if err != nil {
		return result, fmt.Errorf("failed to cast read only components to concrete type: %w", err)
	}
	result.ReadOnlyComponents.ComponentTypes, result.ReadOnlyComponents.IsAllReadOnly = readOnlyComponents.getReadonlyComponentTypes()

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

func toCombinedQueryOptions[QueryOptions iQueryOptions]() (combinedQueryOptions, error) {
	result := combinedQueryOptions{}

	concreteQueryOptions, err := utils.ToConcrete[QueryOptions]()
	if err != nil {
		return result, fmt.Errorf("failed to cast query options to concrete type: %w", err)
	}

	return concreteQueryOptions.getCombinedQueryOptions()
}
