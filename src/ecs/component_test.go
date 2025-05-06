package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Compiler check to verify that `Component` satisfies `IComponent`
var _ IComponent = Component{}

// Compiler check to verify that `Component` satisfies `IComponent`
var _ QueryComponent = Component{}

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
			typesToExclude := toComponentTypes(scenario.components)
			result := getAllRequiredComponents(&typesToExclude, scenario.components)
			assert.Equal(t, scenario.nrExpectedResults, len(result))
		})
	}
}

func TestComponentTypeConversions(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("test that component types differ", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotEqual(
			toComponentType(componentA{}),
			toComponentType(componentB{}),
		)
		assert.NotEqual(
			GetComponentType[componentA](),
			GetComponentType[componentB](),
		)
	})

	t.Run("toComponentDebugType and getComponentDebugType result in the same type", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(
			toComponentDebugType(componentA{}),
			getComponentDebugType[componentA](),
		)
		assert.NotEqual(
			getComponentDebugType[componentA](),
			getComponentDebugType[componentB](),
		)
	})

	t.Run("toComponentType and getComponentType result in the same type", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(
			toComponentType(componentA{}),
			GetComponentType[componentA](),
		)
		assert.NotEqual(
			GetComponentType[componentA](),
			GetComponentType[componentB](),
		)
	})

	t.Run("getting type from an IComponent returns the same type as when passing type param", func(t *testing.T) {
		assert := assert.New(t)

		var iComponent IComponent = componentA{}

		assert.Equal(
			toComponentType(iComponent).String(),
			GetComponentType[componentA]().String(),
		)
	})

	t.Run("return the same result for components and component pointers", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(
			toComponentType(&componentA{}),
			toComponentType(componentA{}),
		)
		assert.Equal(
			GetComponentType[componentA](),
			GetComponentType[*componentA](),
		)
	})
}
