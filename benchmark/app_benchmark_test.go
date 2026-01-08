package benchmark

import (
	"fmt"
	"testing"

	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

func BenchmarkRunApp(b *testing.B) {
	numberOfSystems := []int{250, 1_000, 10_000, 100_000}
	numberOfSChedules := []int{1, 4, 16}

	logger := app.TestLogger{}

	for _, systems := range numberOfSystems {
		for _, schedules := range numberOfSChedules {
			b.Run(fmt.Sprintf("%d-systems-over-%d-schedules", systems, schedules), func(b *testing.B) {
				subApp, err := app.New(&logger, ecs.DefaultWorldConfigs())
				if err != nil {
					b.Fatal(err)
				}
				subApp.UseNTimesRunner(1)

				for i := range schedules {
					schedule := app.Schedule(fmt.Sprintf("schedule.%d", i))
					subApp.AddSchedule(schedule, app.ScheduleOptions{})

					for range systems / schedules {
						subApp.AddSystem(schedule, func() {})
					}
				}

				failOnWarnAndErrLogs(logger, b)

				for b.Loop() {
					subApp.Run(make(chan struct{}), make(chan bool, 1))
				}

				failOnWarnAndErrLogs(logger, b)
			})
		}
	}
}

func failOnWarnAndErrLogs(logger app.TestLogger, b *testing.B) {
	if logger.NumberOfWarnLogs > 0 {
		b.Fatalf("expected 0 warn logs, got %d", logger.NumberOfWarnLogs)
	}
	if logger.NumberOfErrorLogs > 0 {
		b.Fatalf("expected 0 error logs, got %d", logger.NumberOfErrorLogs)
	}
}
