package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Since Get, Get1, Get2 etc. all have the same functionality, we do not extensively test each of them.

// For example, the test "returns an error if a component is exactly like the requested component" in GetTest
// would not change for any of the other functions.
//
// Because if this, the tests will get progressively simpler for each TestGet<number>, until only basic functionality is tested.

type componentA struct {
	value int
	Component
}
type componentB struct {
	value int
	Component
}
type componentC struct {
	value int
	Component
}
type componentD struct {
	value int
	Component
}
type componentE struct {
	value int
	Component
}
type componentF struct {
	value int
	Component
}
type componentG struct {
	value int
	Component
}
type componentH struct {
	value int
	Component
}

const (
	expectedValueA = 101
	expectedValueB = 102
	expectedValueC = 103
	expectedValueD = 104
	expectedValueE = 105
	expectedValueF = 106
	expectedValueG = 107
	expectedValueH = 108
)

func TestGet(t *testing.T) {
	type componentB struct{ Component }
	type componentLikeB struct{ Component }

	type anotherComponent struct{ Component }
	type nonExistingComponent struct{ Component }

	expectedValue := 101

	setup := func(component IComponent) (EntityId, *World, *assert.Assertions) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world, component, &anotherComponent{})
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected component", func(t *testing.T) {
		entity, world, assert := setup(&componentA{value: expectedValue})

		a, err := Get1[componentA](world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValue, (*a).value)
	})

	t.Run("returns an error if the component does not exist on the entity", func(t *testing.T) {
		entity, world, assert := setup(&componentA{value: expectedValue})

		_, err := Get1[nonExistingComponent](world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})

	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		_, world, assert := setup(&componentB{})

		_, err := Get1[componentB](world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if a component is exactly like the requested component", func(t *testing.T) {
		entity, world, assert := setup(&componentB{})

		_, err := Get1[componentLikeB](world, entity)
		assert.Error(err)
	})
}

func TestGet2(t *testing.T) {
	type anotherComponent struct{ Component }
	type nonExistingComponent struct{ Component }

	expectedValueA := 101
	expectedValueB := 102

	setup := func() (EntityId, *World, *assert.Assertions) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world, &componentA{value: expectedValueA}, &anotherComponent{}, &componentB{value: expectedValueB})
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected components regardless of the component order", func(t *testing.T) {
		entity, world, assert := setup()

		a, b, err := Get2[componentA, componentB](world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.Equal(expectedValueA, (*a).value)
		assert.Equal(expectedValueB, (*b).value)

		// other way around
		b, a, err = Get2[componentB, componentA](world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.Equal(expectedValueA, (*a).value)
		assert.Equal(expectedValueB, (*b).value)
	})

	t.Run("returns the expected components even if two of the same components are given", func(t *testing.T) {
		entity, world, assert := setup()

		a, alsoA, err := Get2[componentA, componentA](world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(alsoA)
		assert.Equal(expectedValueA, (*alsoA).value)
	})

	t.Run("returns error if a components was not found, regardless of the component order", func(t *testing.T) {
		entity, world, assert := setup()
		_, _, err := Get2[nonExistingComponent, componentA](world, entity)
		assert.Error(err)
		_, _, err = Get2[componentA, nonExistingComponent](world, entity)
		assert.Error(err)
	})
}

func TestGet3(t *testing.T) {
	type nonExistingComponent struct{ Component }

	setup := func() (EntityId, *World, *assert.Assertions) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
		)
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected components", func(t *testing.T) {
		entity, world, assert := setup()

		a, b, c, err := Get3[componentA, componentB, componentC](world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
	})

	t.Run("returns an error if any of the given components was not found, regardless of the position of the non-existing component", func(t *testing.T) {
		entity, world, assert := setup()

		_, _, _, err := Get3[nonExistingComponent, componentB, componentC](world, entity)
		assert.Error(err)
		_, _, _, err = Get3[componentA, nonExistingComponent, componentC](world, entity)
		assert.Error(err)
		_, _, _, err = Get3[componentA, componentB, nonExistingComponent](world, entity)
		assert.Error(err)
	})
}

func TestGet4(t *testing.T) {
	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
			&componentD{value: expectedValueD},
		)
		assert.NoError(err)

		a, b, c, d, err := Get4[componentA, componentB, componentC, componentD](&world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
		assert.NotNil(d)
		assert.Equal(expectedValueD, (*d).value)
	})
}

func TestGet5(t *testing.T) {
	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
			&componentD{value: expectedValueD},
			&componentE{value: expectedValueE},
		)
		assert.NoError(err)

		a, b, c, d, e, err := Get5[componentA, componentB, componentC, componentD, componentE](&world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
		assert.NotNil(d)
		assert.Equal(expectedValueD, (*d).value)
		assert.NotNil(e)
		assert.Equal(expectedValueE, (*e).value)
	})
}

func TestGet6(t *testing.T) {
	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
			&componentD{value: expectedValueD},
			&componentE{value: expectedValueE},
			&componentF{value: expectedValueF},
		)
		assert.NoError(err)

		a, b, c, d, e, f, err := Get6[componentA, componentB, componentC, componentD, componentE, componentF](&world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
		assert.NotNil(d)
		assert.Equal(expectedValueD, (*d).value)
		assert.NotNil(e)
		assert.Equal(expectedValueE, (*e).value)
		assert.NotNil(f)
		assert.Equal(expectedValueF, (*f).value)
	})
}

func TestGet7(t *testing.T) {
	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
			&componentD{value: expectedValueD},
			&componentE{value: expectedValueE},
			&componentF{value: expectedValueF},
			&componentG{value: expectedValueG},
		)
		assert.NoError(err)

		a, b, c, d, e, f, g, err := Get7[componentA, componentB, componentC, componentD, componentE, componentF, componentG](&world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
		assert.NotNil(d)
		assert.Equal(expectedValueD, (*d).value)
		assert.NotNil(e)
		assert.Equal(expectedValueE, (*e).value)
		assert.NotNil(f)
		assert.Equal(expectedValueF, (*f).value)
		assert.NotNil(g)
		assert.Equal(expectedValueG, (*g).value)
	})
}

func TestGet8(t *testing.T) {
	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()
		entity, err := Spawn(&world,
			&componentA{value: expectedValueA},
			&componentB{value: expectedValueB},
			&componentC{value: expectedValueC},
			&componentD{value: expectedValueD},
			&componentE{value: expectedValueE},
			&componentF{value: expectedValueF},
			&componentG{value: expectedValueG},
			&componentH{value: expectedValueH},
		)
		assert.NoError(err)

		a, b, c, d, e, f, g, h, err := Get8[componentA, componentB, componentC, componentD, componentE, componentF, componentG, componentH](&world, entity)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValueA, (*a).value)
		assert.NotNil(b)
		assert.Equal(expectedValueB, (*b).value)
		assert.NotNil(c)
		assert.Equal(expectedValueC, (*c).value)
		assert.NotNil(d)
		assert.Equal(expectedValueD, (*d).value)
		assert.NotNil(e)
		assert.Equal(expectedValueE, (*e).value)
		assert.NotNil(f)
		assert.Equal(expectedValueF, (*f).value)
		assert.NotNil(g)
		assert.Equal(expectedValueG, (*g).value)
		assert.NotNil(h)
		assert.Equal(expectedValueH, (*h).value)
	})
}
