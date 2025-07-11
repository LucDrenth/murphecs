package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery0(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("query with default options return the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query0[With[componentA]]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// assert that exec on empty world query returns no results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())

		// spawn 2 entities that are expected to be returned from the query and 1 decoy entity that should be skipped
		expectedEntity1, err := Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}) // decoy component, we should not get this one in the query results
		assert.NoError(err)
		expectedEntity2, err := Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)

		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())

		foundExpectedEntity1 := false
		foundExpectedEntity2 := false
		query.Iter(func(entityId EntityId) {
			switch entityId {
			case expectedEntity1:
				foundExpectedEntity1 = true
			case expectedEntity2:
				foundExpectedEntity2 = true
			default:
				assert.FailNow("returned unexpected entity", entityId)
			}
		})

		assert.True(foundExpectedEntity1)
		assert.True(foundExpectedEntity2)

		// assert that clearing the query results works as expected
		query.Clear()
		assert.Equal(uint(0), query.NumberOfResult())
	})

	t.Run("Query0 satisfies Query", func(t *testing.T) {
		var _ Query = &Query0[Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query0[Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)

		// 2 results
		_, err = Spawn(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

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
		world := NewDefaultWorld()
		query := Query1[componentA, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())

		expectedEntity1, err := Spawn(world, &componentA{value: expectedValue1}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{value: -10}) // decoy component, we should not get this one in the query results
		assert.NoError(err)
		expectedEntity2, err := Spawn(world, &componentA{value: expectedValue2}, &componentB{})
		assert.NoError(err)

		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		query.Iter(func(entityId EntityId, a *componentA) {
			switch entityId {
			case expectedEntity1:
				assert.Equal(expectedValue1, a.value)
			case expectedEntity2:
				assert.Equal(expectedValue2, a.value)
			default:
				assert.FailNow("returned unexpected entity", entityId)
			}
		})

		query.Clear()
		assert.Equal(uint(0), query.NumberOfResult())
	})

	t.Run("Query1 satisfies Query", func(t *testing.T) {
		var _ Query = &Query1[componentA, Default]{}
	})

	t.Run("query with With filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, With[componentB]]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(2), query.NumberOfResult())
	})

	t.Run("query with Without filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, Without[componentB]]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(1), query.NumberOfResult())
	})

	t.Run("query with AND filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		expected, err := Spawn(world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, And[With[componentB], With[componentC]]]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(1), query.NumberOfResult())
		query.Iter(func(entityId EntityId, _ *componentA) {
			assert.Equal(expected, entityId)
		})
	})

	t.Run("query with OR filter returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, Or[With[componentB], With[componentC]]]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(2), query.NumberOfResult())
	})

	t.Run("query with With filter and all optional components returns the expected results", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentC{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{}, &componentC{})
		assert.NoError(err)

		query := Query1[componentA, QueryOptions[With[componentB], Optional1[componentA], NoReadOnly, NotLazy, DefaultWorld]]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(4), query.NumberOfResult())
	})

	t.Run("queried component can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 10
		world := NewDefaultWorld()
		query := Query1[componentA, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)
		_, err = Spawn(world, &componentA{value: 0}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)

		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId EntityId, a *componentA) {
			a.value = expectedValue
		})

		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId EntityId, a *componentA) {
			assert.Equal(expectedValue, a.value)
		})
	})

	t.Run("queried component can not be mutated if is specified as read-only", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 0
		world := NewDefaultWorld()
		query := Query1[componentA, AllReadOnly]{}
		err := query.Prepare(world)
		assert.NoError(err)
		_, err = Spawn(world, &componentA{value: 0}, &componentB{})
		assert.NoError(err)
		_, err = Spawn(world, &componentB{})
		assert.NoError(err)

		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId EntityId, a *componentA) {
			a.value = 10
		})

		err = query.Exec(world)
		assert.NoError(err)
		query.Iter(func(entityId EntityId, a *componentA) {
			assert.Equal(expectedValue, a.value)
		})
	})

	t.Run("query results stops iterating when returning an error", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &componentA{})
		assert.NoError(err)
		_, err = Spawn(world, &componentA{})
		assert.NoError(err)
		query := Query1[componentA, Default]{}
		err = query.Prepare(world)
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)

		assert.Equal(uint(2), query.NumberOfResult())
		numberOfIterations := 0
		err = query.IterUntil(func(_ EntityId, _ *componentA) error {
			numberOfIterations++
			return errors.New("oops")
		})

		assert.Error(err)
		assert.Equal(1, numberOfIterations)
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query1[componentA, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 3})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(3, queryResultComponentA.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 5})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery2(t *testing.T) {
	type componentA struct {
		Component
		value int
	}
	type componentB struct {
		Component
		value int
	}

	t.Run("Query2 satisfies Query", func(t *testing.T) {
		var _ Query = &Query2[componentA, componentB, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query2[componentA, componentB, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 3}, &componentB{value: 30})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(3, queryResultComponentA.value)
		assert.Equal(30, queryResultComponentB.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 5}, &componentB{value: 50})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery3(t *testing.T) {
	type componentA struct {
		Component
		value int
	}
	type componentB struct {
		Component
		value int
	}
	type componentC struct {
		Component
		value int
	}

	t.Run("Query3 satisfies Query", func(t *testing.T) {
		var _ Query = &Query3[componentA, componentB, componentC, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query3[componentA, componentB, componentC, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 3}, &componentB{value: 30}, &componentC{value: 300})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(3, queryResultComponentA.value)
		assert.Equal(30, queryResultComponentB.value)
		assert.Equal(300, queryResultComponentC.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 5}, &componentB{value: 50}, &componentC{value: 500})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery4(t *testing.T) {
	type componentA struct {
		Component
		value int
	}
	type componentB struct {
		Component
		value int
	}
	type componentC struct {
		Component
		value int
	}
	type componentD struct {
		Component
		value int
	}

	t.Run("Query4 satisfies Query", func(t *testing.T) {
		var _ Query = &Query4[componentA, componentB, componentC, componentD, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query4[componentA, componentB, componentC, componentD, Default]{}
		err := query.Prepare(world)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 3}, &componentB{value: 30}, &componentC{value: 300}, &componentD{value: 3000})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(3, queryResultComponentA.value)
		assert.Equal(30, queryResultComponentB.value)
		assert.Equal(300, queryResultComponentC.value)
		assert.Equal(3000, queryResultComponentD.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 5}, &componentB{value: 50}, &componentC{value: 500}, &componentD{value: 5000})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}
