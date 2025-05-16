package ecs

import (
	"testing"

	"github.com/lucdrenth/murphecs/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashComponentIds(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("returns a unique hash for each combination of components", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentHashes := []string{
			hashComponentIds([]ComponentId{}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentB](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentC](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentC](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			}),
		}

		assert.True(utils.IsUnique(componentHashes))
	})

	t.Run("hashing is deterministic", func(t *testing.T) {
		assert := assert.New(t)
		world := DefaultWorld()

		assert.Equal(
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
			}),
		)
		assert.Equal(
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
			}),
		)
		assert.Equal(
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			}),
			hashComponentIds([]ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			}),
		)
		assert.Equal(
			hashComponentIds([]ComponentId{}),
			hashComponentIds([]ComponentId{}),
		)
	})
}

func TestArchetypeIsFromComponents(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	world := DefaultWorld()

	scenarios := []struct {
		description           string
		archeTypeComponentIds []ComponentId
		componentIdsToCheck   []ComponentId
		expected              bool
	}{
		{
			description:           "returns false when archetype does not contain any of given input components",
			archeTypeComponentIds: []ComponentId{},
			componentIdsToCheck: []ComponentId{
				ComponentIdFor[componentA](&world),
			},
			expected: false,
		},
		{
			description: "returns false when archetype has a component that the input does not have",
			archeTypeComponentIds: []ComponentId{
				ComponentIdFor[componentA](&world),
			},
			componentIdsToCheck: []ComponentId{},
			expected:            false,
		},
		{
			description: "archetype components are the same",
			archeTypeComponentIds: []ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			},
			componentIdsToCheck: []ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			},
			expected: true,
		},
		{
			description: "archetype components are the same but in different order",
			archeTypeComponentIds: []ComponentId{
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
				ComponentIdFor[componentC](&world),
			},
			componentIdsToCheck: []ComponentId{
				ComponentIdFor[componentC](&world),
				ComponentIdFor[componentA](&world),
				ComponentIdFor[componentB](&world),
			},
			expected: true,
		},
		{
			description:           "returns true for both no components",
			archeTypeComponentIds: []ComponentId{},
			componentIdsToCheck:   []ComponentId{},
			expected:              true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			assert := assert.New(t)

			world := DefaultWorld()
			archetype, err := newArchetype(&world, scenario.archeTypeComponentIds)
			assert.NoError(err)

			result := archetype.IsFromComponents(scenario.componentIdsToCheck)
			assert.Equal(scenario.expected, result)
		})
	}
}
