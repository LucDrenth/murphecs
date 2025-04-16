package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	type someStruct struct{}
	type anotherStruct struct{}

	t.Run("normal type reflect equals result of TypeOf with same type", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(
			reflect.TypeOf(someStruct{}) == TypeOf[someStruct](),
		)
		assert.False(
			reflect.TypeOf(anotherStruct{}) == TypeOf[someStruct](),
		)
	})
}
