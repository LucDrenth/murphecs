package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCombinedQueryOptions(t *testing.T) {
	// TODO
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
