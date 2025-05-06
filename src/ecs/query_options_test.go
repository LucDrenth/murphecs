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

		queryOptions := QueryOptions[QueryParamFilter, NotAllReadOnly]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := QueryOptions[QueryParamFilter, ReadOnlyComponents]{}
		_, err := queryOptions.getCombinedQueryOptions()
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := Default{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.False(result.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := With[componentA]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
	})

	t.Run("successfully creates the combined query options with all read-only", func(t *testing.T) {
		assert := assert.New(t)

		var queryOptions QueryOption = &AllReadOnly{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.True(result.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with an AND filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := And[With[componentA], With[componentB]]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.False(result.IsAllReadOnly)
	})

	t.Run("successfully creates the combined query options with an OR filter", func(t *testing.T) {
		assert := assert.New(t)

		queryOptions := Or[With[componentA], With[componentB]]{}
		result, err := queryOptions.getCombinedQueryOptions()
		assert.NoError(err)
		assert.Equal(1, len(result.Filters))
		assert.False(result.IsAllReadOnly)
	})
}

func TestValidateCombinedQueryOptions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("returns no error for default query options", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{}
		assert.NoError(options.validateOptions([]QueryComponent{componentA{}, componentB{}}))
	})

	t.Run("returns no error if there is nothing wrong with the query options", func(t *testing.T) {
		assert := assert.New(t)

		options := combinedQueryOptions{
			Filters:       []QueryFilter{queryFilterWith{c: []ComponentType{GetComponentType[componentB]()}}},
			IsAllReadOnly: true,
		}
		assert.NoError(options.validateOptions([]QueryComponent{componentA{}}))
	})
}

func TestMergeQueryOptions(t *testing.T) {
	t.Run("has IsAllReadOnly set to true if there is any AllReadOnly option", func(t *testing.T) {
		assert := assert.New(t)

		result, err := mergeQueryOptions([]QueryOption{
			AllReadOnly{},
		})
		assert.NoError(err)
		assert.True(result.IsAllReadOnly)
	})
}

func TestGetInnerQueryComponent(t *testing.T) {
	type componentA struct{ Component }

	t.Run("component returns itself", func(t *testing.T) {
		assert := assert.New(t)

		result := getInnerQueryComponent(componentA{})
		assert.NotNil(result)
		assert.Equal(
			GetComponentType[componentA](),
			toComponentType(result),
		)
	})

	t.Run("read-only component returns the inner type", func(t *testing.T) {
		assert := assert.New(t)

		result := getInnerQueryComponent(ReadOnly[componentA]{})
		assert.NotNil(result)
		assert.Equal(
			GetComponentType[componentA](),
			toComponentType(result),
		)
	})

	t.Run("optional component returns the inner type", func(t *testing.T) {
		assert := assert.New(t)

		result := getInnerQueryComponent(Optional[componentA]{})
		assert.NotNil(result)
		assert.Equal(
			GetComponentType[componentA](),
			toComponentType(result),
		)
	})

	t.Run("nested read-only + optional component returns the correct type", func(t *testing.T) {
		assert := assert.New(t)

		result := getInnerQueryComponent(ReadOnly[Optional[componentA]]{})
		assert.NotNil(result)
		assert.Equal(
			GetComponentType[componentA](),
			toComponentType(result),
		)
	})
}
