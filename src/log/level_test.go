package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowed(t *testing.T) {
	assert := assert.New(t)

	assert.True(LevelDebug.Allows(LevelDebug))
	assert.True(LevelDebug.Allows(LevelInfo))
	assert.True(LevelDebug.Allows(LevelWarn))
	assert.True(LevelDebug.Allows(LevelError))

	assert.False(LevelInfo.Allows(LevelDebug))
	assert.True(LevelInfo.Allows(LevelInfo))
	assert.True(LevelInfo.Allows(LevelWarn))
	assert.True(LevelInfo.Allows(LevelError))

	assert.False(LevelWarn.Allows(LevelDebug))
	assert.False(LevelWarn.Allows(LevelInfo))
	assert.True(LevelWarn.Allows(LevelWarn))
	assert.True(LevelWarn.Allows(LevelError))

	assert.False(LevelError.Allows(LevelDebug))
	assert.False(LevelError.Allows(LevelInfo))
	assert.False(LevelError.Allows(LevelWarn))
	assert.True(LevelError.Allows(LevelError))
}
