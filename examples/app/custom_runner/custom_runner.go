// Demonstrate how to define a custom runner. This lets you customize when
// the app updates are performed.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const update app.Schedule = "Update"

type customRunner struct {
	app.RunnerBasis
}

// Run systems when pressing enter in the console
func (runner *customRunner) Run(exitChannel <-chan struct{}, executor app.Executor) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Press enter to run systems")
		scanner.Scan()

		select {
		case <-exitChannel:
			return
		default:
		}

		runner.Start()
		executor.Run(runner.CurrentTick())
		runner.Done()
	}
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	runner := customRunner{
		RunnerBasis: app.NewRunnerBasis(myApp),
	}
	myApp.SetRunner(&runner) // <--- Use our custom runner

	myApp.
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddResource(&logger).
		AddSystem(update, func(log app.Logger) {
			log.Info("running system 1!")
		}).
		AddSystem(update, func(log app.Logger) {
			log.Info("running system 2!")
		})

	run.RunApps(myApp)
}
