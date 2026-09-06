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

		queryOptions := QueryOptions[QueryParamFilter, NotLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect lazy option", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, IsQueryLazy, DefaultWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect target world", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NotLazy, TargetWorld]{}
		_, err := queryOptions.GetCombinedQueryOptions(world)
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := Default{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Empty(result.Filters)
		assert.Empty(result.OptionalComponents)
		assert.False(result.isLazy)
		assert.Nil(result.TargetWorld)
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[With[componentA], NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Len(result.Filters, 1)
		assert.Empty(result.OptionalComponents)
	})

	t.Run("successfully creates the combined query options with an and filter", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[And[With[componentA], With[componentB]], NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Len(result.Filters, 1)
		assert.Len(result.OptionalComponents, 0)
	})

	t.Run("successfully creates the combined query options with an or filter", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[Or[With[componentA], With[componentB]], NotLazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.Len(result.Filters, 1)
		assert.Len(result.OptionalComponents, 0)
	})

	t.Run("creates a lazy query", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, Lazy, DefaultWorld]{}
		result, err := queryOptions.GetCombinedQueryOptions(world)
		assert.NoError(err)
		assert.True(result.isLazy)
	})

	t.Run("creates a query with a custom target world", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		queryOptions := QueryOptions[NoFilter, NotLazy, TestCustomTargetWorld]{}
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
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if any optional component is not in the query components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		options := CombinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentA](world)},
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
		}
		assert.NoError(options.validateOptions([]ComponentId{ComponentIdFor[componentA](world), ComponentIdFor[componentB](world)}))
	})
}

func TestMergeQueryOptions(t *testing.T) {
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

	t.Run("has the world set to the custom target world", func(t *testing.T) {
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
