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

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, NoReadOnly]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[QueryParamFilter, NoOptional, ReadOnlyComponents]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect optional components", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[NoFilter, OptionalComponents, NoReadOnly]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[NoFilter, NoOptional, NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[With[componentA], NoOptional, NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of optional components", func(t *testing.T) {
		assert := assert.New(t)

		var queryOptions iQueryOptions = &QueryOptions[NoFilter, Optional1[componentA], NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))

		queryOptions = &QueryOptions[NoFilter, Optional2[componentA, componentB], NoReadOnly]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(2, len(result.OptionalComponents))
	})

	t.Run("successfully creates combined query options with the right amount of read-only components", func(t *testing.T) {
		assert := assert.New(t)

		var queryOptions iQueryOptions = &QueryOptions[NoFilter, NoOptional, NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, AllReadOnly]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)

		queryOptions = &QueryOptions[NoFilter, NoOptional, ReadOnly1[componentA]]{}
		result, err = queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentTypes))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with all options applied", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[Without[componentB], Optional1[componentA], ReadOnly1[componentA]]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentTypes))
	})

	t.Run("successfully creates the combined query options with an and filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[And[With[componentA], With[componentB]], NoOptional, NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
	})

	t.Run("successfully creates the combined query options with an or filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[Or[With[componentA], With[componentB]], NoOptional, NoReadOnly]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
	})
}
