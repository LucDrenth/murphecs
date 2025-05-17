package app

import (
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/stretchr/testify/assert"
)

const testSchedule Schedule = "update"

type testResourceForSubAppA struct{}
type testResourceForSubAppB struct{}

type invalidFeature struct {
	Feature
}

func (f invalidFeature) Init() {
	f.AddResource(&testResourceForSubAppA{})
}

type testFeatureForSubAppA struct {
	Feature
}

func (f *testFeatureForSubAppA) Init() {
	f.
		AddResource(&testResourceForSubAppA{}).
		AddSystem(testSchedule, func() {}).
		AddFeature(&testFeatureForSubAppB{})
}

type testFeatureForSubAppB struct {
	Feature
}

func (f *testFeatureForSubAppB) Init() {
	f.AddResource(&testResourceForSubAppB{})
}

func TestProcessFeatures(t *testing.T) {
	assert := assert.New(t)

	t.Run("handles empty features", func(t *testing.T) {
		logger := testLogger{}
		app, err := NewBasicSubApp(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.processFeatures()
		assert.Equal(uint(0), logger.err)
	})

	t.Run("logs error if a feature its Init method does not have pointer receiver", func(t *testing.T) {
		logger := testLogger{}
		app, err := NewBasicSubApp(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddFeature(&invalidFeature{})
		app.processFeatures()
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(0), app.NumberOfResources())
		assert.Equal(uint(0), app.NumberOfSystems())
	})

	t.Run("adds all resources of the feature and its nested features", func(t *testing.T) {
		logger := testLogger{}
		app, err := NewBasicSubApp(&logger, ecs.DefaultWorldConfigs())
		app.AddSchedule(testSchedule, ScheduleTypeRepeating)

		assert.NoError(err)
		app.AddFeature(&testFeatureForSubAppA{})
		app.processFeatures()
		assert.Equal(uint(0), logger.err)
		assert.Equal(uint(2), app.NumberOfResources())
		assert.Equal(uint(1), app.NumberOfSystems())
	})
}
