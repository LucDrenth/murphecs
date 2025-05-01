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

type combinedReadOnlyComponent struct {
	ComponentTypes []ComponentType
	IsAllReadOnly  bool
}

type CombinedQueryOptions struct {
	Filters            []QueryFilter
	OptionalComponents []ComponentType
	ReadOnlyComponents combinedReadOnlyComponent
}

func (o *CombinedQueryOptions) validateFilters(entityData *EntityData) bool {
	for i := range o.Filters {
		if !o.Filters[i].Validate(entityData) {
			return false
		}
	}

	return true
}

func getCombinedQueryOptions[filters QueryParamFilter, optionals OptionalComponents, readOnly ReadOnlyComponents]() (CombinedQueryOptions, error) {
	result := CombinedQueryOptions{}

	concreteFilters, err := utils.ToConcrete[filters]()
	if err != nil {
		return result, fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	filter, err := getFilterFromQueryOption(concreteFilters)
	if err != nil {
		return result, fmt.Errorf("failed to create filter: %w", err)
	}
	if filter != nil {
		result.Filters = append(result.Filters, filter)
	}

	concreteOptionals, err := utils.ToConcrete[optionals]()
	if err != nil {
		return result, fmt.Errorf("failed to cast optionals to concrete type: %w", err)
	}
	result.OptionalComponents = concreteOptionals.getOptionalComponentTypes()

	readOnlyComponents, err := utils.ToConcrete[readOnly]()
	if err != nil {
		return result, fmt.Errorf("failed to cast optionals to concrete type: %w", err)
	}

	result.ReadOnlyComponents.ComponentTypes, result.ReadOnlyComponents.IsAllReadOnly = readOnlyComponents.getReadonlyComponentTypes()

	return result, nil
}

func getFilterFromQueryOption(filters QueryParamFilter) (QueryFilter, error) {
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

			filterA, err := getFilterFromQueryOption(filterParamA)
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter for a: %w", err)
			}

			filterB, err := getFilterFromQueryOption(filterParamB)
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

			filterA, err := getFilterFromQueryOption(filterParamA)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for a: %w", err)
			}

			filterB, err := getFilterFromQueryOption(filterParamB)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for b: %w", err)
			}

			return queryFilterOr{a: filterA, b: filterB}, nil
		}
	}

	return nil, fmt.Errorf("unhandled filter type: %d", filters.getFilterType())
}
