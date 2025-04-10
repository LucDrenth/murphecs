package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstDuplicate(t *testing.T) {
	scenarios := []struct {
		description     string
		elements        []string
		expectedElement *string
		expectedIndexA  int
		expectedIndexB  int
	}{
		{
			description:     "returns nil regardless of order",
			elements:        []string{"A", "B", "C"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns nil regardless of order",
			elements:        []string{"A", "C", "B"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns nil regardless of order",
			elements:        []string{"B", "A", "C"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns nil regardless of order",
			elements:        []string{"B", "C", "A"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns nil regardless of order",
			elements:        []string{"C", "A", "B"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns nil regardless of order",
			elements:        []string{"C", "B", "A"},
			expectedElement: nil,
			expectedIndexA:  -1,
			expectedIndexB:  -1,
		},
		{
			description:     "returns duplicate regardless of order",
			elements:        []string{"A", "A", "B"},
			expectedElement: PointerTo("A"),
			expectedIndexA:  0,
			expectedIndexB:  1,
		},
		{
			description:     "returns duplicate regardless of order",
			elements:        []string{"A", "B", "A"},
			expectedElement: PointerTo("A"),
			expectedIndexA:  0,
			expectedIndexB:  2,
		},
		{
			description:     "returns duplicate regardless of order",
			elements:        []string{"B", "A", "A"},
			expectedElement: PointerTo("A"),
			expectedIndexA:  1,
			expectedIndexB:  2,
		},
		{
			description:     "returns the first duplicate if there are multiple duplicates",
			elements:        []string{"A", "A", "B", "B"},
			expectedElement: PointerTo("A"),
			expectedIndexA:  0,
			expectedIndexB:  1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			duplicateElement, indexA, indexB := GetFirstDuplicate(scenario.elements)

			if scenario.expectedElement == nil {
				if duplicateElement != nil {
					assert.Fail(t, "expected nil but got ", *duplicateElement)
				}
			} else {
				assert.Equal(t, *scenario.expectedElement, *duplicateElement)
			}

			assert.Equal(t, scenario.expectedIndexA, indexA)
			assert.Equal(t, scenario.expectedIndexB, indexB)
		})
	}
}
