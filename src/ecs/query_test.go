package ecs

import (
	"testing"

	"github.com/lucdrenth/murph_engine/src/utils"
	"github.com/stretchr/testify/assert"
)

// TODO much more tests to write here

func TestRangeQueryResult1(t *testing.T) {
	type componentA struct{ Component }

	t.Run("range works as expected", func(t *testing.T) {
		assert := assert.New(t)

		expectedEntityIds := []EntityId{
			3,
			10,
		}

		result := query1Result[componentA]{
			[]*componentA{
				utils.PointerTo(componentA{}),
				utils.PointerTo(componentA{}),
			},
			expectedEntityIds,
		}

		entityIdResults := []EntityId{}

		for entityId, component := range result.Range() {
			assert.NotNil(component)
			entityIdResults = append(entityIdResults, entityId)
		}

		assert.ElementsMatch(expectedEntityIds, entityIdResults)
	})
}

func TestQuery1(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("not specifying any options results in all entities with the component", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		expected := []EntityId{}

		expectedEntity, err := Spawn(&world, componentA{}, componentB{}, componentC{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)

		expectedEntity, err = Spawn(&world, componentA{}, componentB{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)

		_, err = Spawn(&world, componentA{}, componentC{})
		assert.NoError(err)

		expectedEntity, err = Spawn(&world, componentB{}, componentC{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)

		_, err = Spawn(&world, componentA{})
		assert.NoError(err)

		expectedEntity, err = Spawn(&world, componentB{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)

		_, err = Spawn(&world, componentC{})
		assert.NoError(err)

		result := Query1[componentB](&world)
		resultEntities := []EntityId{}

		// check that Iter works as expected
		err = result.Iter(func(entityId EntityId, b *componentB) error {
			assert.NotNil(b)
			resultEntities = append(resultEntities, entityId)
			return nil
		})

		assert.NoError(err)

		// check both resultEntities and result.entityIds to ensure that Iter() works as expected
		assert.ElementsMatch(expected, result.entityIds)
		assert.ElementsMatch(expected, resultEntities)
	})
}
