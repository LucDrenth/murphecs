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

		_, err := getCombinedQueryOptions[QueryParamFilter, AllRequired, NoReadOnly]()
		assert.Error(err)
	})

	t.Run("returns an error when passing an incorrect read-only option", func(t *testing.T) {
		assert := assert.New(t)

		_, err := getCombinedQueryOptions[QueryParamFilter, AllRequired, ReadOnlyComponents]()
		assert.Error(err)
	})

	t.Run("returns an error when passing incorrect optional components", func(t *testing.T) {
		assert := assert.New(t)

		_, err := getCombinedQueryOptions[NoFilter, OptionalComponents, NoReadOnly]()
		assert.Error(err)
	})

	t.Run("successfully creates the combined query options with default options", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[NoFilter, AllRequired, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.Filters))
		assert.Equal(0, len(result.OptionalComponents))
	})

	t.Run("successfully creates the combined query options with the right amount of filters", func(t *testing.T) {
		assert := assert.New(t)

		result, err := getCombinedQueryOptions[With[componentA], AllRequired, NoReadOnly]()
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

		result, err := getCombinedQueryOptions[NoFilter, AllRequired, NoReadOnly]()
		assert.NoError(err)
		assert.Equal(0, len(result.ReadOnlyComponents.ComponentTypes))
		assert.False(result.ReadOnlyComponents.IsAllReadOnly)

		result, err = getCombinedQueryOptions[NoFilter, AllRequired, AllReadOnly]()
		assert.NoError(err)
		assert.True(result.ReadOnlyComponents.IsAllReadOnly)

		result, err = getCombinedQueryOptions[NoFilter, AllRequired, ReadOnly1[componentA]]()
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
}

func TestQueryFilter(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("queryFilterWith only validates if entry has the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentType]uint{
			GetComponentType[componentA](): 0,
		}}

		filter := queryFilterWith{c: []ComponentType{
			GetComponentType[componentA](),
		}}
		assert.True(filter.Validate(&entityData))

		filter = queryFilterWith{c: []ComponentType{
			GetComponentType[componentB](),
		}}
		assert.False(filter.Validate(&entityData))
	})

	t.Run("queryFilterWithout only validates if entry does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentType]uint{
			GetComponentType[componentA](): 0,
		}}

		filter := queryFilterWithout{c: []ComponentType{
			GetComponentType[componentA](),
		}}
		assert.False(filter.Validate(&entityData))

		filter = queryFilterWithout{c: []ComponentType{
			GetComponentType[componentB](),
		}}
		assert.True(filter.Validate(&entityData))
	})

	t.Run("queryFilterAnd only validates if both sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentType]uint{
			GetComponentType[componentA](): 0,
		}}

		// both are true
		filter := queryFilterAnd{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentA](),
			}},
			b: queryFilterWithout{c: []ComponentType{
				GetComponentType[componentB](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// one is true, 1 is false
		filter = queryFilterAnd{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentA](),
			}},
			b: queryFilterWith{c: []ComponentType{
				GetComponentType[componentB](),
			}},
		}
		assert.False(filter.Validate(&entityData))

		// both are false
		filter = queryFilterAnd{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentB](),
			}},
			b: queryFilterWith{c: []ComponentType{
				GetComponentType[componentC](),
			}},
		}
		assert.False(filter.Validate(&entityData))
	})

	t.Run("queryFilterOr returns true if either one or both of the sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentType]uint{
			GetComponentType[componentA](): 0,
			GetComponentType[componentB](): 0,
		}}

		// both are true
		filter := queryFilterOr{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentA](),
			}},
			b: queryFilterWith{c: []ComponentType{
				GetComponentType[componentB](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// one is true, one is false
		filter = queryFilterOr{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentA](),
			}},
			b: queryFilterWith{c: []ComponentType{
				GetComponentType[componentC](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// both are false
		filter = queryFilterOr{
			a: queryFilterWith{c: []ComponentType{
				GetComponentType[componentC](),
			}},
			b: queryFilterWith{c: []ComponentType{
				GetComponentType[componentD](),
			}},
		}
		assert.False(filter.Validate(&entityData))
	})
}
