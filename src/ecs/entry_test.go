package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComponentFromEntry(t *testing.T) {
	type componentA struct {
		value int
		Component
	}
	const expectedValueA = 101
	type componentB struct {
		value int
		Component
	}
	const expectedValueB = 102

	t.Run("gets the component if its present in entry", func(t *testing.T) {
		assert := assert.New(t)
		entry := entry{components: []IComponent{
			componentA{value: expectedValueA},
			componentB{value: expectedValueB},
		}}

		componentA, _, err := getComponentFromEntry[componentA](&entry)
		assert.NoError(err)
		assert.Equal(expectedValueA, (*componentA).value)

		componentB, _, err := getComponentFromEntry[componentB](&entry)
		assert.NoError(err)
		assert.Equal(expectedValueB, (*componentB).value)
	})

	t.Run("return an error if the entry does not contain the component", func(t *testing.T) {
		assert := assert.New(t)

		entry := entry{components: []IComponent{
			componentA{},
		}}

		componentA, _, err := getComponentFromEntry[componentB](&entry)
		assert.Error(err)
		assert.True(errors.Is(err, ErrComponentNotFound))
		assert.Nil(componentA)
	})
}

func TestEntryContainsComponentType(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	scenarios := []struct {
		description      string
		entry            entry
		componentToCheck IComponent
		expectedResult   bool
	}{
		{
			description:      "returns false for empty entry",
			entry:            entry{components: []IComponent{}},
			componentToCheck: componentA{},
			expectedResult:   false,
		},
		{
			description: "returns false when component not present",
			entry: entry{components: []IComponent{
				componentA{},
			}},
			componentToCheck: componentB{},
			expectedResult:   false,
		},
		{
			description: "returns true if component is present",
			entry: entry{components: []IComponent{
				componentA{}, componentB{},
			}},
			componentToCheck: componentA{},
			expectedResult:   true,
		},
		{
			description: "returns true if component is present",
			entry: entry{components: []IComponent{
				componentA{}, componentB{},
			}},
			componentToCheck: componentB{},
			expectedResult:   true,
		},
		{
			description: "returns true if component is the only component in the list",
			entry: entry{components: []IComponent{
				componentA{},
			}},
			componentToCheck: componentA{},
			expectedResult:   true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			result := scenario.entry.containsComponentType(toComponentType(scenario.componentToCheck))
			assert.Equal(t, scenario.expectedResult, result)
		})
	}
}
