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

func TestGet(t *testing.T) {
	type componentA struct {
		value int
		Component
	}

	type componentB struct{ Component }
	type componentLikeB struct{ Component }

	type anotherComponent struct{ Component }
	type nonExistingComponent struct{ Component }

	expectedValue := 100

	setup := func(component IComponent) (entityId, *world, *assert.Assertions) {
		assert := assert.New(t)
		world := NewWorld()
		entity, err := world.Spawn(component, anotherComponent{})
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected component", func(t *testing.T) {
		entity, world, assert := setup(componentA{value: expectedValue})

		a, err := Get[componentA](entity, world)
		assert.NoError(err)
		assert.NotNil(a)
		assert.Equal(expectedValue, (*a).value)
	})

	t.Run("returns the expected component", func(t *testing.T) {
		_, world, assert := setup(componentA{value: expectedValue})

		_, err := Get[nonExistingComponent](nonExistingEntity, world)
		assert.Error(err)
	})

	t.Run("returns an error if a component is exactly like the requested component", func(t *testing.T) {
		entity, world, assert := setup(componentB{})

		_, err := Get[componentLikeB](entity, world)
		assert.Error(err)
	})
}

func TestGet2(t *testing.T) {
	type componentA struct {
		value int
		Component
	}
	type componentB struct {
		value int
		Component
	}
	type anotherComponent struct{ Component }
	type nonExistingComponent struct{ Component }

	expectedValueA := 100
	expectedValueB := 101

	setup := func() (entityId, *world, *assert.Assertions) {
		assert := assert.New(t)
		world := NewWorld()
		entity, err := world.Spawn(componentA{value: expectedValueA}, anotherComponent{}, componentB{value: expectedValueB})
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected components regardless of the component order", func(t *testing.T) {
		entity, world, assert := setup()

		a, b, err := Get2[componentA, componentB](entity, world)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.Equal(expectedValueA, (*a).value)
		assert.Equal(expectedValueB, (*b).value)

		// other way around
		b, a, err = Get2[componentB, componentA](entity, world)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.Equal(expectedValueA, (*a).value)
		assert.Equal(expectedValueB, (*b).value)
	})

	t.Run("returns error if a components was not found, regardless of the component order", func(t *testing.T) {
		_, world, assert := setup()
		_, _, err := Get2[nonExistingComponent, componentA](nonExistingEntity, world)
		assert.Error(err)
		_, _, err = Get2[componentA, nonExistingComponent](nonExistingEntity, world)
		assert.Error(err)
	})

	t.Run("returns error if a components was not found, regardless of the component order", func(t *testing.T) {
		_, world, assert := setup()

		_, _, err := Get2[componentA, componentA](nonExistingEntity, world)
		assert.Error(err)
	})
}

func TestGet3(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type nonExistingComponent struct{ Component }

	setup := func() (entityId, *world, *assert.Assertions) {
		assert := assert.New(t)
		world := NewWorld()
		entity, err := world.Spawn(componentA{}, componentB{}, componentC{})
		assert.NoError(err)

		return entity, &world, assert
	}

	t.Run("returns the expected components", func(t *testing.T) {
		entity, world, assert := setup()

		a, b, c, err := Get3[componentA, componentB, componentC](entity, world)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.NotNil(c)
	})

	t.Run("returns an error if any of the given components was not found, regardless of the position of the non-existing component", func(t *testing.T) {
		entity, world, assert := setup()

		_, _, _, err := Get3[nonExistingComponent, componentB, componentC](entity, world)
		assert.Error(err)
		_, _, _, err = Get3[componentA, nonExistingComponent, componentC](entity, world)
		assert.Error(err)
		_, _, _, err = Get3[componentA, componentB, nonExistingComponent](entity, world)
		assert.Error(err)
	})
}

func TestGet4(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }

	t.Run("returns the expected components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()
		entity, err := world.Spawn(componentA{}, componentB{}, componentC{}, componentD{})
		assert.NoError(err)

		a, b, c, d, err := Get4[componentA, componentB, componentC, componentD](entity, &world)
		assert.NoError(err)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.NotNil(c)
		assert.NotNil(d)
	})
}
