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
	type componentC struct{ Component }

	t.Run("query with default options return the expected results", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue1 := 10
		expectedValue2 := 20
		world := NewWorld()
		query := Query1[componentA, DefaultQueryOptions]{}
		err := query.Prepare()
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

	t.Run("Query1 satisfies Query", func(t *testing.T) {
		var _ Query = &Query1[componentA, DefaultQueryOptions]{}
	})

	t.Run("query with With filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		_, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, With[componentB]]{}
		err = query.Prepare()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(2), query.Result().NumberOfResult())
	})

	t.Run("query with Without filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		_, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, Without[componentB]]{}
		err = query.Prepare()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(1), query.Result().NumberOfResult())
	})

	t.Run("query with AND filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		expected, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, And[With[componentB], With[componentC]]]{}
		err = query.Prepare()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(1), query.Result().NumberOfResult())
		query.results.Iter(func(entityId EntityId, _ *componentA) error {
			assert.Equal(expected, entityId)
			return nil
		})
	})

	t.Run("query with OR filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		_, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, Or[With[componentB], With[componentC]]]{}
		err = query.Prepare()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(2), query.Result().NumberOfResult())
	})

	t.Run("query with With filter and all optional components returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		_, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, QueryOptions[With[componentB], Optional1[componentA], NoReadOnly]]{}
		err = query.Prepare()
		assert.NoError(err)
		query.Exec(&world)

		assert.Equal(uint(4), query.Result().NumberOfResult())
	})

	t.Run("queried component can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 10
		world := NewWorld()
		query := Query1[componentA, QueryOptions[NoFilter, NoOptional, NoReadOnly]]{}
		err := query.Prepare()
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

	t.Run("queried component can not be mutated if is specified as read-only", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 0
		world := NewWorld()
		query := Query1[componentA, AllReadOnly]{}
		err := query.Prepare()
		assert.NoError(err)
		_, err = Spawn(&world, &componentA{value: 0}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(&world, &componentB{})
		assert.NoError(err)

		query.Exec(&world)
		query.results.Iter(func(entityId EntityId, a *componentA) error {
			a.value = 10
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
		query := Query1[componentA, DefaultQueryOptions]{}
		err = query.Prepare()
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

func TestQuery2(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("Query2 satisfies Query", func(t *testing.T) {
		var _ Query = &Query2[componentA, componentB, DefaultQueryOptions]{}
	})
}

func TestQuery3(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("Query3 satisfies Query", func(t *testing.T) {
		var _ Query = &Query3[componentA, componentB, componentC, DefaultQueryOptions]{}
	})
}

func TestQuery4(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }

	t.Run("Query4 satisfies Query", func(t *testing.T) {
		var _ Query = &Query4[componentA, componentB, componentC, componentD, DefaultQueryOptions]{}
	})
}
