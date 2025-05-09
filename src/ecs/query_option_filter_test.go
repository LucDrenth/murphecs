package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryFilter(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("queryFilterWith only validates if entry has the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentId]uint{
			ComponentIdFor[componentA](): 0,
		}}

		filter := queryFilterWith{c: []ComponentId{
			ComponentIdFor[componentA](),
		}}
		assert.True(filter.Validate(&entityData))

		filter = queryFilterWith{c: []ComponentId{
			ComponentIdFor[componentB](),
		}}
		assert.False(filter.Validate(&entityData))
	})

	t.Run("queryFilterWithout only validates if entry does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentId]uint{
			ComponentIdFor[componentA](): 0,
		}}

		filter := queryFilterWithout{c: []ComponentId{
			ComponentIdFor[componentA](),
		}}
		assert.False(filter.Validate(&entityData))

		filter = queryFilterWithout{c: []ComponentId{
			ComponentIdFor[componentB](),
		}}
		assert.True(filter.Validate(&entityData))
	})

	t.Run("queryFilterAnd only validates if both sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentId]uint{
			ComponentIdFor[componentA](): 0,
		}}

		// both are true
		filter := queryFilterAnd{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](),
			}},
			b: queryFilterWithout{c: []ComponentId{
				ComponentIdFor[componentB](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// one is true, 1 is false
		filter = queryFilterAnd{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](),
			}},
			b: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](),
			}},
		}
		assert.False(filter.Validate(&entityData))

		// both are false
		filter = queryFilterAnd{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](),
			}},
			b: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](),
			}},
		}
		assert.False(filter.Validate(&entityData))
	})

	t.Run("queryFilterOr returns true if either one or both of the sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := EntityData{components: map[ComponentId]uint{
			ComponentIdFor[componentA](): 0,
			ComponentIdFor[componentB](): 0,
		}}

		// both are true
		filter := queryFilterOr{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](),
			}},
			b: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// one is true, one is false
		filter = queryFilterOr{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](),
			}},
			b: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](),
			}},
		}
		assert.True(filter.Validate(&entityData))

		// both are false
		filter = queryFilterOr{
			a: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](),
			}},
			b: queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentD](),
			}},
		}
		assert.False(filter.Validate(&entityData))
	})
}
