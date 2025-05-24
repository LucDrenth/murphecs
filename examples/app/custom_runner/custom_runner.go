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
	world  *ecs.World
	logger app.Logger
}

// Run systems when pressing enter in the console
func (runner *customRunner) Run(exitChannel <-chan struct{}, systems []*app.SystemSet) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Press enter to run systems")
		scanner.Scan()

		select {
		case <-exitChannel:
			return
		default:
		}

		for _, systemSet := range systems {
			errors := systemSet.Exec(runner.world, nil)
			for _, err := range errors {
				runner.logger.Error(fmt.Sprintf("system returned error: %v", err))
			}
		}
	}
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	runner := customRunner{
		world:  myApp.World(),
		logger: logger,
	}
	myApp.SetRunner(&runner) // <--- Use our custom runner

	myApp.AddSchedule(update, app.ScheduleTypeRepeating)
	myApp.AddResource(&logger)
	myApp.AddSystem(update, func(log app.Logger) {
		log.Info("running system 1!")
	})
	myApp.AddSystem(update, func(log app.Logger) {
		log.Info("running system 2!")
	})

	run.RunApp(&myApp)
}
