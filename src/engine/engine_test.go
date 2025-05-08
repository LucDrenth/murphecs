package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEngine(t *testing.T) {
	t.Run("default engine does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		_, err := Default()
		assert.NoError(err)
	})
}
