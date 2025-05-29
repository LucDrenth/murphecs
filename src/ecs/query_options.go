package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murphecs/src/utils"
)

// CombinedQueryOptions is a combination of all possible all query options, parsed for easy use within queries.
type CombinedQueryOptions struct {
	Filters            []QueryFilter
	OptionalComponents []ComponentId
	ReadOnlyComponents combinedReadOnlyComponent
	isLazy             bool
	TargetWorld        *WorldId
}

func (o *CombinedQueryOptions) isArchetypeFilteredOut(archetype *Archetype) bool {
	for i := range o.Filters {
		if !o.Filters[i].ArchetypeMeetsCriteria(archetype) {
			return true
		}
	}

	return false
}

// validateOptions returns an error if there are any invalid or non-logical options.
// If an error is returned, it does not mean that the combinedQueryOptions can not be
// used in a query, thus the error should be treated as a warning.
func (o *CombinedQueryOptions) validateOptions(queryComponents []ComponentId) error {
	if duplicate, _, _ := utils.GetFirstDuplicate(o.OptionalComponents); duplicate != nil {
		return fmt.Errorf("optional component %s is given multiple times", (*duplicate).DebugString())
	}

	for _, optional := range o.OptionalComponents {
		if !slices.Contains(queryComponents, optional) {
			return fmt.Errorf("optional component %s is not in query", optional.DebugString())
		}
	}

	if duplicate, _, _ := utils.GetFirstDuplicate(o.ReadOnlyComponents.ComponentIds); duplicate != nil {
		return fmt.Errorf("read-only component %s is given multiple times", (*duplicate).DebugString())
	}

	for _, readOnly := range o.ReadOnlyComponents.ComponentIds {
		if !slices.Contains(queryComponents, readOnly) {
			return fmt.Errorf("read-only component %s is not in query", readOnly.DebugString())
		}
	}

	if o.ReadOnlyComponents.IsAllReadOnly && len(o.ReadOnlyComponents.ComponentIds) > 0 {
		return fmt.Errorf("should not have specific read-only components together with IsAllReadOnly")
	}

	return nil
}

func (o *CombinedQueryOptions) optimize(queryComponents []ComponentId) {
	// If IsAllReadOnly is not set, and all query components are individually marked as ReadOnly, we can
	// set IsAllReadOnly to true. This saves a little bit of memory and slightly increase performance.
	if !o.ReadOnlyComponents.IsAllReadOnly {
		setAllReadOnly := true
		for _, component := range queryComponents {
			if !slices.ContainsFunc(o.ReadOnlyComponents.ComponentIds, func(c ComponentId) bool {
				return component.id == c.id
			}) {
				setAllReadOnly = false
			}
		}

		if setAllReadOnly {
			o.ReadOnlyComponents.IsAllReadOnly = true
		}
	}

	if o.ReadOnlyComponents.IsAllReadOnly {
		// free up memory
		o.ReadOnlyComponents.ComponentIds = nil
	}
}

type combinedReadOnlyComponent struct {
	ComponentIds  []ComponentId
	IsAllReadOnly bool
}

type QueryOption interface {
	GetCombinedQueryOptions(*World) (CombinedQueryOptions, error)
}

// default query options: NoFilter, NoOptional, NoReadonly, NotLazy
type Default struct{}
type QueryOptionsAllReadOnly struct{}
type QueryOptions[_ QueryParamFilter, _ OptionalComponents, _ ReadOnlyComponents, _ IsQueryLazy, _ TargetWorld] struct{}

func (Default) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return CombinedQueryOptions{}, nil
}

func (o QueryOptionsAllReadOnly) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return CombinedQueryOptions{
		ReadOnlyComponents: combinedReadOnlyComponent{IsAllReadOnly: true},
	}, nil
}

func (o QueryOptions[QueryParamFilter, OptionalComponents, ReadOnlyComponents, IsQueryLazy, TargetWorld]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	result := CombinedQueryOptions{}

	concreteFilters, err := utils.ToConcrete[QueryParamFilter]()
	if err != nil {
		return result, fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	filter, err := getFilterFromConcreteQueryParamFilter(concreteFilters, world)
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
	result.OptionalComponents = concreteOptionals.getOptionalComponentIds(world)

	readOnlyComponents, err := utils.ToConcrete[ReadOnlyComponents]()
	if err != nil {
		return result, fmt.Errorf("failed to cast read only components to concrete type: %w", err)
	}
	result.ReadOnlyComponents.ComponentIds, result.ReadOnlyComponents.IsAllReadOnly = readOnlyComponents.getReadonlyComponentIds(world)

	queryOptionLazy, err := utils.ToConcrete[IsQueryLazy]()
	if err != nil {
		return result, fmt.Errorf("failed to cast IsQueryLazy to concrete type: %w", err)
	}
	result.isLazy = queryOptionLazy.isLazy()

	queryOptionTargetWorld, err := utils.ToConcrete[TargetWorld]()
	if err != nil {
		return result, fmt.Errorf("failed to cast TargetWorld to concrete type: %w", err)
	}
	result.TargetWorld = queryOptionTargetWorld.GetWorldId()

	return result, nil
}

func getFilterFromConcreteQueryParamFilter(filters QueryParamFilter, world *World) (QueryFilter, error) {
	switch concreteFilterType := filters.getFilterType(); concreteFilterType {
	case filterTypeWith:
		return &queryFilterWith{c: filters.getComponents(world)}, nil
	case filterTypeWithout:
		return &queryFilterWithout{c: filters.getComponents(world)}, nil
	case filterTypeNone:
		return nil, nil
	case filterTypeAnd:
		{
			filterParamA, filterParamB, err := filters.getNestedFilters()
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter: %w", err)
			}

			filterA, err := getFilterFromConcreteQueryParamFilter(filterParamA, world)
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter for a: %w", err)
			}

			filterB, err := getFilterFromConcreteQueryParamFilter(filterParamB, world)
			if err != nil {
				return nil, fmt.Errorf("failed to create AND filter for b: %w", err)
			}

			return &queryFilterAnd{a: filterA, b: filterB}, nil
		}
	case filterTypeOr:
		{
			filterParamA, filterParamB, err := filters.getNestedFilters()
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter: %w", err)
			}

			filterA, err := getFilterFromConcreteQueryParamFilter(filterParamA, world)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for a: %w", err)
			}

			filterB, err := getFilterFromConcreteQueryParamFilter(filterParamB, world)
			if err != nil {
				return nil, fmt.Errorf("failed to create OR filter for b: %w", err)
			}

			return &queryFilterOr{a: filterA, b: filterB}, nil
		}
	}

	return nil, fmt.Errorf("unhandled filter type: %d", filters.getFilterType())
}

