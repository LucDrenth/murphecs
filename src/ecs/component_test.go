package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Compiler check to verify that `Component` satisfies `IComponent`
var _ IComponent = Component{}

/*****************************
 * Component without require *
 *****************************/
type componentWithoutRequire struct{ Component }

/***************************************
 * Components that requires each other *
 ***************************************/
type componentThatRequiresEachOtherA struct{ Component }
type componentThatRequiresEachOtherB struct{ Component }

func (a componentThatRequiresEachOtherA) RequiredComponents() []IComponent {
	return []IComponent{componentThatRequiresEachOtherB{}}
}
func (a componentThatRequiresEachOtherB) RequiredComponents() []IComponent {
	return []IComponent{componentThatRequiresEachOtherA{}}
}

/**********************************
 * Component that requires itself *
 **********************************/
type componentThatRequiresItself struct{ Component }

func (a componentThatRequiresItself) RequiredComponents() []IComponent {
	return []IComponent{componentThatRequiresItself{}}
}

/****************************************************
 * Components with a tree structure *
 ****************************************************
 *
 *				    ---> 2A <---
 *			  	  /			   	 \
 *				 /				  \
 *		   ---> 1A ----> 2B -----> 3A
 *		 /			 	  \
 * 0A --			   	    -----> 3B
 *		 \
 *		  \
 *		    --> 1B ----> 2C -----> 3C
 *						/
 *					   /
 * 0B <---------------
 *
 ****************************************************/
type componentTree0A struct{ Component }
type componentTree0B struct{ Component }
type componentTree1A struct{ Component }
type componentTree1B struct{ Component }
type componentTree2A struct{ Component }
type componentTree2B struct{ Component }
type componentTree2C struct{ Component }
type componentTree3A struct{ Component }
type componentTree3B struct{ Component }
type componentTree3C struct{ Component }

func (a componentTree0A) RequiredComponents() []IComponent {
	return []IComponent{componentTree1A{}, componentTree1B{}}
}
func (a componentTree1B) RequiredComponents() []IComponent {
	return []IComponent{componentTree2C{}}
}
func (a componentTree2C) RequiredComponents() []IComponent {
	return []IComponent{componentTree3C{}, componentTree0B{}}
}
func (a componentTree1A) RequiredComponents() []IComponent {
	return []IComponent{componentTree2A{}, componentTree2B{}}
}
func (a componentTree2B) RequiredComponents() []IComponent {
	return []IComponent{componentTree3A{}, componentTree3B{}}
}
func (a componentTree3A) RequiredComponents() []IComponent {
	return []IComponent{componentTree2A{}}
}

func TestGetAllRequiredComponents(t *testing.T) {
	scenarios := []struct {
		description       string
		components        []IComponent
		nrExpectedResults int
	}{
		{
			description:       "handles empty list of components",
			components:        []IComponent{},
			nrExpectedResults: 0,
		},
		{
			description:       "handles component with recursive require",
			components:        []IComponent{componentThatRequiresEachOtherA{}},
			nrExpectedResults: 1,
		},
		{
			description:       "handles components that require each other",
			components:        []IComponent{componentThatRequiresEachOtherA{}, componentThatRequiresEachOtherB{}},
			nrExpectedResults: 0,
		},
		{
			description:       "handles component that requires itself",
			components:        []IComponent{componentThatRequiresItself{}},
			nrExpectedResults: 0,
		},
		{
			description:       "returns empty for component without required components",
			components:        []IComponent{componentWithoutRequire{}},
			nrExpectedResults: 0,
		},
		{
			description:       "handles complex tree of required components",
			components:        []IComponent{componentTree0A{}, componentTree0B{}},
			nrExpectedResults: 8,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			world := NewDefaultWorld()

			typesToExclude := toComponentIds(scenario.components, &world)
			result := getAllRequiredComponents(&typesToExclude, scenario.components, &world)
			assert.Equal(t, scenario.nrExpectedResults, len(result))
		})
	}
}

func TestComponentIdConversions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("test that component IDs differ", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		assert.NotEqual(
			ComponentIdOf(componentA{}, &world),
			ComponentIdOf(componentB{}, &world),
		)
		assert.NotEqual(
			ComponentIdFor[componentA](&world),
			ComponentIdFor[componentB](&world),
		)
	})

	t.Run("toComponentDebugType and getComponentDebugType result in the same type", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(
			ComponentDebugStringOf(componentA{}),
			ComponentDebugStringFor[componentA](),
		)
		assert.NotEqual(
			ComponentDebugStringFor[componentA](),
			ComponentDebugStringFor[componentB](),
		)
	})

	t.Run("toComponentId and getComponentId result in the same type", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		assert.Equal(
			ComponentIdOf(componentA{}, &world),
			ComponentIdFor[componentA](&world),
		)
		assert.NotEqual(
			ComponentIdFor[componentA](&world),
			ComponentIdFor[componentB](&world),
		)
	})

	t.Run("getting type from an IComponent returns the same type as when passing type param", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		var iComponent IComponent = componentA{}

		a := ComponentIdOf(iComponent, &world)
		b := ComponentIdFor[componentA](&world)

		assert.Equal(
			a.DebugString(),
			b.DebugString(),
		)
	})

	t.Run("return the same result for components and component pointers", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		assert.Equal(
			ComponentIdOf(&componentA{}, &world),
			ComponentIdOf(componentA{}, &world),
		)
		assert.Equal(
			ComponentIdFor[componentA](&world),
			ComponentIdFor[*componentA](&world),
		)
	})
}

func TestComponentIdRegistry(t *testing.T) {
	assert := assert.New(t)

	type componentA struct{ Component }
	type componentB struct{ Component }

	componentIdRegistry := componentRegistry{
		components: map[reflect.Type]uint{},
		currentId:  0,
	}

	assert.Equal(uint(1), componentIdRegistry.getId(reflect.TypeFor[componentA]()))
	assert.Equal(uint(1), componentIdRegistry.getId(reflect.TypeFor[componentA]()))
	assert.Equal(uint(2), componentIdRegistry.getId(reflect.TypeFor[componentB]()))
}
