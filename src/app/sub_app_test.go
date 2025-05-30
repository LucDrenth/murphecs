package app

import (
	"fmt"
	"strconv"
	"testing"
	"time"

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

func TestAddResource(t *testing.T) {
	type resourceA struct{}

	t.Run("logs an error when resource is not valid", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddResource(&resourceA{})
		app.AddResource(&resourceA{}) // <- resource already added so its not valid
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(1), app.NumberOfResources())
	})

	t.Run("successfully adds a resource", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddResource(&resourceA{})
		assert.Equal(uint(0), logger.err)
		assert.Equal(uint(1), app.NumberOfResources())
	})
}

func TestAddSchedule(t *testing.T) {
	const schedule Schedule = "schedule"

	t.Run("logs an error when the schedule type is not valid", func(t *testing.T) {
		assert := assert.New(t)

		var invalidScheduleType scheduleType = 100
		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddSchedule(schedule, invalidScheduleType)
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(0), app.NumberOfSchedules())
	})

	t.Run("logs an error when schedule not valid", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddSchedule(schedule, ScheduleTypeStartup)
		app.AddSchedule(schedule, ScheduleTypeStartup) // <-- already added, so not valid
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(1), app.NumberOfSchedules())
	})

	t.Run("successfully adds a schedule", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddSchedule(schedule, ScheduleTypeStartup)
		assert.Equal(uint(0), logger.err)
		assert.Equal(uint(1), app.NumberOfSchedules())
	})
}

func TestAddSystemToSubApp(t *testing.T) {
	const schedule Schedule = "schedule"

	t.Run("logs error when schedule does not exist", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddSystem(schedule, func() {})
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(0), app.NumberOfSystems())
	})

	t.Run("logs error when system is not valid", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)
		app.AddSchedule(schedule, ScheduleTypeStartup)

		app.AddSystem(schedule, "a string is not a valid system")
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(0), app.NumberOfSystems())
	})

	t.Run("successfully adds system", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)
		app.AddSchedule(schedule, ScheduleTypeStartup)

		app.AddSystem(schedule, func() {})
		assert.Equal(uint(0), logger.err)
		assert.Equal(uint(1), app.NumberOfSystems())
	})
}

func TestAddFeature(t *testing.T) {
	t.Run("logs error if a feature its Init method does not have pointer receiver", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.AddFeature(&invalidFeature{})
		assert.Equal(uint(1), logger.err)
		assert.Equal(uint(0), app.NumberOfResources())
		assert.Equal(uint(0), app.NumberOfSystems())
	})

	t.Run("adds all resources of the feature and its nested features", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		app.AddSchedule(testSchedule, ScheduleTypeRepeating)

		assert.NoError(err)
		app.AddFeature(&testFeatureForSubAppA{})
		assert.Equal(uint(0), logger.err)
		assert.Equal(uint(2), app.NumberOfResources())
		assert.Equal(uint(1), app.NumberOfSystems())
	})
}

func TestSetRunner(t *testing.T) {
	t.Run("logs an error when passing nil runner", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.SetRunner(nil)
		assert.Equal(uint(1), logger.err)
	})

	t.Run("logs no error when passing a proper runner", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.SetRunner(&fixedRunner{})
		assert.Equal(uint(0), logger.err)
	})
}

