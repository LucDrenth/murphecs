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

		_, err := getCombinedQueryOptions[QueryParamFilter, NoOptional, NoReadOnly]()
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		_, err := getCombinedQueryOptions[QueryParamFilter, NoOptional, ReadOnlyComponents]()
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect optional components", func(t *testing.T) {
		assert := assert.New(t)

		_, err := getCombinedQueryOptions[NoFilter, OptionalComponents, NoReadOnly]()
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[NoFilter, NoOptional, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[With[componentA], NoOptional, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of optional components", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[NoFilter, Optional1[componentA], NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))

		result, err = getCombinedQueryOptions[NoFilter, Optional2[componentA, componentB], NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(2, len(result.OptionalComponents))
	})

	t.Run("successfully creates combined query options with the right amount of read-only components", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[NoFilter, NoOptional, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)

		result, err = getCombinedQueryOptions[NoFilter, NoOptional, AllReadOnly]()
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)

		result, err = getCombinedQueryOptions[NoFilter, NoOptional, ReadOnly1[componentA]]()
		assert.NoError(err)
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentTypes))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with all options applied", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[Without[componentB], Optional1[componentA], ReadOnly1[componentA]]()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(1, len(result.OptionalComponents))
		assert.Equal(1, len(result.ReadOnlyComponents.ComponentTypes))
	})

	t.Run("successfully creates the combined query options with an and filter", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[And[With[componentA], With[componentB]], NoOptional, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
	})

	t.Run("successfully creates the combined query options with an or filter", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[Or[With[componentA], With[componentB]], NoOptional, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
	})
}
