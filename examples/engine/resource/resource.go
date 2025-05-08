// Demonstrate how to use app resources
package main

import (
	"fmt"

	engine "github.com/lucdrenth/murph_engine/src/engine"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type counter struct {
	value int
}

func main() {
	e, err := engine.Default()
	if err != nil {
		panic(err)
	}

	// We add the counter resource to the app.
	// We have to pass it as a reference, or it will result in an error.
	e.App(engine.AppIDCore).
		AddResource(&counter{value: 10}).
		AddSystem(schedule.Update, incrementCounter).
		AddSystem(schedule.Update, logCounter)

	e.Run()
}

// Define the counter as a parameter of this system to use it.
// We can set any resource as a parameter as long as its added to the app.
// Note that we use a pointer to counter so that we can mutate it.
func incrementCounter(counter *counter) {
	counter.value++
}

// We define the counter param as a value (not a pointer to the counter), resulting in
// us getting a copy of the actual counter.
// This allows this system to be ran in parallel with other systems.
//
// The log.logger resource is a special resource that is added to every app by default.
// We do not need to add it ourselves.
func logCounter(counter counter, log log.Logger) {
	log.Info(fmt.Sprintf("counter value: %d", counter.value))
}
