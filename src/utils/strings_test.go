package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringStartsWithUpper(t *testing.T) {
	scenarios := []struct {
		description    string
		input          string
		expectedResult bool
	}{
		{
			input:          "",
			expectedResult: false,
		},
		{
			input:          "lower",
			expectedResult: false,
		},
		{
			input:          "lOwer",
			expectedResult: false,
		},
		{
			input:          "loweR",
			expectedResult: false,
		},
		{
			input:          "Upper",
			expectedResult: true,
		},
		{
			input:          "UPPER",
			expectedResult: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run("returns the expected result", func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(scenario.expectedResult, StringStartsWithUpper(scenario.input))
		})
	}
}
