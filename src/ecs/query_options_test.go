package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCombinedQuery(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }
	type componentE struct{ Component }
	type componentF struct{ Component }

	t.Run("returns error when a component is specified as optional multiple times", func(t *testing.T) {
		assert := assert.New(t)

		result, err := createCombinedQueryOptions([]queryOption{
			Optional[componentA](),
			Optional[componentB](),
			Optional[componentA](),
		})

		assert.ErrorIs(err, ErrDuplicateComponent)
		assert.Len(result.optionalComponents, 2)
	})

	t.Run("succeeds with all types of options", func(t *testing.T) {
		assert := assert.New(t)

		result, err := createCombinedQueryOptions([]queryOption{
			Optional[componentA](),
			Optional[componentB](),
			With[componentA](),
			Without[componentC](),
			Or(
				With[componentD](),
				And(
					Without[componentE](),
					Without[componentF](),
				),
			),
		})

		assert.NoError(err)
		assert.Len(result.optionalComponents, 2)
		assert.Len(result.filters, 3) // the nested filters are not counted here
	})
}

func TestQueryFilter(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("queryFilterWith only validates if entry has the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := entityData{components: map[componentType]uint{
			getComponentType[componentA](): 0,
		}}

		filter := With[componentA]()
		assert.True(filter.validate(&entityData))

		filter = With[componentB]()
		assert.False(filter.validate(&entityData))
	})

	t.Run("queryFilterWithout only validates if entry does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		entityData := entityData{components: map[componentType]uint{
			getComponentType[componentA](): 0,
		}}

		filter := Without[componentA]()
		assert.False(filter.validate(&entityData))

		filter = Without[componentB]()
		assert.True(filter.validate(&entityData))
	})

	t.Run("queryFilterAnd only validates if both sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := entityData{components: map[componentType]uint{
			getComponentType[componentA](): 0,
		}}

		// both are true
		filter := And(
			With[componentA](),
			Without[componentB](),
		)
		assert.True(filter.validate(&entityData))

		// one is true, 1 is false
		filter = And(
			With[componentA](),
			With[componentB](),
		)
		assert.False(filter.validate(&entityData))

		// both are false
		// one is true, 1 is false
		filter = And(
			With[componentB](),
			With[componentC](),
		)
		assert.False(filter.validate(&entityData))
	})

	t.Run("queryFilterOr returns true if either one or both of the sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		entityData := entityData{components: map[componentType]uint{
			getComponentType[componentA](): 0,
			getComponentType[componentB](): 0,
		}}

		// both are true
		filter := Or(
			With[componentA](),
			With[componentB](),
		)
		assert.True(filter.validate(&entityData))

		// one is true, one is false
		filter = Or(
			With[componentA](),
			With[componentC](),
		)
		assert.True(filter.validate(&entityData))

		// both are false
		filter = Or(
			With[componentC](),
			With[componentD](),
		)
		assert.False(filter.validate(&entityData))
	})
}
