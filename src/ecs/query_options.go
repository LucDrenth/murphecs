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

func getCombinedQueryOptions[filters QueryParamFilter, optionals OptionalComponents, readOnly ReadOnlyComponents]() (CombinedQueryOptions, error) {
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

	readOnlyComponents, err := utils.ToConcrete[readOnly]()
	if err != nil {
		return result, fmt.Errorf("failed to cast optionals to concrete type: %w", err)
	}

	result.ReadOnlyComponents.ComponentTypes, result.ReadOnlyComponents.IsAllReadOnly = readOnlyComponents.getReadonlyComponentTypes()

	return result, nil
}
