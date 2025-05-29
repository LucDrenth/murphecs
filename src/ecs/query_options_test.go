package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCombinedQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("returns an error when passing an incorrect query param filter", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, NoReadOnly, NotLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, ReadOnlyComponents, NotLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect optional components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, OptionalComponents, NoReadOnly, NotLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect lazy option", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly, IsQueryLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect target world", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly, NotLazy, TargetWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := Default{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.False(result.isLazy)
		assert.Nil(result.TargetWorld)
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[With[componentA], NoOptional, NoReadOnly, NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of optional components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		var queryOptions QueryOption = &QueryOptions[NoFilter, Optional1[componentA], NoReadOnly, NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))

		queryOptions = &QueryOptions[NoFilter, Optional2[componentA, componentB], NoReadOnly, NotLazy, DefaultWorld]{}
		result, err = queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(2, len(result.OptionalComponents))
	})

	t.Run("successfully creates combined query options with the right amount of read-only components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		var queryOptions QueryOption = &Default{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, AllReadOnly, NotLazy, DefaultWorld]{}
		result, err = queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, ReadOnly1[componentA], NotLazy, DefaultWorld]{}
		result, err = queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentIds))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with all options applied", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[Without[componentB], Optional1[componentA], ReadOnly1[componentA], NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("successfully creates the combined query options with an and filter", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[And[With[componentA], With[componentB]], NoOptional, NoReadOnly, NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("successfully creates the combined query options with an or filter", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[Or[With[componentA], With[componentB]], NoOptional, NoReadOnly, NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("creates a lazy query", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly, Lazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.True(result.isLazy)
	})

	t.Run("creates a query with a custom target world", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly, NotLazy, TestCustomTargetWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Equal(TestCustomTargetWorldId, *result.TargetWorld)
	})
}

func TestValidateCombinedQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("returns an error if the same component is given multiple times as optional component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentA](world), ComponentIdFor[componentA](world)},
			ReadOnlyComponents: combinedReadOnlyComponent{},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if the same component is given multiple times as read-only component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{ComponentIdFor[componentA](world), ComponentIdFor[componentA](world)},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if there are read-only components while AllReadOnly is set to true", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds:  []ComponentId{ComponentIdFor[componentA](world)},
				IsAllReadOnly: true,
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if any read-only component is not in the query components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{ComponentIdFor[componentA](world)},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if any optional component is not in the query components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentA](world)},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns no error for default query options", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{}
		assert.NoError(options.validateOptions([]ComponentId{ComponentIdFor[componentA](world), ComponentIdFor[componentB](world)}))
	})

	t.Run("returns no error if there is nothing wrong", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentB](world)},
			ReadOnlyComponents: combinedReadOnlyComponent{
				IsAllReadOnly: true,
			},
		}
		assert.NoError(options.validateOptions([]ComponentId{ComponentIdFor[componentA](world), ComponentIdFor[componentB](world)}))
	})
}

func TestMergeQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("has IsAllReadOnly set to true if there is any AllReadOnly option", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		result, err := mergeQueryOptions([]QueryOption{
			ReadOnly1[componentA]{},
			AllReadOnly{},
			ReadOnly1[componentB]{},
		}, world)
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("has isLazy set to true if there is any Lazy option", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		result, err := mergeQueryOptions([]QueryOption{
			Lazy{},
			NotLazy{},
		}, world)
		assert.NoError(err)
		assert.True(result.isLazy)
	})

	t.Run("has the set to the custom target world", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		result, err := mergeQueryOptions([]QueryOption{
			DefaultWorld{},
			TestCustomTargetWorld{},
		}, world)
		assert.NoError(err)
		assert.Equal(TestCustomTargetWorldId, *result.TargetWorld)

		result, err = mergeQueryOptions([]QueryOption{
			TestCustomTargetWorld{},
			DefaultWorld{},
		}, world)
		assert.NoError(err)
		assert.Equal(TestCustomTargetWorldId, *result.TargetWorld)
	})
}

func TestOptimizeQueryOptions(t *testing.T) {
	t.Run("does not set IsAllReadOnly to true when not all components are individually set to ReadOnly", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentIdA := ComponentIdFor[emptyComponentA](world)
		componentIdB := ComponentIdFor[emptyComponentB](world)
		componentIdC := ComponentIdFor[emptyComponentC](world)

		// same component order
		{
			options := CombinedQueryOptions{
				ReadOnlyComponents: combinedReadOnlyComponent{
					IsAllReadOnly: false,
					ComponentIds:  []ComponentId{componentIdA, componentIdB},
				},
			}
			options.optimize([]ComponentId{componentIdA, componentIdB, componentIdC})
			assert.False(options.ReadOnlyComponents.IsAllReadOnly)
		}

		// different component order
		{
			options := CombinedQueryOptions{
				ReadOnlyComponents: combinedReadOnlyComponent{
					IsAllReadOnly: false,
					ComponentIds:  []ComponentId{componentIdA, componentIdB},
				},
			}
			options.optimize([]ComponentId{componentIdB, componentIdA, componentIdC})
			assert.False(options.ReadOnlyComponents.IsAllReadOnly)
		}
	})

	t.Run("when all components are individually set to ReadOnly, set IsAllReadOnly to true", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentIdA := ComponentIdFor[emptyComponentA](world)
		componentIdB := ComponentIdFor[emptyComponentB](world)

		// same component order
		{
			options := CombinedQueryOptions{
				ReadOnlyComponents: combinedReadOnlyComponent{
					IsAllReadOnly: false,
					ComponentIds:  []ComponentId{componentIdA, componentIdB},
				},
			}
			options.optimize([]ComponentId{componentIdA, componentIdB})
			assert.True(options.ReadOnlyComponents.IsAllReadOnly)
		}

		// different component order
		{
			options := CombinedQueryOptions{
				ReadOnlyComponents: combinedReadOnlyComponent{
					IsAllReadOnly: false,
					ComponentIds:  []ComponentId{componentIdA, componentIdB},
				},
			}
			options.optimize([]ComponentId{componentIdB, componentIdA})
			assert.True(options.ReadOnlyComponents.IsAllReadOnly)
		}
	})
}
