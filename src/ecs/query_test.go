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
		err := query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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
		err = query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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

	t.Run("Uses the target world", func(t *testing.T) {
		assert := assert.New(t)

		world1 := NewDefaultWorld()
		_, err := Spawn(world1, &componentA{}) // distraction component
		assert.NoError(err)
		_, err = Spawn(world1, &componentB{value: 5})
		assert.NoError(err)
		otherWorlds := &map[WorldId]*World{
			TestCustomTargetWorldId: world1,
		}

		world2 := NewDefaultWorld()
		query := Query1[componentB, TestCustomTargetWorld]{}
		err = query.Prepare(world2, otherWorlds)
		assert.NoError(err)
		err = query.Exec(world1)
		assert.NoError(err)

		_, item, err := query.Single()
		assert.NoError(err)
		assert.Equal(5, item.value)
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
		err := query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
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
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery5(t *testing.T) {
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
	type componentE struct {
		Component
		value int
	}

	t.Run("Query5 satisfies Query", func(t *testing.T) {
		var _ Query = &Query5[componentA, componentB, componentC, componentD, componentE, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query5[componentA, componentB, componentC, componentD, componentE, Default]{}
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4}, &componentE{value: 5})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, queryResultComponentE, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)
		assert.Equal(5, queryResultComponentE.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104}, &componentE{value: 105})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery6(t *testing.T) {
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
	type componentE struct {
		Component
		value int
	}
	type componentF struct {
		Component
		value int
	}

	t.Run("Query6 satisfies Query", func(t *testing.T) {
		var _ Query = &Query6[componentA, componentB, componentC, componentD, componentE, componentF, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query6[componentA, componentB, componentC, componentD, componentE, componentF, Default]{}
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4}, &componentE{value: 5}, &componentF{value: 6})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, queryResultComponentE, queryResultComponentF, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)
		assert.Equal(5, queryResultComponentE.value)
		assert.Equal(6, queryResultComponentF.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104}, &componentE{value: 105}, &componentF{value: 106})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery7(t *testing.T) {
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
	type componentE struct {
		Component
		value int
	}
	type componentF struct {
		Component
		value int
	}
	type componentG struct {
		Component
		value int
	}

	t.Run("Query7 satisfies Query", func(t *testing.T) {
		var _ Query = &Query7[componentA, componentB, componentC, componentD, componentE, componentF, componentG, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query7[componentA, componentB, componentC, componentD, componentE, componentF, componentG, Default]{}
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4}, &componentE{value: 5}, &componentF{value: 6}, &componentG{value: 7})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, queryResultComponentE, queryResultComponentF, queryResultComponentG, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)
		assert.Equal(5, queryResultComponentE.value)
		assert.Equal(6, queryResultComponentF.value)
		assert.Equal(7, queryResultComponentG.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104}, &componentE{value: 105}, &componentF{value: 106}, &componentG{value: 107})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery8(t *testing.T) {
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
	type componentE struct {
		Component
		value int
	}
	type componentF struct {
		Component
		value int
	}
	type componentG struct {
		Component
		value int
	}
	type componentH struct {
		Component
		value int
	}

	t.Run("Query8 satisfies Query", func(t *testing.T) {
		var _ Query = &Query8[componentA, componentB, componentC, componentD, componentE, componentF, componentG, componentH, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query8[componentA, componentB, componentC, componentD, componentE, componentF, componentG, componentH, Default]{}
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4}, &componentE{value: 5}, &componentF{value: 6}, &componentG{value: 7}, &componentH{value: 8})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, queryResultComponentE, queryResultComponentF, queryResultComponentG, queryResultComponentH, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)
		assert.Equal(5, queryResultComponentE.value)
		assert.Equal(6, queryResultComponentF.value)
		assert.Equal(7, queryResultComponentG.value)
		assert.Equal(8, queryResultComponentH.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104}, &componentE{value: 105}, &componentF{value: 106}, &componentG{value: 107}, &componentH{value: 108})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}

func TestQuery9(t *testing.T) {
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
	type componentE struct {
		Component
		value int
	}
	type componentF struct {
		Component
		value int
	}
	type componentG struct {
		Component
		value int
	}
	type componentH struct {
		Component
		value int
	}
	type componentI struct {
		Component
		value int
	}

	t.Run("Query9 satisfies Query", func(t *testing.T) {
		var _ Query = &Query9[componentA, componentB, componentC, componentD, componentE, componentF, componentG, componentH, componentI, Default]{}
	})

	t.Run("Single() works as expected", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		query := Query9[componentA, componentB, componentC, componentD, componentE, componentF, componentG, componentH, componentI, Default]{}
		err := query.Prepare(world, nil)
		assert.NoError(err)

		// 0 results
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(0), query.NumberOfResult())
		_, _, _, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)

		// 1 result
		entity, err := Spawn(world, &componentA{value: 1}, &componentB{value: 2}, &componentC{value: 3}, &componentD{value: 4}, &componentE{value: 5}, &componentF{value: 6}, &componentG{value: 7}, &componentH{value: 8}, &componentI{value: 9})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(1), query.NumberOfResult())
		queryResultEntity, queryResultComponentA, queryResultComponentB, queryResultComponentC, queryResultComponentD, queryResultComponentE, queryResultComponentF, queryResultComponentG, queryResultComponentH, queryResultComponentI, err := query.Single()
		assert.NoError(err)
		assert.Equal(entity, queryResultEntity)
		assert.Equal(1, queryResultComponentA.value)
		assert.Equal(2, queryResultComponentB.value)
		assert.Equal(3, queryResultComponentC.value)
		assert.Equal(4, queryResultComponentD.value)
		assert.Equal(5, queryResultComponentE.value)
		assert.Equal(6, queryResultComponentF.value)
		assert.Equal(7, queryResultComponentG.value)
		assert.Equal(8, queryResultComponentH.value)
		assert.Equal(9, queryResultComponentI.value)

		// 2 results
		_, err = Spawn(world, &componentA{value: 101}, &componentB{value: 102}, &componentC{value: 103}, &componentD{value: 104}, &componentE{value: 105}, &componentF{value: 106}, &componentG{value: 107}, &componentH{value: 108}, &componentI{value: 109})
		assert.NoError(err)
		err = query.Exec(world)
		assert.NoError(err)
		assert.Equal(uint(2), query.NumberOfResult())
		_, _, _, _, _, _, _, _, _, _, err = query.Single()
		assert.ErrorIs(err, ErrUnexpectedNumberOfQueryResults)
	})
}
