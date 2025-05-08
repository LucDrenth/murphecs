package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// combinedQueryOptions is a combination of all possible all query options, parsed for easy use within queries.
type combinedQueryOptions struct {
	Filters            []QueryFilter
	OptionalComponents []ComponentType
	ReadOnlyComponents combinedReadOnlyComponent
	isLazy             bool
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
func (o *combinedQueryOptions) validateOptions(queryComponents []ComponentType) error {
	if duplicate, _, _ := utils.GetFirstDuplicate(o.OptionalComponents); duplicate != nil {
		return fmt.Errorf("optional component %s is given multiple times", (*duplicate).String())
	}

	for _, optional := range o.OptionalComponents {
		if !slices.Contains(queryComponents, optional) {
			return fmt.Errorf("optional component %s is not in query", optional.String())
		}
	}

	if duplicate, _, _ := utils.GetFirstDuplicate(o.ReadOnlyComponents.ComponentTypes); duplicate != nil {
		return fmt.Errorf("read-only component %s is given multiple times", (*duplicate).String())
	}

	for _, readOnly := range o.ReadOnlyComponents.ComponentTypes {
		if !slices.Contains(queryComponents, readOnly) {
			return fmt.Errorf("read-only component %s is not in query", readOnly.String())
		}
	}

	if o.ReadOnlyComponents.IsAllReadOnly && len(o.ReadOnlyComponents.ComponentTypes) > 0 {
		return fmt.Errorf("can not have specific read-only components together with IsAllReadOnly")
	}

	return nil
}

type combinedReadOnlyComponent struct {
	ComponentTypes []ComponentType
	IsAllReadOnly  bool
}

type QueryOption interface {
	getCombinedQueryOptions() (combinedQueryOptions, error)
}

// default query options: NoFilter, NoOptional, NoReadonly
type Default struct{}
type QueryOptionsAllReadOnly struct{}
type QueryOptions[_ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents, _ IsQueryLazy] struct{}

func (Default) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return combinedQueryOptions{}, nil
}

func (o QueryOptionsAllReadOnly) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return combinedQueryOptions{
		ReadOnlyComponents: combinedReadOnlyComponent{IsAllReadOnly: true},
	}, nil
}

func (o QueryOptions[QueryParamFilter, OptionalComponents, ReadOnlyComponents, IsQueryLazy]) getCombinedQueryOptions() (combinedQueryOptions, error) {
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

	queryOptionLazy, err := utils.ToConcrete[IsQueryLazy]()
	if err != nil {
		return result, fmt.Errorf("failed to cast IsQueryLazy to concrete type: %w", err)
	}
	result.isLazy = queryOptionLazy.isLazy()

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
		result.OptionalComponents = append(result.OptionalComponents, result.OptionalComponents...)

		result.ReadOnlyComponents.ComponentTypes = append(result.ReadOnlyComponents.ComponentTypes, options.ReadOnlyComponents.ComponentTypes...)

		if !result.ReadOnlyComponents.IsAllReadOnly && options.ReadOnlyComponents.IsAllReadOnly {
			result.ReadOnlyComponents.IsAllReadOnly = true
		}

		if !result.isLazy && options.isLazy {
			result.isLazy = true
		}
	}

	return result, nil
}

func QueryWithOptional[C IComponent](query Query) error {
	options := query.getOptions()

	componentType := GetComponentType[C]()
	options.OptionalComponents = append(options.OptionalComponents, componentType)

	return query.Validate()
}

func QueryWithReadOnly[C IComponent](query Query) error {
	options := query.getOptions()

	componentType := GetComponentType[C]()
	options.ReadOnlyComponents.ComponentTypes = append(options.ReadOnlyComponents.ComponentTypes, componentType)

	return query.Validate()
}

func QueryWithAllReadOnly(query Query) {
	query.getOptions().ReadOnlyComponents.IsAllReadOnly = true
}

func QueryWith[C IComponent](query Query) error {
	componentType := GetComponentType[C]()
	options := query.getOptions()
	options.Filters = append(options.Filters, queryFilterWith{c: []ComponentType{componentType}})

	return query.Validate()
}

func QueryWithout[C IComponent](query Query) error {
	componentType := GetComponentType[C]()
	options := query.getOptions()
	options.Filters = append(options.Filters, queryFilterWithout{c: []ComponentType{componentType}})

	return query.Validate()
}

func QueryWithFilters[Filters QueryParamFilter](query Query) error {
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

	return query.Validate()
}
