package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstDuplicate(t *testing.T) {
	scenarios := []struct {
		description    string
		elements       []string
		expectedResult *string
	}{
		{
			description:    "returns nil regardless of order",
			elements:       []string{"A", "B", "C"},
			expectedResult: nil,
		},
		{
			description:    "returns nil regardless of order",
			elements:       []string{"A", "C", "B"},
			expectedResult: nil,
		},
		{
			description:    "returns nil regardless of order",
			elements:       []string{"B", "A", "C"},
			expectedResult: nil,
		},
		{
			description:    "returns nil regardless of order",
			elements:       []string{"B", "C", "A"},
			expectedResult: nil,
		},
		{
			description:    "returns nil regardless of order",
			elements:       []string{"C", "A", "B"},
			expectedResult: nil,
		},
		{
			description:    "returns nil regardless of order",
			elements:       []string{"C", "B", "A"},
			expectedResult: nil,
		},
		{
			description:    "returns duplicate regardless of order",
			elements:       []string{"A", "A", "B"},
			expectedResult: PointerTo("A"),
		},
		{
			description:    "returns duplicate regardless of order",
			elements:       []string{"A", "B", "A"},
			expectedResult: PointerTo("A"),
		},
		{
			description:    "returns duplicate regardless of order",
			elements:       []string{"B", "A", "A"},
			expectedResult: PointerTo("A"),
		},
		{
			description:    "returns the first duplicate if there are multiple duplicates",
			elements:       []string{"A", "A", "B", "B"},
			expectedResult: PointerTo("A"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			result := GetFirstDuplicate(scenario.elements)

			if scenario.expectedResult == nil {
				if result != nil {
					assert.Fail(t, "expected nil but got ", *result)
				}
			} else {
				assert.Equal(t, *scenario.expectedResult, *result)
			}
		})
	}
}
