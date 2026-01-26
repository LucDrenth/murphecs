package ecs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Compiler check to verify that `Component` satisfies `IComponent`
var _ AnyComponent = Component{}

/*****************************
 * Component without require *
 *****************************/
type componentWithoutRequire struct{ Component }

/***************************************
 * Components that requires each other *
 ***************************************/
type componentThatRequiresEachOtherA struct{ Component }
type componentThatRequiresEachOtherB struct{ Component }

func (a componentThatRequiresEachOtherA) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentThatRequiresEachOtherB{}}
}
func (a componentThatRequiresEachOtherB) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentThatRequiresEachOtherA{}}
}

/**********************************
 * Component that requires itself *
 **********************************/
type componentThatRequiresItself struct{ Component }

func (a componentThatRequiresItself) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentThatRequiresItself{}}
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

func (a componentTree0A) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree1A{}, componentTree1B{}}
}
func (a componentTree1B) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree2C{}}
}
func (a componentTree2C) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree3C{}, componentTree0B{}}
}
func (a componentTree1A) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree2A{}, componentTree2B{}}
}
func (a componentTree2B) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree3A{}, componentTree3B{}}
}
func (a componentTree3A) RequiredComponents() []AnyComponent {
	return []AnyComponent{componentTree2A{}}
}

func TestGetAllRequiredComponents(t *testing.T) {
	scenarios := []struct {
		description       string
		components        []AnyComponent
		nrExpectedResults int
	}{
		{
			description:       "handles empty list of components",
			components:        []AnyComponent{},
			nrExpectedResults: 0,
		},
		{
			description:       "handles component with recursive require",
			components:        []AnyComponent{componentThatRequiresEachOtherA{}},
			nrExpectedResults: 1,
		},
		{
			description:       "handles components that require each other",
			components:        []AnyComponent{componentThatRequiresEachOtherA{}, componentThatRequiresEachOtherB{}},
			nrExpectedResults: 0,
		},
		{
			description:       "handles component that requires itself",
			components:        []AnyComponent{componentThatRequiresItself{}},
			nrExpectedResults: 0,
		},
		{
			description:       "returns empty for component without required components",
			components:        []AnyComponent{componentWithoutRequire{}},
			nrExpectedResults: 0,
		},
		{
			description:       "handles complex tree of required components",
			components:        []AnyComponent{componentTree0A{}, componentTree0B{}},
			nrExpectedResults: 8,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			assert := assert.New(t)

			world := NewDefaultWorld()

			typesToExclude := toComponentIds(scenario.components, world)
			result := getAllRequiredComponents(&typesToExclude, scenario.components, world)
			assert.Len(result, scenario.nrExpectedResults)
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
			ComponentIdOf(componentA{}, world),
			ComponentIdOf(componentB{}, world),
		)
		assert.NotEqual(
			ComponentIdFor[componentA](world),
			ComponentIdFor[componentB](world),
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
			ComponentIdOf(componentA{}, world),
			ComponentIdFor[componentA](world),
		)
		assert.NotEqual(
			ComponentIdFor[componentA](world),
			ComponentIdFor[componentB](world),
		)
	})

	t.Run("getting type from an IComponent returns the same type as when passing type param", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		var iComponent AnyComponent = componentA{}

		a := ComponentIdOf(iComponent, world)
		b := ComponentIdFor[componentA](world)

		assert.Equal(
			a.DebugString(),
			b.DebugString(),
		)
	})

	t.Run("return the same result for components and component pointers", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		assert.Equal(
			ComponentIdOf(&componentA{}, world),
			ComponentIdOf(componentA{}, world),
		)
		assert.Equal(
			ComponentIdFor[componentA](world),
			ComponentIdFor[*componentA](world),
		)
	})
}

func TestComponentIdRegistry(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }
	type componentE struct{ Component }

	t.Run("returns the right ID's", func(t *testing.T) {
		assert := assert.New(t)

		componentIdRegistry := newComponentRegistry()

		check := func() {
			assert.Equal(uint(1), componentIdRegistry.getId(reflect.TypeFor[componentA]()))
			assert.Equal(uint(1), componentIdRegistry.getId(reflect.TypeFor[componentA]()))
			assert.Equal(uint(2), componentIdRegistry.getId(reflect.TypeFor[componentB]()))
		}

		check()
		componentIdRegistry.processComponentIdRegistries()
		componentIdRegistry.concurrencySafeComponents = nil // set to nil to confirm that this map is not used anymore
		check()
	})

	t.Run("handles concurrency", func(t *testing.T) {
		// no asserts needed. Panic occurs when concurrency does not work as expected.

		componentIdRegistry := newComponentRegistry()

		check := func() {
			var waitGroup sync.WaitGroup

			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentA]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentB]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentC]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentD]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentE]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentF]()) })

			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentA]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentB]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentC]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentD]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentE]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentF]()) })

			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentA]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentB]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentC]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentD]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentE]()) })
			waitGroup.Go(func() { componentIdRegistry.getId(reflect.TypeFor[componentF]()) })

			waitGroup.Wait()
		}

		check()
		componentIdRegistry.processComponentIdRegistries()
		componentIdRegistry.concurrencySafeComponents = nil // set to nil to confirm that this map is not used anymore
		check()
	})
}
