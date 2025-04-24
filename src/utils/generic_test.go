package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type interfaceForToConcrete interface {
	getValue() int
}

type structForToConcrete1 struct{}
type structForToConcrete2 struct{}
type nestedStructForToConcrete[t interfaceForToConcrete] struct{}

func (s structForToConcrete1) getValue() int {
	return 1
}

func (s structForToConcrete2) getValue() int {
	return 2
}

func (s nestedStructForToConcrete[T]) getValue() int {
	return 3
}

func TestToConcrete(t *testing.T) {
	t.Run("returns error when trying to make interface a concrete type", func(t *testing.T) {
		assert := assert.New(t)

		v, err := ToConcrete[interfaceForToConcrete]()
		assert.Error(err)
		assert.Nil(v)
	})

	t.Run("returns a concrete type", func(t *testing.T) {
		assert := assert.New(t)

		value1, err := ToConcrete[structForToConcrete1]()
		assert.NoError(err)
		assert.Equal(1, value1.getValue())

		value2, err := ToConcrete[structForToConcrete2]()
		assert.NoError(err)
		assert.Equal(2, value2.getValue())
	})

	t.Run("handles nested concrete types", func(t *testing.T) {
		assert := assert.New(t)

		value, err := ToConcrete[nestedStructForToConcrete[structForToConcrete2]]()
		assert.NoError(err)
		assert.Equal(3, value.getValue())
	})
}