func TestRun(t *testing.T) {
	t.Run("runs all systems once and then exists", func(t *testing.T) {
		assert := assert.New(t)

		const (
			startup Schedule = "startup"
			update  Schedule = "startup"
			cleanup Schedule = "cleanup"
		)

		numberOfSystemRuns := 0

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.
			AddSchedule(startup, ScheduleTypeStartup).
			AddSchedule(update, ScheduleTypeRepeating).
			AddSchedule(cleanup, ScheduleTypeCleanup)

		app.
			AddSystem(startup, func() { numberOfSystemRuns++ }).
			AddSystem(update, func() { numberOfSystemRuns++ }).
			AddSystem(cleanup, func() { numberOfSystemRuns++ })

		runner := app.NewNTimesRunner(1)
		app.SetRunner(&runner)

		isDoneChannel := make(chan bool)
		go app.Run(make(chan struct{}), isDoneChannel)
		<-isDoneChannel

		assert.Equal(uint(0), logger.err)
		assert.Equal(3, numberOfSystemRuns)
	})

	t.Run("fixed runner stops when closing exit channel", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		exitChannel := make(chan struct{})
		isDoneChannel := make(chan bool)

		go app.Run(exitChannel, isDoneChannel)
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(exitChannel)
		}()

		<-isDoneChannel
		assert.Equal(uint(0), logger.err)
	})

	t.Run("uncapped runner stops when closing exit channel", func(t *testing.T) {
		assert := assert.New(t)

		logger := testLogger{}
		app, err := New(&logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)

		app.UseUncappedRunner()

		exitChannel := make(chan struct{})
		isDoneChannel := make(chan bool)

		go app.Run(exitChannel, isDoneChannel)

		// simulate an exit signal
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(exitChannel)
		}()

		<-isDoneChannel
		assert.Equal(uint(0), logger.err)
	})
}

func TestConcurrency(t *testing.T) {
	const (
		startup Schedule = "Startup"
		update  Schedule = "Update"
	)

	type component struct {
		ecs.Component
		data map[string]int
	}

	t.Run("between-world query that competes for a resource", func(t *testing.T) {
		// Run two subApps that both compete for `component.data`. When concurrency is not
		// handled, this test will panic.
		// It is not a guarantee that it will panic because of the nature of data races but,
		// but increasing numberOfRuns will increase the likelihood.
		numberOfRuns := 10_000

		assert := assert.New(t)

		logger := &testLogger{}
		worldConfigs := ecs.DefaultWorldConfigs()
		worldConfigs.Id = &ecs.TestCustomTargetWorldId
		subAppA, err := New(logger, worldConfigs)
		assert.NoError(err)
		subAppA.AddSchedule(startup, ScheduleTypeStartup)
		subAppA.AddSchedule(update, ScheduleTypeRepeating)
		runner := subAppA.NewNTimesRunner(numberOfRuns)
		subAppA.SetRunner(&runner)

		subAppB, err := New(logger, ecs.DefaultWorldConfigs())
		assert.NoError(err)
		subAppB.AddSchedule(startup, ScheduleTypeStartup)
		subAppB.AddSchedule(update, ScheduleTypeRepeating)
		err = subAppB.RegisterOuterWorld(ecs.TestCustomTargetWorldId, subAppA.World())
		assert.NoError(err)
		runner = subAppB.NewNTimesRunner(numberOfRuns)
		subAppB.SetRunner(&runner)

		subAppA.
			AddSystem(startup, func(world *ecs.World) error {
				_, err := ecs.Spawn(world, &component{
					data: map[string]int{},
				})
				return err
			}).
			AddSystem(update, func(query *ecs.Query1[component, ecs.Default]) {
				query.Result().Iter(func(entityId ecs.EntityId, a *component) {
					target := len(a.data)
					a.data[strconv.Itoa(target)] = target
				})
			})

		subAppB.
			AddResource(logger).
			AddSystem(update, func(log *testLogger, query *ecs.Query1[component, ecs.TestCustomTargetWorld]) {
				query.Result().Iter(func(entityId ecs.EntityId, a *component) {
					for k, v := range a.data {
						log.Info(fmt.Sprintf("%s=%d\t", k, v))
					}
				})
			})

		exitChannel := make(chan struct{})
		isDoneChannelA := make(chan bool)
		isDoneChannelB := make(chan bool)
		go subAppA.Run(exitChannel, isDoneChannelA)
		go subAppB.Run(exitChannel, isDoneChannelB)
		<-isDoneChannelA
		<-isDoneChannelB

		assert.Equal(uint(0), logger.err)
	})
}
