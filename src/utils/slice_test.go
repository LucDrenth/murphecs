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

func TestIsUnique(t *testing.T) {
	scenarios := []struct {
		description string
		elements    []string
		expected    bool
	}{
		{
			description: "returns true regardless of order",
			elements:    []string{"A", "B", "C"},
			expected:    true,
		},
		{
			description: "returns true regardless of order",
			elements:    []string{"A", "C", "B"},
			expected:    true,
		},
		{
			description: "returns true regardless of order",
			elements:    []string{"B", "A", "C"},
			expected:    true,
		},
		{
			description: "returns true regardless of order",
			elements:    []string{"B", "C", "A"},
			expected:    true,
		},
		{
			description: "returns true regardless of order",
			elements:    []string{"C", "A", "B"},
			expected:    true,
		},
		{
			description: "returns true regardless of order",
			elements:    []string{"C", "B", "A"},
			expected:    true,
		},
		{
			description: "returns false if there is a duplicate, no matter the place",
			elements:    []string{"A", "A", "B"},
			expected:    false,
		},
		{
			description: "returns false if there is a duplicate, no matter the place",
			elements:    []string{"A", "B", "A"},
			expected:    false,
		},
		{
			description: "returns false if there is a duplicate, no matter the place",
			elements:    []string{"B", "A", "A"},
			expected:    false,
		},
		{
			description: "returns false if there are multiple duplicates",
			elements:    []string{"A", "A", "B", "B"},
			expected:    false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(scenario.expected, IsUnique(scenario.elements))
		})
	}
}

func TestRemoveFromSlice(t *testing.T) {
	scenarios := []struct {
		description    string
		slice          []int
		indexToRemove  int
		expectedResult []int
	}{
		{
			description:    "remove the first element",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  0,
			expectedResult: []int{1, 2, 3, 4},
		},
		{
			description:    "remove the last element",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  4,
			expectedResult: []int{0, 1, 2, 3},
		},
		{
			description:    "remove an element from the middle",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  2,
			expectedResult: []int{0, 1, 3, 4},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			RemoveFromSlice(&scenario.slice, scenario.indexToRemove)

			assert.ElementsMatch(t, scenario.expectedResult, scenario.slice)
		})
	}
}

func TestRemoveFromSliceAndMaintainOrder(t *testing.T) {
	scenarios := []struct {
		description    string
		slice          []int
		indexToRemove  int
		expectedResult []int
	}{
		{
			description:    "remove the first element",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  0,
			expectedResult: []int{1, 2, 3, 4},
		},
		{
			description:    "remove the last element",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  4,
			expectedResult: []int{0, 1, 2, 3},
		},
		{
			description:    "remove an element from the middle",
			slice:          []int{0, 1, 2, 3, 4},
			indexToRemove:  2,
			expectedResult: []int{0, 1, 3, 4},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			RemoveFromSliceAndMaintainOrder(&scenario.slice, scenario.indexToRemove)

			assert.Equal(t, scenario.expectedResult, scenario.slice)
		})
	}
}
