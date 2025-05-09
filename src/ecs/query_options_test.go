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

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, NoReadOnly, NotLazy]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, ReadOnlyComponents, NotLazy]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect optional components", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[NoFilter, OptionalComponents, NoReadOnly, NotLazy]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := Default{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.False(result.isLazy)
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[With[componentA], NoOptional, NoReadOnly, NotLazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of optional components", func(t *testing.T) {
		assert := assert.New(t)

		var queryOptions QueryOption = &QueryOptions[NoFilter, Optional1[componentA], NoReadOnly, NotLazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))

		queryOptions = &QueryOptions[NoFilter, Optional2[componentA, componentB], NoReadOnly, NotLazy]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(2, len(result.OptionalComponents))
	})

	t.Run("successfully creates combined query options with the right amount of read-only components", func(t *testing.T) {
		assert := assert.New(t)

		var queryOptions QueryOption = &Default{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, AllReadOnly, NotLazy]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, ReadOnly1[componentA], NotLazy]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentIds))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with all options applied", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[Without[componentB], Optional1[componentA], ReadOnly1[componentA], NotLazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("successfully creates the combined query options with an and filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[And[With[componentA], With[componentB]], NoOptional, NoReadOnly, NotLazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("successfully creates the combined query options with an or filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[Or[With[componentA], With[componentB]], NoOptional, NoReadOnly, NotLazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentIds))
	})

	t.Run("creates a lazy query", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly, Lazy]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.True(result.isLazy)
	})
}

func TestValidateCombinedQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("returns an error if the same component is given multiple times as optional component", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentA](), ComponentIdFor[componentA]()},
			ReadOnlyComponents: combinedReadOnlyComponent{},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if the same component is given multiple times as read-only component", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{ComponentIdFor[componentA](), ComponentIdFor[componentA]()},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if there are read-only components while AllReadOnly is set to true", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds:  []ComponentId{ComponentIdFor[componentA]()},
				IsAllReadOnly: true,
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if any read-only component is not in the query components", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{ComponentIdFor[componentA]()},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns an error if any optional component is not in the query components", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentA]()},
			ReadOnlyComponents: combinedReadOnlyComponent{
				ComponentIds: []ComponentId{},
			},
		}
		assert.Error(options.validateOptions([]ComponentId{}))
	})

	t.Run("returns no error for default query options", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{}
		assert.NoError(options.validateOptions([]ComponentId{ComponentIdFor[componentA](), ComponentIdFor[componentB]()}))
	})

	t.Run("returns no error if there is nothing wrong", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:            []QueryFilter{},
			OptionalComponents: []ComponentId{ComponentIdFor[componentB]()},
			ReadOnlyComponents: combinedReadOnlyComponent{
				IsAllReadOnly: true,
			},
		}
		assert.NoError(options.validateOptions([]ComponentId{ComponentIdFor[componentA](), ComponentIdFor[componentB]()}))
	})
}

func TestMergeQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("has IsAllReadOnly set to true if there is any AllReadOnly option", func(t *testing.T) {
		assert := assert.New(t)

		result, err := mergeQueryOptions([]QueryOption{
			ReadOnly1[componentA]{},
			AllReadOnly{},
			ReadOnly1[componentB]{},
		})
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("has isLazy set to true if there is any Lazy option", func(t *testing.T) {
		assert := assert.New(t)

		result, err := mergeQueryOptions([]QueryOption{
			Lazy{},
			NotLazy{},
		})
		assert.NoError(err)
		assert.True(result.isLazy)
	})
}
