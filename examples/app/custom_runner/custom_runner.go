// Demonstrate how to define a custom runner. This lets you customize when
// the app updates are performed.
//
// Support for between-world queries is left out of here for simplicity. Pass
// subApp.OuterWorlds to a runner param and pass it on to the systemSet.Exec
// function to support them.
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
	world          *ecs.World
	eventStorage   *app.EventStorage
	logger         app.Logger
	onFirstRunDone func() // do not set this yourself
	onRunDone      func() // do not set this yourself
	currentTick    *uint  // do not update this yourself
}

// Run systems when pressing enter in the console
func (runner *customRunner) Run(exitChannel <-chan struct{}, systems []*app.SystemSet) {
	scanner := bufio.NewScanner(os.Stdin)
	isFirstRun := true

	for {
		fmt.Print("Press enter to run systems")
		scanner.Scan()

		select {
		case <-exitChannel:
			return
		default:
		}

		for _, systemSet := range systems {
			errors := systemSet.Exec(runner.world, nil, runner.eventStorage, *runner.currentTick)
			for _, err := range errors {
				runner.logger.Error(fmt.Sprintf("system returned error: %v", err))
			}
		}

		if isFirstRun {
			runner.onFirstRunDone()
			isFirstRun = false
		}
		runner.onRunDone()
	}
}

// This method will be called by the SubApp before running the repeated schedules
func (runner *customRunner) SetOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

// This method will be called by the SubApp before running the repeated schedules
func (runner *customRunner) SetOnRunDone(handler func()) {
	runner.onRunDone = handler
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	runner := customRunner{
		world:        myApp.World(),
		eventStorage: myApp.EventStorage(),
		currentTick:  myApp.GetCurrentTick(),
		logger:       logger,
	}
	myApp.SetRunner(&runner) // <--- Use our custom runner

	myApp.
		AddSchedule(update, app.ScheduleTypeRepeating).
		AddResource(&logger).
		AddSystem(update, func(log app.Logger) {
			log.Info("running system 1!")
		}).
		AddSystem(update, func(log app.Logger) {
			log.Info("running system 2!")
		})

	run.RunApps(myApp)
}
