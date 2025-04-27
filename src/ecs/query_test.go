package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery1(t *testing.T) {
	type componentA struct {
		Component
		value int
	}
	type componentB struct {
		Component
		value int
	}

	t.Run("query with default options return the expected results", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue1 := 10
		expectedValue2 := 20
		world := NewWorld()
		query := Query1[componentA, NoFilter, AllRequired]{}
		err := query.PrepareOptions()
		assert.NoError(err)

		query.Exec(&world)
		assert.Equal(uint(0), query.Result().NumberOfResult())

		expectedEntity1, err := Spawn(&world, &componentA{value: expectedValue1}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{value: -10}) // decoy component, we should not get this one in the query results
		assert.NoError(err)
		expectedEntity2, err := Spawn(&world, &componentA{value: expectedValue2}, &componentB{})
		assert.NoError(err)

		query.Exec(&world)
		assert.Equal(uint(2), query.Result().NumberOfResult())
		query.results.Iter(func(entityId EntityId, a *componentA) error {
			if entityId == expectedEntity1 {
				assert.Equal(expectedValue1, a.value)
			} else if entityId == expectedEntity2 {
				assert.Equal(expectedValue2, a.value)
			} else {
				assert.FailNow("returned unexpected entity", entityId)
			}
			return nil
		})

		query.Result().Clear()
		assert.Equal(uint(0), query.results.NumberOfResult())
	})

	t.Run("queried component can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 10
		world := NewWorld()
		query := Query1[componentA, NoFilter, AllRequired]{}
		err := query.PrepareOptions()
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{value: 0}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)

		query.Exec(&world)
		query.results.Iter(func(entityId EntityId, a *componentA) error {
			a.value = expectedValue
			return nil
		})

		query.Exec(&world)
		query.results.Iter(func(entityId EntityId, a *componentA) error {
			assert.Equal(expectedValue, a.value)
			return nil
		})
	})

	t.Run("query results stops iterating when returning an error", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		_, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		query := Query1[componentA, NoFilter, AllRequired]{}
		err = query.PrepareOptions()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(2), query.Result().NumberOfResult())
		numberOfIterations := 0
		query.Result().Iter(func(_ EntityId, _ *componentA) error {
			numberOfIterations++
			return errors.New("oops")
		})

		assert.Equal(1, numberOfIterations)
	})
}

// TODO more tests for Query2, Query3, etc.
