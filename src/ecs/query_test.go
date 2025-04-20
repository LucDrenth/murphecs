package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery1(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentWithValue struct {
		val int
		Component
	}

	t.Run("returns the right entities", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		expected := []EntityId{}

		// spawn entities with components with different orders to ensure that the results
		// are not spawn-order dependent.
		expectedEntity, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)
		expectedEntity, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)
		_, err = Spawn(&world, &componentA{}, &componentC{})
		assert.NoError(err)
		expectedEntity, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)

		expectedEntity, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		expected = append(expected, expectedEntity)

		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)

		result := Query1[componentB](&world)

		assert.ElementsMatch(expected, result.entityIds)
	})

	t.Run("can iter over query results", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)

		numberOfResults := 0

		queryResult := Query1[componentA](&world)
		queryResult.Iter(func(entityId EntityId, a *componentA) error {
			numberOfResults++
			return nil
		})

		assert.Equal(2, numberOfResults)
	})

	t.Run("iter stops iterating when returning an error", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)

		numberOfIterations := 0

		queryResult := Query1[componentA](&world)
		queryResult.Iter(func(entityId EntityId, a *componentA) error {
			numberOfIterations++
			return errors.New("")
		})

		assert.Equal(1, numberOfIterations)
	})

	t.Run("can range over query results", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)

		numberOfIterations := 0

		queryResult := Query1[componentA](&world)
		for range queryResult.Range() {
			numberOfIterations++
		}

		assert.Equal(2, numberOfIterations)
	})

	t.Run("query result implements QueryResult", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		q := Query1[componentA](&world)
		var result QueryResult = &q
		assert.Equal(uint(0), result.NumberOfResult())

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)

		q = Query1[componentA](&world)
		result = &q
		assert.Equal(uint(1), result.NumberOfResult())
	})

	t.Run("queried components can be mutated", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentWithValue{val: 10})
		assert.NoError(err)

		queryResult := Query1[componentWithValue](&world)

		// mutate with Iter
		queryResult.Iter(func(_ EntityId, component *componentWithValue) error {
			component.val++
			return nil
		})

		// mutate again with Range
		for component := range queryResult.Range() {
			component.val++
		}

		queryResult = Query1[componentWithValue](&world)
		queryResult.Iter(func(_ EntityId, component *componentWithValue) error {
			assert.Equal(12, component.val)
			return nil
		})
	})

	t.Run("components are not nil if not marked as optional", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentA{})
		assert.NoError(err)

		queryResults := Query1[componentA](&world)
		queryResults.Iter(func(_ EntityId, a *componentA) error {
			assert.NotNil(a)
			return nil
		})
	})

	t.Run("components can be nil if marked as optional", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentA{})
		assert.NoError(err)

		numberOfNilResults := 0
		numberOfNotNilResults := 0

		queryResults := Query1[componentA](&world, Optional[componentA]())
		queryResults.Iter(func(_ EntityId, a *componentA) error {
			if a == nil {
				numberOfNilResults++
			} else {
				numberOfNotNilResults++
			}

			return nil
		})

		assert.Equal(2, numberOfNilResults)
		assert.Equal(3, numberOfNotNilResults)
	})
}
