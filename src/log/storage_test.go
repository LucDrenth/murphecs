package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	assert := assert.New(t)

	storage := NewStorage()
	message1 := "some message"
	message2 := "another message"
	caller1 := "some/place.go:15"
	caller2 := "another/place/in/code.go:1"

	assert.False(storage.Exists(LevelDebug, message1, caller1))

	storage.Insert(LevelDebug, message1, caller1)
	assert.True(storage.Exists(LevelDebug, message1, caller1))
	assert.False(storage.Exists(LevelDebug, message1, caller2))
	assert.False(storage.Exists(LevelDebug, message2, caller1))
	assert.False(storage.Exists(LevelDebug, message2, caller2))
	assert.False(storage.Exists(LevelInfo, message1, caller1))

	storage.Clear()
	assert.False(storage.Exists(LevelDebug, message1, caller1))
}
