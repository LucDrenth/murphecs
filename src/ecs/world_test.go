package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEntity(t *testing.T) {
	assert := assert.New(t)

	world := NewWorld()
	entity1 := world.createEntity()
	entity2 := world.createEntity()

	assert.NotEqual(entity1, entity2)
}

func TestGetComponentRegistry(t *testing.T) {
	type componentA struct{ Component }

	assert := assert.New(t)

	// create component registry if it is not present yet
	world := NewWorld()
	componentRegistry := world.getComponentRegistry(getComponentType[componentA]())
	assert.NotNil(componentRegistry)

	// get the same component registry if it is already present
	componentRegistry2 := world.getComponentRegistry(getComponentType[componentA]())
	assert.NotNil(componentRegistry)
	assert.Equal(componentRegistry, componentRegistry2)
}
