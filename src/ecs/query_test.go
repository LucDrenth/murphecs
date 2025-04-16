package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO much more tests to write here

func TestQuery1(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("not specifying any options results in all entities with the component", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		expected := []entityId{}

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
		resultEntities := []entityId{}

		for entityId, b := range result.Iter() {
			assert.NotNil(b)
			resultEntities = append(resultEntities, entityId)
		}

		assert.NoError(err)

		// check both resultEntities and result.entityIds to ensure that Iter() works as expected
		assert.ElementsMatch(expected, result.entityIds)
		assert.ElementsMatch(expected, resultEntities)
	})
}
