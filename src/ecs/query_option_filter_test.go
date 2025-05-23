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

		world := NewDefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		entityData := world.entities[entity]

		filter := queryFilterWith{c: []ComponentId{
			ComponentIdFor[componentA](&world),
		}}
		assert.True(filter.EntityMeetsCriteria(entityData))

		filter = queryFilterWith{c: []ComponentId{
			ComponentIdFor[componentB](&world),
		}}
		assert.False(filter.EntityMeetsCriteria(entityData))
	})

	t.Run("queryFilterWithout only validates if entry does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		entityData := world.entities[entity]

		filter := queryFilterWithout{c: []ComponentId{
			ComponentIdFor[componentA](&world),
		}}
		assert.False(filter.EntityMeetsCriteria(entityData))

		filter = queryFilterWithout{c: []ComponentId{
			ComponentIdFor[componentB](&world),
		}}
		assert.True(filter.EntityMeetsCriteria(entityData))
	})

	t.Run("queryFilterAnd only validates if both sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		entityData := world.entities[entity]

		// both are true
		filter := queryFilterAnd{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](&world),
			}},
			b: &queryFilterWithout{c: []ComponentId{
				ComponentIdFor[componentB](&world),
			}},
		}
		assert.True(filter.EntityMeetsCriteria(entityData))

		// one is true, 1 is false
		filter = queryFilterAnd{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](&world),
			}},
			b: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](&world),
			}},
		}
		assert.False(filter.EntityMeetsCriteria(entityData))

		// both are false
		filter = queryFilterAnd{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](&world),
			}},
			b: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](&world),
			}},
		}
		assert.False(filter.EntityMeetsCriteria(entityData))
	})

	t.Run("queryFilterOr returns true if either one or both of the sub-filters are true", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		entityData := world.entities[entity]

		// both are true
		filter := queryFilterOr{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](&world),
			}},
			b: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentB](&world),
			}},
		}
		assert.True(filter.EntityMeetsCriteria(entityData))

		// one is true, one is false
		filter = queryFilterOr{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentA](&world),
			}},
			b: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](&world),
			}},
		}
		assert.True(filter.EntityMeetsCriteria(entityData))

		// both are false
		filter = queryFilterOr{
			a: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentC](&world),
			}},
			b: &queryFilterWith{c: []ComponentId{
				ComponentIdFor[componentD](&world),
			}},
		}
		assert.False(filter.EntityMeetsCriteria(entityData))
	})
}