func toCombinedQueryOptions[QueryOptions QueryOption](world *World) (CombinedQueryOptions, error) {
	result := CombinedQueryOptions{}

	concreteQueryOptions, err := utils.ToConcrete[QueryOptions]()
	if err != nil {
		return result, fmt.Errorf("failed to cast query options to concrete type: %w", err)
	}

	return concreteQueryOptions.GetCombinedQueryOptions(world)
}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions2[A, B QueryOption] struct{}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions3[A, B, C QueryOption] struct{}

// Multiple query options that will be combined. Multiple filters will result in them being used with an AND operator.
type QueryOptions4[A, B, C, D QueryOption] struct{}

func (o QueryOptions2[A, B]) GetCombinedQueryOptions(world *World) (result CombinedQueryOptions, err error) {
	a, err := utils.ToConcrete[A]()
	if err != nil {
		return result, err
	}

	b, err := utils.ToConcrete[B]()
	if err != nil {
		return result, err
	}

	return mergeQueryOptions([]QueryOption{a, b}, world)
}
func (o QueryOptions3[A, B, C]) GetCombinedQueryOptions(world *World) (result CombinedQueryOptions, err error) {
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

	return mergeQueryOptions([]QueryOption{a, b, c}, world)
}
func (o QueryOptions4[A, B, C, D]) GetCombinedQueryOptions(world *World) (result CombinedQueryOptions, err error) {
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

	return mergeQueryOptions([]QueryOption{a, b, c, d}, world)
}

func mergeQueryOptions(queryOptions []QueryOption, world *World) (result CombinedQueryOptions, err error) {
	for _, queryOption := range queryOptions {
		options, err := queryOption.GetCombinedQueryOptions(world)
		if err != nil {
			return result, err
		}

		result.Filters = append(result.Filters, options.Filters...)
		result.OptionalComponents = append(result.OptionalComponents, result.OptionalComponents...)

		result.ReadOnlyComponents.ComponentIds = append(result.ReadOnlyComponents.ComponentIds, options.ReadOnlyComponents.ComponentIds...)

		if !result.ReadOnlyComponents.IsAllReadOnly && options.ReadOnlyComponents.IsAllReadOnly {
			result.ReadOnlyComponents.IsAllReadOnly = true
		}

		if !result.isLazy && options.isLazy {
			result.isLazy = true
		}

		if result.TargetWorld == nil && options.TargetWorld != nil {
			result.TargetWorld = options.TargetWorld
		}
	}

	return result, nil
}

func QueryWithOptional[C IComponent](world *World, query Query) error {
	options := query.getOptions()

	componentId := ComponentIdFor[C](world)
	options.OptionalComponents = append(options.OptionalComponents, componentId)

	return query.Validate()
}

func QueryWithReadOnly[C IComponent](world *World, query Query) error {
	options := query.getOptions()

	componentId := ComponentIdFor[C](world)
	options.ReadOnlyComponents.ComponentIds = append(options.ReadOnlyComponents.ComponentIds, componentId)

	return query.Validate()
}

func QueryWithAllReadOnly(query Query) {
	query.getOptions().ReadOnlyComponents.IsAllReadOnly = true
}

func QueryWith[C IComponent](world *World, query Query) error {
	componentId := ComponentIdFor[C](world)
	options := query.getOptions()
	options.Filters = append(options.Filters, &queryFilterWith{c: []ComponentId{componentId}})

	return query.Validate()
}

func QueryWithout[C IComponent](world *World, query Query) error {
	componentId := ComponentIdFor[C](world)
	options := query.getOptions()
	options.Filters = append(options.Filters, &queryFilterWithout{c: []ComponentId{componentId}})

	return query.Validate()
}

func QueryWithFilters[Filters QueryParamFilter](world *World, query Query) error {
	concreteFilters, err := utils.ToConcrete[Filters]()
	if err != nil {
		return fmt.Errorf("failed to cast filter to concrete type: %w", err)
	}

	filter, err := getFilterFromConcreteQueryParamFilter(concreteFilters, world)
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
