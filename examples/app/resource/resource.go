// Demonstrate how to use app resources. Resources can be seen
// as singletons that can be used a system params.
package main

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const update app.Schedule = "Update"

type counter interface {
	getValue() int
	increaseCount()
}

type doublingCounter struct {
	value int
}

func (c *doublingCounter) getValue() int {
	return c.value
}

func (c *doublingCounter) increaseCount() {
	c.value *= 2
}

type incrementalCounter struct {
	value int
}

func (c *incrementalCounter) getValue() int {
	return c.value
}

func (c *incrementalCounter) increaseCount() {
	c.value += 1
}

func main() {
	var logger app.Logger = &app.NoOpLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.SetTickRate(time.Second)
	myApp.AddSchedule(update, app.ScheduleTypeRepeating)

	// Add a struct resource to use it in system parameters
	counterA := doublingCounter{value: 1}
	myApp.
		AddResource(&counterA).
		AddSystem(update, func(c *doublingCounter) {
			c.increaseCount()
			fmt.Println("A1: ", c.getValue())
			// --> A1: 2
			// --> A1: 4
			// --> A1: 8
			// ...
		}).
		AddSystem(update, func(c doublingCounter) {
			fmt.Println("A2: ", c.getValue())
			// --> A2: 2
			// --> A2: 4
			// --> A2: 8
			// ...

			// increaseCount does not persist between system runs because we use it by value.
			for range 3 {
				c.increaseCount()
			}
		})

	// Add a resource by an interface reference to use it in system parameters as an interface
	var counterB counter = &incrementalCounter{}
	myApp.
		AddResource(&counterB).
		AddSystem(update, func(c counter) {
			c.increaseCount()
			fmt.Println("B: ", c.getValue())
			// --> B: 1
			// --> B: 2
			// --> B: 3
			// ...
		})

	// Add a resource by an interface value to use it in system parameters as the struct implementation
	var counterC counter = &incrementalCounter{}
	myApp.
		AddResource(counterC).
		AddSystem(update, func(c *incrementalCounter) {
			c.increaseCount()
			fmt.Println("C1: ", c.getValue())
			// --> C1: 1
			// --> C1: 2
			// --> C1: 3
			// ...
		}).
		AddSystem(update, func(c incrementalCounter) {
			fmt.Println("C2: ", c.getValue())
			// --> C2: 1
			// --> C2: 2
			// --> C2: 3
			// ...

			// increaseCount does not persist between system runs because we use it by value.
			for range 3 {
				c.increaseCount()
			}
		})

	myApp.AddSystem(update, func() { fmt.Println() })

	run.RunApps(myApp)
}
